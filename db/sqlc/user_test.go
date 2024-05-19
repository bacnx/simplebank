package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bacnx/simplebank/util"
	"github.com/stretchr/testify/require"
)

// When pass argument
//
//	first arg is string to custom account.currency
func createRandomUser(t *testing.T) User {
	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	args := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, users, 5)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.PasswordChangedAt, newUser.PasswordChangedAt)
	require.WithinDuration(t, oldUser.PasswordChangedAt, newUser.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)

	require.NotEqual(t, oldUser.PasswordChangedAt, newUser.PasswordChangedAt)
	require.WithinDuration(t, time.Now(), newUser.PasswordChangedAt, time.Second)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.WithinDuration(t, oldUser.CreatedAt, newUser.CreatedAt, time.Second)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, newHashedPassword, newUser.HashedPassword)

	require.NotEqual(t, oldUser.PasswordChangedAt, newUser.PasswordChangedAt)
	require.WithinDuration(t, time.Now(), newUser.PasswordChangedAt, time.Second)

	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, newFullName, newUser.FullName)

	require.Equal(t, oldUser.Username, newUser.Username)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.WithinDuration(t, oldUser.CreatedAt, newUser.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, user2.Username)

	_, err = testQueries.GetUser(context.Background(), user.Username)
	require.ErrorIs(t, err, sql.ErrNoRows)
}
