package config

import (
	"reece.start/internal/configuration"
)

// CreateTestConfig creates a test configuration with sensible defaults for testing.
// This provides a consistent test configuration across all test files.
func CreateTestConfig() *configuration.Config {
	return &configuration.Config{
		Test:                   true,
		Host:                   "localhost",
		Port:                   "8080",
		FrontendUrl:            "http://localhost:3000",
		JwtSecret:              "test-secret",
		JwtIssuer:              "test-issuer",
		JwtAudience:            "test-audience",
		JwtExpirationTime:      3600,
		StorageEndpoint:        "localhost:9000",
		StorageAccessKeyId:     "minioadmin",
		StorageSecretAccessKey: "minioadmin",
		StorageUseSSL:          false,
		EnableEmail:            false,
		ResendApiKey:           "test-key",
	}
}
