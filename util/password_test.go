package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(8)

	hasedPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = CheckPassword(password, hasedPassword)
	require.NoError(t, err)

	hasedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, hasedPassword, hasedPassword2)
}
