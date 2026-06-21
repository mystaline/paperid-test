package helper

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type envelope struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func newEnvelope(code int, message string) envelope {
	return envelope{Status: code, Message: message}
}

func (e envelope) withData(data any) envelope {
	e.Data = data
	return e
}

func reply(c fiber.Ctx, statusCode int, data any, message string) error {
	return c.Status(statusCode).JSON(
		newEnvelope(statusCode, message).withData(data),
	)
}

func ReplyOK(c fiber.Ctx, message string, data any) error {
	return reply(c, fiber.StatusOK, data, message)
}

func ReplyError(c fiber.Ctx, e error) error {
	var httpError *Error
	if errors.As(e, &httpError) {
		return reply(c, httpError.Code, nil, httpError.Message)
	}

	return reply(c, fiber.ErrInternalServerError.Code, nil, "Something went wrong with server")
}
