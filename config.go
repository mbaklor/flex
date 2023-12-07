package main

import (
	"bytes"
	"fmt"
	"mbaklor/flex/device"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

func flexConfig(ctx *cli.Context) error {
	dev, err := CheckForDevice(ctx)
	if err != nil {
		return err
	}
	if !ctx.Args().Present() {
		return cli.Exit("Required config file in arguments", 1)
	}
	confFile := ctx.Args().First()
	config, err := getConfigFile(confFile)
	if err != nil {
		return cli.Exit(err, 1)
	}
	r, err := CreateConfigForm(config, dev)
	if err != nil {
		return cli.Exit(err, 1)
	}
	res, err := dev.SendToDevice(r)
	if err != nil {
		return cli.Exit(err, 1)
	}
	println(res)
	println(dev.Address, dev.User, dev.Password)
	return nil
}

func getConfigFile(filename string) ([]byte, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config error: %v", err)
	}
	return file, nil
}

func CreateConfigForm(file []byte, dev device.Device) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadfileconf", "config.json")
	if err != nil {
		return nil, err
	}
	part.Write(file)
	writer.Close()

	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/Flexa_upload.cgi", dev.Address.String()), body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	return r, err
}
