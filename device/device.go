package device

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/icholy/digest"
)

type Device struct {
	Address  string `json:"address"`
	User     string `json:"username"`
	Password string `json:"password"`
}

func NewDevice(address, user, pass string) Device {
	if user == "" {
		user = "admin"
	}
	return Device{address, user, pass}
}

func NewDeviceFromFile(filename string) (Device, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return Device{}, err
	}
	var devices Device
	err = json.Unmarshal(file, &devices)
	return devices, err
}

func (d Device) SendToDevice(r *http.Request) (string, error) {
	client := &http.Client{
		Transport: &digest.Transport{
			Username: d.User,
			Password: d.Password,
		},
	}
	res, err := client.Do(r)
	return res.Status, err
}
