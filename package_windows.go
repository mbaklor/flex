//go:build windows

package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func ReadFile(path string) ([]byte, error) {
	var file []byte
	var err error
	file, err = os.ReadFile(path)
	if utf8.ValidString(string(file)) {
		fileStr := string(file)
		if strings.IndexRune(fileStr, '\r') > -1 {
			fmt.Println("\t\tConverting CRLF to LF")
			fileStr = strings.Join(strings.Split(fileStr, "\r"), "")
			file = []byte(fileStr)
		}
	}
	return file, err
}
