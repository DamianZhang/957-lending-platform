package token

import (
	"errors"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrTokenExpired = errors.New("token has expired")
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific userID and duration
	CreateToken(userID string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
