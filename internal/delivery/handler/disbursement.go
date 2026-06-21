package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mystaline/paperid-test/internal/dto"
	"github.com/mystaline/paperid-test/internal/service"
	"github.com/mystaline/paperid-test/pkg/helper"
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
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	body := dto.DisbursementRequest{}
	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	userId, err := strconv.ParseInt(body.UserID, 10, 64)
	if err != nil {
		return helper.ReplyError(c, helper.ErrBadRequest("Invalid user id"))
	}

	param := service.DisbursementParam{
		Ctx:    ctx,
		Amount: body.Amount,
		UserID: userId,
	}

	res := &dto.DisbursementResponse{}
	res, err = h.DisbursementService.Invoke(param)
	if err != nil {
		return helper.ReplyError(c, err)
	}

	return helper.ReplyOK(c, "Successfully disburse balance", res)
}
