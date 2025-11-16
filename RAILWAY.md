# Railway Deployment Setup

This document describes how to set up Railway deployment for the application.

---

## Overview

![Railway project](/assets/railway-overview.png)

The Railway deployment exposes the backend and frontend through a Caddy
reverse proxy, using Railway templates for infrastructure services.

### Favicon links

- **Frontend**: <https://svelte.dev/favicon.png>
- **Backend**: <https://echo.labstack.com/img/favicon.ico>
- **Caddy**: <https://caddyserver.com/resources/images/favicon.png>

---

## Project configuration (shared variables)

1. Open your Railway project.
2. Click **Settings** in the top-right corner.
3. Navigate to **Shared Variables**.
4. Add the variables from the table below (these are shared across services):

| Variable Name                            | Value                     |
| ---------------------------------------- | ------------------------- |
| `ENABLE_EMAIL`                           | `true` or `false`         |
| `GOOGLE_OAUTH_CLIENT_ID`                 | `your_client_._id`        |
| `GOOGLE_OAUTH_CLIENT_SECRET`             | `your_client_secret`      |
| `JWT_AUDIENCE`                           | `https://your.domain.com` |
| `RESEND_API_KEY`                         | `your_api_key`            |
| `STRIPE_ACCOUNT_WEBHOOK_SECRET`          | `your_account_webhook_secret` |
| `STRIPE_BILLING_PORTAL_CONFIGURATION_ID` | `your_configuration_id`       |
| `STRIPE_CONNECT_WEBHOOK_SECRET`          | `your_connect_webhook_secret` |
| `STRIPE_PRO_PLAN_PRICE_ID`               | `your_price_id`               |
| `STRIPE_PRO_PLAN_PRODUCT_ID`             | `your_product_id`             |
| `STRIPE_SECRET_KEY`                      | `your_secret_key`             |

---

## Postgres configuration

1. In the Railway project dashboard, click **Create** (top-right).
2. Select **Database**.
3. Choose **Postgres** and provision the database.

---

## MinIO configuration

1. In the Railway project dashboard, click **Create**.
2. Select **Template**.
3. Search for **"MinIO"**.
4. Select **"MinIO" by "Brody's Projects"** and provision it.

---

## Caddy service configuration

This service acts as the reverse proxy in front of the backend and frontend.

1. Click **Create** and select **GitHub Repo**.
2. Select this repository.
3. Click the newly added **Caddy** service.
4. Go to the **Settings** tab.
5. Under **Source**, click **Add route directory** and enter `/caddy`.
6. Go to the **Variables** tab and add the variables from the env block below.
7. Click **Deploy**.

```env
API_HOST="your.domain.com"
API_URL="${{backend.RAILWAY_PRIVATE_DOMAIN}}:8080"
FRONTEND_HOST="your.domain.com"
FRONTEND_URL="${{frontend.RAILWAY_PRIVATE_DOMAIN}}:8080"
```

---

## Frontend service configuration

1. Click **Create** and select **GitHub Repo**.
2. Select this repository.
3. Click the newly added **Frontend** service.
4. Go to the **Settings** tab.
5. Under **Source**, click **Add route directory** and enter `/frontend`.
6. Go to the **Variables** tab and add the variables from the env block below.
7. Click **Deploy**.

```env
BODY_SIZE_LIMIT="64M"
PUBLIC_GOOGLE_OAUTH_CLIENT_ID="${{shared.GOOGLE_OAUTH_CLIENT_ID}}"
RAILPACK_NO_SPA="1"
RAILPACK_NODE_VERSION="22.20"
```

---

## Backend service configuration

1. Click **Create** and select **GitHub Repo**.
2. Select this repository.
3. Click the newly added **Backend** service.
4. Go to the **Settings** tab.
5. Under **Source**, click **Add route directory** and enter `/backend`.
6. Go to the **Variables** tab and add the variables from the env block below.
7. Click **Deploy**.

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
STRIPE_ACCOUNT_WEBHOOK_SECRET="${{shared.STRIPE_ACCOUNT_WEBHOOK_SECRET}}"
STRIPE_BILLING_PORTAL_CONFIGURATION_ID="${{shared.STRIPE_BILLING_PORTAL_CONFIGURATION_ID}}"
STRIPE_CONNECT_WEBHOOK_SECRET="${{shared.STRIPE_CONNECT_WEBHOOK_SECRET}}"
STRIPE_PRO_PLAN_PRICE_ID="${{shared.STRIPE_PRO_PLAN_PRICE_ID}}"
STRIPE_PRO_PLAN_PRODUCT_ID="${{shared.STRIPE_PRO_PLAN_PRODUCT_ID}}"
STRIPE_SECRET_KEY="${{shared.STRIPE_SECRET_KEY}}"
```

---

## Networking configuration (custom domain)

1. Open the **Caddy** service from the **Architecture** view.
2. Go to the **Settings** tab.
3. Open the **Networking** section and click **Custom Domain**.
4. Enter your domain and follow Railway’s instructions to update your DNS.

---

## Stripe configuration

1. In the Stripe Dashboard, create or select the account you want to use for this project.
2. Go to **Developers → API keys** and create a restricted key (or use the secret key)
   with permissions for billing, subscriptions, and Connect as required.
3. Copy the secret key and set it as the `STRIPE_SECRET_KEY` **shared variable** in Railway
   (see the **Project configuration** table above).
4. Go to **Developers → Webhooks** and click **Add an endpoint**.
5. Set the webhook URL to your deployed backend URL for account/snapshot events, for example:
   `https://your.domain.com/api/webhooks/stripe/account/snapshot`.
6. Under **Select events**, subscribe to **only** the following events:
   - `invoice.paid`
   - `invoice.payment_failed`
   - `invoice.payment_action_required`
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
7. Save the webhook endpoint, then copy the **Signing secret** and set it as
   the `STRIPE_ACCOUNT_WEBHOOK_SECRET` **shared variable** in Railway.
8. Create a second webhook endpoint for **connect/thin** events with URL
   `https://your.domain.com/api/webhooks/stripe/connect/thin` and subscribe to **only**
   the following Connect Account v2 events:
   - `v2.core.account.created`
   - `v2.core.account.updated`
   - `v2.core.account.closed`
   - `v2.core.account_person.updated`
   - `v2.core.account[identity].updated`
   - `v2.core.account[configuration.customer].capability_status_updated`
   - `v2.core.account[configuration.merchant].capability_status_updated`
   - `v2.core.account[configuration.recipient].capability_status_updated`
   - `v2.core.account[requirements].updated`
   - `v2.core.account.updated`
9. Save the webhook endpoint, then copy the **Signing secret** and set it as
   the `STRIPE_CONNECT_WEBHOOK_SECRET` **shared variable** in Railway.
10. Go to **Billing → Prices / Products** and create a subscription product and price
    for your Pro plan.
11. Copy the Pro plan **Product ID** and **Price ID** and set them as
    `STRIPE_PRO_PLAN_PRODUCT_ID` and `STRIPE_PRO_PLAN_PRICE_ID` **shared variables**
    in Railway.
12. Go to **Billing → Customer portal** and configure the billing portal as desired.
    Copy the **Configuration ID** and set it as `STRIPE_BILLING_PORTAL_CONFIGURATION_ID`
    in Railway.
13. Verify that all Stripe-related shared variables from the **Project configuration**
    table are set, then redeploy the backend service so the new configuration is picked up.
