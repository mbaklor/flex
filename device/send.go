package device

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func SendToDev(dev Device, body *bytes.Buffer, contentType string) error {
	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cgi-bin/Flexa_upload.cgi", dev.Address.String()), body)
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", contentType)

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("cyan")
	s.Suffix = color.GreenString(" Sending to %s", dev.Address.String())
	s.Start()
	res, err := dev.SendToDevice(r)
	s.Stop()
	if err != nil {
		color.Red("Failed to send to %s: %v", dev.Address, err)
	} else {
		color.Green("Successfully sent to %s! Got reply of %s", dev.Address.String(), res)
	}
	return nil
}

func SendToDevs(devs []Device, body *bytes.Buffer, contentType string) error {
	for _, dev := range devs {
		err := SendToDev(dev, body, contentType)
		if err != nil {
			return err
		}
	}
	return nil
}
