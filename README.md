# reece-start

Starter kit for building B2B applications using Stripe connect. It is based on SvelteKit + Go, and adopts a (kindof) cloud-native approach for deployments.

## Features

- Organization-Member-User authentication
  - Organization invitation links & emails
  - RBAC using scopes
  - Sudo users
  - Impersonation
  - Google OAuth
  - User & organization logos
  - Settings pages
- Landing, pricing, and faq page templates
- Email notification API
- Organization onboarding
- Billing with free and paid plans
- Stripe connect support
- i18n support
- Basic CI checks
  - Tests
    - Backend
    - Frontend
      - Unit
      - Component
  - Typing
  - Buildability
  - Linting
  - Formatting
- Railway deployment template (TODO)
- Structured logging (json) (TODO)
- Posthog integration (TODO)
- Sentry integration (TODO)

## Getting started

```
# terminal 1
make infra-start
make backend-dev

# terminal 2
make frontend-dev
```

## Commands

Run `make` or `make help` to see all available commands.

### Infrastructure

- `make infra-start` - Start the docker containers required for deployment
- `make infra-stop` - Stop the docker containers

### Frontend

- `make frontend-dev` - Start the dev frontend server w/ auto-reload
- `make frontend-build` - Build the frontend for production
- `make frontend-lint` - Lint frontend code
- `make frontend-lint-fix` - Lint and fix frontend code
- `make frontend-format` - Format frontend code
- `make frontend-typecheck` - Run TypeScript type checking
- `make frontend-test` - Run frontend tests
- `make frontend-test-watch` - Run frontend tests in watch mode

### Backend

- `make backend-dev` - Start the dev backend server w/ auto-reload
- `make backend-build` - Build the backend for production
- `make backend-format` - Format backend code
- `make backend-test` - Run backend tests (requires Docker for testcontainers)
- `make backend-test-verbose` - Run backend tests with verbose output
- `make backend-test-coverage` - Run backend tests with coverage report

### Stripe

- `make stripe-listen` - Start the Stripe webhook event listener

## Project structure

### `backend/`

An echo project which uses Gorm, River, Postgres, and Minio. All endpoints follow the JSON API spec and use the go validator library to validate incoming requests. The backend is designed under the assumption of a Organization-Member-User authentication model, whereby there is a many-many relationship between Organizations and Users. You can remove the `Organization*` models to get a simpler authentication model with only users and no organizations.

### `frontend/`

A SvelteKit project which uses DaisyUI. Load functions and form actions are used to take advantage of sveltes SSR + CSR features. Mutations are made directly against the backend when it makes sense. Otherwise, form actions are preferred as this makes the site more accessible.

### `caddy/`

A reverse proxy used in local and cloud deployments. Used to ensure that the all API requests from the browser are made to the same origin so that cookies are preserved and no CORS preflight requests are needed. Caddy is used as the reverse proxy and is built as a standalone docker container.

## Deployment

The project is intended to be deployed to railway, however you can easily adapt the docker compose file to deploy the app using docker compose.

## Relevant docs

- https://svelte.dev/
- https://daisyui.com/
- https://lucide.dev/
- https://gorm.io/
- https://riverqueue.com/
- https://echo.labstack.com/
- https://resend.com/
- https://posthog.com/
- https://stripe.com/
