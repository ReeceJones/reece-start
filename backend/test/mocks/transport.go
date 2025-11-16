package mocks

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Mock clients for external services
//
// These mocks prevent actual API calls to external services during tests:
// - Stripe: Intercepts HTTP calls at the transport level to prevent network requests
// - Resend: Intercepts HTTP calls at the transport level to prevent network requests
//   (Email sending is also disabled in tests via EnableEmail: false as a double safeguard)

var (
	// originalTransport stores the original http.DefaultTransport so we can restore it
	originalTransport http.RoundTripper
	// internalTransport is used for localhost/internal requests
	internalTransport http.RoundTripper
	transportMutex    sync.Mutex
	transportReplaced bool
)

// MockHTTPTransport intercepts HTTP requests and returns mock responses
// This prevents actual API calls to external services during tests
type MockHTTPTransport struct{}

// RoundTrip implements http.RoundTripper and returns mock responses
func (m *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	host := req.URL.Host

	// Allow localhost/internal requests to pass through (for test server, database connections, etc.)
	if host == "" || host == "localhost" || strings.HasPrefix(host, "127.0.0.1") ||
		strings.HasPrefix(host, "::1") || strings.Contains(host, "localhost") {
		// Use a separate transport for internal requests to avoid circular reference
		if internalTransport == nil {
			internalTransport = &http.Transport{}
		}
		return internalTransport.RoundTrip(req)
	}

	// Handle Stripe API calls - return mock responses
	if strings.Contains(url, "api.stripe.com") {
		return m.handleStripeRequest(req)
	}

	// Handle Resend API calls - return mock responses
	if strings.Contains(url, "api.resend.com") {
		return m.handleResendRequest(req)
	}

	// Block all other external API calls
	return nil, errors.New("mock HTTP transport: actual API calls are blocked in tests - attempted call to " + url)
}

// handleStripeRequest returns mock responses for Stripe API calls
func (m *MockHTTPTransport) handleStripeRequest(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	method := req.Method

	// Handle POST /v2/core/accounts (create Stripe Connect account)
	if method == "POST" && strings.Contains(url, "/v2/core/accounts") {
		accountID := "acct_" + uuid.New().String()[:24]

		// Read request body to extract display name if needed
		var displayName string
		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var requestData map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &requestData); err == nil {
				if dn, ok := requestData["display_name"].(string); ok {
					displayName = dn
				}
			}
		}

		if displayName == "" {
			displayName = "Test Account"
		}

		// Create mock account response
		account := map[string]interface{}{
			"id":           accountID,
			"display_name": displayName,
			"object":       "account",
			"identity": map[string]interface{}{
				"entity_type": "individual",
				"country":     "US",
			},
			"defaults": map[string]interface{}{
				"currency": "usd",
			},
			"configuration": map[string]interface{}{
				"customer": map[string]interface{}{
					"capabilities": map[string]interface{}{
						"automatic_indirect_tax": map[string]interface{}{
							"status": "active",
						},
					},
				},
				"merchant": map[string]interface{}{
					"capabilities": map[string]interface{}{
						"card_payments": map[string]interface{}{
							"status": "active",
						},
					},
				},
				"recipient": map[string]interface{}{
					"capabilities": map[string]interface{}{
						"stripe_balance": map[string]interface{}{
							"stripe_transfers": map[string]interface{}{
								"status": "active",
							},
							"payouts": map[string]interface{}{
								"status": "active",
							},
						},
					},
				},
			},
			"requirements": map[string]interface{}{
				"currently_due": []interface{}{},
			},
		}

		responseBody, _ := json.Marshal(account)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			Request:    req,
		}, nil
	}

	// Handle GET /v2/core/accounts/{id} (retrieve Stripe Connect account)
	if method == "GET" && strings.Contains(url, "/v2/core/accounts/") {
		// Extract account ID from URL
		parts := strings.Split(url, "/v2/core/accounts/")
		accountID := parts[len(parts)-1]
		if strings.Contains(accountID, "?") {
			accountID = strings.Split(accountID, "?")[0]
		}

		account := map[string]interface{}{
			"id":           accountID,
			"display_name": "Test Account",
			"object":       "account",
			"identity": map[string]interface{}{
				"entity_type": "individual",
				"country":     "US",
			},
			"defaults": map[string]interface{}{
				"currency": "usd",
			},
			"configuration": map[string]interface{}{
				"customer": map[string]interface{}{
					"capabilities": map[string]interface{}{
						"automatic_indirect_tax": map[string]interface{}{
							"status": "active",
						},
					},
				},
				"merchant": map[string]interface{}{
					"capabilities": map[string]interface{}{
						"card_payments": map[string]interface{}{
							"status": "active",
						},
					},
				},
				"recipient": map[string]interface{}{
					"capabilities": map[string]interface{}{
						"stripe_balance": map[string]interface{}{
							"stripe_transfers": map[string]interface{}{
								"status": "active",
							},
							"payouts": map[string]interface{}{
								"status": "active",
							},
						},
					},
				},
			},
			"requirements": map[string]interface{}{
				"currently_due": []interface{}{},
			},
		}

		responseBody, _ := json.Marshal(account)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			Request:    req,
		}, nil
	}

	// Handle POST /v2/core/account_links (create onboarding link)
	if method == "POST" && strings.Contains(url, "/v2/core/account_links") {
		linkID := "link_" + uuid.New().String()[:24]
		link := map[string]interface{}{
			"id":  linkID,
			"url": "https://connect.stripe.com/setup/test/" + linkID,
		}

		responseBody, _ := json.Marshal(link)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			Request:    req,
		}, nil
	}

	// Handle POST /v1/checkout/sessions (create checkout session)
	if method == "POST" && strings.Contains(url, "/v1/checkout/sessions") {
		sessionID := "cs_test_" + uuid.New().String()[:24]
		session := map[string]interface{}{
			"id":  sessionID,
			"url": "https://checkout.stripe.com/test/" + sessionID,
		}

		responseBody, _ := json.Marshal(session)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			Request:    req,
		}, nil
	}

	// Handle POST /v1/billing_portal/sessions (create billing portal session)
	if method == "POST" && strings.Contains(url, "/v1/billing_portal/sessions") {
		sessionID := "bps_test_" + uuid.New().String()[:24]
		session := map[string]interface{}{
			"id":  sessionID,
			"url": "https://billing.stripe.com/test/" + sessionID,
		}

		responseBody, _ := json.Marshal(session)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			Request:    req,
		}, nil
	}

	// Handle GET /v1/subscriptions/{id} (retrieve subscription)
	if method == "GET" && strings.Contains(url, "/v1/subscriptions/") {
		// Extract subscription ID from URL
		parts := strings.Split(url, "/v1/subscriptions/")
		subscriptionID := parts[len(parts)-1]
		if strings.Contains(subscriptionID, "?") {
			subscriptionID = strings.Split(subscriptionID, "?")[0]
		}

		subscription := map[string]interface{}{
			"id":     subscriptionID,
			"status": "active",
			"metadata": map[string]interface{}{
				"organization_id": "1",
			},
			"billing_cycle_anchor": time.Now().Unix(),
			"items": map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"price": map[string]interface{}{
							"product": map[string]interface{}{
								"id": "prod_test_123",
							},
							"unit_amount": 1000,
						},
					},
				},
			},
		}

		responseBody, _ := json.Marshal(subscription)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			Request:    req,
		}, nil
	}

	// For other Stripe endpoints, return a generic error
	return nil, errors.New("mock HTTP transport: unhandled Stripe endpoint - " + url)
}

// handleResendRequest returns mock responses for Resend API calls
func (m *MockHTTPTransport) handleResendRequest(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	method := req.Method

	// Handle POST /emails (send email)
	if method == "POST" && strings.Contains(url, "/emails") {
		emailID := uuid.New().String()
		response := map[string]interface{}{
			"id": emailID,
		}

		responseBody, _ := json.Marshal(response)
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			Request:    req,
		}, nil
	}

	// For other Resend endpoints, return a generic error
	return nil, errors.New("mock HTTP transport: unhandled Resend endpoint - " + url)
}

// ReplaceDefaultTransport replaces http.DefaultTransport with a mock transport
// This intercepts all HTTP calls including those from Stripe SDK and Resend SDK
func ReplaceDefaultTransport() {
	transportMutex.Lock()
	defer transportMutex.Unlock()

	if !transportReplaced {
		originalTransport = http.DefaultTransport
		http.DefaultTransport = &MockHTTPTransport{}
		transportReplaced = true
	}
}

// ReplaceDefaultTransportWithCleanup replaces http.DefaultTransport with a mock transport
// and registers a cleanup function with the test to automatically restore it when the test finishes.
// This is the preferred way to use the mock transport in tests as it ensures cleanup even if the test panics.
func ReplaceDefaultTransportWithCleanup(t *testing.T) {
	ReplaceDefaultTransport()
	t.Cleanup(func() {
		RestoreDefaultTransport()
	})
}

// RestoreDefaultTransport restores the original http.DefaultTransport
func RestoreDefaultTransport() {
	transportMutex.Lock()
	defer transportMutex.Unlock()

	if transportReplaced {
		http.DefaultTransport = originalTransport
		transportReplaced = false
	}
}
