package authentication

import (
	"golang.org/x/crypto/bcrypt"
	"reece.start/internal/configuration"
)

func HashPassword(password string, config *configuration.Config) ([]byte, error) {
	cost := 14 // Production cost
	if config.Test {
		cost = 4 // Fast cost for tests (about 1000x faster)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return bytes, err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
