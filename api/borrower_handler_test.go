package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	"github.com/DamianZhang/957-lending-platform/service"
	mocksvc "github.com/DamianZhang/957-lending-platform/service/mock"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSignUpAPI(t *testing.T) {
	expectedReq := expectedSignUpRequest()
	expectedRsp := SignUpResponse{
		Email:    expectedReq.Email,
		LineID:   expectedReq.LineID,
		Nickname: expectedReq.Nickname,
		Role:     util.BorrowerRole,
	}

	testCases := []struct {
		name         string
		reqBody      SignUpRequest
		setReqHeader func(req *http.Request)
		buildStubs   func(svc *mocksvc.MockBorrowerService)
		checkRsp     func(t *testing.T, rsp *http.Response)
	}{
		{
			name:    "OK",
			reqBody: expectedReq,
			setReqHeader: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				input := &service.SignUpInput{
					Email:    expectedReq.Email,
					Password: expectedReq.Password,
					LineID:   expectedReq.LineID,
					Nickname: expectedReq.Nickname,
				}
				output := &service.SignUpOutput{
					Borrower: db.User{
						Email:    expectedRsp.Email,
						LineID:   expectedRsp.LineID,
						Nickname: expectedRsp.Nickname,
						Role:     expectedRsp.Role,
					},
				}
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(input)).
					Times(1).
					Return(output, nil)
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusCreated, rsp.StatusCode)
				requireRspMatchExpectedRsp(t, rsp, expectedRsp)
			},
		},
		{
			name:         "NoneContentType",
			reqBody:      expectedReq,
			setReqHeader: func(req *http.Request) {},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusBadRequest, rsp.StatusCode)
			},
		},
		{
			name: "InvalidEmail",
			reqBody: SignUpRequest{
				Email:    "invalid-email",
				Password: expectedReq.Password,
				LineID:   expectedReq.LineID,
				Nickname: expectedReq.Nickname,
			},
			setReqHeader: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusBadRequest, rsp.StatusCode)
			},
		},
		{
			name:    "ServiceInternalFailure",
			reqBody: expectedReq,
			setReqHeader: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, service.NewError(service.ErrInternalFailure, nil))
			},
			checkRsp: func(t *testing.T, rsp *http.Response) {
				require.Equal(t, fiber.StatusInternalServerError, rsp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocksvc.NewMockBorrowerService(ctrl)
			tc.buildStubs(svc)

			server := newTestServer(t, svc)

			data, err := json.Marshal(tc.reqBody)
			require.NoError(t, err)

			url := "/api/v1/borrowers/sign_up"
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			tc.setReqHeader(req)

			rsp, err := server.app.Test(req)
			require.NoError(t, err)
			tc.checkRsp(t, rsp)
		})
	}
}

func expectedSignUpRequest() SignUpRequest {
	return SignUpRequest{
		Email:    util.RandomEmail(),
		Password: util.RandomString(6),
		LineID:   util.RandomString(6),
		Nickname: util.RandomString(6),
	}
}

func requireRspMatchExpectedRsp(t *testing.T, rsp *http.Response, expectedRsp SignUpResponse) {
	data, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	require.NoError(t, err)

	var actualRsp SignUpResponse
	err = json.Unmarshal(data, &actualRsp)
	require.NoError(t, err)

	require.Equal(t, expectedRsp, actualRsp)
}
