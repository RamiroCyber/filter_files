package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"read_files/util"
	"read_files/util/constants"
)

func LoadEnvironment() {
	err := godotenv.Load(".env")
	if err != nil {
		util.CustomLogger(constants.Error, fmt.Sprintf(".env: %v", err))
	}
}
