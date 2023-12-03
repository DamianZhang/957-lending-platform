package api

import (
	"fmt"
	"strings"

	"github.com/DamianZhang/957-lending-platform/token"
	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a fiber middleware for authorization
func authMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			errMsg := "authorization header is not provided"
			return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			errMsg := "invalid authorization header format"
			return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			errMsg := fmt.Sprintf("unsupported authorization type %s", authorizationType)
			return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			errMsg := fmt.Sprintf("failed to verify access token: %s", err.Error())
			return errorResponse(ctx, fiber.StatusUnauthorized, errMsg)
		}

		ctx.Locals(authorizationPayloadKey, payload)
		return ctx.Next()
	}
}
