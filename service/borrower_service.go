package service

import (
	"context"
)

type BorrowerService interface {
	SignUp(ctx context.Context, input *SignUpInput) (*SignUpOutput, error)
	SignIn(ctx context.Context, input *SignInInput) (*SignInOutput, error)
	RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error)
}
