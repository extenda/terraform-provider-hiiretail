package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOIDCDiscoveryResponse_Validate tests the validation logic for OIDC discovery responses
func TestOIDCDiscoveryResponse_Validate(t *testing.T) {
	tests := []struct {
		name          string
		response      *OIDCDiscoveryResponse
		expectedError string
	}{
		{
			name: "valid_response",
			response: &OIDCDiscoveryResponse{
				Issuer:                   "https://auth.retailsvc.com",
				TokenEndpoint:            "https://auth.retailsvc.com/oauth2/token",
				AuthorizationEndpoint:    "https://auth.retailsvc.com/oauth2/authorize",
				JwksUri:                  "https://auth.retailsvc.com/.well-known/jwks.json",
				GrantTypesSupported:      []string{"client_credentials", "authorization_code"},
				TokenEndpointAuthMethods: []string{"client_secret_basic", "client_secret_post"},
				ResponseTypesSupported:   []string{"code", "token"},
				ScopesSupported:          []string{"iam:read", "iam:write"},
			},
			expectedError: "",
		},
		{
			name: "missing_issuer",
			response: &OIDCDiscoveryResponse{
				TokenEndpoint:            "https://auth.retailsvc.com/oauth2/token",
				GrantTypesSupported:      []string{"client_credentials"},
				TokenEndpointAuthMethods: []string{"client_secret_basic"},
			},
			expectedError: "missing required field: issuer",
		},
		{
			name: "missing_token_endpoint",
			response: &OIDCDiscoveryResponse{
				Issuer:              "https://auth.retailsvc.com",
				GrantTypesSupported: []string{"client_credentials"},
			},
			expectedError: "missing required field: token_endpoint",
		},
		{
			name: "unsupported_grant_types",
			response: &OIDCDiscoveryResponse{
				Issuer:                   "https://auth.retailsvc.com",
				TokenEndpoint:            "https://auth.retailsvc.com/oauth2/token",
				GrantTypesSupported:      []string{"authorization_code", "implicit"},
				TokenEndpointAuthMethods: []string{"client_secret_basic"},
			},
			expectedError: "client_credentials grant type not supported",
		},
		{
			name: "invalid_token_endpoint_url",
			response: &OIDCDiscoveryResponse{
				Issuer:              "https://auth.retailsvc.com",
				TokenEndpoint:       "not-a-valid-url",
				GrantTypesSupported: []string{"client_credentials"},
			},
			expectedError: "invalid token_endpoint URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.response.Validate()
			if tt.expectedError == "" {
				assert.NoError(t, err, "Expected no validation error")
			} else {
				assert.Error(t, err, "Expected validation error")
				assert.Contains(t, err.Error(), tt.expectedError, "Error message should contain expected text")
			}
		})
	}
}

// TestDiscoveryClient_FetchDiscovery tests OAuth2 discovery endpoint fetching
func TestDiscoveryClient_FetchDiscovery(t *testing.T) {
	// Mock discovery response
	validDiscoveryResponse := &OIDCDiscoveryResponse{
		Issuer:                   "https://auth.retailsvc-test.com",
		TokenEndpoint:            "https://auth.retailsvc-test.com/oauth2/token",
		AuthorizationEndpoint:    "https://auth.retailsvc-test.com/oauth2/authorize",
		JwksUri:                  "https://auth.retailsvc-test.com/.well-known/jwks.json",
		GrantTypesSupported:      []string{"client_credentials", "authorization_code"},
		TokenEndpointAuthMethods: []string{"client_secret_basic", "client_secret_post"},
		ResponseTypesSupported:   []string{"code", "token"},
		ScopesSupported:          []string{"iam:read", "iam:write", "iam:admin"},
	}

	t.Run("successful_discovery", func(t *testing.T) {
		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/.well-known/openid-configuration", r.URL.Path)
			assert.Equal(t, "application/json", r.Header.Get("Accept"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			responseBytes, _ := json.Marshal(validDiscoveryResponse)
			w.Write(responseBytes)
		}))
		defer server.Close()

		client := NewDiscoveryClient(server.URL, 30*time.Second)
		ctx := context.Background()

		response, err := client.FetchDiscovery(ctx)
		require.NoError(t, err, "FetchDiscovery should succeed")
		assert.Equal(t, validDiscoveryResponse.Issuer, response.Issuer)
		assert.Equal(t, validDiscoveryResponse.TokenEndpoint, response.TokenEndpoint)
		assert.Contains(t, response.GrantTypesSupported, "client_credentials")
	})

	t.Run("discovery_endpoint_not_found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}))
		defer server.Close()

		client := NewDiscoveryClient(server.URL, 30*time.Second)
		ctx := context.Background()

		_, err := client.FetchDiscovery(ctx)
		require.Error(t, err, "FetchDiscovery should fail with 404")
		assert.Contains(t, err.Error(), "discovery endpoint not available")
	})
}
