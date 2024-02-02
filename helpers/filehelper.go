package helpers

import (
	"os"
)

func LoadFile(fileName string) (string, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
