package impl

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	mockdb "github.com/DamianZhang/957-lending-platform/db/mock"
	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/golang/mock/gomock"
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

	e.arg.HashedPassword = gotArg.HashedPassword
	return reflect.DeepEqual(e.arg, gotArg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestSignUp(t *testing.T) {
	borrower, password := expectedBorrower(t)

	testCases := []struct {
		name        string
		input       *service.SignUpInput
		buildStubs  func(store *mockdb.MockStore)
		checkOutput func(output *service.SignUpOutput, err error)
	}{
		{
			name: "OK",
			input: &service.SignUpInput{
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

					HashedPassword: borrower.HashedPassword,
					Role:           borrower.Role,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(borrower, nil)
			},
			checkOutput: func(output *service.SignUpOutput, err error) {
				require.NoError(t, err)
				require.Equal(t, borrower, output.Borrower)
			},
		},
		{
			name: "TooLongPassword",
			input: &service.SignUpInput{
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
			checkOutput: func(output *service.SignUpOutput, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrInternalFailure)
				require.Nil(t, output)
			},
		},
		{
			name: "DBErrConnDone",
			input: &service.SignUpInput{
				Email:    borrower.Email,
				Password: password,
				LineID:   borrower.LineID,
				Nickname: borrower.Nickname,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrConnDone)
			},
			checkOutput: func(output *service.SignUpOutput, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrInternalFailure)
				require.Nil(t, output)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			borrowerService := NewBorrowerServiceImpl(store, nil)

			output, err := borrowerService.SignUp(context.Background(), tc.input)
			tc.checkOutput(output, err)
		})
	}
}

func TestSignIn(t *testing.T) {
	borrower, password := expectedBorrower(t)

	testCases := []struct {
		name        string
		input       *service.SignInInput
		buildStubs  func(store *mockdb.MockStore)
		checkOutput func(output *service.SignInOutput, err error)
	}{
		{
			name: "OK",
			input: &service.SignInInput{
				Email:    borrower.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), borrower.Email).
					Times(1).
					Return(borrower, nil)
			},
			checkOutput: func(output *service.SignInOutput, err error) {
				require.NoError(t, err)
				require.Equal(t, borrower, output.Borrower)
			},
		},
		{
			name: "DBErrRecordNotFound",
			input: &service.SignInInput{
				Email:    "RecordNotFound",
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq("RecordNotFound")).
					Times(1).
					Return(db.User{}, db.ErrRecordNotFound)
			},
			checkOutput: func(output *service.SignInOutput, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrRecordNotFound)
				require.Nil(t, output)
			},
		},
		{
			name: "DBErrConnDone",
			input: &service.SignInInput{
				Email:    borrower.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, db.ErrConnDone)
			},
			checkOutput: func(output *service.SignInOutput, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrInternalFailure)
				require.Nil(t, output)
			},
		},
		{
			name: "WrongPassword",
			input: &service.SignInInput{
				Email:    borrower.Email,
				Password: "WrongPassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), borrower.Email).
					Times(1).
					Return(borrower, nil)
			},
			checkOutput: func(output *service.SignInOutput, err error) {
				var svcError service.Error
				require.ErrorAs(t, err, &svcError)
				require.ErrorIs(t, svcError.SvcErr(), service.ErrUnauthorized)
				require.Nil(t, output)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			borrowerService := NewBorrowerServiceImpl(store, nil)

			output, err := borrowerService.SignIn(context.Background(), tc.input)
			tc.checkOutput(output, err)
		})
	}
}

func expectedBorrower(t *testing.T) (borrower db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	borrower = db.User{
		Email:    util.RandomEmail(),
		LineID:   util.RandomString(6),
		Nickname: util.RandomString(6),

		HashedPassword: hashedPassword,
		Role:           util.BorrowerRole,
	}
	return
}
