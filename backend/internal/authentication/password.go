package authentication

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

var (
	// testMode enables fast password hashing for tests (uses cost 4 instead of 14)
	testMode = os.Getenv("TEST_MODE") == "true"
)

// SetTestMode enables fast password hashing for tests
// This should be called from test setup to speed up password operations
func SetTestMode(enabled bool) {
	testMode = enabled
}

func HashPassword(password string) ([]byte, error) {
	cost := 14 // Production cost
	if testMode {
		cost = 4 // Fast cost for tests (about 1000x faster)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return bytes, err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
