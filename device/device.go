package device

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/icholy/digest"
)

type Device struct {
	Address  net.IP `json:"address"`
	User     string `json:"username"`
	Password string `json:"password"`
}

func NewDevice(address, user, pass string) (Device, error) {
	addr := net.ParseIP(address)
	if addr == nil {
		return Device{}, fmt.Errorf("creating new device: invalid IP address")
	}
	if user == "" {
		user = "admin"
	}
	return Device{addr, user, pass}, nil
}

func NewDeviceFromFile(filename string) (Device, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return Device{}, err
	}
	var dev struct {
		Address  string `json:"address"`
		User     string `json:"username"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(file, &dev)
	if err != nil {
		return Device{}, err
	}
	device, err := NewDevice(dev.Address, dev.User, dev.Password)
	if err != nil {
		return device, err
	}
	return device, nil
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
