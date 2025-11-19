package authentication

import (
	"testing"

	"github.com/stretchr/testify/require"
	testconfig "reece.start/test/config"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123!"
	config := testconfig.CreateTestConfig()

	hashedPassword, err := HashPassword(password, config)

	require.NoError(t, err)
	require.NotNil(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword123!"
	config := testconfig.CreateTestConfig()

	hashedPassword, err := HashPassword(password, config)

	require.NoError(t, err)
	require.NotNil(t, hashedPassword)
	require.True(t, CheckPasswordHash(password, string(hashedPassword)))
}
