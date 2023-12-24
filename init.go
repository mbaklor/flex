package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

//go:embed template
var Template embed.FS

func flexaInit(ctx *cli.Context) error {
	init, err := GetInitInfo(ctx)

	color.Green("Initializing project: %s", init.Name)
	err = CreateInitDir(init.Name)
	if err != nil {
		return cli.Exit(err, 1)
	}
	err = walkTemplate(init.IsWeb)
	if err != nil {
		color.Yellow("Error encountered, cleaning up project directory")
		os.Chdir("..")
		e := os.Remove(init.Name)
		if e != nil {
			return cli.Exit(e, 1)
		}
		return cli.Exit(err, 1)
	}
	err = WriteManifest(init)
	if err != nil {
		color.Yellow("Error encountered, cleaning up project directory")
		os.Chdir("..")
		e := os.Remove(init.Name)
		if e != nil {
			return cli.Exit(e, 1)
		}
		return cli.Exit(err, 1)
	}
	if init.IsGit {
		color.Green("Initializing git repository in project folder")
		err = InitGit()
		if err != nil {
			return cli.Exit(err, 1)
		}
	}
	color.Green("Created project %s successfully", init.Name)
	return nil
}

func CreateInitDir(name string) error {
	_, err := os.Stat(name)
	if !os.IsNotExist(err) {
		return fmt.Errorf("path '%s' exists, can't create flexa package", name)
	}

	err = os.Mkdir(name, 755)
	if err != nil {
		return err
	}
	fmt.Printf("\tCreated project folder in '%s'\n", name)
	err = os.Chdir(name)
	if err != nil {
		return err
	}
	return nil
}

func walkTemplate(web bool) error {
	dir := "template"
	err := fs.WalkDir(Template, dir, func(path string, d fs.DirEntry, err error) error {
		if path == dir {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := getRelPath(dir, path)
		if relPath == "" {
			return nil
		}
		if d.IsDir() {
			return os.Mkdir(relPath, 755)
		}
		if relPath == "web_ui/menu.json" {
			if !web {
				return nil
			}
		}
		err = writeFileFromTemplate(path, relPath)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func writeFileFromTemplate(path, relPath string) error {
	fmt.Printf("\tcreating: %s\n", relPath)
	file, err := os.Create(relPath)
	if err != nil {
		return err
	}
	fileBytes, err := Template.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		return err
	}
	return nil
}

func WriteManifest(init initInfo) error {
	manifest := CreateManifest(init.Name, init.AppLog)
	fileBytes, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	file, err := os.Create("manifest.json")
	if err != nil {
		return err
	}
	_, err = file.Write(fileBytes)
	if err != nil {
		return err
	}
	return nil
}

func InitGit() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("can't find git binary, make sure you have git installed on this system")
	}
	err = exec.Command("git", "init").Run()
	if err != nil {
		return err
	}
	return nil
}
