package organizations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"reece.start/internal/access"
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

// Organization Membership API Types
type OrganizationMembershipAttributes struct {
	Role string `json:"role" validate:"required,oneof=admin member"`
}

type UpdateOrganizationMembershipAttributes struct {
	Role *string `json:"role,omitempty" validate:"omitempty,oneof=admin member"`
}

type UserRelationshipDataObject struct {
	Id   string            `json:"id" validate:"required"`
	Type constants.ApiType `json:"type" validate:"required,oneof=user"`
}

type OrganizationRelationshipDataObject struct {
	Id   string            `json:"id" validate:"required"`
	Type constants.ApiType `json:"type" validate:"required,oneof=organization"`
}

type UserRelationshipData struct {
	Data UserRelationshipDataObject `json:"data" validate:"required"`
}

type OrganizationRelationshipData struct {
	Data OrganizationRelationshipDataObject `json:"data" validate:"required"`
}

type OrganizationMembershipRelationships struct {
	User         UserRelationshipData         `json:"user"`
	Organization OrganizationRelationshipData `json:"organization"`
}

type CreateOrganizationMembershipRelationships struct {
	User         UserRelationshipData         `json:"user" validate:"required"`
	Organization OrganizationRelationshipData `json:"organization" validate:"required"`
}

type OrganizationMembershipData struct {
	Id            string                              `json:"id"`
	Type          constants.ApiType                   `json:"type"`
	Attributes    OrganizationMembershipAttributes    `json:"attributes"`
	Relationships OrganizationMembershipRelationships `json:"relationships"`
}

type CreateOrganizationMembershipRequest struct {
	Data struct {
		Attributes    OrganizationMembershipAttributes           `json:"attributes"`
		Relationships CreateOrganizationMembershipRelationships `json:"relationships"`
	} `json:"data"`
}

type CreateOrganizationMembershipResponse struct {
	Data OrganizationMembershipData `json:"data"`
}

type UpdateOrganizationMembershipRequest struct {
	Data struct {
		Attributes UpdateOrganizationMembershipAttributes `json:"attributes"`
	} `json:"data"`
}

type UpdateOrganizationMembershipResponse struct {
	Data OrganizationMembershipData `json:"data"`
}

type GetOrganizationMembershipsResponse struct {
	Data     []OrganizationMembershipData `json:"data"`
	Included []interface{}                `json:"included,omitempty"`
}

type GetOrganizationMembershipResponse struct {
	Data     OrganizationMembershipData `json:"data"`
	Included []interface{}              `json:"included,omitempty"`
}

type GetOrganizationMembershipsQuery struct {
	OrganizationID uint `query:"organizationId" validate:"required,min=1"`
}

type GetOrganizationInvitationsQuery struct {
	OrganizationID uint `query:"organizationId" validate:"required,min=1"`
}

type InviteToOrganizationRelationships struct {
	Organization OrganizationRelationshipData `json:"organization" validate:"required"`
}

type InviteToOrganizationRequest struct {
	Data struct {
		Type          constants.ApiType                   `json:"type" validate:"required,oneof=organization-invitation"`
		Attributes    InviteToOrganizationAttributes    `json:"attributes"`
		Relationships InviteToOrganizationRelationships `json:"relationships"`
	} `json:"data"`
}

// User data for included section
type UserIncludedAttributes struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserIncludedMeta struct {
	LogoDistributionUrl string `json:"logoDistributionUrl,omitempty"`
}

type UserIncludedData struct {
	Id         string                 `json:"id"`
	Type       constants.ApiType      `json:"type"`
	Attributes UserIncludedAttributes `json:"attributes"`
	Meta       UserIncludedMeta       `json:"meta,omitempty"`
}

// Organization data for included section
type OrganizationIncludedAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OrganizationIncludedMeta struct {
	LogoDistributionUrl string `json:"logoDistributionUrl,omitempty"`
}

type OrganizationIncludedData struct {
	Id         string                          `json:"id"`
	Type       constants.ApiType               `json:"type"`
	Attributes OrganizationIncludedAttributes `json:"attributes"`
	Meta       OrganizationIncludedMeta       `json:"meta,omitempty"`
}

type OrganizationInvitationAttributes struct {
	Email string `json:"email" validate:"required,email"`
	Role string `json:"role" validate:"required,oneof=admin member"`
	Status string `json:"status" validate:"required,oneof=pending accepted declined expired revoked"`
}

type OrganizationInvitationRelationships struct {
	Organization OrganizationRelationshipData `json:"organization" validate:"required"`
	InvitingUser UserRelationshipData `json:"invitingUser" validate:"required"`
}

type OrganizationInvitationData struct {
	Id         string                `json:"id"`
	Type       constants.ApiType     `json:"type" validate:"required,oneof=organization-invitation"`
	Attributes OrganizationInvitationAttributes `json:"attributes"`
	Relationships OrganizationInvitationRelationships `json:"relationships"`
}

type InviteToOrganizationAttributes struct {
	Email string `json:"email" validate:"required,email"`
	Role string `json:"role" validate:"required,oneof=admin member"`
}

type InviteToOrganizationResponse struct {
	Data OrganizationInvitationData `json:"data"`
}

type GetOrganizationInvitationsResponse struct {
	Data []OrganizationInvitationData `json:"data"`
}

type GetOrganizationInvitationResponse struct {
	Data     OrganizationInvitationData `json:"data"`
	Included []interface{}              `json:"included,omitempty"`
}

type OrganizationInvitationIdentifier struct {
	Id   string            `json:"id" validate:"required"`
	Type constants.ApiType `json:"type" validate:"required,oneof=organization-invitation"`
}

type AcceptOrganizationInvitationRequest struct {
	Data OrganizationInvitationIdentifier `json:"data" validate:"required"`
}

type DeclineOrganizationInvitationRequest struct {
	Data OrganizationInvitationIdentifier `json:"data" validate:"required"`
}

type AcceptOrganizationInvitationResponse struct {
	Data     OrganizationInvitationData `json:"data"`
	Included []interface{}              `json:"included,omitempty"`
}

type DeclineOrganizationInvitationResponse struct {
	Data     OrganizationInvitationData `json:"data"`
	Included []interface{}              `json:"included,omitempty"`
}

func CreateOrganizationEndpoint(c echo.Context, req CreateOrganizationRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
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
			return err
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
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	organizations, err := getOrganizationsByUserID(GetOrganizationsByUserIDServiceRequest{
		UserID: userID,
		Tx:     db,
		MinioClient: minioClient,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, organizationsToResponse(organizations))
}

func GetOrganizationEndpoint(c echo.Context) error {
	// Parse the organization ID from the URL parameter
	paramOrgID, err := api.ParseOrganizationIDFromParams(c)
	if err != nil {
		return err
	}

	if err := access.HasAccessToOrganization(c, paramOrgID); err != nil {
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	organization, err := getOrganizationByID(GetOrganizationByIDServiceRequest{
		OrganizationID: paramOrgID,
		Tx:             db,
		MinioClient:    minioClient,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetOrganizationResponse{
		Data: mapOrganizationToResponse(organization),
	})
}

func UpdateOrganizationEndpoint(c echo.Context, req UpdateOrganizationRequest) error {
	paramOrgID, err := api.ParseOrganizationIDFromParams(c)
	if err != nil {
		return err
	}

	if err := access.HasAdminAccessToOrganization(c, paramOrgID); err != nil {
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	var response UpdateOrganizationResponse

	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
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
			return err
		}

		response = UpdateOrganizationResponse{
			Data: mapOrganizationToResponse(organization),
		}

		return nil
	})

	if err != nil {
		return err // Middleware will handle all error types
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteOrganizationEndpoint(c echo.Context) error {
	paramOrgID, err := api.ParseOrganizationIDFromParams(c)
	if err != nil {
		return err
	}

	if err := access.HasAdminAccessToOrganization(c, paramOrgID); err != nil {
		return err
	}

	db := middleware.GetDB(c)

	err = deleteOrganization(DeleteOrganizationServiceRequest{
		OrganizationID: uint(paramOrgID),
		Tx:             db,
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// Organization Membership Endpoints
func GetOrganizationMembershipsEndpoint(c echo.Context, query GetOrganizationMembershipsQuery) error {
	if err := access.HasAccessToOrganization(c, query.OrganizationID); err != nil {
		return err
	}

	db := middleware.GetDB(c)

	memberships, err := getOrganizationMemberships(GetOrganizationMembershipsServiceRequest{
		OrganizationID: query.OrganizationID,
		Tx:             db,
		MinioClient:    middleware.GetMinioClient(c),
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapMembershipsToResponseWithIncluded(memberships))
}

func GetOrganizationMembershipEndpoint(c echo.Context) error {
	paramMembershipID, err := api.ParseMembershipIDFromParams(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	// First get the membership to check the organization
	membership, err := getOrganizationMembershipByID(GetOrganizationMembershipByIDServiceRequest{
		MembershipID: uint(paramMembershipID),
		Tx:           db,
		MinioClient:  minioClient,
	})

	if err != nil {
		return err
	}

	if err := access.HasAccessToOrganization(c, membership.Membership.OrganizationID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetOrganizationMembershipResponse{
		Data:     mapMembershipToResponse(membership),
		Included: []interface{}{mapUserToIncludedData(membership)},
	})
}

func CreateOrganizationMembershipEndpoint(c echo.Context, req CreateOrganizationMembershipRequest) error {
	orgId, err := strconv.ParseUint(req.Data.Relationships.Organization.Data.Id, 10, 32)
	if err != nil {
		return err
	}

	if err := access.HasAdminAccessToOrganization(c, uint(orgId)); err != nil {
		return err
	}

	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)

	var response CreateOrganizationMembershipResponse

	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		membership, err := createOrganizationMembership(CreateOrganizationMembershipServiceRequest{
			Params: CreateOrganizationMembershipParams{
				UserID:         uint(userID),
				OrganizationID: uint(orgId),
				Role:           req.Data.Attributes.Role,
			},
			Tx: tx,
		})

		if err != nil {
			return err
		}

		response = CreateOrganizationMembershipResponse{
			Data: mapMembershipToResponse(membership),
		}

		return nil
	})

	if err != nil {
		return err // Middleware will handle all error types
	}

	return c.JSON(http.StatusCreated, response)
}

func UpdateOrganizationMembershipEndpoint(c echo.Context, req UpdateOrganizationMembershipRequest) error {
	paramMembershipID, err := api.ParseMembershipIDFromParams(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)

	var response UpdateOrganizationMembershipResponse

	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		// First get the membership to check the organization
		membership, err := getOrganizationMembershipByID(GetOrganizationMembershipByIDServiceRequest{
			MembershipID: uint(paramMembershipID),
			Tx:           tx,
			MinioClient:  nil, // Not needed for update operation
		})

		if err != nil {
			return api.ErrMembershipNotFound
		}

		if err := access.HasAdminAccessToOrganization(c, membership.Membership.OrganizationID); err != nil {
			return err
		}

		updatedMembership, err := updateOrganizationMembership(UpdateOrganizationMembershipServiceRequest{
			Params: UpdateOrganizationMembershipParams{
				MembershipID: uint(paramMembershipID),
				Role:         req.Data.Attributes.Role,
			},
			Tx: tx,
		})

		if err != nil {
			return err
		}

		response = UpdateOrganizationMembershipResponse{
			Data: mapMembershipToResponse(updatedMembership),
		}

		return nil
	})

	if err != nil {
		return err // Middleware will handle all error types
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteOrganizationMembershipEndpoint(c echo.Context) error {
	paramMembershipID, err := api.ParseMembershipIDFromParams(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)

	// First get the membership to check the organization
	membership, err := getOrganizationMembershipByID(GetOrganizationMembershipByIDServiceRequest{
		MembershipID: uint(paramMembershipID),
		Tx:           db,
		MinioClient:  nil, // Not needed for delete operation
	})

	if err != nil {
		return err
	}

	if err := access.HasAdminAccessToOrganization(c, membership.Membership.OrganizationID); err != nil {
		return err
	}

	err = deleteOrganizationMembership(DeleteOrganizationMembershipServiceRequest{
		MembershipID: uint(paramMembershipID),
		Tx:           db,
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func InviteToOrganizationEndpoint(c echo.Context, req InviteToOrganizationRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	// Parse the organization ID from the relationships
	paramOrgID, err := api.ParseOrganizationIDFromString(req.Data.Relationships.Organization.Data.Id)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)
	riverClient := middleware.GetRiverClient(c)

	var response InviteToOrganizationResponse

	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		if err := access.HasAdminAccessToOrganization(c, uint(paramOrgID)); err != nil {
			return err
		}

		invitation, err := createOrganizationInvitation(CreateOrganizationInvitationServiceRequest{
			Params: CreateOrganizationInvitationParams{
				Email:          req.Data.Attributes.Email,
				Role:           req.Data.Attributes.Role,
				OrganizationID: uint(paramOrgID),
				InvitingUserID: userID,
			},
			Tx:          tx,
			RiverClient: riverClient,
		})

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return api.ErrInvitationAlreadyExists
			}

			return err
		}

		response = InviteToOrganizationResponse{
			Data: mapInvitationToResponse(invitation),
		}

		return nil
	})

	if err != nil {
		return err // Middleware will handle all error types
	}

	return c.JSON(http.StatusCreated, response)
}

func GetOrganizationInvitationsEndpoint(c echo.Context, query GetOrganizationInvitationsQuery) error {
	if err := access.HasAdminAccessToOrganization(c, query.OrganizationID); err != nil {
		return err
	}

	db := middleware.GetDB(c)

	invitations, err := getOrganizationInvitations(GetOrganizationInvitationsServiceRequest{
		OrganizationID: query.OrganizationID,
		Tx:             db,
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapInvitationsToResponse(invitations))
}

func GetOrganizationInvitationEndpoint(c echo.Context) error {
	// Parse the invitation ID from the URL parameter
	paramInvitationID, err := api.ParseOrganizationInvitationIDFromParams(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	// First get the invitation to check the organization
	invitation, err := getOrganizationInvitationByID(GetOrganizationInvitationByIDServiceRequest{
		InvitationID: paramInvitationID,
		Tx:           db,
		MinioClient:  minioClient,
	})

	if err != nil {
		return err
	}


	included := []interface{}{
		mapOrganizationToIncludedData(invitation.Organization),
		mapInvitingUserToIncludedData(invitation.InvitingUser),
	}

	return c.JSON(http.StatusOK, GetOrganizationInvitationResponse{
		Data:     mapInvitationToResponse(invitation),
		Included: included,
	})
}

func DeleteOrganizationInvitationEndpoint(c echo.Context) error {
	// Parse the invitation ID from the URL parameter
	paramInvitationID, err := api.ParseOrganizationInvitationIDFromParams(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	// First get the invitation to check the organization
	invitation, err := getOrganizationInvitationByID(GetOrganizationInvitationByIDServiceRequest{
		InvitationID: paramInvitationID,
		Tx:           db,
		MinioClient:  minioClient,
	})

	if err != nil {
		return err
	}

	if err := access.HasAdminAccessToOrganization(c, invitation.Invitation.OrganizationID); err != nil {
		return err
	}

	err = deleteOrganizationInvitation(DeleteOrganizationInvitationServiceRequest{
		InvitationID: paramInvitationID,
		Tx:           db,
	})

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func AcceptOrganizationInvitationEndpoint(c echo.Context, req AcceptOrganizationInvitationRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	// Parse the invitation ID from the URL parameter
	paramInvitationID, err := api.ParseOrganizationInvitationIDFromParams(c)
	if err != nil {
		return err
	}

	// Validate that the request body ID matches the URL parameter
	if req.Data.Id != paramInvitationID.String() {
		return api.ErrInvalidInvitationID
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	var response AcceptOrganizationInvitationResponse

	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		invitation, err := acceptOrganizationInvitation(AcceptOrganizationInvitationServiceRequest{
			InvitationID: paramInvitationID,
			UserID:       userID,
			Tx:           tx,
			MinioClient:  minioClient,
		})

		if err != nil {
			return err
		}

		included := []interface{}{
			mapOrganizationToIncludedData(invitation.Organization),
			mapInvitingUserToIncludedData(invitation.InvitingUser),
		}

		response = AcceptOrganizationInvitationResponse{
			Data:     mapInvitationToResponse(invitation),
			Included: included,
		}

		return nil
	})

	if err != nil {
		return err // Middleware will handle all error types
	}

	return c.JSON(http.StatusOK, response)
}

func DeclineOrganizationInvitationEndpoint(c echo.Context, req DeclineOrganizationInvitationRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	// Parse the invitation ID from the URL parameter
	paramInvitationID, err := api.ParseOrganizationInvitationIDFromParams(c)
	if err != nil {
		return err
	}

	// Validate that the request body ID matches the URL parameter
	if req.Data.Id != paramInvitationID.String() {
		return api.ErrInvalidInvitationID
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)

	var response DeclineOrganizationInvitationResponse

	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		invitation, err := declineOrganizationInvitation(DeclineOrganizationInvitationServiceRequest{
			InvitationID: paramInvitationID,
			UserID:       userID,
			Tx:           tx,
			MinioClient:  minioClient,
		})

		if err != nil {
			return err
		}

		included := []interface{}{
			mapOrganizationToIncludedData(invitation.Organization),
			mapInvitingUserToIncludedData(invitation.InvitingUser),
		}

		response = DeclineOrganizationInvitationResponse{
			Data:     mapInvitationToResponse(invitation),
			Included: included,
		}

		return nil
	})

	if err != nil {
		return err // Middleware will handle all error types
	}

	return c.JSON(http.StatusOK, response)
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

func mapMembershipToResponse(membershipDto *OrganizationMembershipDto) OrganizationMembershipData {
	return OrganizationMembershipData{
		Id:   strconv.FormatUint(uint64(membershipDto.Membership.ID), 10),
		Type: constants.ApiTypeOrganizationMembership,
		Attributes: OrganizationMembershipAttributes{
			Role: membershipDto.Membership.Role,
		},
		Relationships: OrganizationMembershipRelationships{
			User: UserRelationshipData{
				Data: UserRelationshipDataObject{
					Id:   strconv.FormatUint(uint64(membershipDto.User.ID), 10),
					Type: constants.ApiTypeUser,
				},
			},
			Organization: OrganizationRelationshipData{
				Data: OrganizationRelationshipDataObject{
					Id:   strconv.FormatUint(uint64(membershipDto.Organization.ID), 10),
					Type: constants.ApiTypeOrganization,
				},
			},
		},
	}
}

func mapUserToIncludedData(membershipDto *OrganizationMembershipDto) UserIncludedData {
	return UserIncludedData{
		Id:   strconv.FormatUint(uint64(membershipDto.User.ID), 10),
		Type: constants.ApiTypeUser,
		Attributes: UserIncludedAttributes{
			Name:  membershipDto.User.Name,
			Email: membershipDto.User.Email,
		},
		Meta: UserIncludedMeta{
			LogoDistributionUrl: membershipDto.UserLogoDistributionUrl,
		},
	}
}

func mapOrganizationToIncludedData(organizationDto *OrganizationDto) OrganizationIncludedData {
	return OrganizationIncludedData{
		Id:   strconv.FormatUint(uint64(organizationDto.Organization.ID), 10),
		Type: constants.ApiTypeOrganization,
		Attributes: OrganizationIncludedAttributes{
			Name:        organizationDto.Organization.Name,
			Description: organizationDto.Organization.Description,
		},
		Meta: OrganizationIncludedMeta{
			LogoDistributionUrl: organizationDto.LogoDistributionUrl,
		},
	}
}

func mapInvitingUserToIncludedData(invitingUserDto *InvitingUserDto) UserIncludedData {
	return UserIncludedData{
		Id:   strconv.FormatUint(uint64(invitingUserDto.User.ID), 10),
		Type: constants.ApiTypeUser,
		Attributes: UserIncludedAttributes{
			Name:  invitingUserDto.User.Name,
			Email: invitingUserDto.User.Email,
		},
		Meta: UserIncludedMeta{
			LogoDistributionUrl: invitingUserDto.UserLogoDistributionUrl,
		},
	}
}

func mapMembershipsToResponseWithIncluded(membershipDtos []*OrganizationMembershipDto) GetOrganizationMembershipsResponse {
	data := []OrganizationMembershipData{}
	included := []interface{}{}
	userMap := make(map[string]bool) // To avoid duplicate users in included section

	for _, membershipDto := range membershipDtos {
		data = append(data, mapMembershipToResponse(membershipDto))

		// Add user to included section if not already added
		userID := strconv.FormatUint(uint64(membershipDto.User.ID), 10)
		if !userMap[userID] {
			userIncluded := mapUserToIncludedData(membershipDto)
			included = append(included, userIncluded)
			userMap[userID] = true
		}
	}

	return GetOrganizationMembershipsResponse{
		Data:     data,
		Included: included,
	}
}

func mapInvitationToResponse(invitationDto *OrganizationInvitationDto) OrganizationInvitationData {
	return OrganizationInvitationData{
		Id:   invitationDto.Invitation.ID.String(),
		Type: constants.ApiTypeOrganizationInvitation,
		Attributes: OrganizationInvitationAttributes{
			Email: invitationDto.Invitation.Email,
			Role: invitationDto.Invitation.Role,
			Status: invitationDto.Invitation.Status,
		},
		Relationships: OrganizationInvitationRelationships{
			Organization: OrganizationRelationshipData{
				Data: OrganizationRelationshipDataObject{
					Id:   strconv.FormatUint(uint64(invitationDto.Invitation.OrganizationID), 10),
					Type: constants.ApiTypeOrganization,
				},
			},
			InvitingUser: UserRelationshipData{
				Data: UserRelationshipDataObject{
					Id:   strconv.FormatUint(uint64(invitationDto.Invitation.InvitingUserID), 10),
					Type: constants.ApiTypeUser,
				},
			},
		},
	}
}

func mapInvitationsToResponse(invitations []*OrganizationInvitationDto) GetOrganizationInvitationsResponse {
	data := []OrganizationInvitationData{}
	for _, invitation := range invitations {
		data = append(data, mapInvitationToResponse(invitation))
	}
	return GetOrganizationInvitationsResponse{Data: data}
}
