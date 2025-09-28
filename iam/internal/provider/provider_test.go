package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestHiiRetailIamProvider(t *testing.T) {
	t.Run("Provider metadata", func(t *testing.T) {
		p := &HiiRetailIamProvider{version: "test"}
		resp := &provider.MetadataResponse{}
		p.Metadata(context.Background(), provider.MetadataRequest{}, resp)

		if resp.TypeName != "hiiretail_iam" {
			t.Errorf("Expected TypeName to be 'hiiretail_iam', got %s", resp.TypeName)
		}
		if resp.Version != "test" {
			t.Errorf("Expected Version to be 'test', got %s", resp.Version)
		}
	})

	t.Run("Provider schema", func(t *testing.T) {
		p := &HiiRetailIamProvider{}
		resp := &provider.SchemaResponse{}
		p.Schema(context.Background(), provider.SchemaRequest{}, resp)

		// Check that all required attributes are present
		expectedAttrs := []string{"tenant_id", "base_url", "client_id", "client_secret"}
		for _, attr := range expectedAttrs {
			if _, exists := resp.Schema.Attributes[attr]; !exists {
				t.Errorf("Expected attribute %s to be present in schema", attr)
			}
		}

		// Check required attributes
		if !resp.Schema.Attributes["tenant_id"].IsRequired() {
			t.Error("tenant_id should be required")
		}
		if !resp.Schema.Attributes["client_id"].IsRequired() {
			t.Error("client_id should be required")
		}
		if !resp.Schema.Attributes["client_secret"].IsRequired() {
			t.Error("client_secret should be required")
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

func TestHiiRetailIamProvider_Configure(t *testing.T) {
	testCases := []struct {
		name          string
		config        map[string]tftypes.Value
		expectedError string
	}{
		{
			name: "Valid configuration with all required fields",
			config: map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, "test-tenant"),
				"client_id":     tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret": tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":      tftypes.NewValue(tftypes.String, "https://test-api.example.com"),
			},
			expectedError: "",
		},
		{
			name: "Valid configuration with default base_url",
			config: map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, "test-tenant"),
				"client_id":     tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret": tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":      tftypes.NewValue(tftypes.String, nil),
			},
			expectedError: "",
		},
		{
			name: "Missing tenant_id",
			config: map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, nil),
				"client_id":     tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret": tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":      tftypes.NewValue(tftypes.String, nil),
			},
			expectedError: "Missing tenant_id",
		},
		{
			name: "Missing client_id",
			config: map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, "test-tenant"),
				"client_id":     tftypes.NewValue(tftypes.String, nil),
				"client_secret": tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":      tftypes.NewValue(tftypes.String, nil),
			},
			expectedError: "Missing client_id",
		},
		{
			name: "Missing client_secret",
			config: map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, "test-tenant"),
				"client_id":     tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret": tftypes.NewValue(tftypes.String, nil),
				"base_url":      tftypes.NewValue(tftypes.String, nil),
			},
			expectedError: "Missing client_secret",
		},
		{
			name: "Invalid base_url",
			config: map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, "test-tenant"),
				"client_id":     tftypes.NewValue(tftypes.String, "test-client-id"),
				"client_secret": tftypes.NewValue(tftypes.String, "test-client-secret"),
				"base_url":      tftypes.NewValue(tftypes.String, "not-a-valid-url"),
			},
			expectedError: "Invalid base_url",
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
					found := false
					for _, diag := range resp.Diagnostics.Errors() {
						if contains(diag.Summary(), tc.expectedError) {
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

func TestHiiRetailIamProvider_OIDCConfiguration(t *testing.T) {
	testCases := []struct {
		name         string
		clientId     string
		clientSecret string
		baseUrl      string
		expectError  bool
	}{
		{
			name:         "Valid OIDC configuration",
			clientId:     "test-client",
			clientSecret: "test-secret",
			baseUrl:      "https://test-api.example.com",
			expectError:  false,
		},
		{
			name:         "Valid OIDC with default URL",
			clientId:     "test-client",
			clientSecret: "test-secret",
			baseUrl:      "",
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &HiiRetailIamProvider{}

			// Create schema
			schemaResp := &provider.SchemaResponse{}
			p.Schema(context.Background(), provider.SchemaRequest{}, schemaResp)

			// Create configuration
			configMap := map[string]tftypes.Value{
				"tenant_id":     tftypes.NewValue(tftypes.String, "test-tenant"),
				"client_id":     tftypes.NewValue(tftypes.String, tc.clientId),
				"client_secret": tftypes.NewValue(tftypes.String, tc.clientSecret),
			}

			if tc.baseUrl != "" {
				configMap["base_url"] = tftypes.NewValue(tftypes.String, tc.baseUrl)
			} else {
				configMap["base_url"] = tftypes.NewValue(tftypes.String, nil)
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
				if apiClient, ok := resp.ResourceData.(*APIClient); ok {
					if apiClient.TenantID != "test-tenant" {
						t.Errorf("Expected TenantID to be 'test-tenant', got %s", apiClient.TenantID)
					}

					expectedBaseURL := tc.baseUrl
					if expectedBaseURL == "" {
						expectedBaseURL = "https://iam-api.retailsvc-test.com"
					}
					if apiClient.BaseURL != expectedBaseURL {
						t.Errorf("Expected BaseURL to be '%s', got %s", expectedBaseURL, apiClient.BaseURL)
					}

					if apiClient.HTTPClient == nil {
						t.Error("Expected HTTPClient to be configured")
					}
				} else {
					t.Error("Expected ResourceData to be *APIClient")
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	providerFunc := New("1.0.0")
	provider := providerFunc()

	if hiiProvider, ok := provider.(*HiiRetailIamProvider); ok {
		if hiiProvider.version != "1.0.0" {
			t.Errorf("Expected version to be '1.0.0', got %s", hiiProvider.version)
		}
	} else {
		t.Error("Expected provider to be *HiiRetailIamProvider")
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
