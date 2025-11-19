package users

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/authentication"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

func createUser(request CreateUserServiceRequest) (*UserDto, error) {
	tx := request.Tx
	params := request.Params
	config := request.Config
	posthogClient := request.PostHogClient

	hashedPassword, err := authentication.HashPassword(params.Password, config)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	if err := tx.Create(&user).Error; err != nil {
		// Check if this is a unique constraint violation (duplicate email)
		if api.IsUniqueConstraintViolation(err) {
			return nil, api.ErrUserEmailAlreadyExists
		}
		return nil, err
	}

	// Log user created event to PostHog
	posthogClient.Capture(
		fmt.Sprintf("%d", user.ID),
		"user created",
		map[string]any{
			"user_id":       user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"signup_method": "email",
		},
	)

	// Generate JWT token for the new user
	token, err := createAuthenticatedUserToken(CreateAuthenticatedUserTokenServiceRequest{
		Params: CreateAuthenticatedUserTokenParams{
			UserId: user.ID,
		},
		Tx:     tx,
		Config: config,
	})

	if err != nil {
		return nil, err
	}

	return &UserDto{
		User:  user,
		Token: token,
	}, nil
}

func createAuthenticatedUserToken(request CreateAuthenticatedUserTokenServiceRequest) (string, error) {
	tx := request.Tx
	config := request.Config

	var user models.User
	err := tx.First(&user, request.Params.UserId).Error
	if err != nil {
		return "", err
	}

	var selectMembershipRole SelectMembershipRole = SelectMembershipRole{}
	scopes := make([]constants.UserScope, 0)

	if request.Params.OrganizationId != nil && *request.Params.OrganizationId != 0 {
		err := tx.Model(&models.OrganizationMembership{}).Where("user_id = ? AND organization_id = ?", request.Params.UserId, request.Params.OrganizationId).Select("role").First(&selectMembershipRole).Error
		if err != nil {
			return "", err
		}

		organizationScopes := constants.OrganizationRoleToScopes[constants.OrganizationRole(*selectMembershipRole.Role)]
		scopes = append(scopes, organizationScopes...)
	}

	if user.Role != "" {
		userScopes := constants.UserRoleToScopes[constants.UserRole(user.Role)]
		scopes = append(scopes, userScopes...)
	}

	userRole := constants.UserRole(user.Role)

	isImpersonating := false
	var impersonatingUserId *string
	if request.Params.ImpersonatingUserId != nil && *request.Params.ImpersonatingUserId != 0 {
		isImpersonating = true
		impersonatingUserId = &[]string{fmt.Sprintf("%d", *request.Params.ImpersonatingUserId)}[0]
	}

	jwtOptions := authentication.JwtOptions{
		UserId:              request.Params.UserId,
		OrganizationId:      request.Params.OrganizationId,
		OrganizationRole:    selectMembershipRole.Role,
		Scopes:              &scopes,
		Role:                &userRole,
		IsImpersonating:     &isImpersonating,
		ImpersonatingUserId: impersonatingUserId,
	}

	// Use custom expiry if provided (for OAuth tokens)
	if request.Params.CustomExpiry != nil {
		jwtOptions.CustomExpiry = request.Params.CustomExpiry
	}

	token, err := authentication.CreateJWT(config, jwtOptions)

	return token, err
}

func loginUser(request LoginUserServiceRequest) (*UserDto, error) {
	tx := request.Tx
	params := request.Params
	config := request.Config
	minioClient := request.MinioClient

	var user models.User
	err := tx.Where("email = ?", params.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrUnauthorizedInvalidLogin
		}
		return nil, err
	}

	// Check if password matches
	if !authentication.CheckPasswordHash(params.Password, string(user.HashedPassword)) {
		return nil, api.ErrUnauthorizedInvalidLogin
	}

	// Generate the JWT token for the user
	token, err := createAuthenticatedUserToken(CreateAuthenticatedUserTokenServiceRequest{
		Params: CreateAuthenticatedUserTokenParams{
			UserId: user.ID,
		},
		Tx:     tx,
		Config: config,
	})
	if err != nil {
		return nil, err
	}

	// Get the logo distribution URL for the user
	logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
		UserID:      user.ID,
		Tx:          tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &UserDto{
		User:                &user,
		Token:               token,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}

func getUserByID(request GetUserByIDServiceRequest) (*UserDto, error) {
	tx := request.Tx
	userID := request.UserID
	minioClient := request.MinioClient

	var user models.User
	err := tx.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, api.ErrUserNotFound
		}
		return nil, err
	}

	// get the logo distribution URL for the user
	logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
		UserID:      user.ID,
		Tx:          tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &UserDto{
		User:                &user,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}

func updateUser(request UpdateUserServiceRequest) (*UserDto, error) {
	tx := request.Tx
	params := request.Params
	minioClient := request.MinioClient
	config := request.Config

	// Get the existing user
	var user models.User
	if err := tx.First(&user, params.UserID).Error; err != nil {
		return nil, err
	}

	// Update fields if provided
	if params.Name != "" {
		user.Name = params.Name
	}

	if params.Email != "" {
		user.Email = params.Email
	}

	if params.Password != "" {
		hashedPassword, err := authentication.HashPassword(params.Password, config)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = hashedPassword
	}

	if params.Logo != "" {
		// decode the image from base64 to a binary file
		decodedImage, err := base64.StdEncoding.DecodeString(params.Logo)
		if err != nil {
			return nil, err
		}

		slog.Info("Uploading logo for user", "userID", user.ID, "length", len(decodedImage))

		// Get the mime type from the image
		mimeType := http.DetectContentType(decodedImage)

		slog.Info("Detected logo mime type", "mimeType", mimeType)

		objectName := fmt.Sprintf("%d", user.ID)

		// upload the image to minio
		_, err = minioClient.PutObject(context.Background(), string(constants.StorageBucketUserLogos), objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{
			ContentType: mimeType,
		})

		if err != nil {
			return nil, err
		}

		slog.Info("Updated logo for user", "userID", user.ID)

		user.LogoFileStorageKey = objectName
	}

	// Save the updated user
	if err := tx.Save(&user).Error; err != nil {
		return nil, err
	}

	logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
		UserID:      user.ID,
		Tx:          tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &UserDto{
		User:                &user,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}

func GetUserLogoDistributionUrl(request GetUserLogoDistributionUrlServiceRequest) (string, error) {
	tx := request.Tx
	minioClient := request.MinioClient
	userID := request.UserID

	var user models.User
	err := tx.First(&user, userID).Error
	if err != nil {
		return "", err
	}

	objectName := user.LogoFileStorageKey

	// If user has uploaded a logo, return the presigned URL for that
	if objectName != "" {
		presignedUrl, err := minioClient.PresignedGetObject(context.Background(), string(constants.StorageBucketUserLogos), objectName, time.Hour*24, url.Values{})
		if err != nil {
			return "", err
		}
		return presignedUrl.String(), nil
	}

	// If no uploaded logo but Google profile image exists, return that as fallback
	if user.GoogleProfileImage != "" {
		return user.GoogleProfileImage, nil
	}

	return "", nil
}

// normalizePageSize ensures page size is within valid bounds (1-100, default 20)
func normalizePageSize(size int) int {
	if size <= 0 {
		return 20
	}
	if size > 100 {
		return 100
	}
	return size
}

// applySearchFilter applies search filter to a GORM query
func applySearchFilter(query *gorm.DB, search string) *gorm.DB {
	if search == "" {
		return query
	}
	searchPattern := "%" + search + "%"
	return query.Where("name ILIKE ? OR email ILIKE ? OR id::text ILIKE ?",
		searchPattern, searchPattern, searchPattern)
}

// applyPaginationFilter applies cursor-based pagination filter to a GORM query
func applyPaginationFilter(query *gorm.DB, cursor GetUsersCursor) *gorm.DB {
	if cursor == (GetUsersCursor{}) {
		return query.Order("id ASC")
	}
	if cursor.Direction == "next" {
		return query.Where("id > ?", cursor.UserID).Order("id ASC")
	}
	if cursor.Direction == "prev" {
		return query.Where("id < ?", cursor.UserID).Order("id DESC")
	}
	return query.Order("id ASC")
}

// calculatePaginationState determines hasNext and hasPrev based on cursor and result count
func calculatePaginationState(cursor GetUsersCursor, resultCount, pageSize int) (hasNext, hasPrev bool) {
	hasMoreResults := resultCount > pageSize

	if cursor == (GetUsersCursor{}) {
		// No cursor - first page
		return hasMoreResults, false
	}

	if cursor.Direction == "next" {
		// Forward pagination
		return hasMoreResults, true
	}

	if cursor.Direction == "prev" {
		// Backward pagination
		return true, hasMoreResults
	}

	return false, false
}

// createUserDtoWithLogo creates a UserDto with logo distribution URL
func createUserDtoWithLogo(user *models.User, tx *gorm.DB, minioClient *minio.Client) (*UserDto, error) {
	logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
		UserID:      user.ID,
		Tx:          tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &UserDto{
		User:                user,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}

func getUsers(request GetUsersServiceRequest) (*GetUsersServiceResponse, error) {
	tx := request.Tx
	minioClient := request.MinioClient
	cursor := request.Cursor
	size := normalizePageSize(request.Size)
	search := request.Search

	// Parse cursor (user ID) if provided
	var getUsersCursor GetUsersCursor
	if err := api.ParseCursor(cursor, &getUsersCursor); err != nil {
		return nil, err
	}

	// Build query with search and pagination filters
	query := tx.Model(&models.User{})
	query = applySearchFilter(query, search)
	query = applyPaginationFilter(query, getUsersCursor)

	// Get one extra record to determine if there are more pages
	var users []models.User
	err := query.Limit(size + 1).Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Calculate pagination state
	hasNext, hasPrev := calculatePaginationState(getUsersCursor, len(users), size)

	// Remove the extra record if needed
	if hasNext || hasPrev {
		users = users[:size]
	}

	// Convert to DTOs with logo distribution URLs
	var userDtos []*UserDto
	for _, user := range users {
		dto, err := createUserDtoWithLogo(&user, tx, minioClient)
		if err != nil {
			return nil, err
		}
		userDtos = append(userDtos, dto)
	}

	// Generate cursors if needed
	var nextCursor string
	var prevCursor string
	if hasNext && len(userDtos) > 0 {
		nextCursor, err = api.EncodeCursor(GetUsersCursor{
			UserID:    userDtos[len(userDtos)-1].User.ID,
			Direction: "next",
		})
		if err != nil {
			return nil, err
		}
	}

	if hasPrev && len(userDtos) > 0 {
		prevCursor, err = api.EncodeCursor(GetUsersCursor{
			UserID:    userDtos[0].User.ID,
			Direction: "prev",
		})
		if err != nil {
			return nil, err
		}
	}

	return &GetUsersServiceResponse{
		Users:      userDtos,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}, nil
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func googleOAuthCallback(request GoogleOAuthCallbackServiceRequest) (*UserDto, error) {
	tx := request.Tx
	config := request.Config
	minioClient := request.MinioClient
	posthogClient := request.PostHogClient
	params := request.Params

	// Configure OAuth
	oauth2Config := &oauth2.Config{
		ClientID:     config.GoogleOAuthClientId,
		ClientSecret: config.GoogleOAuthClientSecret,
		RedirectURL:  params.RedirectUri,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Exchange code for token
	token, err := oauth2Config.Exchange(context.Background(), params.Code)
	if err != nil {
		return nil, api.ErrUnauthorizedInvalidLogin
	}

	// Get user info from Google
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, api.ErrUnauthorizedInvalidLogin
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, api.ErrUnauthorizedInvalidLogin
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, api.ErrUnauthorizedInvalidLogin
	}

	// Check if user exists by Google ID
	var user models.User
	err = tx.Where("google_id = ?", googleUser.ID).First(&user).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// If user doesn't exist by Google ID, check by email
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = tx.Where("email = ?", googleUser.Email).First(&user).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		// If user exists by email, link Google account
		if err == nil {
			user.GoogleId = googleUser.ID
			if user.GoogleProfileImage == "" && googleUser.Picture != "" {
				user.GoogleProfileImage = googleUser.Picture
			}
			if err := tx.Save(&user).Error; err != nil {
				return nil, err
			}
		} else {
			// Create new user
			user = models.User{
				Name:               googleUser.Name,
				Email:              googleUser.Email,
				GoogleId:           googleUser.ID,
				GoogleProfileImage: googleUser.Picture,
				HashedPassword:     nil, // OAuth users don't have passwords
			}

			if err := tx.Create(&user).Error; err != nil {
				return nil, err
			}

			// Log user created event to PostHog (OAuth signup)
			posthogClient.Capture(
				fmt.Sprintf("%d", user.ID),
				"user created",
				map[string]any{
					"user_id":       user.ID,
					"email":         user.Email,
					"name":          user.Name,
					"signup_method": "oauth_google",
				},
			)
		}
	}

	// Generate JWT token with OAuth token expiration
	jwtToken, err := createAuthenticatedUserToken(CreateAuthenticatedUserTokenServiceRequest{
		Params: CreateAuthenticatedUserTokenParams{
			UserId:       user.ID,
			CustomExpiry: &token.Expiry,
		},
		Tx:     tx,
		Config: config,
	})

	if err != nil {
		return nil, err
	}

	// Get the logo distribution URL (prefer uploaded logo over Google profile image)
	logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
		UserID:      user.ID,
		Tx:          tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	// If no uploaded logo and Google profile image exists, use Google profile image as fallback
	if logoDistributionUrl == "" && user.GoogleProfileImage != "" {
		logoDistributionUrl = user.GoogleProfileImage
	}

	return &UserDto{
		User:                &user,
		Token:               jwtToken,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}
