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

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	args := UpdateUserParams{
		Username:       user.Username,
		HashedPassword: "changed_secret",
	}

	user2, err := testQueries.UpdateUser(context.Background(), args)
	require.NoError(t, err)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.FullName, user2.FullName)
	require.WithinDuration(t, user.CreatedAt, user2.CreatedAt, time.Second)

	require.NotEqual(t, user.HashedPassword, user2.HashedPassword)
	require.NotEqual(t, user.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.DeleteUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, user2.Username)

	_, err = testQueries.GetUser(context.Background(), user.Username)
	require.ErrorIs(t, err, sql.ErrNoRows)
}
