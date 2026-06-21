// Fiber app instantiate
package app

import (
	"cmp"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mystaline/paperid-test/internal/config"
	"github.com/mystaline/paperid-test/internal/delivery/routes"
	"github.com/mystaline/paperid-test/internal/repository"
	"github.com/mystaline/paperid-test/internal/service"
	"github.com/mystaline/paperid-test/pkg/db"
)

type App struct {
	app  *fiber.App
	pool *pgxpool.Pool
}

func New(appConfig config.AppConfig) *App {
	db.InitIDGenerator()

	pool, err := db.NewPool(appConfig)
	if err != nil {
		log.Fatalf("Failed to instantiate database pool connection: %v", err)
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	transactionManager := repository.NewTransactionManager(pool)

	transactionLogRepository := repository.NewTransactionLogRepository()
	userRepository := repository.NewUserRepository()
	walletRepository := repository.NewWalletRepository()
	disbursementService := service.NewDisbursementService(
		transactionManager,
		transactionLogRepository,
		userRepository,
		walletRepository,
	)

	setupRoutes(app, disbursementService)

	return &App{app: app, pool: pool}
}

func setupRoutes(app *fiber.App, disbursementService *service.DisbursementService) {
	routes.SetupRoutes(app, disbursementService)
}

func (a *App) Shutdown() error {
	a.pool.Close()
	return a.app.Shutdown()
}

func (a *App) Run() error {
	port := cmp.Or(os.Getenv("PORT"), "8080")
	return a.app.Listen(":" + port)
}
