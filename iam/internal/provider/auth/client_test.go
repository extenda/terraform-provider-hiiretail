package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// TestAuthClient_TokenAcquisition tests OAuth2 client credentials token acquisition
func TestAuthClient_TokenAcquisition(t *testing.T) {
	validTokenResponse := &oauth2.Token{
		AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test.token",
		TokenType:    "Bearer",
		RefreshToken: "",
		Expiry:       time.Now().Add(1 * time.Hour),
	}

	t.Run("successful_token_acquisition", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/oauth2/token", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

			// Verify request body
			body := parseForm(t, r)
			assert.Equal(t, "client_credentials", body.Get("grant_type"))
			assert.Contains(t, body.Get("scope"), "iam:read")
			assert.Contains(t, body.Get("scope"), "iam:write")

			// Verify authorization header or body credentials
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Credentials in body
				assert.Equal(t, "test-client-123", body.Get("client_id"))
				assert.Equal(t, "test-secret-456", body.Get("client_secret"))
			} else {
				// Basic auth
				assert.True(t, strings.HasPrefix(authHeader, "Basic "))
			}

			// Return valid token response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"access_token": validTokenResponse.AccessToken,
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        "iam:read iam:write",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read", "iam:write"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()
		token, err := client.GetToken(ctx)
		require.NoError(t, err, "GetToken should succeed")
		assert.Equal(t, validTokenResponse.AccessToken, token.AccessToken)
		assert.Equal(t, "Bearer", token.TokenType)
		assert.True(t, token.Valid(), "Token should be valid")
	})

	t.Run("invalid_credentials", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			response := map[string]interface{}{
				"error":             "invalid_client",
				"error_description": "Client authentication failed",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "invalid-client",
			ClientSecret: "invalid-secret",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()
		_, err = client.GetToken(ctx)
		require.Error(t, err, "GetToken should fail with invalid credentials")

		var authErr *AuthError
		assert.ErrorAs(t, err, &authErr)
		assert.Equal(t, AuthErrorCredentials, authErr.Type)
		assert.Contains(t, authErr.Message, "Client authentication failed")
		assert.False(t, authErr.Retryable)
	})

	t.Run("invalid_scope", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)

			response := map[string]interface{}{
				"error":             "invalid_scope",
				"error_description": "The requested scope is invalid, unknown, or malformed",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"invalid:scope"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()
		_, err = client.GetToken(ctx)
		require.Error(t, err, "GetToken should fail with invalid scope")

		var authErr *AuthError
		assert.ErrorAs(t, err, &authErr)
		assert.Equal(t, AuthErrorCredentials, authErr.Type)
		assert.Contains(t, authErr.Message, "invalid scope")
	})

	t.Run("rate_limiting", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)

			response := map[string]interface{}{
				"error":             "rate_limited",
				"error_description": "Too many token requests",
				"retry_after":       60,
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()
		_, err = client.GetToken(ctx)
		require.Error(t, err, "GetToken should fail with rate limiting")

		var authErr *AuthError
		assert.ErrorAs(t, err, &authErr)
		assert.Equal(t, AuthErrorRateLimit, authErr.Type)
		assert.True(t, authErr.Retryable)
		assert.Equal(t, 60*time.Second, authErr.RetryAfter)
	})

	t.Run("server_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			response := map[string]interface{}{
				"error":             "server_error",
				"error_description": "The authorization server encountered an unexpected condition",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()
		_, err = client.GetToken(ctx)
		require.Error(t, err, "GetToken should fail with server error")

		var authErr *AuthError
		assert.ErrorAs(t, err, &authErr)
		assert.Equal(t, AuthErrorServerError, authErr.Type)
		assert.True(t, authErr.Retryable)
	})

	t.Run("network_timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate slow response
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      10 * time.Millisecond, // Very short timeout
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()
		_, err = client.GetToken(ctx)
		require.Error(t, err, "GetToken should fail with timeout")

		var authErr *AuthError
		assert.ErrorAs(t, err, &authErr)
		assert.Equal(t, AuthErrorNetwork, authErr.Type)
		assert.True(t, authErr.Retryable)
	})
}

// TestAuthClient_TokenRefresh tests token refresh and expiration handling
func TestAuthClient_TokenRefresh(t *testing.T) {
	t.Run("automatic_token_refresh", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++

			var expiresIn int
			if callCount == 1 {
				// First call - return token that expires quickly
				expiresIn = 1
			} else {
				// Subsequent calls - return longer-lived token
				expiresIn = 3600
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"access_token": fmt.Sprintf("token-call-%d", callCount),
				"token_type":   "Bearer",
				"expires_in":   expiresIn,
				"scope":        "iam:read iam:write",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read", "iam:write"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()

		// First token acquisition
		token1, err := client.GetToken(ctx)
		require.NoError(t, err, "First GetToken should succeed")
		assert.Equal(t, "token-call-1", token1.AccessToken)

		// Wait for token to expire
		time.Sleep(2 * time.Second)

		// Second token acquisition should trigger refresh
		token2, err := client.GetToken(ctx)
		require.NoError(t, err, "Second GetToken should succeed")
		assert.Equal(t, "token-call-2", token2.AccessToken)
		assert.NotEqual(t, token1.AccessToken, token2.AccessToken)
		assert.Equal(t, 2, callCount, "Should have made 2 token requests")
	})

	t.Run("token_validation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"access_token": "valid-token-123",
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        "iam:read",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()

		// Get token
		token, err := client.GetToken(ctx)
		require.NoError(t, err, "GetToken should succeed")

		// Validate token
		valid, err := client.ValidateToken(ctx, token)
		require.NoError(t, err, "ValidateToken should succeed")
		assert.True(t, valid, "Token should be valid")

		// Test with expired token
		expiredToken := &oauth2.Token{
			AccessToken: "expired-token",
			TokenType:   "Bearer",
			Expiry:      time.Now().Add(-1 * time.Hour), // Expired
		}

		valid, err = client.ValidateToken(ctx, expiredToken)
		require.NoError(t, err, "ValidateToken should succeed")
		assert.False(t, valid, "Expired token should be invalid")
	})
}

// TestAuthClient_ConcurrentAccess tests thread-safe token management
func TestAuthClient_ConcurrentAccess(t *testing.T) {
	t.Run("concurrent_token_acquisition", func(t *testing.T) {
		callCount := 0
		var mutex sync.Mutex

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mutex.Lock()
			callCount++
			currentCall := callCount
			mutex.Unlock()

			// Simulate some processing time
			time.Sleep(10 * time.Millisecond)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"access_token": fmt.Sprintf("concurrent-token-%d", currentCall),
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        "iam:read",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()
		numGoroutines := 5
		var wg sync.WaitGroup
		tokens := make([]*oauth2.Token, numGoroutines)

		// Concurrent token acquisition
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				token, err := client.GetToken(ctx)
				require.NoError(t, err, "Concurrent GetToken should succeed")
				tokens[index] = token
			}(i)
		}

		wg.Wait()

		// Verify all goroutines got tokens
		for i, token := range tokens {
			assert.NotNil(t, token, "Token %d should not be nil", i)
			assert.True(t, token.Valid(), "Token %d should be valid", i)
		}

		// With proper caching, should have made only one actual HTTP request
		mutex.Lock()
		assert.Equal(t, 1, callCount, "Should have made only one token request due to caching")
		mutex.Unlock()
	})

	t.Run("concurrent_token_refresh", func(t *testing.T) {
		callCount := 0
		var mutex sync.Mutex

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mutex.Lock()
			callCount++
			currentCall := callCount
			mutex.Unlock()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"access_token": fmt.Sprintf("refresh-token-%d", currentCall),
				"token_type":   "Bearer",
				"expires_in":   1, // Short expiration to force refresh
				"scope":        "iam:read",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     server.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()

		// Get initial token
		_, err = client.GetToken(ctx)
		require.NoError(t, err, "Initial GetToken should succeed")

		// Wait for token to expire
		time.Sleep(2 * time.Second)

		numGoroutines := 3
		var wg sync.WaitGroup
		tokens := make([]*oauth2.Token, numGoroutines)

		// Concurrent token refresh
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				token, err := client.GetToken(ctx)
				require.NoError(t, err, "Concurrent GetToken after expiry should succeed")
				tokens[index] = token
			}(i)
		}

		wg.Wait()

		// Verify all goroutines got valid tokens
		for i, token := range tokens {
			assert.NotNil(t, token, "Token %d should not be nil", i)
			assert.True(t, token.Valid(), "Token %d should be valid", i)
		}

		// Should have made exactly 2 requests (initial + one refresh)
		mutex.Lock()
		assert.Equal(t, 2, callCount, "Should have made exactly 2 token requests (initial + refresh)")
		mutex.Unlock()
	})
}

// TestAuthClient_HTTPClientIntegration tests HTTP client integration with tokens
func TestAuthClient_HTTPClientIntegration(t *testing.T) {
	t.Run("authenticated_http_client", func(t *testing.T) {
		// OAuth2 token server
		tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"access_token": "test-access-token-123",
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        "iam:read",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer tokenServer.Close()

		// API server that requires authentication
		apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			assert.Equal(t, "Bearer test-access-token-123", authHeader)
			assert.Equal(t, "test-tenant-123", r.Header.Get("X-Tenant-ID"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"message": "Authenticated API call successful",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer apiServer.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     tokenServer.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()

		// Get authenticated HTTP client
		httpClient, err := client.HTTPClient(ctx)
		require.NoError(t, err, "HTTPClient should succeed")

		// Make API call with authenticated client
		resp, err := httpClient.Get(apiServer.URL + "/api/test")
		require.NoError(t, err, "API call should succeed")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("token_refresh_on_401", func(t *testing.T) {
		tokenCallCount := 0
		tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCallCount++

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			response := map[string]interface{}{
				"access_token": fmt.Sprintf("refreshed-token-%d", tokenCallCount),
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        "iam:read",
			}
			json.NewEncoder(w).Encode(response)
		}))
		defer tokenServer.Close()

		apiCallCount := 0
		apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiCallCount++

			if apiCallCount == 1 {
				// First call - return 401 to simulate expired token
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)

				response := map[string]interface{}{
					"error":             "invalid_token",
					"error_description": "The access token expired",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			// Subsequent calls - check for new token and succeed
			authHeader := r.Header.Get("Authorization")
			assert.Equal(t, "Bearer refreshed-token-2", authHeader)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "Success after token refresh"}`))
		}))
		defer apiServer.Close()

		config := &AuthClientConfig{
			TenantID:     "test-tenant-123",
			ClientID:     "test-client-123",
			ClientSecret: "test-secret-456",
			TokenURL:     tokenServer.URL + "/oauth2/token",
			Scopes:       []string{"iam:read"},
			Timeout:      30 * time.Second,
			MaxRetries:   3,
		}

		client, err := NewAuthClient(config)
		require.NoError(t, err, "NewAuthClient should succeed")

		ctx := context.Background()

		// Get authenticated HTTP client with retry capability
		httpClient, err := client.HTTPClientWithRetry(ctx)
		require.NoError(t, err, "HTTPClientWithRetry should succeed")

		// Make API call - should automatically retry with new token after 401
		resp, err := httpClient.Get(apiServer.URL + "/api/test")
		require.NoError(t, err, "API call should succeed after token refresh")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 2, tokenCallCount, "Should have made 2 token requests")
		assert.Equal(t, 2, apiCallCount, "Should have made 2 API calls")
	})
}

// Mock structures and helper functions for testing

type AuthClientConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	BaseURL      string
	TokenURL     string
	Scopes       []string
	Timeout      time.Duration
	MaxRetries   int
}

type AuthClient struct {
	config      *AuthClientConfig
	oauth2cfg   *clientcredentials.Config
	tokenSource oauth2.TokenSource
	httpClient  *http.Client
	mutex       sync.RWMutex
}

type AuthErrorType int

const (
	AuthErrorUnknown AuthErrorType = iota
	AuthErrorConfiguration
	AuthErrorDiscovery
	AuthErrorCredentials
	AuthErrorNetwork
	AuthErrorServerError
	AuthErrorRateLimit
	AuthErrorTokenExpired
)

type AuthError struct {
	Type       AuthErrorType
	Message    string
	Underlying error
	Retryable  bool
	RetryAfter time.Duration
}

func (e *AuthError) Error() string {
	return e.Message
}

// Mock implementations for testing (these will be replaced by actual implementations)

func NewAuthClient(config *AuthClientConfig) (*AuthClient, error) {
	// This is a mock implementation for testing
	return &AuthClient{config: config}, nil
}

func (c *AuthClient) GetToken(ctx context.Context) (*oauth2.Token, error) {
	// This is a mock implementation for testing
	return nil, fmt.Errorf("not implemented - this is a test mock")
}

func (c *AuthClient) ValidateToken(ctx context.Context, token *oauth2.Token) (bool, error) {
	// Simple validation for testing
	return token != nil && token.Valid(), nil
}

func (c *AuthClient) HTTPClient(ctx context.Context) (*http.Client, error) {
	// This is a mock implementation for testing
	return &http.Client{}, nil
}

func (c *AuthClient) HTTPClientWithRetry(ctx context.Context) (*http.Client, error) {
	// This is a mock implementation for testing
	return &http.Client{}, nil
}

// Helper function to parse form data from request
func parseForm(t *testing.T, r *http.Request) url.Values {
	err := r.ParseForm()
	require.NoError(t, err, "Should be able to parse form data")
	return r.Form
}
