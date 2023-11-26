package api

import (
	"fmt"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/gofiber/fiber/v2"
)

type BorrowerHandler struct {
	borrowerService service.BorrowerService
}

func NewBorrowerHandler(borrowerService service.BorrowerService) *BorrowerHandler {
	return &BorrowerHandler{
		borrowerService: borrowerService,
	}
}

func (handler *BorrowerHandler) Route(app *fiber.App) {
	router := app.Group("/api/v1/borrowers")
	router.Post("/sign_up", handler.SignUp)
}

func (handler *BorrowerHandler) SignUp(ctx *fiber.Ctx) error {
	var (
		req service.SignUpRequest
		err = ctx.BodyParser(&req)
	)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Message: fmt.Sprintf("failed to parse request: %s", err.Error()),
		})
	}

	rsp, err := handler.borrowerService.SignUp(ctx.Context(), &req)
	if err != nil {
		apiError := FromSvcError(err)
		return ctx.Status(apiError.Status).JSON(ErrorResponse{
			Message: apiError.Message,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(GeneralResponse{
		Data: rsp,
	})
}
