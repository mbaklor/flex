package main

import (
	"bytes"
	"fmt"
	"io"
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
	conf, err := getConfigFile()
	defer conf.Close()
	if err != nil {
		return err
	}
	println(confFile)
	println(dev.Address, dev.User, dev.Password)
	return nil
}

func getConfigFile() (*os.File, error) {
	filename := "config.json"
	file, err := os.Open(filename)
	if err != nil {
		return nil, cli.Exit(err, 1)
	}
	return file, nil
}

func CreateConfigForm(file *os.File, dev device.Device) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadfileconf", "config.json")
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)
	writer.Close()

	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/Flexa_upload.cgi", dev.Address), body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	return r, err
}
