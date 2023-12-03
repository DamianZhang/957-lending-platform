package service

import (
	"time"

	"github.com/DamianZhang/957-lending-platform/cache"
	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
	"github.com/google/uuid"
)

type SignUpInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	LineID   string `json:"line_id"`
	Nickname string `json:"nickname"`
}

type SignUpOutput struct {
	Borrower db.User `json:"borrower"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInOutput struct {
	Borrower db.User `json:"borrower"`
}

type CreateSessionInput struct {
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Email     string    `json:"email"`
}

type CreateSessionOutput struct {
	Session cache.Session `json:"session"`
}

type RefreshTokenInput struct {
	SessionID string `json:"session_id"`
}

type RefreshTokenOutput struct {
	Session cache.Session `json:"session"`
}

type GetBorrowerByIDInput struct {
	BorrowerID uuid.UUID `json:"borrower_id"`
}

type GetBorrowerByIDOutput struct {
	Borrower db.User `json:"borrower"`
}
