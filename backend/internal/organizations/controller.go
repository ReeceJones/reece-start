package organizations

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/constants"
	"reece.start/internal/middleware"
	"reece.start/internal/models"
)

// API Types
type OrganizationAttributes struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
}

type CreateOrganizationAttributes struct {
	OrganizationAttributes
}

type UpdateOrganizationAttributes struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
}

type OrganizationData struct {
	Id         string                `json:"id"`
	Type       constants.ApiType     `json:"type"`
	Attributes OrganizationAttributes `json:"attributes"`
}

type CreateOrganizationRequest struct {
	Data struct {
		Attributes CreateOrganizationAttributes `json:"attributes"`
	} `json:"data"`
}

type CreateOrganizationResponse struct {
	Data OrganizationData `json:"data"`
}

type UpdateOrganizationRequest struct {
	Data struct {
		Attributes UpdateOrganizationAttributes `json:"attributes"`
	} `json:"data"`
}

type UpdateOrganizationResponse struct {
	Data OrganizationData `json:"data"`
}

type GetOrganizationsResponse struct {
	Data []OrganizationData `json:"data"`
}

type GetOrganizationResponse struct {
	Data OrganizationData `json:"data"`
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

	var response CreateOrganizationResponse
	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		organization, err := createOrganization(CreateOrganizationServiceRequest{
			Params: CreateOrganizationParams{
				Name:        req.Data.Attributes.Name,
				Description: req.Data.Attributes.Description,
				UserID: userID,
			},
			Tx: db,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    constants.ErrorCodeInternalServerError,
				Message: err.Error(),
			})
		}

		response = CreateOrganizationResponse{
			Data: mapOrganizationToApiData(organization),
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

	organizations, err := getOrganizationsByUserID(GetOrganizationsByUserIDServiceRequest{
		UserID: userID,
		Tx:     db,
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
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, api.ApiError{
			Code:    constants.ErrorCodeNotFound,
			Message: "Organization not found",
		})
	}

	return c.JSON(http.StatusOK, mapOrganizationToApiData(organization))
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
			},
			Tx: db,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.ApiError{
				Code:    constants.ErrorCodeInternalServerError,
				Message: err.Error(),
			})
		}

		response = UpdateOrganizationResponse{
			Data: mapOrganizationToApiData(organization),
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
func mapOrganizationToApiData(organization *models.Organization) OrganizationData {
	return OrganizationData{
		Id:   strconv.FormatUint(uint64(organization.ID), 10),
		Type: constants.ApiTypeOrganization,
		Attributes: OrganizationAttributes{
			Name: organization.Name,
			Description: organization.Description,
		},
	}
}

func organizationsToResponse(organizations []models.Organization) GetOrganizationsResponse {
	data := []OrganizationData{}
	for _, org := range organizations {
		data = append(data, mapOrganizationToApiData(&org))
	}
	return GetOrganizationsResponse{Data: data}
}
