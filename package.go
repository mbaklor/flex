package main

import (
	"archive/zip"
	"fmt"
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
	err = ZipPackage(dir)
	if err != nil {
		return err
	}
	println(dev.Address.String())
	return nil
}

func ZipPackage(dir string) error {
	packFile, err := os.Create("package.zip")
	if err != nil {
		return err
	}
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
		fmt.Printf("\t%s\n", relPath)
		zipFile, err := zipper.Create(relPath)
		if err != nil {
			return err
		}
		file, err := os.ReadFile(path)
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
		return err
	}

	err = zipper.Close()
	if err != nil {
		return err
	}
	return nil
}
