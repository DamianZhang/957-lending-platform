package api

import (
	"errors"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/gofiber/fiber/v2"
)

type APIError struct {
	Status  int
	Message string
}

func FromSvcError(err error) APIError {
	var (
		apiError APIError
		svcError service.Error
	)
	if errors.As(err, &svcError) {
		switch svcError.SvcErr() {
		case service.ErrBadRequest:
			apiError.Status = fiber.StatusBadRequest
		case service.ErrInternalFailure:
			apiError.Status = fiber.StatusInternalServerError
		default:
			apiError.Status = fiber.StatusInternalServerError
		}

		apiError.Message = svcError.AppErr().Error()
	}

	return apiError
}
