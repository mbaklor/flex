package main

import (
	"fmt"
	"log"
	"mbaklor/flex/device"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func CreateDeviceFlags(flags ...cli.Flag) []cli.Flag {
	deviceFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a", "addr"},
			Usage:   "device IP address",
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u", "user"},
			Usage:   "device username",
			Value:   "admin",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p", "pass"},
			Usage:   "Device password in plain text (don't judge me)",
		},
		&cli.StringSliceFlag{
			Name:    "device-file",
			Aliases: []string{"f"},
			Usage:   "JSON file that contains device IP, username and password, can be used multiple times",
		},
	}
	return append(flags, deviceFlags...)
}

func ShowHelpAndError(ctx *cli.Context, err error) error {
	cli.ShowSubcommandHelp(ctx)
	return cli.Exit(err, 1)
}

func ExitHandler(ctx *cli.Context, err error) {
	e, is := err.(cli.ExitCoder)
	if is {
		color.Red(e.Error())
		cli.OsExiter(e.ExitCode())
	}
}

func CheckForDevice(ctx *cli.Context) ([]device.Device, error) {
	devFile := ctx.StringSlice("device-file")
	devIP := ctx.String("address")
	devCount := len(devFile)
	if devCount == 0 && devIP == "" {
		return nil, fmt.Errorf("Required either device file or IP address")
	}

	devs := make([]device.Device, devCount)
	if devCount != 0 {
		for idx := range devs {
			dev, err := device.NewDeviceFromFile(devFile[idx])
			if err != nil {
				return nil, err
			}
			devs[idx] = dev
		}
	} else {
		devUser := ctx.String("username")
		devPass := ctx.String("password")
		dev, err := device.NewDevice(devIP, devUser, devPass)
		if err != nil {
			return nil, err
		}
		devs = append(devs, dev)
	}
	return devs, nil
}

func main() {

	app := &cli.App{
		Name:            "Flex",
		Version:         "v0.0.1",
		Usage:           "CLI tool for Flexa development",
		ArgsUsage:       " ",
		HideHelpCommand: true,
		ExitErrHandler:  ExitHandler,
		Commands: []*cli.Command{
			{
				Name:            "init",
				Aliases:         []string{"i"},
				Usage:           "init flexa package",
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "name of the flexa app",
					},
					&cli.StringFlag{
						Name:    "app-log",
						Aliases: []string{"l"},
						Usage:   "filename for the log file created for this app",
					},
					&cli.BoolFlag{
						Name:    "web-log",
						Aliases: []string{"w"},
						Usage:   "add a file to show the app log in the web UI",
					},
					&cli.BoolFlag{
						Name:    "git",
						Aliases: []string{"g"},
						Usage:   "initialize a git repository for package",
					},
					&cli.BoolFlag{
						Name:    "confirm",
						Aliases: []string{"y"},
						Usage:   "confirm all defaults, if used must provide name flag as well",
					},
				},
				Action: flexaInit,
			},
			{
				Name:            "package",
				Aliases:         []string{"p"},
				Usage:           "pack and upload current package",
				HideHelpCommand: true,
				Flags: CreateDeviceFlags(
					&cli.StringFlag{
						Name:    "directory",
						Aliases: []string{"d"},
						Usage:   "directory to pack and send [default: .]",
						Value:   ".",
					},
				),
				Action: flexPack,
			},
			{
				Name:            "config",
				Aliases:         []string{"c"},
				Usage:           "send config json to device",
				ArgsUsage:       "config file [eg: config.json]",
				HideHelpCommand: true,
				Flags:           CreateDeviceFlags(),
				Action:          flexConfig,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
