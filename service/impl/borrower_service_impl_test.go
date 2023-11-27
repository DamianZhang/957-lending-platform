package impl

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	mockdb "github.com/DamianZhang/957-lending-platform/db/mock"
	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	gotArg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(gotArg.HashedPassword, e.password)
	if err != nil {
		return false
	}

	if gotArg.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return false
	}

	if gotArg.Role != util.BorrowerRole {
		return false
	}

	e.arg.HashedPassword = gotArg.HashedPassword
	e.arg.ID = gotArg.ID
	e.arg.Role = gotArg.Role
	return reflect.DeepEqual(e.arg, gotArg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestSignUp(t *testing.T) {
	borrower, password := randomBorrower(t)

	testCases := []struct {
		name        string
		input       *service.SignUpRequest
		buildStubs  func(store *mockdb.MockStore)
		checkOutput func(rsp *service.SignUpResponse, err error)
	}{
		{
			name: "OK",
			input: &service.SignUpRequest{
				Email:    borrower.Email,
				Password: password,
				LineID:   borrower.LineID,
				Nickname: borrower.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Email:    borrower.Email,
					LineID:   borrower.LineID,
					Nickname: borrower.Nickname,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(borrower, nil)
			},
			checkOutput: func(rsp *service.SignUpResponse, err error) {
				require.NoError(t, err)
				requireRspMatchBorrower(t, rsp, borrower)
			},
		},
		{
			name: "InvalidEmail",
			input: &service.SignUpRequest{
				Email:    "invalid-email",
				Password: password,
				LineID:   borrower.LineID,
				Nickname: borrower.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkOutput: func(rsp *service.SignUpResponse, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrBadRequest)
				require.Nil(t, rsp)
			},
		},
		{
			name: "TooLongPassword",
			input: &service.SignUpRequest{
				Email:    borrower.Email,
				Password: "01234567890123456789012345678901234567890123456789012345678901234567890123456789",
				LineID:   borrower.LineID,
				Nickname: borrower.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkOutput: func(rsp *service.SignUpResponse, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrInternalFailure)
				require.Nil(t, rsp)
			},
		},
		{
			name: "InternalError",
			input: &service.SignUpRequest{
				Email:    borrower.Email,
				Password: password,
				LineID:   borrower.LineID,
				Nickname: borrower.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkOutput: func(rsp *service.SignUpResponse, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrInternalFailure)
				require.Nil(t, rsp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			borrowerService := NewBorrowerServiceImpl(store)

			rsp, err := borrowerService.SignUp(context.Background(), tc.input)
			tc.checkOutput(rsp, err)
		})
	}
}

func randomBorrower(t *testing.T) (borrower db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	borrower = db.User{
		Email:    util.RandomEmail(),
		LineID:   util.RandomString(6),
		Nickname: util.RandomString(6),

		HashedPassword: hashedPassword,
		ID:             uuid.New(),
		Role:           util.BorrowerRole,
	}
	return
}

func requireRspMatchBorrower(t *testing.T, rsp *service.SignUpResponse, borrower db.User) {
	require.Equal(t, borrower.Email, rsp.Email)
	require.Equal(t, borrower.LineID, rsp.LineID)
	require.Equal(t, borrower.Nickname, rsp.Nickname)

	require.Equal(t, borrower.ID, rsp.ID)
	require.Equal(t, borrower.Role, rsp.Role)
}
