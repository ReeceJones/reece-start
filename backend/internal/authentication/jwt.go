package authentication

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"reece.start/internal/api"
	"reece.start/internal/configuration"
	"reece.start/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId              string                      `json:"user_id"`
	OrganizationId      *string                     `json:"organization_id"`
	OrganizationRole    *constants.OrganizationRole `json:"organization_role"`
	Scopes              *[]constants.UserScope      `json:"scopes"`
	Role                *constants.UserRole         `json:"role"`
	IsImpersonating     *bool                       `json:"is_impersonating"`
	ImpersonatingUserId *string                     `json:"impersonating_user_id"` // The actual user id of the authenticated user
}

type JwtOptions struct {
	UserId              uuid.UUID
	OrganizationId      *uuid.UUID
	OrganizationRole    *constants.OrganizationRole
	Scopes              *[]constants.UserScope
	Role                *constants.UserRole
	IsImpersonating     *bool
	ImpersonatingUserId *uuid.UUID
	CustomExpiry        *time.Time
}

func CreateJWT(config *configuration.Config, options JwtOptions) (string, error) {
	now := time.Now()
	userIdString := options.UserId.String()

	activeOrganizationId := getActiveOrganizationIdFromOptions(options)
	expiresAt := getExpiryFromOptions(config, options)
	impersonatingUserId := getImpersonatingUserIdFromOptions(options)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		UserId:              userIdString,
		OrganizationId:      activeOrganizationId,
		OrganizationRole:    options.OrganizationRole,
		Scopes:              options.Scopes,
		Role:                options.Role,
		IsImpersonating:     options.IsImpersonating,
		ImpersonatingUserId: impersonatingUserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    config.JwtIssuer,
			Subject:   userIdString,
			Audience:  jwt.ClaimStrings{config.JwtAudience},
		},
	})

	return token.SignedString([]byte(config.JwtSecret))
}

func ValidateJWT(config *configuration.Config, tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		slog.Error("Error parsing JWT token", "error", err, "token", tokenString)
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	slog.Error("Invalid JWT token: claims type assertion failed or token is not valid", "error", api.ErrForbiddenNoAccess)
	return nil, api.ErrForbiddenNoAccess
}

func getActiveOrganizationIdFromOptions(options JwtOptions) *string {
	if options.OrganizationId != nil {
		orgIdString := options.OrganizationId.String()
		return &orgIdString
	}
	return nil
}

func getExpiryFromOptions(config *configuration.Config, options JwtOptions) *jwt.NumericDate {
	if options.CustomExpiry != nil {
		return jwt.NewNumericDate(*options.CustomExpiry)
	}

	now := time.Now()
	return jwt.NewNumericDate(now.Add(time.Duration(config.JwtExpirationTime) * time.Second))
}

func getImpersonatingUserIdFromOptions(options JwtOptions) *string {
	if options.ImpersonatingUserId != nil {
		impersonatingUserIdString := options.ImpersonatingUserId.String()
		return &impersonatingUserIdString
	}
	return nil
}
