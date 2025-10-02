package main
package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/testutils"
)

func TestOAuth2Timeout(t *testing.T) {
	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready
	env.ValidateMockServerReady(t)

	// Configure OAuth2 exactly like the provider does
	config := &clientcredentials.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     fmt.Sprintf("%s/oauth2/token", env.BaseURL),
	}

	// Create HTTP client with OAuth2 configuration and timeout like the provider
	baseHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create context with timeout for OAuth2 operations
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	httpClient := config.Client(context.WithValue(ctx, oauth2.HTTPClient, baseHTTPClient))

	// Try to make a request to test OAuth2 flow
	testURL := fmt.Sprintf("%s/iam/v1/custom-roles", env.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", testURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("X-Tenant-ID", "test-tenant-123")

	t.Logf("Making OAuth2 request to: %s", testURL)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("OAuth2 request failed: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("OAuth2 request succeeded with status: %d", resp.StatusCode)
}