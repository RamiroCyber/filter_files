package main

import (
	"fmt"
	"log"
	"read_files/config"
	"read_files/router"
	"read_files/util"
	"read_files/util/constants"
)

func init() {
	config.LoadEnvironment()
}

func main() {
	app := router.InitializeRoutes()
	if err := app.Listen(":3000"); err != nil {
		util.CustomLogger(constants.Error, fmt.Sprintf("Listen: %v", err))
		log.Panicf("Falha ao iniciar o servidor : %v", err)
	}

}