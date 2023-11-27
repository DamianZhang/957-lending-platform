package impl

import (
	"context"

	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	"github.com/DamianZhang/957-lending-platform/service"
	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewBorrowerServiceImpl(borrowerStore db.Store) service.BorrowerService {
	return &borrowerServiceImpl{
		validate:      validator.New(),
		borrowerStore: borrowerStore,
	}
}

type borrowerServiceImpl struct {
	validate      *validator.Validate
	borrowerStore db.Store
}

func (svc *borrowerServiceImpl) SignUp(ctx context.Context, req *service.SignUpRequest) (*service.SignUpResponse, error) {
	err := svc.validate.Struct(req)
	if err != nil {
		return nil, service.NewError(service.ErrBadRequest, err)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, service.NewError(service.ErrInternalFailure, err)
	}

	arg := db.CreateUserParams{
		Email:    req.Email,
		LineID:   req.LineID,
		Nickname: req.Nickname,

		HashedPassword: hashedPassword,
		ID:             uuid.New(),
		Role:           util.BorrowerRole,
	}

	borrower, err := svc.borrowerStore.CreateUser(ctx, arg)
	if err != nil {
		return nil, service.NewError(service.ErrInternalFailure, err)
	}

	rsp := &service.SignUpResponse{
		ID:       borrower.ID,
		Email:    borrower.Email,
		LineID:   borrower.LineID,
		Nickname: borrower.Nickname,
		Role:     borrower.Role,
	}
	return rsp, nil
}
