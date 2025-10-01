package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// OAuth2IntegrationTestSuite provides end-to-end testing for OAuth2 authentication flow
type OAuth2IntegrationTestSuite struct {
	suite.Suite
	mockOCMSServer   *httptest.Server
	mockAPIServer    *httptest.Server
	authClient       *AuthClient
	discoveryClient  *DiscoveryClient
	testTenantID     string
	testClientID     string
	testClientSecret string
}

// SetupSuite initializes the test suite with mock OCMS and API servers
func (suite *OAuth2IntegrationTestSuite) SetupSuite() {
	suite.testTenantID = "integration-test-tenant-123"
	suite.testClientID = "integration-test-client-456"
	suite.testClientSecret = "integration-test-secret-789-very-secure"

	// Setup mock OCMS server
	suite.setupMockOCMSServer()

	// Setup mock API server
	suite.setupMockAPIServer()

	// Initialize auth client
	config := &AuthClientConfig{
		TenantID:     suite.testTenantID,
		ClientID:     suite.testClientID,
		ClientSecret: suite.testClientSecret,
		BaseURL:      suite.mockAPIServer.URL,
		TokenURL:     suite.mockOCMSServer.URL + "/oauth2/token",
		Scopes:       []string{"iam:read", "iam:write"},
		Timeout:      30 * time.Second,
		MaxRetries:   3,
	}

	var err error
	suite.authClient, err = NewAuthClient(config)
	suite.Require().NoError(err, "Failed to create auth client")

	// Initialize discovery client
	suite.discoveryClient = NewDiscoveryClient(suite.mockOCMSServer.URL, 30*time.Second)
}

// TearDownSuite cleans up test resources
func (suite *OAuth2IntegrationTestSuite) TearDownSuite() {
	if suite.mockOCMSServer != nil {
		suite.mockOCMSServer.Close()
	}
	if suite.mockAPIServer != nil {
		suite.mockAPIServer.Close()
	}
}

func (suite *OAuth2IntegrationTestSuite) setupMockOCMSServer() {
	suite.mockOCMSServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			suite.handleDiscoveryEndpoint(w, r)
		case "/oauth2/token":
			suite.handleTokenEndpoint(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func (suite *OAuth2IntegrationTestSuite) setupMockAPIServer() {
	suite.mockAPIServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authentication header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			suite.respondWithError(w, http.StatusUnauthorized, "missing_authorization", "Authorization header required")
			return
		}

		if !suite.isValidAuthHeader(authHeader) {
			suite.respondWithError(w, http.StatusUnauthorized, "invalid_token", "Invalid or expired access token")
			return
		}

		// Verify tenant ID header
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID != suite.testTenantID {
			suite.respondWithError(w, http.StatusBadRequest, "invalid_tenant", "Invalid tenant ID")
			return
		}

		// Handle different API endpoints
		switch r.URL.Path {
		case "/iam/v1/groups":
			suite.handleGroupsAPI(w, r)
		case "/iam/v1/roles":
			suite.handleRolesAPI(w, r)
		case "/iam/v1/custom-roles":
			suite.handleCustomRolesAPI(w, r)
		case "/iam/v1/role-bindings":
			suite.handleRoleBindingsAPI(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "endpoint_not_found", "message": "API endpoint not found"}`))
		}
	}))
}

func (suite *OAuth2IntegrationTestSuite) handleDiscoveryEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	discoveryResponse := map[string]interface{}{
		"issuer":                                suite.mockOCMSServer.URL,
		"token_endpoint":                        suite.mockOCMSServer.URL + "/oauth2/token",
		"authorization_endpoint":                suite.mockOCMSServer.URL + "/oauth2/authorize",
		"jwks_uri":                              suite.mockOCMSServer.URL + "/.well-known/jwks.json",
		"grant_types_supported":                 []string{"client_credentials", "authorization_code"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
		"response_types_supported":              []string{"code", "token"},
		"scopes_supported":                      []string{"iam:read", "iam:write", "iam:admin"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(discoveryResponse)
}

func (suite *OAuth2IntegrationTestSuite) handleTokenEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		suite.respondWithError(w, http.StatusBadRequest, "invalid_request", "Failed to parse form data")
		return
	}

	// Validate grant type
	grantType := r.Form.Get("grant_type")
	if grantType != "client_credentials" {
		suite.respondWithError(w, http.StatusBadRequest, "unsupported_grant_type", "Only client_credentials grant type is supported")
		return
	}

	// Validate credentials
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

	if clientID != suite.testClientID || clientSecret != suite.testClientSecret {
		suite.respondWithError(w, http.StatusUnauthorized, "invalid_client", "Client authentication failed")
		return
	}

	// Generate mock token
	tokenResponse := map[string]interface{}{
		"access_token": fmt.Sprintf("integration-test-token-%d", time.Now().Unix()),
		"token_type":   "Bearer",
		"expires_in":   3600,
		"scope":        r.Form.Get("scope"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse)
}

func (suite *OAuth2IntegrationTestSuite) isValidAuthHeader(authHeader string) bool {
	return len(authHeader) > 7 && authHeader[:7] == "Bearer " && len(authHeader) > 20
}

func (suite *OAuth2IntegrationTestSuite) respondWithError(w http.ResponseWriter, statusCode int, errorCode, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]interface{}{
		"error":             errorCode,
		"error_description": description,
	}
	json.NewEncoder(w).Encode(errorResponse)
}

func (suite *OAuth2IntegrationTestSuite) handleGroupsAPI(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"groups": []map[string]interface{}{
			{"id": "group-1", "name": "Test Group 1", "description": "Integration test group 1"},
			{"id": "group-2", "name": "Test Group 2", "description": "Integration test group 2"},
		},
		"total": 2,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (suite *OAuth2IntegrationTestSuite) handleRolesAPI(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"roles": []map[string]interface{}{
			{"id": "role-1", "name": "Test Role 1", "permissions": []string{"read", "write"}},
			{"id": "role-2", "name": "Test Role 2", "permissions": []string{"read", "admin"}},
		},
		"total": 2,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (suite *OAuth2IntegrationTestSuite) handleCustomRolesAPI(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"custom_roles": []map[string]interface{}{
			{"id": "custom-role-1", "name": "Custom Test Role", "permissions": []string{"custom:read", "custom:write"}},
		},
		"total": 1,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (suite *OAuth2IntegrationTestSuite) handleRoleBindingsAPI(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"role_bindings": []map[string]interface{}{
			{"id": "binding-1", "role_id": "role-1", "subject_id": "user-1", "subject_type": "user"},
			{"id": "binding-2", "role_id": "role-2", "subject_id": "group-1", "subject_type": "group"},
		},
		"total": 2,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Test OAuth2 discovery flow
func (suite *OAuth2IntegrationTestSuite) TestOAuth2Discovery() {
	ctx := context.Background()

	// Test discovery endpoint
	response, err := suite.discoveryClient.FetchDiscovery(ctx)
	suite.Require().NoError(err, "Discovery should succeed")

	suite.Equal(suite.mockOCMSServer.URL, response.Issuer)
	suite.Equal(suite.mockOCMSServer.URL+"/oauth2/token", response.TokenEndpoint)
	suite.Contains(response.GrantTypesSupported, "client_credentials")
	suite.Contains(response.ScopesSupported, "iam:read")
	suite.Contains(response.ScopesSupported, "iam:write")
}

// Test end-to-end OAuth2 authentication flow
func (suite *OAuth2IntegrationTestSuite) TestEndToEndOAuth2Flow() {
	ctx := context.Background()

	// Step 1: Acquire OAuth2 token
	token, err := suite.authClient.GetToken(ctx)
	suite.Require().NoError(err, "Token acquisition should succeed")
	suite.NotNil(token, "Token should not be nil")
	suite.True(token.Valid(), "Token should be valid")
	suite.Equal("Bearer", token.TokenType)

	// Step 2: Use token for authenticated API calls
	httpClient, err := suite.authClient.HTTPClient(ctx)
	suite.Require().NoError(err, "HTTP client creation should succeed")

	// Test various API endpoints
	endpoints := []string{
		"/iam/v1/groups",
		"/iam/v1/roles",
		"/iam/v1/custom-roles",
		"/iam/v1/role-bindings",
	}

	for _, endpoint := range endpoints {
		resp, err := httpClient.Get(suite.mockAPIServer.URL + endpoint)
		suite.Require().NoError(err, "API call to %s should succeed", endpoint)
		suite.Equal(http.StatusOK, resp.StatusCode, "API call to %s should return 200", endpoint)
		resp.Body.Close()
	}
}

// Test concurrent OAuth2 operations
func (suite *OAuth2IntegrationTestSuite) TestConcurrentOAuth2Operations() {
	ctx := context.Background()
	numGoroutines := 5
	var wg sync.WaitGroup

	results := make([]error, numGoroutines)

	// Test concurrent token acquisition and API calls
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Get token
			token, err := suite.authClient.GetToken(ctx)
			if err != nil {
				results[index] = fmt.Errorf("token acquisition failed: %w", err)
				return
			}

			// Use token for API call
			httpClient, err := suite.authClient.HTTPClient(ctx)
			if err != nil {
				results[index] = fmt.Errorf("HTTP client creation failed: %w", err)
				return
			}

			resp, err := httpClient.Get(suite.mockAPIServer.URL + "/iam/v1/groups")
			if err != nil {
				results[index] = fmt.Errorf("API call failed: %w", err)
				return
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results[index] = fmt.Errorf("API call returned status %d", resp.StatusCode)
				return
			}

			results[index] = nil // Success
		}(i)
	}

	wg.Wait()

	// Verify all concurrent operations succeeded
	for i, err := range results {
		suite.NoError(err, "Concurrent operation %d should succeed", i)
	}
}

// Test token refresh scenarios
func (suite *OAuth2IntegrationTestSuite) TestTokenRefreshScenarios() {
	ctx := context.Background()

	// Create client with short token expiration for testing
	shortExpiryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth2/token" {
			tokenResponse := map[string]interface{}{
				"access_token": fmt.Sprintf("short-lived-token-%d", time.Now().UnixNano()),
				"token_type":   "Bearer",
				"expires_in":   1, // 1 second expiration
				"scope":        "iam:read iam:write",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tokenResponse)
		}
	}))
	defer shortExpiryServer.Close()

	config := &AuthClientConfig{
		TenantID:     suite.testTenantID,
		ClientID:     suite.testClientID,
		ClientSecret: suite.testClientSecret,
		TokenURL:     shortExpiryServer.URL + "/oauth2/token",
		Scopes:       []string{"iam:read", "iam:write"},
		Timeout:      30 * time.Second,
	}

	authClient, err := NewAuthClient(config)
	suite.Require().NoError(err, "Auth client creation should succeed")

	// Get initial token
	token1, err := authClient.GetToken(ctx)
	suite.Require().NoError(err, "Initial token acquisition should succeed")

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Get token again - should trigger refresh
	token2, err := authClient.GetToken(ctx)
	suite.Require().NoError(err, "Token refresh should succeed")

	// Tokens should be different (new token acquired)
	suite.NotEqual(token1.AccessToken, token2.AccessToken, "Refreshed token should be different")
}

// Test error handling scenarios
func (suite *OAuth2IntegrationTestSuite) TestErrorHandlingScenarios() {
	ctx := context.Background()

	testCases := []struct {
		name              string
		serverHandler     http.HandlerFunc
		expectedError     string
		expectedType      AuthErrorType
		expectedRetryable bool
	}{
		{
			name: "invalid_credentials",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				suite.respondWithError(w, http.StatusUnauthorized, "invalid_client", "Client authentication failed")
			}),
			expectedError:     "Client authentication failed",
			expectedType:      AuthErrorCredentials,
			expectedRetryable: false,
		},
		{
			name: "server_error",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				suite.respondWithError(w, http.StatusInternalServerError, "server_error", "Internal server error")
			}),
			expectedError:     "Internal server error",
			expectedType:      AuthErrorServerError,
			expectedRetryable: true,
		},
		{
			name: "rate_limiting",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Retry-After", "60")
				suite.respondWithError(w, http.StatusTooManyRequests, "rate_limited", "Too many requests")
			}),
			expectedError:     "Too many requests",
			expectedType:      AuthErrorRateLimit,
			expectedRetryable: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			errorServer := httptest.NewServer(tc.serverHandler)
			defer errorServer.Close()

			config := &AuthClientConfig{
				TenantID:     suite.testTenantID,
				ClientID:     suite.testClientID,
				ClientSecret: suite.testClientSecret,
				TokenURL:     errorServer.URL + "/oauth2/token",
				Scopes:       []string{"iam:read"},
				Timeout:      30 * time.Second,
			}

			authClient, err := NewAuthClient(config)
			suite.Require().NoError(err, "Auth client creation should succeed")

			_, err = authClient.GetToken(ctx)
			suite.Require().Error(err, "Token acquisition should fail")

			var authErr *AuthError
			suite.Require().ErrorAs(err, &authErr, "Error should be AuthError type")
			suite.Equal(tc.expectedType, authErr.Type, "Error type should match")
			suite.Contains(authErr.Message, tc.expectedError, "Error message should contain expected text")
			suite.Equal(tc.expectedRetryable, authErr.Retryable, "Retryable flag should match")
		})
	}
}

// Test long-running operations with token lifecycle
func (suite *OAuth2IntegrationTestSuite) TestLongRunningOperations() {
	ctx := context.Background()

	// Simulate long-running operation with multiple API calls
	httpClient, err := suite.authClient.HTTPClient(ctx)
	suite.Require().NoError(err, "HTTP client creation should succeed")

	// Make multiple API calls over time
	for i := 0; i < 5; i++ {
		resp, err := httpClient.Get(suite.mockAPIServer.URL + "/iam/v1/groups")
		suite.Require().NoError(err, "API call %d should succeed", i+1)
		suite.Equal(http.StatusOK, resp.StatusCode, "API call %d should return 200", i+1)
		resp.Body.Close()

		// Wait between calls
		time.Sleep(100 * time.Millisecond)
	}
}

// TestOAuth2IntegrationSuite runs the integration test suite
func TestOAuth2IntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	suite.Run(t, new(OAuth2IntegrationTestSuite))
}

// TestRealOCMSIntegration tests against real OCMS endpoints (optional)
func TestRealOCMSIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real OCMS integration test in short mode")
	}

	// Check if integration test environment is available
	testBaseURL := os.Getenv("INTEGRATION_TEST_BASE_URL")
	testTenantID := os.Getenv("INTEGRATION_TEST_TENANT_ID")
	testClientID := os.Getenv("INTEGRATION_TEST_CLIENT_ID")
	testClientSecret := os.Getenv("INTEGRATION_TEST_CLIENT_SECRET")

	if testBaseURL == "" || testTenantID == "" || testClientID == "" || testClientSecret == "" {
		t.Skip("Integration test environment variables not set")
	}

	ctx := context.Background()

	// Test discovery against real OCMS
	discoveryClient := NewDiscoveryClient("https://auth.retailsvc-test.com", 30*time.Second)

	response, err := discoveryClient.FetchDiscovery(ctx)
	require.NoError(t, err, "Real OCMS discovery should succeed")
	assert.Equal(t, "https://auth.retailsvc-test.com", response.Issuer)
	assert.Contains(t, response.GrantTypesSupported, "client_credentials")

	// Test token acquisition with real credentials
	config := &AuthClientConfig{
		TenantID:     testTenantID,
		ClientID:     testClientID,
		ClientSecret: testClientSecret,
		BaseURL:      testBaseURL,
		Scopes:       []string{"iam:read"},
		Timeout:      30 * time.Second,
	}

	authClient, err := NewAuthClient(config)
	require.NoError(t, err, "Real auth client creation should succeed")

	token, err := authClient.GetToken(ctx)
	require.NoError(t, err, "Real token acquisition should succeed")
	assert.True(t, token.Valid(), "Real token should be valid")

	// Test API call with real token
	httpClient, err := authClient.HTTPClient(ctx)
	require.NoError(t, err, "Real HTTP client creation should succeed")

	resp, err := httpClient.Get(testBaseURL + "/iam/v1/groups")
	if err == nil {
		defer resp.Body.Close()
		// Real API might return different status codes based on permissions
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusForbidden,
			"Real API call should return 200 OK or 403 Forbidden")
	}
}
