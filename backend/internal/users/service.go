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
	"reece.start/internal/authentication"
	"reece.start/internal/configuration"
	"reece.start/internal/constants"
	"reece.start/internal/models"
)

type CreateUserParams struct {
	Name     string
	Email    string
	Password string
	Timezone string
}

type CreateUserServiceRequest struct {
	Params CreateUserParams
	Tx     *gorm.DB
	Config *configuration.Config
}

type LoginUserParams struct {
	Email    string
	Password string
}

type LoginUserServiceRequest struct {
	Params LoginUserParams
	Tx     *gorm.DB
	Config *configuration.Config
	MinioClient *minio.Client
}

type GetUserByIDServiceRequest struct {
	UserID uint
	Tx     *gorm.DB
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
	Params UpdateUserParams
	Tx     *gorm.DB
	MinioClient *minio.Client
}

type GetUserLogoDistributionUrlServiceRequest struct {
	UserID uint
	Tx     *gorm.DB
	MinioClient *minio.Client
}

type UserDto struct {
	User *models.User
	Token string
	LogoDistributionUrl string
}

type CreateAuthenticatedUserTokenServiceRequest struct {
	Params CreateAuthenticatedUserTokenParams
	Tx     *gorm.DB
	Config *configuration.Config
}

type CreateAuthenticatedUserTokenParams struct {
	UserId uint
	OrganizationId *uint
}

type SelectMembershipRole struct {
	Role *constants.OrganizationRole
}

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
	token, err := authentication.CreateJWT(config, authentication.JwtOptions{
		UserId: user.ID,
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

	var selectMembershipRole SelectMembershipRole = SelectMembershipRole{}
	var organizationScopes *[]constants.OrganizationScope

	if request.Params.OrganizationId != nil {
		err := tx.Model(&models.OrganizationMembership{}).Where("user_id = ? AND organization_id = ?", request.Params.UserId, request.Params.OrganizationId).Select("role").First(&selectMembershipRole).Error
		if err != nil {
			return "", err
		}
		
		scopes := constants.OrganizationRoleToScopes[constants.OrganizationRole(*selectMembershipRole.Role)]
		organizationScopes = &scopes
	}

	token, err := authentication.CreateJWT(config, authentication.JwtOptions{
		UserId: request.Params.UserId,
		OrganizationId: request.Params.OrganizationId,
		OrganizationRole: selectMembershipRole.Role,
		OrganizationScopes: organizationScopes,
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
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check if password matches
	if !authentication.CheckPasswordHash(params.Password, string(user.HashedPassword)) {
		return nil, errors.New("invalid email or password")
	}

	// Generate the JWT token for the user
	token, err := authentication.CreateJWT(config, authentication.JwtOptions{
		UserId: user.ID,
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
			return nil, errors.New("user not found")
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