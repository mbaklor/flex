//go:build !windows

package main

import (
	"os"
)

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
