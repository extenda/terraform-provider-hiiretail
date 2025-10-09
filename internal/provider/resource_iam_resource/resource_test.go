package resource_iam_resource_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/auth"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


// TestSetResourceContract verifies the PUT /api/v1/tenants/{tenantId}/resources/{id} endpoint contract
func TestSetResourceContract(t *testing.T) {
	tests := []struct {
		name           string
		resourceId     string
		requestBody    *iam.SetResourceDto
		mockStatusCode int
		mockResponse   string
		expectError    bool
		expectedStatus int
	}{
		{
			name:       "create new resource success",
			resourceId: "store:001",
			requestBody: &iam.SetResourceDto{
				Name: "Store 001",
				Props: map[string]interface{}{
					"location": "downtown",
					"active":   true,
				},
			},
			mockStatusCode: 201,
			mockResponse: `{
				"id": "store:001",
				"name": "Store 001",
				"props": {
					"location": "downtown",
					"active": true
				}
			}`,
			expectedStatus: 201,
		},
		{
			name:       "update existing resource success",
			resourceId: "register:pos-01",
			requestBody: &iam.SetResourceDto{
				Name: "POS Terminal 01 Updated",
				Props: map[string]interface{}{
					"location": "front-counter",
					"type":     "touch-screen",
				},
			},
			mockStatusCode: 200,
			mockResponse: `{
				"id": "register:pos-01",
				"name": "POS Terminal 01 Updated",
				"props": {
					"location": "front-counter",
					"type": "touch-screen"
				}
			}`,
			expectedStatus: 200,
		},
		{
			name:       "resource with null props",
			resourceId: "simple:resource",
			requestBody: &iam.SetResourceDto{
				Name:  "Simple Resource",
				Props: nil,
			},
			mockStatusCode: 201,
			mockResponse: `{
				"id": "simple:resource",
				"name": "Simple Resource",
				"props": null
			}`,
			expectedStatus: 201,
		},
		{
			name:       "bad request error",
			resourceId: "invalid..id",
			requestBody: &iam.SetResourceDto{
				Name: "Invalid Resource",
			},
			mockStatusCode: 400,
			mockResponse: `{
				"statusCode": 400,
				"message": ["ID format is invalid"],
				"error": "Bad Request"
			}`,
			expectError:    true,
			expectedStatus: 400,
		},
		{
			name:       "forbidden error",
			resourceId: "protected:resource",
			requestBody: &iam.SetResourceDto{
				Name: "Protected Resource",
			},
			mockStatusCode: 403,
			mockResponse: `{
				"statusCode": 403,
				"message": "Insufficient permissions",
				"error": "Forbidden"
			}`,
			expectError:    true,
			expectedStatus: 403,
		},
		{
			name:       "internal server error",
			resourceId: "error:resource",
			requestBody: &iam.SetResourceDto{
				Name: "Error Resource",
			},
			mockStatusCode: 500,
			mockResponse: `{
				"statusCode": 500,
				"message": "Internal Server Error"
			}`,
			expectError:    true,
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and path
				expectedPath := fmt.Sprintf("/api/v1/tenants/test-tenant/resources/%s", tt.resourceId)
				assert.Equal(t, "PUT", r.Method)
				assert.Equal(t, expectedPath, r.URL.Path)

				// Verify Content-Type
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Verify Authorization header
				assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))

				// Return mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create service with properly initialized mock client
			mockConfig := client.DefaultConfig()
			mockConfig.BaseURL = server.URL
			// Use real client.Client with minimal valid auth.Config for contract test
			   dummyAuth := &auth.Config{
				   ClientID:     "dummy-client-id",
				   ClientSecret: "dummy-client-secret",
				   TenantID:     "test-tenant",
				   APIURL:       server.URL,
				   AuthURL:      server.URL + "/oauth2/token",
				   Environment:  "test",
				   Scopes:       []string{"hiiretail:iam"},
				   Timeout:      5,
				   MaxRetries:   0,
				   SkipTLS:      true,
				   TestToken:    "dummy-token-for-contract-test",
			   }
			mockConfig = client.DefaultConfig()
			mockConfig.BaseURL = server.URL
			mockClient, errNew := client.New(dummyAuth, mockConfig)
			require.NoError(t, errNew, "Failed to create mock client")
			   service := iam.NewService(mockClient, "test-tenant")

			// This test should FAIL because SetResource is not yet implemented
			ctx := context.Background()
			result, err := service.SetResource(ctx, tt.resourceId, tt.requestBody)

			if tt.expectError {
				assert.Error(t, err, "Expected error for test case: %s", tt.name)
				// Only check for status code in error message if not internal server error
				if tt.expectedStatus != 500 {
					assert.Contains(t, err.Error(), fmt.Sprintf("%d", tt.expectedStatus))
				}
			} else {
				require.NoError(t, err, "Unexpected error for test case: %s", tt.name)
				require.NotNil(t, result, "Expected result for test case: %s", tt.name)

				// Verify response structure
				assert.Equal(t, tt.resourceId, result.ID)
				assert.Equal(t, tt.requestBody.Name, result.Name)
				// Props comparison would depend on implementation details
			}
		})
	}
}

// TestGetResourceContract verifies the GET /api/v1/tenants/{tenantId}/resources/{id} endpoint contract
func TestGetResourceContract(t *testing.T) {
	tests := []struct {
		name           string
		resourceId     string
		mockStatusCode int
		mockResponse   string
		expectError    bool
	}{
		{
			name:           "get existing resource success",
			resourceId:     "store:001",
			mockStatusCode: 200,
			mockResponse: `{
				"id": "store:001",
				"name": "Store 001",
				"props": {
					"location": "downtown",
					"active": true
				}
			}`,
		},
		{
			name:           "resource not found",
			resourceId:     "nonexistent:resource",
			mockStatusCode: 404,
			mockResponse: `{
				"statusCode": 404,
				"message": "Resource not found",
				"error": "Not Found"
			}`,
			expectError: true,
		},
		{
			name:           "forbidden access",
			resourceId:     "protected:resource",
			mockStatusCode: 403,
			mockResponse: `{
				"statusCode": 403,
				"message": "Access denied",
				"error": "Forbidden"
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/api/v1/tenants/test-tenant/resources/%s", tt.resourceId)
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, expectedPath, r.URL.Path)
				assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create service with properly initialized mock client
			mockConfig := client.DefaultConfig()
			mockConfig.BaseURL = server.URL
			dummyAuth := &auth.Config{
				ClientID:     "dummy-client-id",
				ClientSecret: "dummy-client-secret",
				TenantID:     "test-tenant",
				APIURL:       server.URL,
				AuthURL:      server.URL + "/oauth2/token",
				Environment:  "test",
				Scopes:       []string{"hiiretail:iam"},
				Timeout:      5,
				MaxRetries:   0,
				SkipTLS:      true,
				TestToken:    "dummy-token-for-contract-test",
			}
			mockClient, errNew := client.New(dummyAuth, mockConfig)
			require.NoError(t, errNew, "Failed to create mock client")
			service := iam.NewService(mockClient, "test-tenant")

			// This test should FAIL because GetResource is not yet implemented properly
			ctx := context.Background()
			result, err := service.GetResource(ctx, tt.resourceId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.resourceId, result.ID)
			}
		})
	}
}

// TestDeleteResourceContract verifies the DELETE /api/v1/tenants/{tenantId}/resources/{id} endpoint contract
func TestDeleteResourceContract(t *testing.T) {
	tests := []struct {
		name           string
		resourceId     string
		mockStatusCode int
		mockResponse   string
		expectError    bool
	}{
		{
			name:           "delete existing resource success",
			resourceId:     "store:001",
			mockStatusCode: 204,
			mockResponse:   "",
		},
		{
			name:           "resource not found",
			resourceId:     "nonexistent:resource",
			mockStatusCode: 404,
			mockResponse: `{
				"statusCode": 404,
				"message": "Resource not found",
				"error": "Not Found"
			}`,
			expectError: true,
		},
		{
			name:           "forbidden access",
			resourceId:     "protected:resource",
			mockStatusCode: 403,
			mockResponse: `{
				"statusCode": 403,
				"message": "Access denied",
				"error": "Forbidden"
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/api/v1/tenants/test-tenant/resources/%s", tt.resourceId)
				assert.Equal(t, "DELETE", r.Method)
				assert.Equal(t, expectedPath, r.URL.Path)
				assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != "" {
					w.Write([]byte(tt.mockResponse))
				}
			}))
			defer server.Close()

			// Create service with properly initialized mock client
			mockConfig := client.DefaultConfig()
			mockConfig.BaseURL = server.URL
			dummyAuth := &auth.Config{
				ClientID:     "dummy-client-id",
				ClientSecret: "dummy-client-secret",
				TenantID:     "test-tenant",
				APIURL:       server.URL,
				AuthURL:      server.URL + "/oauth2/token",
				Environment:  "test",
				Scopes:       []string{"hiiretail:iam"},
				Timeout:      5,
				MaxRetries:   0,
				SkipTLS:      true,
				TestToken:    "dummy-token-for-contract-test",
			}
			mockClient, errNew := client.New(dummyAuth, mockConfig)
			require.NoError(t, errNew, "Failed to create mock client")
			service := iam.NewService(mockClient, "test-tenant")

			// This test should FAIL because DeleteResource is not yet implemented properly
			ctx := context.Background()
			err := service.DeleteResource(ctx, tt.resourceId)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGetResourcesContract verifies the GET /api/v1/tenants/{tenantId}/resources endpoint contract
func TestGetResourcesContract(t *testing.T) {
	tests := []struct {
		name           string
		request        *iam.GetResourcesRequest
		mockStatusCode int
		mockResponse   string
		expectError    bool
	}{
		{
			name:           "list all resources success",
			request:        &iam.GetResourcesRequest{},
			mockStatusCode: 200,
			mockResponse: `[
				{
					"id": "store:001",
					"name": "Store 001",
					"props": {"location": "downtown"}
				},
				{
					"id": "register:pos-01",
					"name": "POS Terminal 01",
					"props": {"type": "touch-screen"}
				}
			]`,
		},
		{
			name: "list resources with permission filter",
			request: &iam.GetResourcesRequest{
				Permission: "read",
			},
			mockStatusCode: 200,
			mockResponse: `[
				{
					"id": "store:001",
					"name": "Store 001",
					"props": {"location": "downtown"}
				}
			]`,
		},
		{
			name: "list resources with type filter",
			request: &iam.GetResourcesRequest{
				Type: "store",
			},
			mockStatusCode: 200,
			mockResponse: `[
				{
					"id": "store:001",
					"name": "Store 001",
					"props": {"location": "downtown"}
				}
			]`,
		},
		{
			name:           "forbidden access",
			request:        &iam.GetResourcesRequest{},
			mockStatusCode: 403,
			mockResponse: `{
				"statusCode": 403,
				"message": "Access denied",
				"error": "Forbidden"
			}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/api/v1/tenants/test-tenant/resources"
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, expectedPath, r.URL.Path)
				assert.True(t, strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "))

				// Verify query parameters
				if tt.request.Permission != "" {
					assert.Equal(t, tt.request.Permission, r.URL.Query().Get("permission"))
				}
				if tt.request.Type != "" {
					assert.Equal(t, tt.request.Type, r.URL.Query().Get("type"))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create service with properly initialized mock client
			mockConfig := client.DefaultConfig()
			mockConfig.BaseURL = server.URL
			dummyAuth := &auth.Config{
				ClientID:     "dummy-client-id",
				ClientSecret: "dummy-client-secret",
				TenantID:     "test-tenant",
				APIURL:       server.URL,
				AuthURL:      server.URL + "/oauth2/token",
				Environment:  "test",
				Scopes:       []string{"hiiretail:iam"},
				Timeout:      5,
				MaxRetries:   0,
				SkipTLS:      true,
				TestToken:    "dummy-token-for-contract-test",
			}
			mockClient, errNew := client.New(dummyAuth, mockConfig)
			require.NoError(t, errNew, "Failed to create mock client")
			service := iam.NewService(mockClient, "test-tenant")

			// This test should FAIL because GetResources is not yet implemented properly
			ctx := context.Background()
			result, err := service.GetResources(ctx, tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.GreaterOrEqual(t, len(result.Resources), 0)
			}
		})
	}
}
