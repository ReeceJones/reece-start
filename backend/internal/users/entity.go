package users

import (
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/configuration"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/internal/posthog"
)

// API Types
type UserAttributes struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type CreateUserAttributes struct {
	UserAttributes
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserAttributes struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserAttributes struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8"`
	Logo     string `json:"logo,omitempty" validate:"omitempty,base64"`
}

type GoogleOAuthCallbackAttributes struct {
	Code        string `json:"code" validate:"required"`
	State       string `json:"state" validate:"required"`
	RedirectUri string `json:"redirectUri" validate:"required,url"`
}

type UserTokenRevocation struct {
	LastIssuedAt *time.Time `json:"lastIssuedAt,omitempty"`
	CanRefresh   bool       `json:"canRefresh,omitempty"`
}

type UserMeta struct {
	Token               string              `json:"token,omitempty"`
	LogoDistributionUrl string              `json:"logoDistributionUrl,omitempty"`
	TokenRevocation     UserTokenRevocation `json:"tokenRevocation,omitempty"`
}

type UserData struct {
	Id         string            `json:"id"`
	Type       constants.ApiType `json:"type"`
	Attributes UserAttributes    `json:"attributes"`
}

type UserDataWithMeta struct {
	UserData
	Meta UserMeta `json:"meta"`
}

type CreateUserRequest struct {
	Data struct {
		Attributes CreateUserAttributes `json:"attributes"`
	} `json:"data"`
}

type OrganizationRelationshipData struct {
	Id   string `json:"id" validate:"required"`
	Type string `json:"type" validate:"required,oneof=organization"`
}

type OrganizationRelationship struct {
	Data OrganizationRelationshipData `json:"data" validate:"required"`
}

type UserRelationshipData struct {
	Id   string `json:"id" validate:"required"`
	Type string `json:"type" validate:"required,oneof=user"`
}

type UserRelationship struct {
	Data UserRelationshipData `json:"data" validate:"required"`
}

type CreateAuthenticatedUserTokenRelationships struct {
	Organization     *OrganizationRelationship `json:"organization"`
	ImpersonatedUser *UserRelationship         `json:"impersonatedUser"`
}

type CreateAuthenticatedUserTokenRequestMeta struct {
	StopImpersonating bool `json:"stopImpersonating"`
}

type CreateAuthenticatedUserTokenData struct {
	Type          constants.ApiType                         `json:"type" validate:"oneof=token"`
	Relationships CreateAuthenticatedUserTokenRelationships `json:"relationships"`
	Meta          CreateAuthenticatedUserTokenRequestMeta   `json:"meta"`
}

type CreateAuthenticatedUserTokenRequest struct {
	Data CreateAuthenticatedUserTokenData `json:"data"`
}

type CreateAuthenticatedUserTokenResponseMeta struct {
	Token string `json:"token"`
}

type CreateAuthenticatedUserTokenResponseData struct {
	Type          constants.ApiType                         `json:"type" validate:"oneof=token"`
	Relationships CreateAuthenticatedUserTokenRelationships `json:"relationships"`
	Meta          CreateAuthenticatedUserTokenResponseMeta  `json:"meta"`
}

type CreateAuthenticatedUserTokenResponse struct {
	Data CreateAuthenticatedUserTokenResponseData `json:"data"`
}

type GetUsersQuery struct {
	Cursor string `query:"page[cursor]"`
	Size   int    `query:"page[size]" validate:"min=1,max=100"`
	Search string `query:"search"`
}

type GetUsersResponse struct {
	Data  []UserDataWithMeta  `json:"data"`
	Links api.PaginationLinks `json:"links"`
}

type UserResponse struct {
	Data UserDataWithMeta `json:"data"`
}

type LoginUserRequest struct {
	Data struct {
		Attributes LoginUserAttributes `json:"attributes"`
	} `json:"data"`
}

type UpdateUserRequest struct {
	Data struct {
		Attributes UpdateUserAttributes `json:"attributes"`
	} `json:"data"`
}

type GoogleOAuthCallbackRequest struct {
	Data struct {
		Attributes GoogleOAuthCallbackAttributes `json:"attributes"`
	} `json:"data"`
}

// Service-layer types
type CreateUserParams struct {
	Name     string
	Email    string
	Password string
	Timezone string
}

type GoogleOAuthUserParams struct {
	Name            string
	Email           string
	GoogleId        string
	ProfileImageUrl string
	TokenExpiry     time.Time
}

type CreateUserServiceRequest struct {
	Params        CreateUserParams
	Tx            *gorm.DB
	Config        *configuration.Config
	PostHogClient *posthog.Client
}

type LoginUserParams struct {
	Email    string
	Password string
}

type LoginUserServiceRequest struct {
	Params      LoginUserParams
	Tx          *gorm.DB
	Config      *configuration.Config
	MinioClient *minio.Client
}

type GetUserByIDServiceRequest struct {
	UserID      uint
	Tx          *gorm.DB
	MinioClient *minio.Client
}

type UpdateUserParams struct {
	UserID   uint
	Name     string
	Email    string
	Password string
	Logo     string
}

type UpdateUserServiceRequest struct {
	Params      UpdateUserParams
	Tx          *gorm.DB
	MinioClient *minio.Client
	Config      *configuration.Config
}

type GetUserLogoDistributionUrlServiceRequest struct {
	UserID      uint
	Tx          *gorm.DB
	MinioClient *minio.Client
}

type GetUsersServiceRequest struct {
	Cursor      string
	Size        int
	Search      string
	Tx          *gorm.DB
	MinioClient *minio.Client
}

type GetUsersServiceResponse struct {
	Users      []*UserDto
	NextCursor string
	PrevCursor string
	HasNext    bool
	HasPrev    bool
}

type UserDto struct {
	User                *models.User
	Token               string
	LogoDistributionUrl string
}

type CreateAuthenticatedUserTokenServiceRequest struct {
	Params CreateAuthenticatedUserTokenParams
	Tx     *gorm.DB
	Config *configuration.Config
}

type CreateAuthenticatedUserTokenParams struct {
	UserId              uint
	OrganizationId      *uint
	ImpersonatingUserId *uint
	CustomExpiry        *time.Time
}

type GoogleOAuthCallbackServiceRequest struct {
	Params        GoogleOAuthCallbackParams
	Tx            *gorm.DB
	Config        *configuration.Config
	MinioClient   *minio.Client
	PostHogClient *posthog.Client
}

type GoogleOAuthCallbackParams struct {
	Code        string
	State       string
	RedirectUri string
}

type SelectMembershipRole struct {
	Role *constants.OrganizationRole
}

type GetUsersCursor struct {
	UserID    uint
	Direction string
}

// Type mappers
func mapUserToResponse(params *UserDto) UserResponse {
	return UserResponse{
		Data: UserDataWithMeta{
			UserData: UserData{
				Id:   strconv.FormatUint(uint64(params.User.ID), 10),
				Type: constants.ApiTypeUser,
				Attributes: UserAttributes{
					Name:  params.User.Name,
					Email: params.User.Email,
				},
			},
			Meta: UserMeta{
				Token:               params.Token,
				LogoDistributionUrl: params.LogoDistributionUrl,
				TokenRevocation: UserTokenRevocation{
					LastIssuedAt: params.User.Revocation.LastValidIssuedAt,
					CanRefresh:   params.User.Revocation.CanRefresh,
				},
			},
		},
	}
}

func mapCreateAuthenticatedUserTokenToResponse(req CreateAuthenticatedUserTokenRequest, token string) CreateAuthenticatedUserTokenResponse {
	return CreateAuthenticatedUserTokenResponse{
		Data: CreateAuthenticatedUserTokenResponseData{
			Type:          constants.ApiTypeToken,
			Relationships: req.Data.Relationships,
			Meta: CreateAuthenticatedUserTokenResponseMeta{
				Token: token,
			},
		},
	}
}
