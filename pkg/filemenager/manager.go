package filemenager

import (
	"os"
)

func CreateDirIfNotExist(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.Mkdir(filePath, 0755)
	}
	return nil
}
