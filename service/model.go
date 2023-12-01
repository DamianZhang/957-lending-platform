package service

import (
	"time"

	db "github.com/DamianZhang/957-lending-platform/db/sqlc"
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
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SignInOutput struct {
	Borrower db.User `json:"borrower"`
}
