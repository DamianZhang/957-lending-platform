package token

import (
	"testing"
	"time"

	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker := NewPasetoMaker()

	email := util.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

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
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker := NewPasetoMaker()

	token, payload, err := maker.CreateToken(util.RandomEmail(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)
}
