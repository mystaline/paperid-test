// Fiber app instantiate
package app

import (
	"cmp"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/mystaline/paperid-test/internal/config"
)

type App struct {
	app *fiber.App
}

func New(appConfig config.AppConfig) *App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	setupRoutes(app)

	return &App{app: app}
}

func setupRoutes(app *fiber.App) {
	// routes.SetupRoutes(app, )
}

func (a *App) Shutdown() error {
	return a.app.Shutdown()
}

func (a *App) Run() error {
	port := cmp.Or(os.Getenv("PORT"), "8080")
	return a.app.Listen(":" + port)
}
