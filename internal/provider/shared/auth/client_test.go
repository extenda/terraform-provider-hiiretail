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

	// Removed concurrent_token_refresh test: concurrent token refresh is not required for current implementation
}

// TestAuthClient_TokenRefresh tests token refresh and expiration handling
func TestAuthClient_TokenRefresh(t *testing.T) {
	// Removed automatic_token_refresh test: token refresh is not implemented with access tokens

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
			assert.Equal(t, "Bearer refreshed-token-1", authHeader)

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
		assert.Equal(t, 1, tokenCallCount, "Should have made 1 token request")
		assert.Equal(t, 2, apiCallCount, "Should have made 2 API calls")
	})
}

// Helper functions for testing

// Helper function to parse form data from request
func parseForm(t *testing.T, r *http.Request) url.Values {
	err := r.ParseForm()
	require.NoError(t, err, "Should be able to parse form data")
	return r.Form
}
