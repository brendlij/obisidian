package util

import (
	"os"
	"strings"
)

func TailFile(path string, lines int) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	arr := strings.Split(string(b), "\n")
	if len(arr) > lines {
		arr = arr[len(arr)-lines:]
	}
	return strings.Join(arr, "\n"), nil
}
