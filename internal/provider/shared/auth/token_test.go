package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTokenFetch is a simple test to verify OAuth2 token acquisition works
func TestTokenFetch(t *testing.T) {
	// Test credentials
	testTenantID := "test-tenant-123"
	testClientID := "test-client-456"
	testClientSecret := "test-secret-789-very-secure"

	// Setup mock OAuth2 server
	var mockServer *httptest.Server
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			// Mock discovery response
			discoveryResponse := map[string]interface{}{
				"issuer":                                mockServer.URL,
				"token_endpoint":                        mockServer.URL + "/oauth2/token",
				"authorization_endpoint":                mockServer.URL + "/oauth2/authorize",
				"jwks_uri":                              mockServer.URL + "/.well-known/jwks.json",
				"grant_types_supported":                 []string{"client_credentials", "authorization_code"},
				"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
				"response_types_supported":              []string{"code", "token"},
				"scopes_supported":                      []string{"iam:read", "iam:write", "iam:admin"},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(discoveryResponse)

		case "/oauth2/token":
			// Validate request
			err := r.ParseForm()
			require.NoError(t, err, "Should parse form data")

			grantType := r.Form.Get("grant_type")
			assert.Equal(t, "client_credentials", grantType, "Should use client_credentials grant type")

			clientID := r.Form.Get("client_id")
			clientSecret := r.Form.Get("client_secret")

			// Also check Basic auth if credentials not in form
			if clientID == "" || clientSecret == "" {
				username, password, ok := r.BasicAuth()
				if ok {
					clientID = username
					clientSecret = password
				}
			}

			if clientID != testClientID || clientSecret != testClientSecret {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":             "invalid_client",
					"error_description": "Client authentication failed",
				})
				return
			}

			// Generate successful token response
			tokenResponse := map[string]interface{}{
				"access_token": fmt.Sprintf("test-access-token-%d", time.Now().Unix()),
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        r.Form.Get("scope"),
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tokenResponse)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Create auth client configuration
	config := &AuthClientConfig{
		TenantID:     testTenantID,
		ClientID:     testClientID,
		ClientSecret: testClientSecret,
		BaseURL:      mockServer.URL,
		Scopes:       []string{"iam:read", "iam:write"},
		Timeout:      30 * time.Second,
		MaxRetries:   3,
	}

	// Create auth client
	authClient, err := NewAuthClient(config)
	require.NoError(t, err, "Should create auth client successfully")
	defer authClient.Close()

	// Test token acquisition
	ctx := context.Background()
	token, err := authClient.GetToken(ctx)
	require.NoError(t, err, "Should acquire token successfully")
	require.NotNil(t, token, "Token should not be nil")

	// Verify token properties
	assert.True(t, token.Valid(), "Token should be valid")
	assert.Equal(t, "Bearer", token.TokenType, "Token type should be Bearer")
	assert.NotEmpty(t, token.AccessToken, "Access token should not be empty")
	assert.True(t, len(token.AccessToken) > 10, "Access token should be reasonable length")

	fmt.Printf("✅ Token fetch test successful!\n")
	fmt.Printf("   Token Type: %s\n", token.TokenType)
	fmt.Printf("   Access Token: %s...\n", token.AccessToken[:20])
	fmt.Printf("   Expires: %v\n", token.Expiry)
	fmt.Printf("   Valid: %v\n", token.Valid())

	// Test authenticated HTTP client
	httpClient, err := authClient.HTTPClient(ctx)
	require.NoError(t, err, "Should create authenticated HTTP client")
	require.NotNil(t, httpClient, "HTTP client should not be nil")

	fmt.Printf("✅ Authenticated HTTP client created successfully!\n")
}

// TestTokenFetchError tests error scenarios
func TestTokenFetchError(t *testing.T) {
	// Setup mock server that returns errors
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":             "invalid_client",
				"error_description": "Client authentication failed",
			})
		}
	}))
	defer mockServer.Close()

	// Create auth client with wrong credentials
	config := &AuthClientConfig{
		TenantID:     "test-tenant",
		ClientID:     "wrong-client-id",
		ClientSecret: "wrong-secret",
		TokenURL:     mockServer.URL + "/oauth2/token",
		Scopes:       []string{"iam:read"},
		Timeout:      30 * time.Second,
		MaxRetries:   1, // Single attempt
	}

	authClient, err := NewAuthClient(config)
	require.NoError(t, err, "Should create auth client successfully")
	defer authClient.Close()

	// Test token acquisition failure
	ctx := context.Background()
	_, err = authClient.GetToken(ctx)
	require.Error(t, err, "Should fail to acquire token with wrong credentials")

	fmt.Printf("✅ Token fetch error test successful - correctly failed with: %v\n", err)
}
