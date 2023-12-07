package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"mbaklor/flex/device"
)

func flexaInit(ctx *cli.Context) error {
	name := ctx.String("name")
	println("init project", name)
	return nil
}

func CreateDeviceFlags(flags ...cli.Flag) []cli.Flag {
	deviceFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "address",
			Aliases:  []string{"a", "addr"},
			Usage:    "device IP address",
			Category: "Device Info",
		},
		&cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u", "user"},
			Usage:    "device username",
			Category: "Device Info",
			Value:    "admin",
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p", "pass"},
			Usage:    "Device password in plain text (don't judge me)",
			Category: "Device Info",
		},
		&cli.StringFlag{
			Name:     "device-file",
			Aliases:  []string{"f"},
			Usage:    "name of json file that contains device IP, username and password",
			Category: "Device Info",
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

func CheckForDevice(ctx *cli.Context) (device.Device, error) {
	devFile := ctx.String("device-file")
	devIP := ctx.String("address")
	if devFile == "" && devIP == "" {
		return device.Device{}, fmt.Errorf("Required either device file or IP address")
	}

	var dev device.Device
	var err error
	if devFile != "" {
		dev, err = device.NewDeviceFromFile(devFile)
		if err != nil {
			return device.Device{}, err
		}
	} else {
		devUser := ctx.String("username")
		devPass := ctx.String("password")
		dev, err = device.NewDevice(devIP, devUser, devPass)
		if err != nil {
			return device.Device{}, err
		}
	}
	return dev, nil
}

func main() {

	app := &cli.App{
		Name:            "Flex",
		Usage:           "CLI tool for Flexa development",
		Version:         "v0.0.1",
		ArgsUsage:       " ",
		ExitErrHandler:  ExitHandler,
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "init flexa package",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Required: true,
					},
				},
				Action: flexaInit,
			},
			{
				Name:    "package",
				Aliases: []string{"p"},
				Usage:   "pack and upload current package",
			},
			{
				Name:      "config",
				Aliases:   []string{"c"},
				Usage:     "send config json to device",
				ArgsUsage: "config file [eg: config.json]",
				Flags:     CreateDeviceFlags(),
				Action:    flexConfig,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
