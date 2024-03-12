package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/mbaklor/flex/device"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func flexPack(ctx *cli.Context) error {
	bundle := ctx.Bool("bundle")
	var devs []device.Device
	if !bundle {
		var err error
		devs, err = CheckForDevice(ctx)
		if err != nil {
			b, e := prompt.New().
				Ask("No device supplied, do you want to bundle without sending to device?").
				Choose(
					[]string{"Yes", "No"},
					choose.WithTheme(choose.ThemeLine),
					choose.WithKeyMap(choose.HorizontalKeyMap),
				)
			if e != nil {
				return e
			}
			if b == "Yes" {
				bundle = true
			} else {
				return ShowHelpAndError(ctx, err)
			}
		}
	}
	dir, err := GetPackageDir(ctx)
	if err != nil {
		return cli.Exit(err, 1)
	}
	manifest, err := GetManifest(dir)
	if err != nil {
		return cli.Exit(err, 1)
	}
	color.Green("creating package zip for %v - %v", manifest.Name, manifest.GetVersionString())
	err = ZipPackage(dir)
	if err != nil {
		return cli.Exit(err, 1)
	}
	body, contentType, err := CreatePackageForm(manifest.GetVersionString())
	if !bundle {
		err = device.SendToDevs(devs, body, contentType)
	}
	return nil
}

func GetPackageDir(ctx *cli.Context) (string, error) {
	dir := ctx.String("directory")
	if dir != "." {
		stat, err := os.Stat(dir)
		if err != nil || !stat.IsDir() {
			return "", fmt.Errorf("Can't access package folder, make sure it's here!")
		}
	}
	return dir, nil
}

func getRelPath(rootDir, path string) string {
	trimPath := path
	if rootDir != "." {
		trimPath = strings.TrimPrefix(path, rootDir)
	}
	relPath := strings.Replace(trimPath, "\\", "/", -1)
	relPath = strings.TrimPrefix(relPath, "/")
	if relPath[0] == byte('.') || strings.Contains(relPath, "/.") {
		return ""
	}
	if relPath == "package.zip" {
		return ""
	}
	return relPath
}

func writeFileToZip(zipper *zip.Writer, path, relPath string, info os.FileInfo) error {
	fmt.Printf("\tadding: %s\n", relPath)
	fh, _ := zip.FileInfoHeader(info)
	fh.Method = zip.Deflate
	fh.Name = relPath
	zipFile, err := zipper.CreateHeader(fh)
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
}

func ZipPackage(dir string) error {
	packFile, err := os.Create("package.zip")
	if err != nil {
		return err
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
		relPath := getRelPath(dir, path)
		if relPath == "" {
			return nil
		}
		err = writeFileToZip(zipper, path, relPath, info)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = zipper.Close()
	if err != nil {
		return err
	}
	return nil
}

func CreateFormZip(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", "application/x-zip-compressed")
	return w.CreatePart(h)
}

func CreatePackageForm(version string) (*bytes.Buffer, string, error) {
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
	ver.Write([]byte(version))

	return body, writer.FormDataContentType(), nil
}
