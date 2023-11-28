package api

import (
	"fmt"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BorrowerHandler struct {
	validate        *validator.Validate
	borrowerService service.BorrowerService
}

func NewBorrowerHandler(borrowerService service.BorrowerService) *BorrowerHandler {
	return &BorrowerHandler{
		validate:        validator.New(),
		borrowerService: borrowerService,
	}
}

func (handler *BorrowerHandler) Route(app *fiber.App) {
	router := app.Group("/api/v1/borrowers")
	router.Post("/sign_up", handler.SignUp)
}

func (handler *BorrowerHandler) SignUp(ctx *fiber.Ctx) error {
	var req SignUpRequest

	if err := ctx.BodyParser(&req); err != nil {
		errMsg := fmt.Sprintf("failed to parse request: %s", err.Error())
		return errorResponse(ctx, fiber.StatusBadRequest, errMsg)
	}

	if err := handler.validate.Struct(req); err != nil {
		errMsg := fmt.Sprintf("invalid request: %s", err.Error())
		return errorResponse(ctx, fiber.StatusBadRequest, errMsg)
	}

	input := &service.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
		LineID:   req.LineID,
		Nickname: req.Nickname,
	}

	output, err := handler.borrowerService.SignUp(ctx.Context(), input)
	if err != nil {
		apiError := FromServiceError(err)
		return errorResponse(ctx, apiError.StatusCode, apiError.Message)
	}

	rsp := SignUpResponse{
		Email:    output.Borrower.Email,
		LineID:   output.Borrower.LineID,
		Nickname: output.Borrower.Nickname,
		Role:     output.Borrower.Role,
	}
	return ctx.Status(fiber.StatusCreated).JSON(rsp)
}
