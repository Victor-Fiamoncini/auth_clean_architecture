package config

import "github.com/gofiber/fiber/v2"

// NewApp func
func NewApp() *fiber.App {
	LoadEnv()
	CreateSchemas()

	app := fiber.New()

	SetupApp(app)
	SetupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	return app
}
