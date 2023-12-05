package device

import (
	"net/http"

	"github.com/icholy/digest"
)

type Device struct {
	Address  string
	User     string
	Password string
}

func NewDevice(address, user, pass string) Device {
	if user == "" {
		user = "admin"
	}
	return Device{address, user, pass}
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
