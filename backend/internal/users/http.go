package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"reece.start/internal/access"
	"reece.start/internal/api"
	"reece.start/internal/constants"
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

	impersonatingUserId, _ := middleware.GetImpersonatingUserIDFromJWT(c)

	if req.Data.Relationships.ImpersonatedUser != nil {
		if impersonatingUserId != 0 {
			// if the impersonating user id is not 0, then the user is already impersonating someone
			return api.ErrForbiddenImpersonationNotAllowed
		}

		if err := access.HasAdminAccess(c, []constants.UserScope{constants.UserScopeAdminUsersImpersonate}); err != nil {
			return err
		}

		// set the user id to the impersonated user id
		impersonatingUserId, err = api.ParseUserIDFromString(req.Data.Relationships.ImpersonatedUser.Data.Id)

		if err != nil {
			return err
		}

		actualUserId := userId
		userId = impersonatingUserId
		impersonatingUserId = actualUserId
	}

	if req.Data.Meta.StopImpersonating {		
		if impersonatingUserId == 0 {
			// if the impersonating user id is 0, then the user is not impersonating anyone
			return api.ErrForbiddenImpersonationNotAllowed
		}

		userId = impersonatingUserId
		impersonatingUserId = 0
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
			ImpersonatingUserId: &impersonatingUserId,
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

func GetUsersEndpoint(c echo.Context, query GetUsersQuery) error {
	// Check admin access
	if err := access.HasAdminAccess(c, []constants.UserScope{constants.UserScopeAdminUsersList}); err != nil {
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	// Set default values if not provided
	cursor := query.Cursor
	size := query.Size
	search := query.Search
	if size <= 0 {
		size = 20
	}

	response, err := getUsers(GetUsersServiceRequest{
		Cursor:      cursor,
		Size:        size,
		Search:      search,
		Tx:          db,
		MinioClient: minioClient,
	})

	if err != nil {
		return err
	}

	// Convert to response format
	userData := make([]UserDataWithMeta, 0, len(response.Users))
	for _, userDto := range response.Users {
		userData = append(userData, UserDataWithMeta{
			UserData: UserData{
				Id:   strconv.FormatUint(uint64(userDto.User.ID), 10),
				Type: constants.ApiTypeUser,
				Attributes: UserAttributes{
					Name:  userDto.User.Name,
					Email: userDto.User.Email,
				},
			},
			Meta: UserMeta{
				LogoDistributionUrl: userDto.LogoDistributionUrl,
				TokenRevocation: UserTokenRevocation{
					LastIssuedAt: userDto.User.Revocation.LastValidIssuedAt,
					CanRefresh:   userDto.User.Revocation.CanRefresh,
				},
			},
		})
	}

	return c.JSON(http.StatusOK, GetUsersResponse{
		Data:  userData,
		Links: api.BuildPaginationLinks(api.BuildPaginationLinksParams{
			PrevCursor: response.PrevCursor,
			NextCursor: response.NextCursor,
			Context: c,
		}),
	})
}

func GoogleOAuthCallbackEndpoint(c echo.Context, req GoogleOAuthCallbackRequest) error {
	db := middleware.GetDB(c)
	config := middleware.GetConfig(c)
	minioClient := middleware.GetMinioClient(c)

	user, err := googleOAuthCallback(GoogleOAuthCallbackServiceRequest{
		Params: GoogleOAuthCallbackParams{
			Code:        req.Data.Attributes.Code,
			State:       req.Data.Attributes.State,
			RedirectUri: req.Data.Attributes.RedirectUri,
		},
		Tx:          db,
		Config:      config,
		MinioClient: minioClient,
	})

	if err != nil {
		return err // Middleware will handle the error response
	}

	return c.JSON(http.StatusOK, mapUserToResponse(user))
}
