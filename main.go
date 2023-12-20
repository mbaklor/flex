package main

import (
	"bytes"
	"fmt"
	"log"
	"mbaklor/flex/device"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func flexaInit(ctx *cli.Context) error {
	name := ctx.String("name")
	println("init project", name)
	return nil
}

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

func SendToDev(dev device.Device, body *bytes.Buffer, contentType string) error {
	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/Flexa_upload.cgi", dev.Address.String()), body)
	r.Header.Add("Content-Type", contentType)

	res, err := dev.SendToDevice(r)
	if err != nil {
		return err
	}
	color.Green("Successfully sent to device! Got reply of %s", res)
	return nil
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
						Name:     "name",
						Aliases:  []string{"n"},
						Required: true,
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
