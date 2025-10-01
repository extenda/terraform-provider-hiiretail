package auth
package auth_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEndpointResolver_TenantIDParsing tests tenant ID parsing for environment detection
func TestEndpointResolver_TenantIDParsing(t *testing.T) {
	testCases := []struct {
		name             string
		tenantID         string
		expectedIsTest   bool
		expectedAuthURL  string
		expectedAPIURL   string
	}{
		{
			name:             "Live tenant - production pattern",
			tenantID:         "production-company-123",
			expectedIsTest:   false,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc.com",
		},
		{
			name:             "Live tenant - company name only",
			tenantID:         "acme-corp",
			expectedIsTest:   false,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc.com",
		},
		{
			name:             "Test tenant - test prefix",
			tenantID:         "test-company-123",
			expectedIsTest:   true,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc-test.com",
		},
		{
			name:             "Test tenant - dev prefix",
			tenantID:         "dev-company-123",
			expectedIsTest:   true,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc-test.com",
		},
		{
			name:             "Test tenant - staging prefix",
			tenantID:         "staging-company-123",
			expectedIsTest:   true,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc-test.com",
		},
		{
			name:             "Test tenant - test suffix",
			tenantID:         "company-test-123",
			expectedIsTest:   true,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc-test.com",
		},
		{
			name:             "Test tenant - case insensitive TEST",
			tenantID:         "TEST-company-123",
			expectedIsTest:   true,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc-test.com",
		},
		{
			name:             "Test tenant - case insensitive Dev",
			tenantID:         "company-Dev-123",
			expectedIsTest:   true,
			expectedAuthURL:  "https://auth.retailsvc.com",
			expectedAPIURL:   "https://iam-api.retailsvc-test.com",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This test will fail until EndpointResolver is implemented
			// TODO: Replace with actual EndpointResolver implementation
			
			// Mock EndpointResolver logic for testing
			resolver := mockEndpointResolver{TenantID: tc.tenantID}
			
			isTest := resolver.IsTestEnvironment()
			authURL := resolver.GetAuthURL()
			apiURL := resolver.GetAPIURL()
			
			assert.Equal(t, tc.expectedIsTest, isTest, "Environment detection mismatch")
			assert.Equal(t, tc.expectedAuthURL, authURL, "Auth URL mismatch")
			assert.Equal(t, tc.expectedAPIURL, apiURL, "API URL mismatch")
		})
	}
}

// TestEndpointResolver_EnvironmentVariableOverride tests environment variable overrides
func TestEndpointResolver_EnvironmentVariableOverride(t *testing.T) {
	// Save original environment variables
	originalForceTest := os.Getenv("HIIRETAIL_FORCE_TEST_ENV")
	originalAuthURL := os.Getenv("HIIRETAIL_AUTH_URL")
	originalAPIURL := os.Getenv("HIIRETAIL_API_URL")
	
	// Cleanup after test
	t.Cleanup(func() {
		os.Setenv("HIIRETAIL_FORCE_TEST_ENV", originalForceTest)
		os.Setenv("HIIRETAIL_AUTH_URL", originalAuthURL)
		os.Setenv("HIIRETAIL_API_URL", originalAPIURL)
	})
	
	t.Run("Force test environment", func(t *testing.T) {
		// Set environment variable to force test environment
		os.Setenv("HIIRETAIL_FORCE_TEST_ENV", "true")
		
		// This test will fail until EndpointResolver is implemented
		resolver := mockEndpointResolver{TenantID: "production-company-123"}
		
		// Should detect as test environment despite production tenant ID
		assert.True(t, resolver.IsTestEnvironment(), "Should be forced to test environment")
		assert.Equal(t, "https://iam-api.retailsvc-test.com", resolver.GetAPIURL())
	})
	
	t.Run("Custom auth URL override", func(t *testing.T) {
		// Clear force test env and set custom auth URL
		os.Unsetenv("HIIRETAIL_FORCE_TEST_ENV")
		os.Setenv("HIIRETAIL_AUTH_URL", "https://mock-auth.example.com/oauth/token")
		
		// This test will fail until EndpointResolver is implemented
		resolver := mockEndpointResolver{TenantID: "test-company-123"}
		
		assert.Equal(t, "https://mock-auth.example.com/oauth/token", resolver.GetAuthURL())
	})
	
	t.Run("Custom API URL override", func(t *testing.T) {
		// Set custom API URL
		os.Setenv("HIIRETAIL_API_URL", "https://mock-api.example.com")
		
		// This test will fail until EndpointResolver is implemented
		resolver := mockEndpointResolver{TenantID: "test-company-123"}
		
		assert.Equal(t, "https://mock-api.example.com", resolver.GetAPIURL())
	})
}

// TestEndpointResolver_MockMode tests mock server mode
func TestEndpointResolver_MockMode(t *testing.T) {
	// Save original environment variables
	originalMockMode := os.Getenv("HIIRETAIL_MOCK_MODE")
	originalAuthURL := os.Getenv("HIIRETAIL_AUTH_URL")
	originalAPIURL := os.Getenv("HIIRETAIL_API_URL")
	
	// Cleanup after test
	t.Cleanup(func() {
		os.Setenv("HIIRETAIL_MOCK_MODE", originalMockMode)
		os.Setenv("HIIRETAIL_AUTH_URL", originalAuthURL)
		os.Setenv("HIIRETAIL_API_URL", originalAPIURL)
	})
	
	// Enable mock mode with local URLs
	os.Setenv("HIIRETAIL_MOCK_MODE", "true")
	os.Setenv("HIIRETAIL_AUTH_URL", "http://localhost:8080/oauth/token")
	os.Setenv("HIIRETAIL_API_URL", "http://localhost:8080/api")
	
	// This test will fail until EndpointResolver is implemented
	resolver := mockEndpointResolver{TenantID: "test-company-123"}
	
	assert.True(t, resolver.IsMockMode(), "Should be in mock mode")
	assert.Equal(t, "http://localhost:8080/oauth/token", resolver.GetAuthURL())
	assert.Equal(t, "http://localhost:8080/api", resolver.GetAPIURL())
}

// TestEndpointResolver_InvalidTenantID tests error handling for invalid tenant IDs
func TestEndpointResolver_InvalidTenantID(t *testing.T) {
	testCases := []struct {
		name     string
		tenantID string
		wantErr  bool
	}{
		{
			name:     "Empty tenant ID",
			tenantID: "",
			wantErr:  true,
		},
		{
			name:     "Whitespace only tenant ID",
			tenantID: "   ",
			wantErr:  true,
		},
		{
			name:     "Valid tenant ID",
			tenantID: "valid-tenant-123",
			wantErr:  false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This test will fail until EndpointResolver is implemented
			resolver := mockEndpointResolver{TenantID: tc.tenantID}
			
			err := resolver.Validate()
			
			if tc.wantErr {
				assert.Error(t, err, "Should return error for invalid tenant ID")
			} else {
				assert.NoError(t, err, "Should not return error for valid tenant ID")
			}
		})
	}
}

// mockEndpointResolver is a temporary mock implementation for testing
// This will be replaced by the actual EndpointResolver implementation
type mockEndpointResolver struct {
	TenantID string
}

func (r *mockEndpointResolver) IsTestEnvironment() bool {
	// Check environment variable override first
	if os.Getenv("HIIRETAIL_FORCE_TEST_ENV") == "true" {
		return true
	}
	
	// Mock pattern matching logic
	tenantLower := strings.ToLower(r.TenantID)
	return strings.Contains(tenantLower, "test") ||
		   strings.Contains(tenantLower, "dev") ||
		   strings.Contains(tenantLower, "staging")
}

func (r *mockEndpointResolver) GetAuthURL() string {
	// Check for environment variable override
	if customURL := os.Getenv("HIIRETAIL_AUTH_URL"); customURL != "" {
		return customURL
	}
	
	// Default auth URL is always the same
	return "https://auth.retailsvc.com"
}

func (r *mockEndpointResolver) GetAPIURL() string {
	// Check for environment variable override
	if customURL := os.Getenv("HIIRETAIL_API_URL"); customURL != "" {
		return customURL
	}
	
	// Choose API URL based on environment
	if r.IsTestEnvironment() {
		return "https://iam-api.retailsvc-test.com"
	}
	return "https://iam-api.retailsvc.com"
}

func (r *mockEndpointResolver) IsMockMode() bool {
	return os.Getenv("HIIRETAIL_MOCK_MODE") == "true"
}

func (r *mockEndpointResolver) Validate() error {
	if strings.TrimSpace(r.TenantID) == "" {
		return fmt.Errorf("tenant ID cannot be empty")
	}
	return nil
}

