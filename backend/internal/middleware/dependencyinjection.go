package middleware

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
	"reece.start/internal/configuration"
)

type AppDependencies struct {
	Config *configuration.Config
	DB *gorm.DB
	MinioClient *minio.Client
	RiverClient *river.Client[*sql.Tx]
}

// Middleware to inject config and database into context
func DependencyInjectionMiddleware(dependencies AppDependencies) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", dependencies.Config)
			c.Set("db", dependencies.DB)
			c.Set("minioClient", dependencies.MinioClient)
			c.Set("riverClient", dependencies.RiverClient)
			return next(c)
		}
	}
}

// GetConfigAndDB extracts the config and database from the Echo context
func GetConfig(c echo.Context) *configuration.Config {
	return c.Get("config").(*configuration.Config)
}

func GetDB(c echo.Context) *gorm.DB {
	return c.Get("db").(*gorm.DB)
}

func GetMinioClient(c echo.Context) *minio.Client {
	return c.Get("minioClient").(*minio.Client)
}

func GetRiverClient(c echo.Context) *river.Client[*sql.Tx] {
	return c.Get("riverClient").(*river.Client[*sql.Tx])
}
