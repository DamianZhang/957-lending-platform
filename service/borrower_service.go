package service

import (
	"context"
)

type BorrowerService interface {
	SignUp(ctx context.Context, input *SignUpInput) (*SignUpOutput, error)
	SignIn(ctx context.Context, input *SignInInput) (*SignInOutput, error)
	CreateSession(ctx context.Context, input *CreateSessionInput) (*CreateSessionOutput, error)
	RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error)
	GetBorrowerByID(ctx context.Context, input *GetBorrowerByIDInput) (*GetBorrowerByIDOutput, error)
}
