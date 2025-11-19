package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"reece.start/internal/api"
	"reece.start/internal/configuration"
	appMiddleware "reece.start/internal/middleware"
	"reece.start/internal/organizations"
	"reece.start/internal/stripe"
	"reece.start/internal/users"
)

// Register registers all application routes on the provided Echo instance
func Register(e *echo.Echo, config *configuration.Config) {
	// Create authentication middleware
	auth := appMiddleware.JwtAuthMiddleware(config)

	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Public user routes (no authentication required)
	e.POST("/users", api.Validated(users.CreateUserEndpoint))
	e.POST("/users/login", api.Validated(users.LoginEndpoint))

	// Public OAuth routes (no authentication required)
	e.POST("/oauth/google/callback", api.Validated(users.GoogleOAuthCallbackEndpoint))

	// Webhook routes (no authentication required)
	e.POST("/webhooks/stripe/account/snapshot", stripe.StripeSnapshotWebhookEndpoint)
	e.POST("/webhooks/stripe/connect/thin", stripe.StripeThinWebhookEndpoint)

	// Protected user routes
	e.GET("/users/me", users.GetAuthenticatedUserEndpoint, auth)
	e.GET("/users", api.ValidatedQuery(users.GetUsersEndpoint), auth)
	e.POST("/users/me/token", api.Validated(users.CreateAuthenticatedUserTokenEndpoint), auth)
	e.PATCH("/users/:id", api.Validated(users.UpdateUserEndpoint), auth)

	// Protected organization routes
	e.GET("/organizations", organizations.GetOrganizationsEndpoint, auth)
	e.POST("/organizations", api.Validated(organizations.CreateOrganizationEndpoint), auth)
	e.GET("/organizations/:id", organizations.GetOrganizationEndpoint, auth)
	e.PATCH("/organizations/:id", api.Validated(organizations.UpdateOrganizationEndpoint), auth)
	e.DELETE("/organizations/:id", organizations.DeleteOrganizationEndpoint, auth)
	e.POST("/organizations/:id/stripe-onboarding-link", organizations.CreateStripeOnboardingLinkEndpoint, auth)
	e.POST("/organizations/:id/stripe-dashboard-link", organizations.CreateStripeDashboardLinkEndpoint, auth)
	e.GET("/organizations/:id/subscription", stripe.GetSubscriptionEndpoint, auth)
	e.POST("/organizations/:id/checkout-session", api.Validated(stripe.CreateCheckoutSessionEndpoint), auth)
	e.POST("/organizations/:id/billing-portal-session", api.Validated(stripe.CreateBillingPortalSessionEndpoint), auth)

	// Protected organization membership routes
	e.GET("/organization-memberships", api.ValidatedQuery(organizations.GetOrganizationMembershipsEndpoint), auth)
	e.GET("/organization-memberships/:id", organizations.GetOrganizationMembershipEndpoint, auth)
	e.POST("/organization-memberships", api.Validated(organizations.CreateOrganizationMembershipEndpoint), auth)
	e.PATCH("/organization-memberships/:id", api.Validated(organizations.UpdateOrganizationMembershipEndpoint), auth)
	e.DELETE("/organization-memberships/:id", organizations.DeleteOrganizationMembershipEndpoint, auth)

	// Protected organization invitation routes
	e.POST("/organization-invitations", api.Validated(organizations.InviteToOrganizationEndpoint), auth)
	e.POST("/organization-invitations/:id/accept", api.Validated(organizations.AcceptOrganizationInvitationEndpoint), auth)
	e.POST("/organization-invitations/:id/decline", api.Validated(organizations.DeclineOrganizationInvitationEndpoint), auth)
	e.GET("/organization-invitations", api.ValidatedQuery(organizations.GetOrganizationInvitationsEndpoint), auth)
	e.GET("/organization-invitations/:id", organizations.GetOrganizationInvitationEndpoint, auth)
	e.DELETE("/organization-invitations/:id", organizations.DeleteOrganizationInvitationEndpoint, auth)
}
