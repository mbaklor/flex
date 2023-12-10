package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func flexPack(ctx *cli.Context) error {
	dev, err := CheckForDevice(ctx)
	if err != nil {
		return ShowHelpAndError(ctx, err)
	}
	dir := ctx.String("directory")
	if dir != "." {
		stat, err := os.Stat(dir)
		if err != nil || !stat.IsDir() {
			return cli.Exit(fmt.Errorf("Can't access package folder, make sure it's here!"), 1)
		}
	}
	color.Green("creating package zip")
	zipFile, err := ZipPackage(dir)
	if err != nil {
		return err
	}
	body, contentType, err := CreatePackageForm(zipFile)
	println(contentType)

	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/Flexa_upload.cgi", dev.Address.String()), body)
	r.Header.Add("Content-Type", contentType)

	color.Green("Sending package to %s\n", dev.Address.String())

	res, err := dev.SendToDevice(r)
	if err != nil {
		return cli.Exit(err, 1)
	}
	color.Green("Successfully send to device! Got reply of %s", res)
	return nil
}

func ZipPackage(dir string) (*os.File, error) {
	packFile, err := os.Create("package.zip")
	if err != nil {
		return nil, err
	}
	// packFile := new(bytes.Buffer)
	zipper := zip.NewWriter(packFile)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		trimPath := path
		if dir != "." {
			trimPath = strings.TrimPrefix(path, dir)
		}
		relPath := strings.Replace(strings.TrimPrefix(trimPath, string(filepath.Separator)), "\\", "/", -1)
		if relPath[0] == byte('.') || strings.Contains(relPath, "/.") {
			return nil
		}
		if relPath == "package.zip" {
			return nil
		}
		fmt.Printf("\t%s\n", relPath)
		zipFile, err := zipper.Create(relPath)
		if err != nil {
			return err
		}
		file, err := ReadFile(path)
		if err != nil {
			return err
		}
		_, err = zipFile.Write(file)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = zipper.Close()
	if err != nil {
		return nil, err
	}
	return packFile, nil
}

func CreateFormZip(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", "application/x-zip-compressed")
	return w.CreatePart(h)
}

func CreatePackageForm(c *os.File) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	content, err := os.ReadFile("package.zip")
	if err != nil {
		return nil, "", err
	}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	packPart, err := CreateFormZip(writer, "uploadfileapp", "package.zip")
	if err != nil {
		return nil, "", err
	}
	// io.Copy(packPart, content)
	packPart.Write(content)
	ver, err := writer.CreateFormField("upload_app_version")
	if err != nil {
		return nil, "", err
	}
	ver.Write([]byte("0.0.1+pack4"))

	return body, writer.FormDataContentType(), nil
}
