package api

import (
	"fmt"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/token"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BorrowerHandler struct {
	config          util.Config
	validate        *validator.Validate
	borrowerService service.BorrowerService
	tokenMaker      token.Maker
}

func NewBorrowerHandler(config util.Config, borrowerService service.BorrowerService, tokenMaker token.Maker) *BorrowerHandler {
	return &BorrowerHandler{
		config:          config,
		validate:        validator.New(),
		borrowerService: borrowerService,
		tokenMaker:      tokenMaker,
	}
}

func (handler *BorrowerHandler) Route(app *fiber.App) {
	router := app.Group("/api/v1/borrowers")
	router.Post("/sign_up", handler.SignUp)
	router.Post("/sign_in", handler.SignIn)
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

func (handler *BorrowerHandler) SignIn(ctx *fiber.Ctx) error {
	var req SignInRequest

	if err := ctx.BodyParser(&req); err != nil {
		errMsg := fmt.Sprintf("failed to parse request: %s", err.Error())
		return errorResponse(ctx, fiber.StatusBadRequest, errMsg)
	}

	if err := handler.validate.Struct(req); err != nil {
		errMsg := fmt.Sprintf("invalid request: %s", err.Error())
		return errorResponse(ctx, fiber.StatusBadRequest, errMsg)
	}

	input := &service.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := handler.borrowerService.SignIn(ctx.Context(), input)
	if err != nil {
		apiError := FromServiceError(err)
		return errorResponse(ctx, apiError.StatusCode, apiError.Message)
	}

	refreshToken, err := handler.tokenMaker.CreateToken(
		output.Borrower.Email,
		handler.config.RefreshTokenDuration,
	)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create refresh token: %s", err.Error())
		return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
	}

	rsp := SignInResponse{
		RefreshToken: refreshToken,
	}
	return ctx.Status(fiber.StatusOK).JSON(rsp)
}
