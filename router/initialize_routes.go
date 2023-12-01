package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"read_files/config"
	"read_files/router/handler"
	"read_files/util/constants"
)

func InitializeRoutes() *fiber.App {
	app := config.ConfigsRoutes()

	api := app.Group(constants.API)

	v1 := api.Group(fmt.Sprint("/", constants.V1), func(c *fiber.Ctx) error {
		c.Set(constants.VERSION, constants.V1)
		return c.Next()
	})

	v1.Get("/health", handler.HealthCheck)

	v1.Post("/upload", handler.Upload)

	return app
}
