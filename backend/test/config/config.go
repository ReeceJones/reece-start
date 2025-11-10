package config

import (
	"reece.start/internal/configuration"
)

// CreateTestConfig creates a test configuration with sensible defaults for testing.
// This provides a consistent test configuration across all test files.
func CreateTestConfig() *configuration.Config {
	return &configuration.Config{
		JwtSecret:         "test-secret-key",
		JwtIssuer:         "test-issuer",
		JwtAudience:       "test-audience",
		JwtExpirationTime: 3600, // 1 hour in seconds
	}
}
