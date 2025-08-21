package users

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"reece.start/internal/authentication"
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
}

type LoginUserParams struct {
	Email    string
	Password string
}

type LoginUserServiceRequest struct {
	Params LoginUserParams
	Tx     *gorm.DB
}

type GetUserByIDServiceRequest struct {
	UserID uint
	Tx     *gorm.DB
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

func createUser(request CreateUserServiceRequest) (*models.User, error) {
	tx := request.Tx
	params := request.Params

	hashedPassword, err := authentication.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	err = tx.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func loginUser(request LoginUserServiceRequest) (*models.User, error) {
	tx := request.Tx
	params := request.Params

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

	return &user, nil
}

func getUserByID(request GetUserByIDServiceRequest) (*models.User, error) {
	tx := request.Tx
	userID := request.UserID

	var user models.User
	err := tx.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func updateUser(request UpdateUserServiceRequest) (*models.User, error) {
	tx := request.Tx
	params := request.Params
	minioClient := request.MinioClient

	// Get the existing user
	user, err := getUserByID(GetUserByIDServiceRequest{
		UserID: params.UserID,
		Tx:     tx,
	})
	if err != nil {
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
	}

	// Save the updated user
	err = tx.Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}