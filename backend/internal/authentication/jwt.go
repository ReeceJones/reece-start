package authentication

import (
	"fmt"
	"log"
	"time"

	"reece.start/internal/configuration"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId               string  `json:"user_id"`
	ActiveOrganizationId *string `json:"active_organization_id"`
}

type JwtOptions struct {
	UserId uint
}

func CreateJWT(config *configuration.Config, options JwtOptions) (string, error) {
	now := time.Now()
	userIdString := fmt.Sprintf("%d", options.UserId)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId: userIdString,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(config.JwtExpirationTime) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.JwtIssuer,
			Subject:   userIdString,
			Audience:  jwt.ClaimStrings{config.JwtAudience},
		},
	})

	return token.SignedString([]byte(config.JwtSecret))
}

func UpdateJWT(config *configuration.Config, claims *JwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtSecret))
}

func ValidateJWT(config *configuration.Config, tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	} else {
		log.Printf("Error parsing JWT token: %v", err)
		return nil, err
	}
}
