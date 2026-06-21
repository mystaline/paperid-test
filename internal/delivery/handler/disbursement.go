package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/mystaline/paperid-test/internal/dto"
	"github.com/mystaline/paperid-test/internal/service"
	"github.com/mystaline/paperid-test/pkg/interfaces"
)

type DisbursementHandler struct {
	DisbursementService interfaces.Service[service.DisbursementParam, *dto.DisbursementResponse]
}

func NewDisbursementHandler(
	disbursementService interfaces.Service[service.DisbursementParam, *dto.DisbursementResponse],
) *DisbursementHandler {
	return &DisbursementHandler{
		DisbursementService: disbursementService,
	}
}

func (h *DisbursementHandler) HandleDisbursement(c fiber.Ctx) (err error) {
	body := dto.DisbursementRequest{}
	if err := c.Bind().Body(&body); err != nil {
		return err
	}
	// log.Default().Printf("Received disbursement request: %+v", body)

	userId, err := strconv.ParseInt(body.UserID, 10, 64)
	if err != nil {
		// log.Default().Printf("Error parsing userId: %v", err)
		return c.Status(400).JSON(map[string]any{
			"status": fiber.ErrBadRequest.Code,
			"error":  "Invalid userId parameter",
		})
	}

	param := service.DisbursementParam{
		Amount: body.Amount,
		UserID: userId,
	}

	res := &dto.DisbursementResponse{}
	res, err = h.DisbursementService.Invoke(param)
	if err != nil {
		return c.Status(500).JSON(map[string]any{
			"status": fiber.ErrInternalServerError.Code,
			"error":  err.Error(),
		})
	}

	return c.Status(200).JSON(res)
}
