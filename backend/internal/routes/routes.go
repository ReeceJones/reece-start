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
	// Health check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// Public user routes (no authentication required)
	e.POST("/users", api.Validated(users.CreateUserEndpoint))
	e.POST("/users/login", api.Validated(users.LoginEndpoint))

	// Public OAuth routes (no authentication required)
	e.POST("/oauth/google/callback", api.Validated(users.GoogleOAuthCallbackEndpoint))

	// Webhook routes
	e.POST("/webhooks/stripe/snapshot", stripe.StripeSnapshotWebhookEndpoint)
	e.POST("/webhooks/stripe/thin", stripe.StripeThinWebhookEndpoint)

	// Protected routes (authentication required)
	protected := e.Group("")
	protected.Use(appMiddleware.JwtAuthMiddleware(config))

	protected.GET("/users/me", users.GetAuthenticatedUserEndpoint)
	protected.GET("/users", api.ValidatedQuery(users.GetUsersEndpoint))
	protected.POST("/users/me/token", api.Validated(users.CreateAuthenticatedUserTokenEndpoint))
	protected.PATCH("/users/:id", api.Validated(users.UpdateUserEndpoint))

	protected.GET("/organizations", organizations.GetOrganizationsEndpoint)
	protected.POST("/organizations", api.Validated(organizations.CreateOrganizationEndpoint))
	protected.GET("/organizations/:id", organizations.GetOrganizationEndpoint)
	protected.PATCH("/organizations/:id", api.Validated(organizations.UpdateOrganizationEndpoint))
	protected.DELETE("/organizations/:id", organizations.DeleteOrganizationEndpoint)
	protected.POST("/organizations/:id/stripe-onboarding-link", organizations.CreateStripeOnboardingLinkEndpoint)
	protected.POST("/organizations/:id/stripe-dashboard-link", organizations.CreateStripeDashboardLinkEndpoint)
	protected.GET("/organizations/:id/subscription", stripe.GetSubscriptionEndpoint)
	protected.POST("/organizations/:id/checkout-session", api.Validated(stripe.CreateCheckoutSessionEndpoint))
	protected.POST("/organizations/:id/billing-portal-session", api.Validated(stripe.CreateBillingPortalSessionEndpoint))

	protected.GET("/organization-memberships", api.ValidatedQuery(organizations.GetOrganizationMembershipsEndpoint))
	protected.GET("/organization-memberships/:id", organizations.GetOrganizationMembershipEndpoint)
	protected.POST("/organization-memberships", api.Validated(organizations.CreateOrganizationMembershipEndpoint))
	protected.PATCH("/organization-memberships/:id", api.Validated(organizations.UpdateOrganizationMembershipEndpoint))
	protected.DELETE("/organization-memberships/:id", organizations.DeleteOrganizationMembershipEndpoint)

	protected.POST("/organization-invitations", api.Validated(organizations.InviteToOrganizationEndpoint))
	protected.POST("/organization-invitations/:id/accept", api.Validated(organizations.AcceptOrganizationInvitationEndpoint))
	protected.POST("/organization-invitations/:id/decline", api.Validated(organizations.DeclineOrganizationInvitationEndpoint))
	protected.GET("/organization-invitations", api.ValidatedQuery(organizations.GetOrganizationInvitationsEndpoint))
	protected.GET("/organization-invitations/:id", organizations.GetOrganizationInvitationEndpoint)
	protected.DELETE("/organization-invitations/:id", organizations.DeleteOrganizationInvitationEndpoint)
}
