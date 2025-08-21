package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/authentication"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
	"reece.start/internal/models"
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
	Token string `json:"token"`
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

type CreateUserResponse struct {
	Data UserDataWithMeta `json:"data"`
}

type LoginUserRequest struct {
	Data struct {
		Attributes LoginUserAttributes `json:"attributes"`
	} `json:"data"`
}

type LoginUserResponse struct {
	Data UserDataWithMeta `json:"data"`
}

type UpdateUserRequest struct {
	Data struct {
		Attributes UpdateUserAttributes `json:"attributes"`
	} `json:"data"`
}

type UpdateUserResponse struct {
	Data UserData `json:"data"`
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
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
	}

	// Generate JWT token for the new user
	token, err := authentication.CreateJWT(config, authentication.JwtOptions{
		UserId: user.ID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: "Failed to generate authentication token",
		})
	}

	return c.JSON(http.StatusCreated, userToResponse(user, token))
}

func LoginEndpoint(c echo.Context, req LoginUserRequest) error {
	config := middleware.GetConfig(c)
	db := middleware.GetDB(c)

	user, err := loginUser(LoginUserServiceRequest{
		Params: LoginUserParams{
			Email:    req.Data.Attributes.Email,
			Password: req.Data.Attributes.Password,
		},
		Tx: db,
	})

	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid email or password",
		})
	}

	// Generate JWT token for the authenticated user
	token, err := authentication.CreateJWT(config, authentication.JwtOptions{
		UserId: user.ID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: "Failed to generate authentication token",
		})
	}

	return c.JSON(http.StatusOK, userToResponse(user, token))
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

	user, err := getUserByID(GetUserByIDServiceRequest{
		UserID: userID,
		Tx:     db,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, api.ApiError{
			Code:    constants.ErrorCodeNotFound,
			Message: "User not found",
		})
	}

	return c.JSON(http.StatusOK, userToResponseWithoutToken(user))
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
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, userToResponseWithoutToken(user))
}

// Type mappers
func userToResponse(user *models.User, token string) CreateUserResponse {
	return CreateUserResponse{
		Data: UserDataWithMeta{
			UserData: UserData{
				Id: strconv.FormatUint(uint64(user.ID), 10),
				Type: constants.ApiTypeUser,
				Attributes: UserAttributes{
					Name: user.Name,
					Email: user.Email,
				},
			},
			Meta: UserMeta{
				Token: token,
			},
		},
	}
}

func userToResponseWithoutToken(user *models.User) UpdateUserResponse {
	return UpdateUserResponse{
		Data: UserData{
			Id: strconv.FormatUint(uint64(user.ID), 10),
			Type: constants.ApiTypeUser,
			Attributes: UserAttributes{
				Name: user.Name,
				Email: user.Email,
			},
		},
	}
}
