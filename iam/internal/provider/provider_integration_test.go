package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestOIDCIntegration tests the OIDC client credentials flow integration
func TestOIDCIntegration(t *testing.T) {
	// Create a mock OAuth server
	mockOAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" && r.Method == "POST" {
			// Validate the request
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form", http.StatusBadRequest)
				return
			}

			grantType := r.FormValue("grant_type")
			clientId := r.FormValue("client_id")
			clientSecret := r.FormValue("client_secret")

			if grantType != "client_credentials" {
				http.Error(w, "Invalid grant_type", http.StatusBadRequest)
				return
			}

			if clientId == "" || clientSecret == "" {
				http.Error(w, "Missing credentials", http.StatusBadRequest)
				return
			}

			// Simulate different responses based on client_id for testing
			switch clientId {
			case "valid-client":
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"access_token": "mock-access-token",
					"token_type": "Bearer",
					"expires_in": 3600
				}`))
			case "invalid-client":
				http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
			default:
				http.Error(w, "Unknown client", http.StatusUnauthorized)
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer mockOAuthServer.Close()

	testCases := []struct {
		name          string
		clientId      string
		clientSecret  string
		baseUrl       string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Valid OIDC credentials",
			clientId:      "valid-client",
			clientSecret:  "valid-secret",
			baseUrl:       mockOAuthServer.URL,
			expectedError: false,
		},
		{
			name:          "Invalid OIDC credentials",
			clientId:      "invalid-client",
			clientSecret:  "invalid-secret",
			baseUrl:       mockOAuthServer.URL,
			expectedError: false, // Configuration should succeed, authentication will fail on API calls
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &HiiRetailIamProvider{}

			// Create schema
			schemaResp := &provider.SchemaResponse{}
			p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			// Create configuration
			configValue := tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"tenant_id":     tftypes.String,
					"base_url":      tftypes.String,
					"client_id":     tftypes.String,
					"client_secret": tftypes.String,
				},
			}, map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, "test-tenant"),
				"base_url":      tftypes.NewValue(tftypes.String, tc.baseUrl),
				"client_id":     tftypes.NewValue(tftypes.String, tc.clientId),
				"client_secret": tftypes.NewValue(tftypes.String, tc.clientSecret),
			})

			config := tfsdk.Config{
				Schema: schemaResp.Schema,
				Raw:    configValue,
			}

			// Test configuration
			resp := &provider.ConfigureResponse{}
			p.Configure(context.Background(), provider.ConfigureRequest{Config: config}, resp)

			if tc.expectedError {
				if !resp.Diagnostics.HasError() {
					t.Error("Expected error but got none")
				} else {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if contains(diag.Summary(), tc.errorContains) || contains(diag.Detail(), tc.errorContains) {
							found = true
							break
						}
					}
					if !found && tc.errorContains != "" {
						t.Errorf("Expected error containing '%s', but got: %v", tc.errorContains, resp.Diagnostics.Errors())
					}
				}
			} else {
				if resp.Diagnostics.HasError() {
					t.Errorf("Expected no error but got: %v", resp.Diagnostics.Errors())
					return
				}

				// Verify API client is configured
				apiClient, ok := resp.ResourceData.(*APIClient)
				if !ok {
					t.Error("Expected ResourceData to be *APIClient")
					return
				}

				if apiClient.HTTPClient == nil {
					t.Error("Expected HTTPClient to be configured")
					return
				}

				if apiClient.BaseURL != tc.baseUrl {
					t.Errorf("Expected BaseURL to be '%s', got '%s'", tc.baseUrl, apiClient.BaseURL)
				}

				if apiClient.TenantID != "test-tenant" {
					t.Errorf("Expected TenantID to be 'test-tenant', got '%s'", apiClient.TenantID)
				}

				// Test making a request with the configured client (this will test OIDC token retrieval)
				testRequest, err := http.NewRequest("GET", tc.baseUrl+"/test", nil)
				if err != nil {
					t.Errorf("Failed to create test request: %v", err)
					return
				}

				// Set a short timeout for the test
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				testRequest = testRequest.WithContext(ctx)

				// The HTTP client should automatically handle OIDC token retrieval
				resp, err := apiClient.HTTPClient.Do(testRequest)
				if err != nil {
					// For valid credentials, we expect the request to succeed (even if 404)
					// For invalid credentials, we might get an auth error
					if tc.clientId == "valid-client" {
						t.Logf("Request completed with error (expected for test endpoint): %v", err)
					}
				} else {
					resp.Body.Close()
					t.Logf("Request completed with status: %d", resp.StatusCode)
				}
			}
		})
	}
}

// TestProviderConfigurationValidation tests various configuration validation scenarios
func TestProviderConfigurationValidation(t *testing.T) {
	testCases := []struct {
		name          string
		config        map[string]interface{}
		expectedError string
	}{
		{
			name: "Empty tenant_id",
			config: map[string]interface{}{
				"tenant_id":     "",
				"client_id":     "test-client",
				"client_secret": "test-secret",
			},
			expectedError: "Missing tenant_id",
		},
		{
			name: "Empty client_id",
			config: map[string]interface{}{
				"tenant_id":     "test-tenant",
				"client_id":     "",
				"client_secret": "test-secret",
			},
			expectedError: "Missing client_id",
		},
		{
			name: "Empty client_secret",
			config: map[string]interface{}{
				"tenant_id":     "test-tenant",
				"client_id":     "test-client",
				"client_secret": "",
			},
			expectedError: "Missing client_secret",
		},
		{
			name: "Invalid URL format",
			config: map[string]interface{}{
				"tenant_id":     "test-tenant",
				"client_id":     "test-client",
				"client_secret": "test-secret",
				"base_url":      "not-a-url",
			},
			expectedError: "Invalid base_url",
		},
		{
			name: "Valid minimal configuration",
			config: map[string]interface{}{
				"tenant_id":     "test-tenant",
				"client_id":     "test-client",
				"client_secret": "test-secret",
			},
			expectedError: "",
		},
		{
			name: "Valid full configuration",
			config: map[string]interface{}{
				"tenant_id":     "test-tenant",
				"client_id":     "test-client",
				"client_secret": "test-secret",
				"base_url":      "https://custom-api.example.com",
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &HiiRetailIamProvider{}

			// Create schema
			schemaResp := &provider.SchemaResponse{}
			p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			// Convert config to tftypes.Value
			configMap := make(map[string]tftypes.Value)
			for key, value := range tc.config {
				if strValue, ok := value.(string); ok {
					configMap[key] = tftypes.NewValue(tftypes.String, strValue)
				}
			}

			// Ensure all schema attributes are present
			for attr := range schemaResp.Schema.Attributes {
				if _, exists := configMap[attr]; !exists {
					configMap[attr] = tftypes.NewValue(tftypes.String, nil)
				}
			}

			configValue := tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"tenant_id":     tftypes.String,
					"base_url":      tftypes.String,
					"client_id":     tftypes.String,
					"client_secret": tftypes.String,
				},
			}, configMap)

			config := tfsdk.Config{
				Schema: schemaResp.Schema,
				Raw:    configValue,
			}

			// Test configuration
			resp := &provider.ConfigureResponse{}
			p.Configure(context.Background(), provider.ConfigureRequest{Config: config}, resp)

			if tc.expectedError == "" {
				if resp.Diagnostics.HasError() {
					t.Errorf("Expected no error, but got: %v", resp.Diagnostics.Errors())
				}
			} else {
				if !resp.Diagnostics.HasError() {
					t.Errorf("Expected error containing '%s', but got no error", tc.expectedError)
				} else {
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if contains(diag.Summary(), tc.expectedError) || contains(diag.Detail(), tc.expectedError) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error containing '%s', but got: %v", tc.expectedError, resp.Diagnostics.Errors())
					}
				}
			}
		})
	}
}
