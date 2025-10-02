package acceptance_tests

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/testutils"
	"golang.org/x/oauth2/clientcredentials"
)

func TestOAuth2Flow_Isolated(t *testing.T) {
	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Validate mock server is ready
	env.ValidateMockServerReady(t)

	// Configure OAuth2 like the provider does
	config := &clientcredentials.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     env.BaseURL + "/oauth2/token",
	}

	// Create context with a short timeout to see if it hangs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Log("Creating OAuth2 client...")
	httpClient := config.Client(ctx)

	// Make a simple test request to trigger token acquisition
	t.Log("Making test request to trigger OAuth2 token acquisition...")
	resp, err := httpClient.Get(env.BaseURL + "/iam/v1/custom-roles")
	if err != nil {
		t.Fatalf("OAuth2 flow failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body to debug what the API returns
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	t.Logf("OAuth2 flow successful, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
}
