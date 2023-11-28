package config

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"read_files/util/constants"
	"strings"
)

func ConfigsRoutes() *fiber.App {
	c := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
		}, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		BodyLimit:   10 * 1024 * 1024,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/zip")
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", constants.FileName))
		return c.Next()
	})

	app.Use(c)

	return app
}
