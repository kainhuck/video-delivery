package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ExistFile(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func TrimFilename(filePath string) string {
	return strings.Split(filepath.Base(filePath), ".")[0]
}
