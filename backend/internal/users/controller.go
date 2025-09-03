package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
)

// API Types
type UserAttributes struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
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
	Logo string `json:"logo,omitempty" validate:"omitempty,base64"`
}

type UserMeta struct {
	Token string `json:"token,omitempty"`
	LogoDistributionUrl string `json:"logoDistributionUrl,omitempty"`
}

type UserData struct {
	Id string `json:"id"`
	Type constants.ApiType `json:"type"`
	Attributes UserAttributes `json:"attributes"`
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
	Id string `json:"id" validate:"required"`
	Type string `json:"type" validate:"required,oneof=organization"`
}

type OrganizationRelationship struct {
	Data OrganizationRelationshipData `json:"data" validate:"required"`
}

type CreateAuthenticatedUserTokenRelationships struct {
	Organization *OrganizationRelationship `json:"organization"`
}

type CreateAuthenticatedUserTokenData struct {
	Type constants.ApiType `json:"type" validate:"oneof=token"`
	Relationships CreateAuthenticatedUserTokenRelationships `json:"relationships"`
}

type CreateAuthenticatedUserTokenRequest struct {
	Data CreateAuthenticatedUserTokenData `json:"data"`
}

type CreateAuthenticatedUserTokenMeta struct {
	Token string `json:"token"`
}

type CreateAuthenticatedUserTokenResponseData struct {
	Type constants.ApiType `json:"type" validate:"oneof=token"`
	Relationships CreateAuthenticatedUserTokenRelationships `json:"relationships"`
	Meta CreateAuthenticatedUserTokenMeta `json:"meta"`
}

type CreateAuthenticatedUserTokenResponse struct {
	Data CreateAuthenticatedUserTokenResponseData `json:"data"`
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


func CreateUserEndpoint(c echo.Context, req CreateUserRequest) error {
	config := middleware.GetConfig(c)
	db := middleware.GetDB(c)

	user, err := createUser(CreateUserServiceRequest{
		Params: CreateUserParams{
			Name:     req.Data.Attributes.Name,
			Email:    req.Data.Attributes.Email,
			Password: req.Data.Attributes.Password,
		},
		Tx: db,
		Config: config,
	})

	if err != nil {
		return err // Middleware will handle the error response
	}

	return c.JSON(http.StatusCreated, mapUserToResponse(user))
}

func LoginEndpoint(c echo.Context, req LoginUserRequest) error {
	config := middleware.GetConfig(c)
	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	user, err := loginUser(LoginUserServiceRequest{
		Params: LoginUserParams{
			Email:    req.Data.Attributes.Email,
			Password: req.Data.Attributes.Password,
		},
		Tx: db,
		Config: config,
		MinioClient: minioClient,
	})

	if err != nil {
		return api.ErrUnauthorizedInvalidLogin // Middleware will handle the error response
	}

	return c.JSON(http.StatusOK, mapUserToResponse(user))
}

func GetAuthenticatedUserEndpoint(c echo.Context) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err // Middleware will handle the error response
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	user, err := getUserByID(GetUserByIDServiceRequest{
		UserID: userID,
		Tx:     db,
		MinioClient: minioClient,
	})

	if err != nil {
		return err // Middleware will handle the error response
	}

	return c.JSON(http.StatusOK, mapUserToResponse(user))
}

func CreateAuthenticatedUserTokenEndpoint(c echo.Context, req CreateAuthenticatedUserTokenRequest) error {
	userId, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	var organizationId uint
	if req.Data.Relationships.Organization != nil {
		organizationId, err = api.ParseOrganizationIDFromString(req.Data.Relationships.Organization.Data.Id)
		if err != nil {
			return err
		}
	}

	tx := middleware.GetDB(c)
	config := middleware.GetConfig(c)

	token, err := createAuthenticatedUserToken(CreateAuthenticatedUserTokenServiceRequest{
		Params: CreateAuthenticatedUserTokenParams{
			UserId: userId,
			OrganizationId: &organizationId,
		},
		Tx: tx,
		Config: config,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapCreateAuthenticatedUserTokenToResponse(req, token))
}

func UpdateUserEndpoint(c echo.Context, req UpdateUserRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err // Middleware will handle the error response
	}

	// Parse the user ID from the URL parameter
	paramUserID, err := api.ParseUserIDFromString(c.Param("id"))
	if err != nil {
		return err // Middleware will handle the error response
	}

	// Ensure users can only update their own profile
	if paramUserID != userID {
		return api.ErrForbiddenOwnProfileOnly // Middleware will handle the error response
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	user, err := updateUser(UpdateUserServiceRequest{
		Params: UpdateUserParams{
			UserID:   userID,
			Name:     req.Data.Attributes.Name,
			Email:    req.Data.Attributes.Email,
			Password: req.Data.Attributes.Password,
			Logo: req.Data.Attributes.Logo,
		},
		Tx: db,
		MinioClient: minioClient,
	})

	if err != nil {
		log.Error(err)
		return err // Middleware will handle the error response
	}

	return c.JSON(http.StatusOK, mapUserToResponse(user))
}

// Type mappers
func mapUserToResponse(params *UserDto) UserResponse {
	return UserResponse{
		Data: UserDataWithMeta{
			UserData: UserData{
				Id: strconv.FormatUint(uint64(params.User.ID), 10),
				Type: constants.ApiTypeUser,
				Attributes: UserAttributes{
					Name: params.User.Name,
					Email: params.User.Email,
				},
			},
			Meta: UserMeta{
				Token: params.Token,
				LogoDistributionUrl: params.LogoDistributionUrl,
			},
		},
	}
}

func mapCreateAuthenticatedUserTokenToResponse(req CreateAuthenticatedUserTokenRequest, token string) CreateAuthenticatedUserTokenResponse {
	return CreateAuthenticatedUserTokenResponse{
		Data: CreateAuthenticatedUserTokenResponseData{
			Type: constants.ApiTypeToken,
			Relationships: req.Data.Relationships,
			Meta: CreateAuthenticatedUserTokenMeta{
				Token: token,
			},
		},
	}
}