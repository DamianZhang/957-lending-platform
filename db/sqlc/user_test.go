package db

import (
	"context"
	"testing"
	"time"

	"github.com/DamianZhang/957-lending-platform/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	wanted := CreateUserParams{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		LineID:         util.RandomString(6),
		Nickname:       util.RandomString(6),
		Role:           util.RandomRole(),
	}

	got, err := testStore.CreateUser(context.Background(), wanted)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, wanted.Email, got.Email)
	require.Equal(t, wanted.HashedPassword, got.HashedPassword)
	require.Equal(t, wanted.LineID, got.LineID)
	require.Equal(t, wanted.Nickname, got.Nickname)
	require.Equal(t, wanted.Role, got.Role)
	require.False(t, got.IsEmailVerified)
	require.NotZero(t, got.ID)
	require.NotZero(t, got.CreatedAt)
	require.NotZero(t, got.UpdatedAt)

	return got
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	wanted := createRandomUser(t)
	got, err := testStore.GetUserByEmail(context.Background(), wanted.Email)
	require.NoError(t, err)
	require.NotEmpty(t, got)

	require.Equal(t, wanted.ID, got.ID)
	require.Equal(t, wanted.Email, got.Email)
	require.Equal(t, wanted.HashedPassword, got.HashedPassword)
	require.Equal(t, wanted.LineID, got.LineID)
	require.Equal(t, wanted.Nickname, got.Nickname)
	require.Equal(t, wanted.IsEmailVerified, got.IsEmailVerified)
	require.Equal(t, wanted.Role, got.Role)
	require.WithinDuration(t, wanted.CreatedAt, got.CreatedAt, time.Second)
	require.WithinDuration(t, wanted.UpdatedAt, got.UpdatedAt, time.Second)
}

func TestGetUsers(t *testing.T) {
	wanted := 5
	for i := 0; i < wanted; i++ {
		createRandomUser(t)
	}

	users, err := testStore.GetUsers(context.Background(), GetUsersParams{
		Limit:  int32(wanted),
		Offset: 0,
	})
	require.NoError(t, err)
	require.Len(t, users, wanted)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestUpdateUserByEmailOnlyHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newUpdatedAt := time.Now()
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUserByEmail(context.Background(), UpdateUserByEmailParams{
		Email:     oldUser.Email,
		UpdatedAt: newUpdatedAt,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, oldUser.UpdatedAt, updatedUser.UpdatedAt)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)

	require.Equal(t, oldUser.ID, updatedUser.ID)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.LineID, updatedUser.LineID)
	require.Equal(t, oldUser.Nickname, updatedUser.Nickname)
	require.Equal(t, oldUser.IsEmailVerified, updatedUser.IsEmailVerified)
	require.Equal(t, oldUser.Role, updatedUser.Role)
	require.WithinDuration(t, oldUser.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.WithinDuration(t, newUpdatedAt, updatedUser.UpdatedAt, time.Second)
}

func TestUpdateUserByEmailOnlyIsEmailVerified(t *testing.T) {
	oldUser := createRandomUser(t)

	newUpdatedAt := time.Now()
	newIsEmailVerified := true

	updatedUser, err := testStore.UpdateUserByEmail(context.Background(), UpdateUserByEmailParams{
		Email:     oldUser.Email,
		UpdatedAt: newUpdatedAt,
		IsEmailVerified: pgtype.Bool{
			Bool:  newIsEmailVerified,
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, oldUser.UpdatedAt, updatedUser.UpdatedAt)
	require.NotEqual(t, oldUser.IsEmailVerified, updatedUser.IsEmailVerified)

	require.Equal(t, oldUser.ID, updatedUser.ID)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.LineID, updatedUser.LineID)
	require.Equal(t, oldUser.Nickname, updatedUser.Nickname)
	require.Equal(t, newIsEmailVerified, updatedUser.IsEmailVerified)
	require.Equal(t, oldUser.Role, updatedUser.Role)
	require.WithinDuration(t, oldUser.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.WithinDuration(t, newUpdatedAt, updatedUser.UpdatedAt, time.Second)
}

func TestUpdateUserByEmailAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newUpdatedAt := time.Now()
	newHashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	newLineID := util.RandomString(6)
	newNickname := util.RandomString(6)
	newIsEmailVerified := true

	updatedUser, err := testStore.UpdateUserByEmail(context.Background(), UpdateUserByEmailParams{
		Email:     oldUser.Email,
		UpdatedAt: newUpdatedAt,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
		LineID: pgtype.Text{
			String: newLineID,
			Valid:  true,
		},
		Nickname: pgtype.Text{
			String: newNickname,
			Valid:  true,
		},
		IsEmailVerified: pgtype.Bool{
			Bool:  newIsEmailVerified,
			Valid: true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.NotEqual(t, oldUser.UpdatedAt, updatedUser.UpdatedAt)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.LineID, updatedUser.LineID)
	require.NotEqual(t, oldUser.Nickname, updatedUser.Nickname)
	require.NotEqual(t, oldUser.IsEmailVerified, updatedUser.IsEmailVerified)

	require.Equal(t, oldUser.ID, updatedUser.ID)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newLineID, updatedUser.LineID)
	require.Equal(t, newNickname, updatedUser.Nickname)
	require.Equal(t, newIsEmailVerified, updatedUser.IsEmailVerified)
	require.Equal(t, oldUser.Role, updatedUser.Role)
	require.WithinDuration(t, oldUser.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.WithinDuration(t, newUpdatedAt, updatedUser.UpdatedAt, time.Second)
}
