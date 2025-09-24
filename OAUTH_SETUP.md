# Google OAuth Setup

This document describes how to set up Google OAuth authentication for the Easy Apparel application.

## Backend Configuration

The following environment variables need to be set:

```bash
GOOGLE_OAUTH_CLIENT_ID=your_google_client_id_here
GOOGLE_OAUTH_CLIENT_SECRET=your_google_client_secret_here
```

These are already configured in `backend/internal/configuration/env.go`.

## Frontend Configuration

The frontend needs the Google OAuth Client ID as an environment variable:

```bash
PUBLIC_GOOGLE_OAUTH_CLIENT_ID=your_google_client_id_here
```

## Google Cloud Console Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google People API
4. Go to "Credentials" > "Create Credentials" > "OAuth 2.0 Client IDs"
5. Configure the OAuth consent screen
6. Set the authorized redirect URIs:
   - For development: `http://localhost:4040/oauth/google/callback`
   - For production: `https://yourdomain.com/oauth/google/callback`

## How It Works

### Flow Overview

1. User clicks "Sign in with Google" button on `/signin` or `/signup` pages
2. User is redirected to Google's OAuth consent page
3. After granting permission, Google redirects to `/oauth/google/callback`
4. The callback page submits a form action which validates the oauth state parameter.
5. The server-side form action exchanges the authorization code for user information
6. Backend creates or links the user account and issues a JWT token
7. Server sets an HTTP-only cookie and redirects to the dashboard

### Database Schema

The User model includes these additional fields for Google OAuth:

```go
type User struct {
    // ... other fields ...
    GoogleId           string // Google user ID for account linking
    GoogleProfileImage string // URL to Google profile image
}
```

You can follow this pattern when adding new OAuth providers.
