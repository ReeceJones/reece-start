package configuration

import (
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Host string `env:"HOST" envDefault:"localhost"`
	Port string `env:"PORT" envDefault:"8080"`
	
	FrontendUrl string `env:"FRONTEND_URL" envDefault:"https://localhost:4040"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:3000"`

	DatabaseUri string `env:"DATABASE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/postgres"`

	JwtSecret   string `env:"JWT_SECRET" envDefault:"secret"`
	JwtIssuer   string `env:"JWT_ISSUER" envDefault:"reece-start"`
	JwtAudience string `env:"JWT_AUDIENCE" envDefault:"https://reece.start"`
	JwtExpirationTime int `env:"JWT_EXPIRATION_TIME" envDefault:"86400"`

	StorageEndpoint        string `env:"STORAGE_ENDPOINT" envDefault:"localhost:9000"`
	StorageAccessKeyId     string `env:"STORAGE_ACCESS_KEY_ID" envDefault:"minioadmin"`
	StorageSecretAccessKey string `env:"STORAGE_SECRET_ACCESS_KEY" envDefault:"minioadmin"`
	StorageUseSSL          bool   `env:"STORAGE_USE_SSL" envDefault:"false"`

	ResendApiKey string `env:"RESEND_API_KEY" envDefault:""`
	EnableEmail bool `env:"ENABLE_EMAIL" envDefault:"false"`

	GoogleOAuthClientId string `env:"GOOGLE_OAUTH_CLIENT_ID" envDefault:""`
	GoogleOAuthClientSecret string `env:"GOOGLE_OAUTH_CLIENT_SECRET" envDefault:""`

	StripeWebhookSecret string `env:"STRIPE_WEBHOOK_SECRET" envDefault:""`
	StripeSecretKey string `env:"STRIPE_SECRET_KEY" envDefault:""`
	StripeProPlanPriceId string `env:"STRIPE_PRO_PLAN_PRICE_ID" envDefault:""`
	StripeProPlanProductId string `env:"STRIPE_PRO_PLAN_PRODUCT_ID" envDefault:""`
	StripeBillingPortalConfigurationId string `env:"STRIPE_BILLING_PORTAL_CONFIGURATION_ID" envDefault:""`
	StripeEnableACHDebitPayments bool `env:"STRIPE_ENABLE_ACH_DEBIT_PAYMENTS" envDefault:"false"`
}

func LoadEnvironmentVariables() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Detected .env file, loading environment variables")
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file, %s", err)
		}
	}

	var config Config
	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("Error parsing environment variables, %s", err)
	}

	return &config, nil
}
