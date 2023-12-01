package impl

import (
	"context"

	"github.com/DamianZhang/957-lending-platform/cache"
	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/util"
)

func NewBorrowerServiceImpl(borrowerStore db.Store, borrowerCacher cache.Cacher) service.BorrowerService {
	return &borrowerServiceImpl{
		borrowerStore:  borrowerStore,
		borrowerCacher: borrowerCacher,
	}
}

type borrowerServiceImpl struct {
	borrowerStore  db.Store
	borrowerCacher cache.Cacher
}

func (svc *borrowerServiceImpl) SignUp(ctx context.Context, input *service.SignUpInput) (*service.SignUpOutput, error) {
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return nil, service.NewError(service.ErrInternalFailure, err)
	}

	arg := db.CreateUserParams{
		Email:    input.Email,
		LineID:   input.LineID,
		Nickname: input.Nickname,

		HashedPassword: hashedPassword,
		Role:           util.BorrowerRole,
	}

	borrower, err := svc.borrowerStore.CreateUser(ctx, arg)
	if err != nil {
		return nil, service.NewError(service.ErrInternalFailure, err)
	}

	output := &service.SignUpOutput{
		Borrower: borrower,
	}
	return output, nil
}

func (svc *borrowerServiceImpl) SignIn(ctx context.Context, input *service.SignInInput) (*service.SignInOutput, error) {
	borrower, err := svc.borrowerStore.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, service.NewError(service.ErrUnauthorized, err)
	}

	err = util.CheckPassword(borrower.HashedPassword, input.Password)
	if err != nil {
		return nil, service.NewError(service.ErrUnauthorized, err)
	}

	output := &service.SignInOutput{
		Borrower: borrower,
	}
	return output, nil
}
