package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"reece.start/internal/api"
	"reece.start/internal/middleware"
)


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
