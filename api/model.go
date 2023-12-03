package api

import (
	"time"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	LineID   string `json:"line_id" validate:"required,min=1,max=20"`
	Nickname string `json:"nickname" validate:"required,alphanum,min=1,max=20"`
}

type SignUpResponse struct {
	Email    string `json:"email"`
	LineID   string `json:"line_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignInResponse struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type GetBorrowerResponse struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Email           string    `json:"email"`
	LineID          string    `json:"line_id"`
	Nickname        string    `json:"nickname"`
	IsEmailVerified bool      `json:"is_email_verified"`
	Role            string    `json:"role"`
}
