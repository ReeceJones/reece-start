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
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
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
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid email or password",
		})
	}

	return c.JSON(http.StatusOK, mapUserToResponse(user))
}

func GetAuthenticatedUserEndpoint(c echo.Context) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid user token",
		})
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	user, err := getUserByID(GetUserByIDServiceRequest{
		UserID: userID,
		Tx:     db,
		MinioClient: minioClient,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, api.ApiError{
			Code:    constants.ErrorCodeNotFound,
			Message: "User not found",
		})
	}

	return c.JSON(http.StatusOK, mapUserToResponse(user))
}

func UpdateUserEndpoint(c echo.Context, req UpdateUserRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid user token",
		})
	}

	// Parse the user ID from the URL parameter
	paramUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    constants.ErrorCodeBadRequest,
			Message: "Invalid user ID",
		})
	}

	// Ensure users can only update their own profile
	if uint(paramUserID) != userID {
		return c.JSON(http.StatusForbidden, api.ApiError{
			Code:    constants.ErrorCodeForbidden,
			Message: "You can only update your own profile",
		})
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
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
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