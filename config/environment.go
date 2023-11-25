package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"read_files/util"
)

func LoadEnvironment() {
	err := godotenv.Load(".env")
	if err != nil {
		util.CustomLogger("ERROR", fmt.Sprintf(".env: %v", err))
	}
}
