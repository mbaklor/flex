//go:build windows

package main

import (
	"bufio"
	"os"
	"strings"
)

func ReadFile(path string) ([]byte, error) {
	var file []byte
	var err error
	if strings.HasSuffix(path, ".cgi") {
		file, err = ReadFileOnlyLF(path)
	} else {
		file, err = os.ReadFile(path)
	}
	return file, err
}

func ReadFileOnlyLF(path string) ([]byte, error) {
	osfile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewScanner(osfile)
	file := make([]byte, 0, 256)
	for buf.Scan() {
		file = append(file, buf.Bytes()...)
		file = append(file, byte('\n'))
	}
	osfile.Close()
	return file, nil
}
