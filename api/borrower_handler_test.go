package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DamianZhang/957-lending-platform/service"
	mocksvc "github.com/DamianZhang/957-lending-platform/service/mock"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	req := randomSignUpRequest()

	testCases := []struct {
		name          string
		reqBody       *service.SignUpRequest
		setReqHeader  func(req *http.Request)
		buildStubs    func(svc *mocksvc.MockBorrowerService)
		checkResponse func(rsp *http.Response)
	}{
		{
			name:    "OK",
			reqBody: req,
			setReqHeader: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				rsp := &service.SignUpResponse{
					Email:    req.Email,
					LineID:   req.LineID,
					Nickname: req.Nickname,
					ID:       uuid.New(),
					Role:     util.BorrowerRole,
				}
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(req)).
					Times(1).
					Return(rsp, nil)
			},
			checkResponse: func(rsp *http.Response) {
				require.Equal(t, fiber.StatusCreated, rsp.StatusCode)
				requireRspMatchReq(t, rsp, req)
			},
		},
		{
			name:    "NoneContentType",
			reqBody: req,
			setReqHeader: func(req *http.Request) {
			},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp *http.Response) {
				require.Equal(t, fiber.StatusBadRequest, rsp.StatusCode)
			},
		},
		{
			name:    "ServiceBadRequest",
			reqBody: req,
			setReqHeader: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(req)).
					Times(1).
					Return(nil, service.NewError(service.ErrBadRequest, service.ErrBadRequest))
			},
			checkResponse: func(rsp *http.Response) {
				require.Equal(t, fiber.StatusBadRequest, rsp.StatusCode)
			},
		},
		{
			name:    "ServiceInternalFailure",
			reqBody: req,
			setReqHeader: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			buildStubs: func(svc *mocksvc.MockBorrowerService) {
				svc.EXPECT().
					SignUp(gomock.Any(), gomock.Eq(req)).
					Times(1).
					Return(nil, service.NewError(service.ErrInternalFailure, service.ErrInternalFailure))
			},
			checkResponse: func(rsp *http.Response) {
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
			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			tc.setReqHeader(request)

			rsp, err := server.app.Test(request)
			require.NoError(t, err)
			tc.checkResponse(rsp)
		})
	}
}

func randomSignUpRequest() *service.SignUpRequest {
	return &service.SignUpRequest{
		Email:    util.RandomEmail(),
		Password: util.RandomString(6),
		LineID:   util.RandomString(6),
		Nickname: util.RandomString(6),
	}
}

func requireRspMatchReq(t *testing.T, rsp *http.Response, req *service.SignUpRequest) {
	data, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	require.NoError(t, err)

	var gotRsp service.SignUpResponse
	err = json.Unmarshal(data, &gotRsp)
	require.NoError(t, err)

	require.Equal(t, req.Email, gotRsp.Email)
	require.Equal(t, req.LineID, gotRsp.LineID)
	require.Equal(t, req.Nickname, gotRsp.Nickname)

	require.NotZero(t, gotRsp.ID)
	require.Equal(t, util.BorrowerRole, gotRsp.Role)
}
