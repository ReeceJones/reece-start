# reece-start

This is a repo that I use as a starting point for new projects. It is based on SvelteKit + Go, and adopts a (mostly) cloud-native approach for deployments.

## Features

- Organization-Member-User authentication model
  - TODO: Refactor auth check functions, add scopes
- User settings page
- Organization settings page
- Organization-based billing with free and paid membership tiers (TODO)
- Stripe connect support (TODO)
- Landing, pricing, and feature page templates (TODO)
- Email and SMS notification APIs (WIP)
  - TODO: SMS
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

- `make infra-start` - Start the docker containers required for deployment
- `make infra-stop` - Stop the docker containers
- `make frontend-dev` - Start the dev frontend server w/ auto-reload
- `make backend-dev` - Start the dev backend server w/ auto-reload

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
