package token

import (
	"testing"
	"time"

	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker := NewJWTMaker()

	email := util.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker := NewJWTMaker()

	token, payload, err := maker.CreateToken(util.RandomEmail(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	fakePayload, err := NewPayload(util.RandomEmail(), time.Minute)
	require.NoError(t, err)

	fakeClaims := MyCustomClaims{
		ID:    fakePayload.ID,
		Email: fakePayload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(fakePayload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(fakePayload.IssuedAt),
			NotBefore: jwt.NewNumericDate(fakePayload.IssuedAt),
		},
	}
	fakeJWTToken := jwt.NewWithClaims(jwt.SigningMethodNone, fakeClaims)

	fakeToken, err := fakeJWTToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker := NewJWTMaker()

	payload, err := maker.VerifyToken(fakeToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
