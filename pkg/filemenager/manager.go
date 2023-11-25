package filemenager

import (
	"os"
	"read_files/util/constants"
)

func CreateDirIfNotExist() error {
	if _, err := os.Stat(constants.TargetDir); os.IsNotExist(err) {
		return os.Mkdir(constants.TargetDir, 0755)
	}
	return nil
}
