package mocks

import (
	"reece.start/internal/configuration"
	"reece.start/internal/posthog"
)

func NewMockPosthogClient() *posthog.Client {
	return posthog.NewClient(&configuration.Config{})
}
