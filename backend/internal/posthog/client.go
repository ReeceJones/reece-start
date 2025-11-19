package posthog

import (
	"log/slog"

	"github.com/posthog/posthog-go"
	"reece.start/internal/configuration"
)

// Client wraps the PostHog client with no-op support when not configured
type Client struct {
	client  posthog.Client
	enabled bool
}

// NewClient creates a new PostHog client. Returns a no-op client if PostHog is not configured.
func NewClient(config *configuration.Config) *Client {
	if config.PostHogApiKey == "" {
		slog.Info("PostHog not configured, using no-op client")
		return &Client{
			client:  nil,
			enabled: false,
		}
	}

	client, err := posthog.NewWithConfig(config.PostHogApiKey, posthog.Config{
		Endpoint: config.PostHogHost,
	})
	if err != nil {
		slog.Warn("Failed to initialize PostHog client, using no-op client", "error", err)
		return &Client{
			client:  nil,
			enabled: false,
		}
	}

	slog.Info("PostHog client initialized", "host", config.PostHogHost)
	return &Client{
		client:  client,
		enabled: true,
	}
}

// Capture logs an event to PostHog. No-op if PostHog is not configured.
func (c *Client) Capture(distinctID string, event string, properties map[string]any) {
	if !c.enabled || c.client == nil {
		return
	}

	err := c.client.Enqueue(posthog.Capture{
		DistinctId: distinctID,
		Event:      event,
		Properties: properties,
	})
	if err != nil {
		slog.Warn("Failed to capture PostHog event", "event", event, "error", err)
	}
}

// Close closes the PostHog client and flushes any pending events.
func (c *Client) Close() error {
	if !c.enabled || c.client == nil {
		return nil
	}
	return c.client.Close()
}
