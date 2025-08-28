package organizations

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
)

// API Types
type OrganizationAttributes struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
}

type CreateOrganizationAttributes struct {
	OrganizationAttributes
	Logo string `json:"logo,omitempty" validate:"omitempty,base64"`
}

type UpdateOrganizationAttributes struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	Logo *string `json:"logo,omitempty" validate:"omitempty,base64"`
}

type OrganizationMeta struct {
	LogoDistributionUrl string `json:"logoDistributionUrl,omitempty"`
}

type OrganizationData struct {
	Id         string                `json:"id"`
	Type       constants.ApiType     `json:"type"`
	Attributes OrganizationAttributes `json:"attributes"`
}

type OrganizationDataWithMeta struct {
	OrganizationData
	Meta OrganizationMeta `json:"meta"`
}

type CreateOrganizationRequest struct {
	Data struct {
		Attributes CreateOrganizationAttributes `json:"attributes"`
	} `json:"data"`
}

type CreateOrganizationResponse struct {
	Data OrganizationDataWithMeta `json:"data"`
}

type UpdateOrganizationRequest struct {
	Data struct {
		Attributes UpdateOrganizationAttributes `json:"attributes"`
	} `json:"data"`
}

type UpdateOrganizationResponse struct {
	Data OrganizationDataWithMeta `json:"data"`
}

type GetOrganizationsResponse struct {
	Data []OrganizationDataWithMeta `json:"data"`
}

type GetOrganizationResponse struct {
	Data OrganizationDataWithMeta `json:"data"`
}

func CreateOrganizationEndpoint(c echo.Context, req CreateOrganizationRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid user token",
		})
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	var response CreateOrganizationResponse
	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		organization, err := createOrganization(CreateOrganizationServiceRequest{
			Params: CreateOrganizationParams{
				Name:        req.Data.Attributes.Name,
				Description: req.Data.Attributes.Description,
				UserID: userID,
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

		response = CreateOrganizationResponse{
			Data: mapOrganizationToResponse(organization),
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}

func GetOrganizationsEndpoint(c echo.Context) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid user token",
		})
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	organizations, err := getOrganizationsByUserID(GetOrganizationsByUserIDServiceRequest{
		UserID: userID,
		Tx:     db,
		MinioClient: minioClient,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, organizationsToResponse(organizations))
}

func GetOrganizationEndpoint(c echo.Context) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid user token",
		})
	}

	// Parse the organization ID from the URL parameter
	paramOrgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    constants.ErrorCodeBadRequest,
			Message: "Invalid organization ID",
		})
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	// Check if user has access to this organization
	hasAccess, err := checkUserOrganizationAccess(CheckUserOrganizationAccessServiceRequest{
		UserID:         userID,
		OrganizationID: uint(paramOrgID),
		Tx:             db,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
	}

	if !hasAccess {
		return c.JSON(http.StatusForbidden, api.ApiError{
			Code:    constants.ErrorCodeForbidden,
			Message: "You don't have access to this organization",
		})
	}

	organization, err := getOrganizationByID(GetOrganizationByIDServiceRequest{
		OrganizationID: uint(paramOrgID),
		Tx:             db,
		MinioClient:    minioClient,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, api.ApiError{
			Code:    constants.ErrorCodeNotFound,
			Message: "Organization not found",
		})
	}

	return c.JSON(http.StatusOK, GetOrganizationResponse{
		Data: mapOrganizationToResponse(organization),
	})
}

func UpdateOrganizationEndpoint(c echo.Context, req UpdateOrganizationRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid user token",
		})
	}

	// Parse the organization ID from the URL parameter
	paramOrgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    constants.ErrorCodeBadRequest,
			Message: "Invalid organization ID",
		})
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	var response UpdateOrganizationResponse

	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		// Check if user has admin access to this organization
		hasAdminAccess, err := checkUserOrganizationAdminAccess(CheckUserOrganizationAdminAccessServiceRequest{
			UserID:         userID,
			OrganizationID: uint(paramOrgID),
			Tx:             db,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    constants.ErrorCodeInternalServerError,
				Message: err.Error(),
			})
		}

		if !hasAdminAccess {
			return c.JSON(http.StatusForbidden, api.ApiError{
				Code:    constants.ErrorCodeForbidden,
				Message: "You don't have admin access to this organization",
			})
		}

		organization, err := updateOrganization(UpdateOrganizationServiceRequest{
			Params: UpdateOrganizationParams{
				OrganizationID: uint(paramOrgID),
				Name:           req.Data.Attributes.Name,
				Description:    req.Data.Attributes.Description,
				Logo:           req.Data.Attributes.Logo,
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

		response = UpdateOrganizationResponse{
			Data: mapOrganizationToResponse(organization),
		}

		return nil
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteOrganizationEndpoint(c echo.Context) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, api.ApiError{
			Code:    constants.ErrorCodeUnauthorized,
			Message: "Invalid user token",
		})
	}

	// Parse the organization ID from the URL parameter
	paramOrgID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.ApiError{
			Code:    constants.ErrorCodeBadRequest,
			Message: "Invalid organization ID",
		})
	}

	db := middleware.GetDB(c)

	// Check if user has admin access to this organization
	hasAdminAccess, err := checkUserOrganizationAdminAccess(CheckUserOrganizationAdminAccessServiceRequest{
		UserID:         userID,
		OrganizationID: uint(paramOrgID),
		Tx:             db,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
	}

	if !hasAdminAccess {
		return c.JSON(http.StatusForbidden, api.ApiError{
			Code:    constants.ErrorCodeForbidden,
			Message: "You don't have admin access to this organization",
		})
	}

	err = deleteOrganization(DeleteOrganizationServiceRequest{
		OrganizationID: uint(paramOrgID),
		Tx:             db,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.ApiError{
			Code:    constants.ErrorCodeInternalServerError,
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Type mappers
func mapOrganizationToResponse(params *OrganizationDto) OrganizationDataWithMeta {
	return OrganizationDataWithMeta{
		OrganizationData: OrganizationData{
			Id:   strconv.FormatUint(uint64(params.Organization.ID), 10),
			Type: constants.ApiTypeOrganization,
			Attributes: OrganizationAttributes{
				Name: params.Organization.Name,
				Description: params.Organization.Description,
			},
		},
		Meta: OrganizationMeta{
			LogoDistributionUrl: params.LogoDistributionUrl,
		},
	}
}

func organizationsToResponse(organizations []*OrganizationDto) GetOrganizationsResponse {
	data := []OrganizationDataWithMeta{}
	for _, org := range organizations {
		data = append(data, mapOrganizationToResponse(org))
	}
	return GetOrganizationsResponse{Data: data}
}
