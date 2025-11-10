package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	appMiddleware "reece.start/internal/middleware"
	"reece.start/internal/routes"
)

// NewEcho creates and configures a new Echo instance with all middleware and routes
func NewEcho(deps appMiddleware.AppDependencies) *echo.Echo {
	e := echo.New()

	// Add logging middleware
	e.Use(middleware.Logger())

	// Add error handling middleware
	e.Use(appMiddleware.ErrorHandlingMiddleware)

	// Add CORS middleware
	e.Use(middleware.CORS())

	// Add content type middleware
	e.Use(appMiddleware.ContentTypeMiddleware)

	// Add dependency injection middleware
	e.Use(appMiddleware.DependencyInjectionMiddleware(deps))

	// Register all application routes
	routes.Register(e, deps.Config)

	return e
}
