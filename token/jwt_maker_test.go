package token

import (
	"testing"
	"time"

	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker := NewJWTMaker(util.RandomString(32))

	userID := util.RandomString(10)
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker := NewJWTMaker(util.RandomString(32))

	token, payload, err := maker.CreateToken(util.RandomString(10), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	fakePayload, err := NewPayload(util.RandomString(10), time.Minute)
	require.NoError(t, err)

	fakeClaims := MyCustomClaims{
		ID:     fakePayload.ID,
		UserID: fakePayload.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(fakePayload.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(fakePayload.IssuedAt),
			NotBefore: jwt.NewNumericDate(fakePayload.IssuedAt),
		},
	}
	fakeJWTToken := jwt.NewWithClaims(jwt.SigningMethodNone, fakeClaims)

	fakeToken, err := fakeJWTToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker := NewJWTMaker(util.RandomString(32))

	payload, err := maker.VerifyToken(fakeToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
