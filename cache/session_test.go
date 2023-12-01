package cache

import (
	"context"
	"testing"

	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) Session {
	expected := CreateSessionParams{
		ID:    util.RandomString(10),
		Email: util.RandomEmail(),
	}

	actual, err := testCacher.CreateSession(context.Background(), expected)
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Email, actual.Email)
	require.False(t, actual.IsBlocked)

	return actual
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSessionByID(t *testing.T) {
	expected := createRandomSession(t)
	actual, err := testCacher.GetSessionByID(context.Background(), expected.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Email, actual.Email)
	require.Equal(t, expected.IsBlocked, actual.IsBlocked)
}
