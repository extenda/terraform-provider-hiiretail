package acceptance_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/testutils"
	"golang.org/x/oauth2/clientcredentials"
)

func TestAPIResponse_CreateVsRead(t *testing.T) {
	// Set up test environment with mock server
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)
	env.ValidateMockServerReady(t)

	// Set up OAuth2 client
	config := &clientcredentials.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TokenURL:     env.BaseURL + "/oauth2/token",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	httpClient := config.Client(ctx)

	// Test CREATE operation
	createReq := map[string]interface{}{
		"id":   "test-role-001",
		"name": "Test Custom Role",
		"permissions": []map[string]interface{}{
			{
				"id": "pos.payment.create",
			},
		},
	}

	jsonData, _ := json.Marshal(createReq)

	// POST to create role with OAuth2
	resp, err := httpClient.Post(env.BaseURL+"/iam/v1/custom-roles", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		t.Fatalf("CREATE request failed: %v", err)
	}
	defer resp.Body.Close()

	createBody, _ := io.ReadAll(resp.Body)
	t.Logf("CREATE response status: %d, body: %s", resp.StatusCode, string(createBody))

	// Test READ operation with OAuth2
	readResp, err := httpClient.Get(env.BaseURL + "/iam/v1/custom-roles/test-role-001")
	if err != nil {
		t.Fatalf("READ request failed: %v", err)
	}
	defer readResp.Body.Close()

	readBody, _ := io.ReadAll(readResp.Body)
	t.Logf("READ response status: %d, body: %s", readResp.StatusCode, string(readBody))
}
