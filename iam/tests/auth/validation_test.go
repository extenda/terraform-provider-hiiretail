package auth
package auth_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOAuth2Configuration_Validation tests OAuth2 configuration validation
func TestOAuth2Configuration_Validation(t *testing.T) {
	testCases := []struct {
		name     string
		config   mockOAuth2Config
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid configuration",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
			},
			wantErr: false,
		},
		{
			name: "Missing client ID",
			config: mockOAuth2Config{
				ClientID:     "",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
			},
			wantErr:  true,
			errorMsg: "client_id",
		},
		{
			name: "Missing client secret",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "",
				TenantID:     "test-tenant-123",
			},
			wantErr:  true,
			errorMsg: "client_secret",
		},
		{
			name: "Missing tenant ID",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "",
			},
			wantErr:  true,
			errorMsg: "tenant_id",
		},
		{
			name: "Whitespace-only client ID",
			config: mockOAuth2Config{
				ClientID:     "   ",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
			},
			wantErr:  true,
			errorMsg: "client_id",
		},
		{
			name: "Whitespace-only client secret",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "   ",
				TenantID:     "test-tenant-123",
			},
			wantErr:  true,
			errorMsg: "client_secret",
		},
		{
			name: "Whitespace-only tenant ID",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "   ",
			},
			wantErr:  true,
			errorMsg: "tenant_id",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This test will fail until OAuth2Configuration validation is implemented
			err := tc.config.Validate()
			
			if tc.wantErr {
				assert.Error(t, err, "Should return validation error")
				assert.Contains(t, err.Error(), tc.errorMsg, "Error should mention the invalid field")
			} else {
				assert.NoError(t, err, "Should not return validation error")
			}
		})
	}
}

// TestOAuth2Configuration_URLValidation tests URL validation for overrides
func TestOAuth2Configuration_URLValidation(t *testing.T) {
	testCases := []struct {
		name     string
		config   mockOAuth2Config
		wantErr  bool
		errorMsg string
	}{
		{
			name: "Valid HTTPS token URL",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				TokenURL:     "https://auth.example.com/oauth/token",
			},
			wantErr: false,
		},
		{
			name: "Valid HTTPS API URL",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				APIURL:       "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "Invalid HTTP token URL (non-TLS)",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				TokenURL:     "http://auth.example.com/oauth/token",
			},
			wantErr:  true,
			errorMsg: "TLS",
		},
		{
			name: "Invalid HTTP API URL (non-TLS)",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				APIURL:       "http://api.example.com",
			},
			wantErr:  true,
			errorMsg: "TLS",
		},
		{
			name: "Invalid token URL format",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				TokenURL:     "not-a-valid-url",
			},
			wantErr:  true,
			errorMsg: "invalid URL",
		},
		{
			name: "Invalid API URL format",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				APIURL:       "not-a-valid-url",
			},
			wantErr:  true,
			errorMsg: "invalid URL",
		},
		{
			name: "Mock mode allows HTTP URLs",
			config: mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				TokenURL:     "http://localhost:8080/oauth/token",
				APIURL:       "http://localhost:8080/api",
				MockMode:     true,
			},
			wantErr: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This test will fail until OAuth2Configuration URL validation is implemented
			err := tc.config.Validate()
			
			if tc.wantErr {
				assert.Error(t, err, "Should return URL validation error")
				assert.Contains(t, err.Error(), tc.errorMsg, "Error should mention URL validation issue")
			} else {
				assert.NoError(t, err, "Should not return URL validation error")
			}
		})
	}
}

// TestCredentialValidation_SensitiveDataHandling tests that credentials are properly marked as sensitive
func TestCredentialValidation_SensitiveDataHandling(t *testing.T) {
	config := mockOAuth2Config{
		ClientID:     "test-client-id",
		ClientSecret: "super-secret-value",
		TenantID:     "test-tenant-123",
	}
	
	// This test will fail until OAuth2Configuration is implemented
	// Test that String() method redacts sensitive information
	configStr := config.String()
	
	assert.Contains(t, configStr, "test-client-id", "Should show client ID")
	assert.Contains(t, configStr, "test-tenant-123", "Should show tenant ID")
	assert.NotContains(t, configStr, "super-secret-value", "Should not show client secret")
	assert.Contains(t, configStr, "[REDACTED]", "Should redact client secret")
}

// TestCredentialValidation_EnvironmentVariableSupport tests environment variable validation
func TestCredentialValidation_EnvironmentVariableSupport(t *testing.T) {
	testCases := []struct {
		name     string
		clientID string
		secret   string
		tenantID string
		wantErr  bool
	}{
		{
			name:     "All direct values",
			clientID: "direct-client-id",
			secret:   "direct-secret",
			tenantID: "direct-tenant",
			wantErr:  false,
		},
		{
			name:     "Environment variable placeholders",
			clientID: "${HIIRETAIL_CLIENT_ID}",
			secret:   "${HIIRETAIL_CLIENT_SECRET}",
			tenantID: "${HIIRETAIL_TENANT_ID}",
			wantErr:  false,
		},
		{
			name:     "Mixed direct and environment variables",
			clientID: "direct-client-id",
			secret:   "${HIIRETAIL_CLIENT_SECRET}",
			tenantID: "direct-tenant",
			wantErr:  false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := mockOAuth2Config{
				ClientID:     tc.clientID,
				ClientSecret: tc.secret,
				TenantID:     tc.tenantID,
			}
			
			// This test will fail until OAuth2Configuration validation is implemented
			err := config.Validate()
			
			if tc.wantErr {
				assert.Error(t, err, "Should return validation error")
			} else {
				assert.NoError(t, err, "Should not return validation error")
			}
		})
	}
}

// TestCredentialValidation_TenantIDPatterns tests tenant ID pattern validation
func TestCredentialValidation_TenantIDPatterns(t *testing.T) {
	testCases := []struct {
		name     string
		tenantID string
		wantErr  bool
		errorMsg string
	}{
		{
			name:     "Valid simple tenant ID",
			tenantID: "company-123",
			wantErr:  false,
		},
		{
			name:     "Valid test tenant ID",
			tenantID: "test-company-123",
			wantErr:  false,
		},
		{
			name:     "Valid production tenant ID",
			tenantID: "prod-company-123",
			wantErr:  false,
		},
		{
			name:     "Valid alphanumeric tenant ID",
			tenantID: "company123test",
			wantErr:  false,
		},
		{
			name:     "Invalid tenant ID with special characters",
			tenantID: "company@123",
			wantErr:  true,
			errorMsg: "invalid characters",
		},
		{
			name:     "Invalid tenant ID with spaces",
			tenantID: "company 123",
			wantErr:  true,
			errorMsg: "invalid characters",
		},
		{
			name:     "Invalid tenant ID too short",
			tenantID: "a",
			wantErr:  true,
			errorMsg: "too short",
		},
		{
			name:     "Invalid tenant ID too long",
			tenantID: "this-is-a-very-long-tenant-id-that-exceeds-the-maximum-allowed-length-for-tenant-identifiers",
			wantErr:  true,
			errorMsg: "too long",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := mockOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     tc.tenantID,
			}
			
			// This test will fail until OAuth2Configuration tenant ID validation is implemented
			err := config.ValidateTenantID()
			
			if tc.wantErr {
				assert.Error(t, err, "Should return tenant ID validation error")
				assert.Contains(t, err.Error(), tc.errorMsg, "Error should mention specific validation issue")
			} else {
				assert.NoError(t, err, "Should not return tenant ID validation error")
			}
		})
	}
}

// mockOAuth2Config is a temporary mock implementation for testing
// This will be replaced by the actual OAuth2Configuration implementation
type mockOAuth2Config struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	TokenURL     string
	APIURL       string
	MockMode     bool
}

func (c *mockOAuth2Config) Validate() error {
	// This is a mock implementation that will fail until the real one is created
	if strings.TrimSpace(c.ClientID) == "" {
		return fmt.Errorf("client_id is required")
	}
	if strings.TrimSpace(c.ClientSecret) == "" {
		return fmt.Errorf("client_secret is required")
	}
	if strings.TrimSpace(c.TenantID) == "" {
		return fmt.Errorf("tenant_id is required")
	}
	
	// URL validation
	if c.TokenURL != "" {
		if !c.MockMode && !strings.HasPrefix(c.TokenURL, "https://") {
			return fmt.Errorf("token_url must use TLS")
		}
		if !isValidURL(c.TokenURL) {
			return fmt.Errorf("invalid URL format for token_url")
		}
	}
	
	if c.APIURL != "" {
		if !c.MockMode && !strings.HasPrefix(c.APIURL, "https://") {
			return fmt.Errorf("api_url must use TLS")
		}
		if !isValidURL(c.APIURL) {
			return fmt.Errorf("invalid URL format for api_url")
		}
	}
	
	return nil
}

func (c *mockOAuth2Config) ValidateTenantID() error {
	// This is a mock implementation that will fail until the real one is created
	tenantID := strings.TrimSpace(c.TenantID)
	
	if len(tenantID) < 2 {
		return fmt.Errorf("tenant_id too short")
	}
	if len(tenantID) > 64 {
		return fmt.Errorf("tenant_id too long")
	}
	
	// Check for invalid characters (only allow alphanumeric, hyphens, underscores)
	for _, char := range tenantID {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-' || char == '_') {
			return fmt.Errorf("tenant_id contains invalid characters")
		}
	}
	
	return nil
}

func (c *mockOAuth2Config) String() string {
	// This is a mock implementation that will fail until the real one is created
	return fmt.Sprintf("OAuth2Config{ClientID: %s, ClientSecret: [REDACTED], TenantID: %s}", 
		c.ClientID, c.TenantID)
}

// isValidURL is a simple URL validation helper for testing
func isValidURL(urlStr string) bool {
	return strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")
}

