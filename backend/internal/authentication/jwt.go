package authentication

import (
	"fmt"
	"log"
	"time"

	"reece.start/internal/configuration"
	"reece.start/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId               string                           `json:"user_id"`
	OrganizationId *string                          `json:"organization_id"`
	OrganizationRole     *constants.OrganizationRole       `json:"organization_role"`
	Scopes   *[]constants.UserScope  `json:"scopes"`
	Role *constants.UserRole `json:"role"`
}

type JwtOptions struct {
	UserId uint
	OrganizationId *uint
	OrganizationRole *constants.OrganizationRole
	Scopes *[]constants.UserScope
	Role *constants.UserRole
}

func CreateJWT(config *configuration.Config, options JwtOptions) (string, error) {
	now := time.Now()
	userIdString := fmt.Sprintf("%d", options.UserId)

	var activeOrganizationId *string
	if options.OrganizationId != nil {
		orgIdString := fmt.Sprintf("%d", *options.OrganizationId)
		activeOrganizationId = &orgIdString
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId:               userIdString,
		OrganizationId: activeOrganizationId,
		OrganizationRole:     options.OrganizationRole,
		Scopes:   options.Scopes,
		Role:   options.Role,
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
