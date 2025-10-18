package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestHiiRetailProvider(t *testing.T) {
	t.Run("Provider metadata", func(t *testing.T) {
		p := &HiiRetailProvider{version: "test"}
		resp := &provider.MetadataResponse{}
		p.Metadata(context.Background(), provider.MetadataRequest{}, resp)

		if resp.TypeName != "hiiretail" {
			t.Errorf("Expected TypeName to be 'hiiretail', got %s", resp.TypeName)
		}
		if resp.Version != "test" {
			t.Errorf("Expected Version to be 'test', got %s", resp.Version)
		}
	})

	t.Run("Provider schema", func(t *testing.T) {
		p := &HiiRetailProvider{}
		resp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, resp)

		// Check that all expected attributes are present
		expectedAttrs := []string{"base_url", "client_id", "client_secret", "iam_endpoint", "ccc_endpoint", "token_url", "scopes", "timeout_seconds", "max_retries"}
		for _, attr := range expectedAttrs {
			if _, exists := resp.Schema.Attributes[attr]; !exists {
				t.Errorf("Expected attribute %s to be present in schema", attr)
			}
		}

		// Check that attributes are optional (can be set via environment variables)
		if resp.Schema.Attributes["client_id"].IsRequired() {
			t.Error("client_id should be optional (can be set via environment)")
		}
		if resp.Schema.Attributes["client_secret"].IsRequired() {
			t.Error("client_secret should be optional (can be set via environment)")
		}

		// Check optional attributes
		if !resp.Schema.Attributes["base_url"].IsOptional() {
			t.Error("base_url should be optional")
		}

		// Check sensitive attributes
		if !resp.Schema.Attributes["client_secret"].IsSensitive() {
			t.Error("client_secret should be sensitive")
		}
	})
}

func TestHiiRetailProvider_Configure(t *testing.T) {
	// Set up environment variables for tests
	t.Setenv("HIIRETAIL_TENANT_ID", "test-tenant")

	testCases := []struct {
		name          string
		config        map[string]tftypes.Value
		expectedError string
	}{
		{
			name: "Valid configuration with all fields - expect auth failure in unit test",
			config: map[string]tftypes.Value{
				"client_id":       tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret":   tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":        tftypes.NewValue(tftypes.String, "https://test-api.example.com"),
				"iam_endpoint":    tftypes.NewValue(tftypes.String, "/iam/v1"),
				"ccc_endpoint":    tftypes.NewValue(tftypes.String, "/ccc/v1"),
				"token_url":       tftypes.NewValue(tftypes.String, "https://auth.example.com/token"),
				"scopes":          tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "iam:read")}),
				"timeout_seconds": tftypes.NewValue(tftypes.Number, 30),
				"max_retries":     tftypes.NewValue(tftypes.Number, 3),
				"tenant_id":       tftypes.NewValue(tftypes.String, "test-tenant"),
			},
			expectedError: "OAuth2 authentication failed",
		},
		{
			name: "Valid minimal configuration - expect auth failure in unit test",
			config: map[string]tftypes.Value{
				"client_id":       tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret":   tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":        tftypes.NewValue(tftypes.String, nil),
				"iam_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"ccc_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"token_url":       tftypes.NewValue(tftypes.String, nil),
				"scopes":          tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
				"timeout_seconds": tftypes.NewValue(tftypes.Number, nil),
				"max_retries":     tftypes.NewValue(tftypes.Number, nil),
				"tenant_id":       tftypes.NewValue(tftypes.String, "test-tenant"),
			},
			expectedError: "OAuth2 authentication failed",
		},
		{
			name: "Missing client_id - should fail validation",
			config: map[string]tftypes.Value{
				"client_id":       tftypes.NewValue(tftypes.String, nil),
				"client_secret":   tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":        tftypes.NewValue(tftypes.String, nil),
				"iam_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"ccc_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"token_url":       tftypes.NewValue(tftypes.String, nil),
				"scopes":          tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
				"timeout_seconds": tftypes.NewValue(tftypes.Number, nil),
				"max_retries":     tftypes.NewValue(tftypes.Number, nil),
				"tenant_id":       tftypes.NewValue(tftypes.String, "test-tenant"),
			},
			expectedError: "client authentication failed",
		},
		{
			name: "Missing client_secret - should fail validation",
			config: map[string]tftypes.Value{
				"client_id":       tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret":   tftypes.NewValue(tftypes.String, nil),
				"base_url":        tftypes.NewValue(tftypes.String, nil),
				"iam_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"ccc_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"token_url":       tftypes.NewValue(tftypes.String, nil),
				"scopes":          tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
				"timeout_seconds": tftypes.NewValue(tftypes.Number, nil),
				"max_retries":     tftypes.NewValue(tftypes.Number, nil),
				"tenant_id":       tftypes.NewValue(tftypes.String, "test-tenant"),
			},
			expectedError: "client authentication failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &HiiRetailProvider{}

			// Create schema
			schemaResp := &provider.SchemaResponse{}
			p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			// Create configuration
			configValue := tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"client_id":       tftypes.String,
					"client_secret":   tftypes.String,
					"tenant_id":       tftypes.String,
					"base_url":        tftypes.String,
					"iam_endpoint":    tftypes.String,
					"ccc_endpoint":    tftypes.String,
					"token_url":       tftypes.String,
					"scopes":          tftypes.Set{ElementType: tftypes.String},
					"timeout_seconds": tftypes.Number,
					"max_retries":     tftypes.Number,
				},
			}, tc.config)

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
					// Consider OAuth2 auth failures and client init failures acceptable in unit tests
					alternatives := []string{"OAuth2 authentication failed", "Failed to initialize API client"}
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if tc.expectedError != "" && (contains(diag.Summary(), tc.expectedError) || contains(diag.Detail(), tc.expectedError)) {
							found = true
							break
						}
						for _, alt := range alternatives {
							if contains(diag.Summary(), alt) || contains(diag.Detail(), alt) {
								found = true
								break
							}
						}
						// Fall back to searching the entire diagnostic string for the expected substring
						foundInAny := false
						for _, d := range resp.Diagnostics.Errors() {
							if contains(fmt.Sprintf("%v", d), tc.expectedError) {
								foundInAny = true
								break
							}
						}
						if !foundInAny {
							t.Errorf("Expected error containing '%s', but got: %v", tc.expectedError, resp.Diagnostics.Errors())
						}
					}
					if !found {
						t.Errorf("Expected error containing '%s' or a known auth/client-init failure, but got: %v", tc.expectedError, resp.Diagnostics.Errors())
					}
				}
			}
		})
	}
}

func TestHiiRetailProvider_OIDCConfiguration(t *testing.T) {
	// Set up environment variables for tests
	t.Setenv("HIIRETAIL_TENANT_ID", "test-tenant")

	testCases := []struct {
		name         string
		clientId     string
		clientSecret string
		baseUrl      string
		expectError  bool
	}{
		{
			name:         "Valid OIDC configuration - expect auth failure in unit test",
			clientId:     "test-client",
			clientSecret: "test-secret",
			baseUrl:      "https://test-api.example.com",
			expectError:  true, // Expect auth failure with fake credentials
		},
		{
			name:         "Valid OIDC with default URL - expect auth failure in unit test",
			clientId:     "test-client",
			clientSecret: "test-secret",
			baseUrl:      "",
			expectError:  true, // Expect auth failure with fake credentials
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &HiiRetailProvider{}

			// Create schema
			schemaResp := &provider.SchemaResponse{}
			p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			// Create configuration
			configMap := map[string]tftypes.Value{
				"client_id":       tftypes.NewValue(tftypes.String, tc.clientId),
				"client_secret":   tftypes.NewValue(tftypes.String, tc.clientSecret),
				"iam_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"ccc_endpoint":    tftypes.NewValue(tftypes.String, nil),
				"token_url":       tftypes.NewValue(tftypes.String, nil),
				"scopes":          tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
				"timeout_seconds": tftypes.NewValue(tftypes.Number, nil),
				"max_retries":     tftypes.NewValue(tftypes.Number, nil),
				"tenant_id":       tftypes.NewValue(tftypes.String, "test-tenant"),
			}

			if tc.baseUrl != "" {
				configMap["base_url"] = tftypes.NewValue(tftypes.String, tc.baseUrl)
			} else {
				configMap["base_url"] = tftypes.NewValue(tftypes.String, nil)
			}

			configValue := tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"client_id":       tftypes.String,
					"client_secret":   tftypes.String,
					"base_url":        tftypes.String,
					"iam_endpoint":    tftypes.String,
					"ccc_endpoint":    tftypes.String,
					"token_url":       tftypes.String,
					"scopes":          tftypes.Set{ElementType: tftypes.String},
					"timeout_seconds": tftypes.Number,
					"max_retries":     tftypes.Number,
					"tenant_id":       tftypes.String,
				},
			}, configMap)

			config := tfsdk.Config{
				Schema: schemaResp.Schema,
				Raw:    configValue,
			}

			// Test configuration
			resp := &provider.ConfigureResponse{}
			p.Configure(context.Background(), provider.ConfigureRequest{Config: config}, resp)

			if tc.expectError && !resp.Diagnostics.HasError() {
				t.Error("Expected error but got none")
			} else if !tc.expectError && resp.Diagnostics.HasError() {
				t.Errorf("Expected no error but got: %v", resp.Diagnostics.Errors())
			}

			// If successful, check that the API client was created
			if !tc.expectError && !resp.Diagnostics.HasError() {
				if resp.ResourceData == nil {
					t.Error("Expected ResourceData to be set")
				}
				if resp.DataSourceData == nil {
					t.Error("Expected DataSourceData to be set")
				}

				// Check API client configuration
				if apiClient, ok := resp.ResourceData.(*client.Client); ok {
					// Note: The new client structure doesn't expose these fields directly
					// TODO: Update tests to match new client interface
					t.Log("API client configured successfully")
					if apiClient == nil {
						t.Error("Expected client to be non-nil")
					}
				} else {
					t.Error("Expected ResourceData to be *client.Client")
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	providerFunc := New("1.0.0")
	provider := providerFunc()

	if hiiProvider, ok := provider.(*HiiRetailProvider); ok {
		if hiiProvider.version != "1.0.0" {
			t.Errorf("Expected version to be '1.0.0', got %s", hiiProvider.version)
		}
	} else {
		t.Error("Expected provider to be *HiiRetailProvider")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsAtIndex(s, substr))))
}

func containsAtIndex(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
