package device

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/fatih/color"
)

func SendToDev(dev Device, body *bytes.Buffer, contentType string) error {
	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/Flexa_upload.cgi", dev.Address.String()), body)
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", contentType)

	res, err := dev.SendToDevice(r)
	if err != nil {
		color.Red("Failed to send to %s: %v", dev.Address, err)
	} else {
		color.Green("Successfully sent to device! Got reply of %s", res)
	}
	return nil
}

func SendToDevs(devs []Device, body *bytes.Buffer, contentType string) error {
	for _, dev := range devs {
		color.Green("Sending to %s\n", dev.Address.String())
		err := SendToDev(dev, body, contentType)
		if err != nil {
			return err
		}
	}
	return nil
}
