package impl

import (
	"context"

	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/util"
)

func NewBorrowerServiceImpl(borrowerStore db.Store) service.BorrowerService {
	return &borrowerServiceImpl{
		borrowerStore: borrowerStore,
	}
}

type borrowerServiceImpl struct {
	borrowerStore db.Store
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
