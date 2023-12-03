package api

import (
	"fmt"

	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/token"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	router.Get("/refresh_token", handler.RefreshToken)
	router.Get("/:id", authMiddleware(handler.tokenMaker), handler.GetBorrower)
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

	signInInput := &service.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	}

	signInOutput, err := handler.borrowerService.SignIn(ctx.Context(), signInInput)
	if err != nil {
		apiError := FromServiceError(err)
		return errorResponse(ctx, apiError.StatusCode, apiError.Message)
	}

	refreshToken, payload, err := handler.tokenMaker.CreateToken(
		signInOutput.Borrower.ID.String(),
		handler.config.RefreshTokenDuration,
	)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create refresh token: %s", err.Error())
		return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
	}

	createSessionInput := &service.CreateSessionInput{
		SessionID: payload.ID.String(),
		ExpiresAt: payload.ExpiresAt,
		Email:     signInOutput.Borrower.Email,
	}

	_, err = handler.borrowerService.CreateSession(ctx.Context(), createSessionInput)
	if err != nil {
		apiError := FromServiceError(err)
		return errorResponse(ctx, apiError.StatusCode, apiError.Message)
	}

	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  payload.ExpiresAt,
		HTTPOnly: true,
	}
	ctx.Cookie(cookie)

	rsp := SignInResponse{
		RefreshToken: refreshToken,
	}
	return ctx.Status(fiber.StatusOK).JSON(rsp)
}

func (handler *BorrowerHandler) RefreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	payload, err := handler.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		errMsg := fmt.Sprintf("failed to verify refresh token: %s", err.Error())
		return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
	}

	input := &service.RefreshTokenInput{
		SessionID: payload.ID.String(),
	}

	_, err = handler.borrowerService.RefreshToken(ctx.Context(), input)
	if err != nil {
		apiError := FromServiceError(err)
		return errorResponse(ctx, apiError.StatusCode, apiError.Message)
	}

	accessToken, _, err := handler.tokenMaker.CreateToken(
		payload.UserID,
		handler.config.AccessTokenDuration,
	)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create access token: %s", err.Error())
		return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
	}

	rsp := RefreshTokenResponse{
		AccessToken: accessToken,
	}
	return ctx.Status(fiber.StatusOK).JSON(rsp)
}

func (handler *BorrowerHandler) GetBorrower(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	borrowerID, err := uuid.Parse(id)
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse borrower_id: %s", err.Error())
		return errorResponse(ctx, fiber.StatusBadRequest, errMsg)
	}

	authPayload := ctx.Locals(authorizationPayloadKey).(*token.Payload)
	if id != authPayload.UserID {
		errMsg := "borrower doesn't belong to the authenticated user"
		return errorResponse(ctx, fiber.StatusForbidden, errMsg)
	}

	input := &service.GetBorrowerByIDInput{
		BorrowerID: borrowerID,
	}

	output, err := handler.borrowerService.GetBorrowerByID(ctx.Context(), input)
	if err != nil {
		apiError := FromServiceError(err)
		return errorResponse(ctx, apiError.StatusCode, apiError.Message)
	}

	rsp := GetBorrowerResponse{
		ID:              output.Borrower.ID,
		CreatedAt:       output.Borrower.CreatedAt,
		UpdatedAt:       output.Borrower.UpdatedAt,
		Email:           output.Borrower.Email,
		LineID:          output.Borrower.LineID,
		Nickname:        output.Borrower.Nickname,
		IsEmailVerified: output.Borrower.IsEmailVerified,
		Role:            output.Borrower.Role,
	}
	return ctx.Status(fiber.StatusOK).JSON(rsp)
}
