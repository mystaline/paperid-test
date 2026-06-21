package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mystaline/paperid-test/internal/delivery/handler"
	"github.com/mystaline/paperid-test/internal/service"
)

func SetupRoutes(app *fiber.App, disbursementService *service.DisbursementService) {
	h := handler.NewDisbursementHandler(disbursementService)

	app.Post(
		"/disburse",
		h.HandleDisbursement,
	)
}
