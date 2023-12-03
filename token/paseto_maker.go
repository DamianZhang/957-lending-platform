package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
	parser       *paseto.Parser
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.NotBeforeNbf())

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create v4 symmetric key: %w", err)
	}

	return &PasetoMaker{
		symmetricKey: key,
		parser:       &parser,
	}, nil
}

// CreateToken creates a new token for a specific email and duration
func (maker *PasetoMaker) CreateToken(email string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(email, duration)
	if err != nil {
		return "", nil, err
	}

	token := paseto.NewToken()
	token.SetString("id", payload.ID.String())
	token.SetString("email", payload.Email)
	token.SetExpiration(payload.ExpiresAt)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetNotBefore(payload.IssuedAt)

	return token.V4Encrypt(maker.symmetricKey, nil), payload, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parsedToken, err := maker.parser.ParseV4Local(maker.symmetricKey, token, nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, ErrTokenExpired
		}

		return nil, ErrInvalidToken
	}

	id, err := parsedToken.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidToken
	}
	email, err := parsedToken.GetString("email")
	if err != nil {
		return nil, ErrInvalidToken
	}
	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}
	expiresAt, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		ID:        uuid,
		Email:     email,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}
	return payload, nil
}
