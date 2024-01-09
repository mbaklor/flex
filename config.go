package main

import (
	"bytes"
	"fmt"
	"github.com/mbaklor/flex/device"
	"mime/multipart"
	"os"

	"github.com/urfave/cli/v2"
)

func flexConfig(ctx *cli.Context) error {
	devs, err := CheckForDevice(ctx)
	if err != nil {
		return ShowHelpAndError(ctx, err)
	}
	if !ctx.Args().Present() {
		return ShowHelpAndError(ctx, fmt.Errorf("Required config file in arguments"))
	}
	confFile := ctx.Args().First()
	config, err := getConfigFile(confFile)
	if err != nil {
		return cli.Exit(err, 1)
	}
	body, contentType, err := CreateConfigForm(config)
	if err != nil {
		return cli.Exit(err, 1)
	}
	err = device.SendToDevs(devs, body, contentType)
	return nil
}

func getConfigFile(filename string) ([]byte, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config error: %v", err)
	}
	return file, nil
}

func CreateConfigForm(file []byte) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, err := writer.CreateFormFile("uploadfileconf", "config.json")
	if err != nil {
		return nil, "", err
	}
	_, err = part.Write(file)
	if err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil

}
