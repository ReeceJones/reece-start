# Caddy Reverse Proxy Configuration

Native Caddy configuration that routes:

- Frontend requests directly to FRONTEND_URL
- API requests to API_URL with path prefix removal (/api/_ → /_)

Uses Caddy's native reverse proxy capabilities with automatic HTTP/2 and header management.
