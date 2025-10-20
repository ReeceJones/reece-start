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


func CreateOrganizationEndpoint(c echo.Context, req CreateOrganizationRequest) error {
	userID, err := middleware.GetUserIDFromJWT(c)
	if err != nil {
		return err
	}

	db := middleware.GetDB(c)
	minioClient := middleware.GetMinioClient(c)
	config := middleware.GetConfig(c)
	stripeClient := middleware.GetStripeClient(c)

	var response CreateOrganizationResponse
	err = db.WithContext(c.Request().Context()).Transaction(func(tx *gorm.DB) error {
		organization, err := createOrganization(CreateOrganizationServiceRequest{
			Params: CreateOrganizationParams{
				Name:        req.Data.Attributes.Name,
				Description: req.Data.Attributes.Description,
				UserID: userID,
				Logo: req.Data.Attributes.Logo,
				ContactEmail: req.Data.Attributes.ContactEmail,
				ContactPhone: req.Data.Attributes.ContactPhone,
				ContactPhoneCountry: req.Data.Attributes.ContactPhoneCountry,
				Locale: req.Data.Attributes.Locale,
				EntityType: req.Data.Attributes.EntityType,
				Address: req.Data.Attributes.Address,
			},
			Tx: tx, // Use the transaction instead of the main db connection
			MinioClient: minioClient,
			Context: c.Request().Context(),
			Config: config,
			StripeClient: stripeClient,
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

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: paramOrgID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationRead},
	}); err != nil {
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

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: paramOrgID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationUpdate},
	}); err != nil {
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

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: paramOrgID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationDelete},
	}); err != nil {
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
	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: query.OrganizationID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationMembershipsList},
	}); err != nil {
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

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: membership.Membership.OrganizationID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationMembershipsList},
	}); err != nil {
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

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: uint(orgId),
		Scopes: []constants.UserScope{constants.UserScopeOrganizationMembershipsCreate},
	}); err != nil {
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

		if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
			OrganizationID: membership.Membership.OrganizationID,
			Scopes: []constants.UserScope{constants.UserScopeOrganizationMembershipsUpdate},
		}); err != nil {
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

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: membership.Membership.OrganizationID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationMembershipsDelete},
	}); err != nil {
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
		if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
			OrganizationID: uint(paramOrgID),
			Scopes: []constants.UserScope{constants.UserScopeOrganizationInvitationsCreate},
		}); err != nil {
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
	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: query.OrganizationID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationInvitationsList},
	}); err != nil {
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

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: invitation.Invitation.OrganizationID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationInvitationsDelete},
	}); err != nil {
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

func CreateStripeOnboardingLinkEndpoint(c echo.Context) error {
	paramOrgID, err := api.ParseOrganizationIDFromString(c.Param("id"))
	if err != nil {
		return err
	}

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: paramOrgID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationStripeUpdate},
	}); err != nil {
		return err
	}

    db := middleware.GetDB(c)
    stripeClient := middleware.GetStripeClient(c)
	config := middleware.GetConfig(c)

	link, err := createStripeOnboardingLink(CreateStripeOnboardingLinkServiceRequest{
        Db:           db,
        StripeClient: stripeClient,
        Context:      c.Request().Context(),
		Config: config,
        Params: CreateStripeOnboardingLinkParams{
            OrganizationID: paramOrgID,
        },
    })
    if err != nil {
        return err
    }

	response := CreateStripeOnboardingLinkResponse{
		Data: StripeAccountLinkData{
			Type: string(constants.ApiTypeStripeAccountLink),
			Attributes: StripeAccountLinkAttributes{
				URL:       link.URL,
				ExpiresAt: link.ExpiresAt,
				Livemode:  link.Livemode,
				AccountID: link.Account,
				CreatedAt: link.Created,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

func CreateStripeDashboardLinkEndpoint(c echo.Context) error {
	paramOrgID, err := api.ParseOrganizationIDFromString(c.Param("id"))
	if err != nil {
		return err
	}

	if err := access.HasOrganizationAccess(c, access.HasOrganizationAccessParams{
		OrganizationID: paramOrgID,
		Scopes: []constants.UserScope{constants.UserScopeOrganizationStripeUpdate},
	}); err != nil {
		return err
	}

    db := middleware.GetDB(c)
    stripeClient := middleware.GetStripeClient(c)

	dashboardURL, err := createStripeDashboardLink(CreateStripeDashboardLinkServiceRequest{
        Db:           db,
        StripeClient: stripeClient,
        Context:      c.Request().Context(),
        Params: CreateStripeDashboardLinkParams{
            OrganizationID: paramOrgID,
        },
    })
    if err != nil {
        return err
    }

	response := CreateStripeDashboardLinkResponse{
		Data: StripeDashboardLinkData{
			Type: string(constants.ApiTypeStripeDashboardLink),
			Attributes: StripeDashboardLinkAttributes{
				URL: dashboardURL,
			},
		},
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
				CommonOrganizationAttributes: CommonOrganizationAttributes{
					Name: params.Organization.Name,
					Description: params.Organization.Description,
					Address: api.Address(params.Organization.Address),
					Locale: params.Organization.Locale,
					ContactEmail: params.Organization.ContactEmail,
					ContactPhone: params.Organization.ContactPhone,
					ContactPhoneCountry: params.Organization.ContactPhoneCountry,
				},
			},
		},
		Meta: OrganizationMeta{
			LogoDistributionUrl: params.LogoDistributionUrl,
			OnboardingStatus: params.Organization.OnboardingStatus,
			Stripe: StripeMeta{
				HasPendingRequirements: params.Organization.Stripe.HasPendingRequirements,
				OnboardingStatus: params.Organization.Stripe.OnboardingStatus,
			},
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
