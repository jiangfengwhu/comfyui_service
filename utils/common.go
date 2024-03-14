package utils

import (
	"errors"
	"github.com/google/uuid"
	"os"
)

func GetUUID() string {
	return uuid.New().String()
}
func CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(err, os.ErrNotExist)
}
