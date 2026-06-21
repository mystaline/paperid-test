// Fiber app instantiate
package app

import (
	"cmp"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mystaline/paperid-test/internal/config"
	"github.com/mystaline/paperid-test/internal/delivery/routes"
	"github.com/mystaline/paperid-test/internal/repository"
	"github.com/mystaline/paperid-test/internal/service"
)

type App struct {
	app *fiber.App
}

func New(appConfig config.AppConfig) *App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	walletRepository := repository.NewWalletRepository(&pgxpool.Pool{})
	disbursementService := service.NewDisbursementService(walletRepository)

	setupRoutes(app, disbursementService)

	return &App{app: app}
}

func setupRoutes(app *fiber.App, disbursementService *service.DisbursementService) {
	routes.SetupRoutes(app, disbursementService)
}

func (a *App) Shutdown() error {
	return a.app.Shutdown()
}

func (a *App) Run() error {
	port := cmp.Or(os.Getenv("PORT"), "8080")
	return a.app.Listen(":" + port)
}
