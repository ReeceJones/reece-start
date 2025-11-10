package authentication

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123!"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotNil(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword123!"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotNil(t, hashedPassword)
	require.True(t, CheckPasswordHash(password, string(hashedPassword)))
}
