package api

import (
	"errors"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/gofiber/fiber/v2"
)

type APIError struct {
	StatusCode int
	Message    string
}

func FromServiceError(err error) APIError {
	var (
		apiError APIError
		svcError service.Error
	)
	if errors.As(err, &svcError) {
		switch svcError.SvcErr() {
		case service.ErrInternalFailure:
			apiError.StatusCode = fiber.StatusInternalServerError
		}

		apiError.Message = svcError.Error()
	}

	return apiError
}

func errorResponse(ctx *fiber.Ctx, statusCode int, errMsg string) error {
	return ctx.Status(statusCode).JSON(ErrorResponse{Message: errMsg})
}
