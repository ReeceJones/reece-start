# Railway Deployment Setup

This document describes how to setup Railway deployment for the application.

## Overview

![Railway project](/assets/railway-overview.png)

The railway deployment works by exposing the backend and frontend through a reverse proxy, using templates for infrastructure deployments.

### Favicon links

- [Frontend](https://svelte.dev/favicon.png)
- [Backend](https://echo.labstack.com/img/favicon.ico)
- [Caddy](https://caddyserver.com/resources/images/favicon.png)

## Project Configuration

1. Click "Settings" on the top right corner
2. Navigate to "Shared Variables"
3. For your environment enter the environment variables from the table below

| Variable Name                            | Value                     |
| ---------------------------------------- | ------------------------- |
| `ENABLE_EMAIL`                           | `true` or `false`         |
| `GOOGLE_OAUTH_CLIENT_ID`                 | `your_client_._id`        |
| `GOOGLE_OAUTH_CLIENT_SECRET`             | `your_client_secret`      |
| `JWT_AUDIENCE`                           | `https://your.domain.com` |
| `RESEND_API_KEY`                         | `your_api_key`            |
| `STRIPE_BILLING_PORTAL_CONFIGURATION_ID` | `your_configuration_id`   |
| `STRIPE_PRO_PLAN_PRICE_ID`               | `your_price_id`           |
| `STRIPE_PRO_PLAN_PRODUCT_ID`             | `your_product_id`         |
| `STRIPE_SECRET_KEY`                      | `your_secret_key`         |
| `STRIPE_WEBHOOK_SECRET`                  | `your_webhook_secret`     |

## Postgres Configuration

1. Click "Create" in the top right corner of the project dashboard
2. Select "Database"
3. Select "Postgres"

## MinIO Configuration

1. Click "Create" in the top right corner of the project dashboard
2. Select "Template"
3. Search for "MinIO"
4. Select "MinIO" by "Brody's Projects"

## Caddy Configuration

1. Click "Create" in the top right corner of the project dashboard
2. Select "GitHub Repo"
3. Select your repo
4. Click on the newly added service
5. Click on the "Settings" tab
6. Navigate to "Source" and select "Add route directory"
7. Enter "/caddy"
8. Click on the "Variables" tab and enter the environment variables from the env file below.
9. Press "Deploy" at the top of the screen

```env
API_HOST="your.domain.com"
API_URL="${{backend.RAILWAY_PRIVATE_DOMAIN}}:8080"
FRONTEND_HOST="your.domain.com"
FRONTEND_URL="${{frontend.RAILWAY_PRIVATE_DOMAIN}}:8080"
```

## Frontend Configuration

1. Click "Create" in the top right corner of the project dashboard
2. Select "GitHub Repo"
3. Select your repo
4. Click on the newly added service
5. Click on the "Settings" tab
6. Navigate to "Source" and select "Add route directory"
7. Enter "/frontend"
8. Click on the "Variables" tab and enter the environment variables from the env file below.
9. Press "Deploy" at the top of the screen

```env
BODY_SIZE_LIMIT="64M"
PUBLIC_GOOGLE_OAUTH_CLIENT_ID="${{shared.GOOGLE_OAUTH_CLIENT_ID}}"
RAILPACK_NO_SPA="1"
RAILPACK_NODE_VERSION="22.20"
```

## Backend Configuration

1. Click "Create" in the top right corner of the project dashboard
2. Select "GitHub Repo"
3. Select your repo
4. Click on the newly added service
5. Click on the "Settings" tab
6. Navigate to "Source" and select "Add route directory"
7. Enter "/backend"
8. Click on the "Variables" tab and enter the environment variables from the env file below.
9. Press "Deploy" at the top of the screen

```env
DATABASE_URI="${{Postgres.DATABASE_URL}}"
ENABLE_EMAIL="${{shared.ENABLE_EMAIL}}"
GOOGLE_OAUTH_CLIENT_ID="${{shared.GOOGLE_OAUTH_CLIENT_ID}}"
GOOGLE_OAUTH_CLIENT_SECRET="${{shared.GOOGLE_OAUTH_CLIENT_SECRET}}"
JWT_AUDIENCE="${{shared.JWT_AUDIENCE}}"
RESEND_API_KEY="${{shared.RESEND_API_KEY}}"
STORAGE_ACCESS_KEY_ID="${{Bucket.MINIO_ROOT_USER}}"
STORAGE_ENDPOINT="${{Bucket.MINIO_PRIVATE_HOST}}:${{Bucket.MINIO_PRIVATE_PORT}}"
STORAGE_SECRET_ACCESS_KEY="${{Bucket.MINIO_ROOT_PASSWORD}}"
STRIPE_BILLING_PORTAL_CONFIGURATION_ID="${{shared.STRIPE_BILLING_PORTAL_CONFIGURATION_ID}}"
STRIPE_PRO_PLAN_PRICE_ID="${{shared.STRIPE_PRO_PLAN_PRICE_ID}}"
STRIPE_PRO_PLAN_PRODUCT_ID="${{shared.STRIPE_PRO_PLAN_PRODUCT_ID}}"
STRIPE_SECRET_KEY="${{shared.STRIPE_SECRET_KEY}}"
STRIPE_WEBHOOK_SECRET="${{shared.STRIPE_WEBHOOK_SECRET}}"
```

## Networking Configuration

1. Click on the caddy service from the "Architecture" view
2. Navigate to the "Settings" tab
3. Navigate to "Networking" and click "Custom Domain"
4. Enter your domain and follow the instructions to update your DNS configuration
