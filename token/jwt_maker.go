package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) Maker {
	return &JWTMaker{
		secretKey: secretKey,
	}
}

type MyCustomClaims struct {
	ID     uuid.UUID `json:"id"`
	UserID string    `json:"user_id"`
	jwt.RegisteredClaims
}

// CreateToken creates a new token for a specific userID and duration
func (maker *JWTMaker) CreateToken(userID string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", nil, err
	}

	claims := MyCustomClaims{
		ID:     payload.ID,
		UserID: payload.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(payload.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			NotBefore: jwt.NewNumericDate(payload.IssuedAt),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// token.Method needs to match with our signing algorithm
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &MyCustomClaims{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*MyCustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		ID:        claims.ID,
		UserID:    claims.UserID,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}
	return payload, nil
}
