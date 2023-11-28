package service

import (
	"context"
)

type BorrowerService interface {
	SignUp(ctx context.Context, input *SignUpInput) (*SignUpOutput, error)
}
