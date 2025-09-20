package users

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
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

	hashedPassword, err := authentication.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	if err := tx.Create(&user).Error; err != nil {
		return nil, err
	}

	// Generate JWT token for the new user
	token, err := createAuthenticatedUserToken(CreateAuthenticatedUserTokenServiceRequest{
		Params: CreateAuthenticatedUserTokenParams{
			UserId: user.ID,
		},
		Tx: tx,
		Config: config,
	})

	if err != nil {
		return nil, err
	}

	return &UserDto{
		User: user,
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

	if request.Params.OrganizationId != nil {
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

	token, err := authentication.CreateJWT(config, authentication.JwtOptions{
		UserId: request.Params.UserId,
		OrganizationId: request.Params.OrganizationId,
		OrganizationRole: selectMembershipRole.Role,
		Scopes: &scopes,
		Role: &userRole,
	})

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
			return nil,api.ErrUnauthorizedInvalidLogin
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
		Tx: tx,
		Config: config,
	})
	if err != nil {
		return nil, err
	}

	// Get the logo distribution URL for the user
	logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
		UserID: user.ID,
		Tx:     tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &UserDto{
		User: &user,
		Token: token,
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
		UserID: user.ID,
		Tx:     tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &UserDto{
		User: &user,
		LogoDistributionUrl: logoDistributionUrl,
	}, nil
}

func updateUser(request UpdateUserServiceRequest) (*UserDto, error) {
	tx := request.Tx
	params := request.Params
	minioClient := request.MinioClient

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
		hashedPassword, err := authentication.HashPassword(params.Password)
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

		log.Printf("Uploading logo for user %d of length %d\n", user.ID, len(decodedImage))

		// Get the mime type from the image
		mimeType := http.DetectContentType(decodedImage)

		log.Printf("Detected logo mime type: %s\n", mimeType)

		objectName := fmt.Sprintf("%d", user.ID)

		// upload the image to minio
		_, err = minioClient.PutObject(context.Background(), string(constants.StorageBucketUserLogos), objectName, bytes.NewReader(decodedImage), int64(len(decodedImage)), minio.PutObjectOptions{
			ContentType: mimeType,
		})

		if err != nil {
			return nil, err
		}

		log.Printf("Updated logo for user %d\n", user.ID)

		user.LogoFileStorageKey = objectName
	}

	// Save the updated user
	if err := tx.Save(&user).Error; err != nil {
		return nil, err
	}

	logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
		UserID: user.ID,
		Tx:     tx,
		MinioClient: minioClient,
	})
	if err != nil {
		return nil, err
	}

	return &UserDto{
		User: &user,
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

	if objectName == "" {
		return "", nil
	}

	presignedUrl, err := minioClient.PresignedGetObject(context.Background(), string(constants.StorageBucketUserLogos), objectName, time.Hour*24, url.Values{})
	if err != nil {
		return "", err
	}

	return presignedUrl.String(), nil
}

func getUsers(request GetUsersServiceRequest) (*GetUsersServiceResponse, error) {
	tx := request.Tx
	minioClient := request.MinioClient
	cursor := request.Cursor
	size := request.Size
	search := request.Search

	// Set default values if not provided
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	// Parse cursor (user ID) if provided
	var getUsersCursor GetUsersCursor
	if err := api.ParseCursor(cursor, &getUsersCursor); err != nil {
		return nil, err
	}

	// Get users with cursor-based pagination and search
	var users []models.User
	query := tx
	
	// Apply search filter if provided
	if search != "" {
		// Search by name, email, or ID (case-insensitive)
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR id::text ILIKE ?", 
			searchPattern, searchPattern, searchPattern)
	}
	
	if getUsersCursor != (GetUsersCursor{}) && getUsersCursor.Direction == "next" {
		// Get users after the cursor
		query = query.Where("id > ?", getUsersCursor.UserID).Order("id ASC")
	} else if getUsersCursor != (GetUsersCursor{}) && getUsersCursor.Direction == "prev" {
		query = query.Where("id < ?", getUsersCursor.UserID).Order("id DESC")
	} else {
		query = query.Order("id ASC")
	}
	
	// Get one extra record to determine if there are more pages
	err := query.Limit(size + 1).Find(&users).Error;
	if err != nil {
		return nil, err
	}

	// Check if there are more records
	hasNext := ((getUsersCursor == (GetUsersCursor{}) || getUsersCursor.Direction == "next") && len(users) > size) || (getUsersCursor.Direction == "prev")
	hasPrev := (getUsersCursor.Direction == "prev" && len(users) > size) || (getUsersCursor.Direction == "next")

	// hasNext := len(users) > size
	if hasNext || hasPrev {
		// Remove the extra record
		users = users[:size]
	}

	// hasPrev := getUsersCursor != (GetUsersCursor{})

	// Convert to DTOs with logo distribution URLs
	var userDtos []*UserDto
	for _, user := range users {
		logoDistributionUrl, err := GetUserLogoDistributionUrl(GetUserLogoDistributionUrlServiceRequest{
			UserID:      user.ID,
			Tx:          tx,
			MinioClient: minioClient,
		})
		if err != nil {
			return nil, err
		}

		userDtos = append(userDtos, &UserDto{
			User:                &user,
			LogoDistributionUrl: logoDistributionUrl,
		})
	}

	// Generate next cursor if there are more records
	var nextCursor string
	var prevCursor string
	if hasNext && len(userDtos) > 0 {
		nextCursor, err = api.EncodeCursor(GetUsersCursor{
			UserID: userDtos[len(userDtos)-1].User.ID,
			Direction: "next",
		})
		if err != nil {
			return nil, err
		}
	}

	if hasPrev && len(userDtos) > 0 {
		prevCursor, err = api.EncodeCursor(GetUsersCursor{
			UserID: userDtos[0].User.ID,
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
