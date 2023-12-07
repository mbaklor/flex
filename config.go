package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func flexConfig(ctx *cli.Context) error {
	dev, err := CheckForDevice(ctx)
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
	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/Flexa_upload.cgi", dev.Address.String()), body)
	r.Header.Add("Content-Type", contentType)

	color.Green("Sending %s to %s\n", confFile, dev.Address.String())

	res, err := dev.SendToDevice(r)
	if err != nil {
		return cli.Exit(err, 1)
	}
	color.Green("Successfully send to device! Got reply of %s", res)
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
