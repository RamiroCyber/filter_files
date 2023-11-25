package main

import (
	"fmt"
	"log"
	"os"
	"read_files/config"
	"read_files/router"
	"read_files/util"
)

func init() {
	config.LoadEnvironment()
}

func main() {
	app := router.InitializeRoutes()
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("port_application"))); err != nil {
		util.CustomLogger("ERROR", fmt.Sprintf("Listen: %v", err))
		log.Panicf("Falha ao iniciar o servidor : %v", err)
	}

}
