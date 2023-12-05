package main

import (
	"log"
	"os"

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

func CheckForDevice(ctx *cli.Context) (device.Device, error) {
	devFile := ctx.String("device-file")
	devIP := ctx.String("address")
	if devFile == "" && devIP == "" {
		return device.Device{}, cli.Exit("Required either device file or IP address", 1)
	}

	var dev device.Device
	if devFile != "" {
		var err error
		dev, err = device.NewDeviceFromFile(devFile)
		if err != nil {
			return device.Device{}, cli.Exit(err, 1)
		}
	} else {
		devUser := ctx.String("username")
		devPass := ctx.String("password")
		dev = device.NewDevice(devIP, devUser, devPass)
	}
	return dev, nil
}

func flexConfig(ctx *cli.Context) error {
	dev, err := CheckForDevice(ctx)
	if err != nil {
		return err
	}
	if !ctx.Args().Present() {
		return cli.Exit("Required config file in arguments", 1)
	}
	confFile := ctx.Args().First()
	println(confFile)
	println(dev.Address, dev.User, dev.Password)
	return nil
}

func main() {

	app := &cli.App{
		Name:    "Flex",
		Usage:   "CLI tool for Flexa development",
		Version: "v0.0.1",
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
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "send config json to device",
				Flags:   CreateDeviceFlags(),
				Action:  flexConfig,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
