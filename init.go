package main

import (
	"fmt"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func getName(ctx *cli.Context, confirm bool) (string, error) {
	name := ctx.String("name")
	if name == "" {
		if confirm {
			return "", ShowHelpAndError(ctx, fmt.Errorf("can't use confirm flag without name flag"))
		}
		inp, err := prompt.New().
			Ask("Project name:").
			Input("")
		if err != nil {
			return "", err
		}
		if inp != "" {
			name = inp
		} else {
			return "", fmt.Errorf("No project name provided\nHow do you expect this to work??")
		}
	}
	return name, nil
}

func getLog(ctx *cli.Context, confirm bool) (string, error) {
	logfile := ctx.String("app-log")
	if logfile == "" {
		if confirm {
			return "app_log.log", nil
		}
		inp, err := prompt.New().
			Ask("App log filename:").
			Input("app_log.log")
		if err != nil {
			return "", err
		}
		if inp != "" {
			logfile = inp
		} else {
			logfile = "app_log.log"
		}
	}
	return logfile, nil
}

func getWeb(ctx *cli.Context, confirm bool) (bool, error) {
	isWeb := ctx.Bool("web-log")
	if !isWeb {
		if confirm {
			return false, nil
		}
		inp, err := prompt.New().
			Ask("Show app log page in web UI menu?").
			Choose(
				[]string{"Yes", "No"},
				choose.WithTheme(choose.ThemeLine),
				choose.WithKeyMap(choose.HorizontalKeyMap),
			)
		if err != nil {
			return false, err
		}
		if inp == "Yes" {
			isWeb = true
		}
	}
	return isWeb, nil
}

func getGit(ctx *cli.Context, confirm bool) (bool, error) {
	isGit := ctx.Bool("git")
	if !isGit {
		if confirm {
			return false, nil
		}
		inp, err := prompt.New().
			Ask("Initialize git repository in project?").
			Choose(
				[]string{"Yes", "No"},
				choose.WithTheme(choose.ThemeLine),
				choose.WithKeyMap(choose.HorizontalKeyMap),
			)
		if err != nil {
			return false, err
		}
		if inp == "Yes" {
			isGit = true
		}
	}
	return isGit, nil
}

func flexaInit(ctx *cli.Context) error {
	confirm := ctx.Bool("confirm")
	name, err := getName(ctx, confirm)
	if err != nil {
		return cli.Exit(err, 1)
	}
	logfile, err := getLog(ctx, confirm)
	if err != nil {
		return cli.Exit(err, 1)
	}
	isWeb, err := getWeb(ctx, confirm)
	if err != nil {
		return cli.Exit(err, 1)
	}
	isGit, err := getGit(ctx, confirm)
	if err != nil {
		return cli.Exit(err, 1)
	}

	color.Green("Initializing project: %s with log file %s", name, logfile)
	println(isWeb, isGit)
	return nil
}
