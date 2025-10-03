package organizations

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/riverqueue/river"
	stripeGo "github.com/stripe/stripe-go/v82"
	"gorm.io/gorm"
	"reece.start/internal/api"
	"reece.start/internal/configuration"
	"reece.start/internal/constants"
	"reece.start/internal/models"
	"reece.start/internal/stripe"
)

// API Types
type OrganizationAttributes struct {
	// Basic information
	Name string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	Address api.Address `json:"address"`

	// Localization fields
	Locale string `json:"locale" validate:"required"`

	// Contact information
	ContactEmail string `json:"contactEmail" validate:"omitempty,email"`
	ContactPhone string `json:"contactPhone" validate:"omitempty"`
	WebsiteUrl string `json:"websiteUrl" validate:"omitempty,url"`
}

type CreateOrganizationAttributes struct {
	// Common fields
	OrganizationAttributes
	Logo string `json:"logo,omitempty" validate:"omitempty,base64"`

	// Onboarding - Basic information
	EntityType string `json:"entityType" validate:"required"`
}

type UpdateOrganizationAttributes struct {
	// Basic information
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=255"`
	Logo *string `json:"logo,omitempty" validate:"omitempty,base64"`
	Address *api.Address `json:"address,omitempty" validate:"omitempty"`

	// Localization fields
	Currency *string `json:"currency,omitempty" validate:"omitempty"`
	Locale *string `json:"locale,omitempty" validate:"omitempty"`

	// Contact information
	ContactEmail *string `json:"contactEmail,omitempty" validate:"omitempty,email"`
	ContactPhone *string `json:"contactPhone,omitempty" validate:"omitempty"`
	WebsiteUrl *string `json:"websiteUrl,omitempty" validate:"omitempty,url"`
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

// Service request/response types
type CreateOrganizationParams struct {
	Name   string
	Description string
	UserID uint
	Logo   string
	ContactEmail string
	ContactPhone string
	WebsiteUrl string
	Locale string
	EntityType string
	Address api.Address
}

type CreateOrganizationServiceRequest struct {
	Params CreateOrganizationParams
	Tx     *gorm.DB
	MinioClient *minio.Client
	Config *configuration.Config
	StripeClient *stripe.Client
	Context context.Context
}

type GetOrganizationsByUserIDServiceRequest struct {
	UserID uint
	Tx     *gorm.DB
	MinioClient *minio.Client
}

type GetOrganizationByIDServiceRequest struct {
	OrganizationID uint
	UserID uint
	Tx             *gorm.DB
	MinioClient    *minio.Client
}

type UpdateOrganizationParams struct {
	OrganizationID uint
	Name           *string
	Description    *string
	Logo           *string
	// TODO: need to add a lot more fields here
}

type UpdateOrganizationServiceRequest struct {
	Params UpdateOrganizationParams
	Tx     *gorm.DB
	MinioClient *minio.Client
}

type DeleteOrganizationServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
}

type CheckUserOrganizationAccessServiceRequest struct {
	UserID         uint
	OrganizationID uint
	Tx             *gorm.DB
}

type CheckUserOrganizationAdminAccessServiceRequest struct {
	UserID         uint
	OrganizationID uint
	Tx             *gorm.DB
}

type GetOrganizationLogoDistributionUrlServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
	MinioClient    *minio.Client
}

type OrganizationDto struct {
	Organization        *models.Organization
	LogoDistributionUrl string
}

type OrganizationMembershipDto struct {
	Membership          *models.OrganizationMembership
	User                *models.User
	UserLogoDistributionUrl string
	Organization        *models.Organization
}

// Organization Membership Service Types
type CreateOrganizationMembershipParams struct {
	UserID         uint
	OrganizationID uint
	Role           string
}

type CreateOrganizationMembershipServiceRequest struct {
	Params CreateOrganizationMembershipParams
	Tx     *gorm.DB
}

type GetOrganizationMembershipsServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
	MinioClient    *minio.Client
}

type GetOrganizationMembershipByIDServiceRequest struct {
	MembershipID uint
	Tx           *gorm.DB
	MinioClient  *minio.Client
}

type UpdateOrganizationMembershipParams struct {
	MembershipID uint
	Role         *string
}

type UpdateOrganizationMembershipServiceRequest struct {
	Params UpdateOrganizationMembershipParams
	Tx     *gorm.DB
}

type DeleteOrganizationMembershipServiceRequest struct {
	MembershipID uint
	Tx           *gorm.DB
}

// Organization Invitation Service Types
type CreateOrganizationInvitationParams struct {
	Email          string
	Role           string
	OrganizationID uint
	InvitingUserID uint
}

type CreateOrganizationInvitationServiceRequest struct {
	Params      CreateOrganizationInvitationParams
	Tx          *gorm.DB
	RiverClient *river.Client[*sql.Tx]
}

type InvitingUserDto struct {
	User                    *models.User
	UserLogoDistributionUrl string
}

type OrganizationInvitationDto struct {
	Invitation   *models.OrganizationInvitation
	Organization *OrganizationDto
	InvitingUser *InvitingUserDto
}

type GetOrganizationInvitationsServiceRequest struct {
	OrganizationID uint
	Tx             *gorm.DB
}

type GetOrganizationInvitationByIDServiceRequest struct {
	InvitationID uuid.UUID
	Tx           *gorm.DB
	MinioClient  *minio.Client
}

type DeleteOrganizationInvitationServiceRequest struct {
	InvitationID uuid.UUID
	Tx           *gorm.DB
}

type AcceptOrganizationInvitationServiceRequest struct {
	InvitationID uuid.UUID
	UserID       uint
	Tx           *gorm.DB
	MinioClient  *minio.Client
}

type DeclineOrganizationInvitationServiceRequest struct {
	InvitationID uuid.UUID
	UserID       uint
	Tx           *gorm.DB
	MinioClient  *minio.Client
}

type UpdateOrganizationStripeInformationServiceRequest struct {
	Organization *models.Organization
	StripeAccount stripeGo.V2CoreAccount
}