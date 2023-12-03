package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DamianZhang/957-lending-platform/token"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	req *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	email string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	email := util.RandomEmail()

	testCases := []struct {
		name      string
		setUpAuth func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		checkRsp  func(t *testing.T, rsp *http.Response)
	}{
		{
			name: "OK",
			setUpAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, email, time.Minute)
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusOK, rsp.StatusCode)
			},
		},
		{
			name:      "NoAuthorizationHeader",
			setUpAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, rsp.StatusCode)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setUpAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, "", email, time.Minute)
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, rsp.StatusCode)
			},
		},
		{
			name: "UnsupportedAuthorizationType",
			setUpAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, "unsupported", email, time.Minute)
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, rsp.StatusCode)
			},
		},
		{
			name: "ExpiredToken",
			setUpAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, req, tokenMaker, authorizationTypeBearer, email, -time.Minute)
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, rsp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"
			server.app.Get(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *fiber.Ctx) error {
					return ctx.Status(fiber.StatusOK).SendString("OK")
				},
			)

			req := httptest.NewRequest(http.MethodGet, authPath, nil)
			tc.setUpAuth(t, req, server.tokenMaker)

			rsp, err := server.app.Test(req)
			require.NoError(t, err)
			tc.checkRsp(t, rsp)
		})
	}
}
