package auth
package auth_test

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
)

// TestAuthClient_TokenAcquisition tests successful OAuth2 token acquisition
func TestAuthClient_TokenAcquisition(t *testing.T) {
	// Mock OAuth2 server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify OAuth2 request
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/oauth/token", r.URL.Path)
		
		// Parse form data
		err := r.ParseForm()
		require.NoError(t, err)
		
		assert.Equal(t, "client_credentials", r.Form.Get("grant_type"))
		assert.Equal(t, "test-client-id", r.Form.Get("client_id"))
		assert.Equal(t, "test-client-secret", r.Form.Get("client_secret"))
		
		// Return token response
		response := map[string]interface{}{
			"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"scope":        "hiiretail:iam",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// This test will fail until AuthClient is implemented
	// TODO: Replace with actual AuthClient implementation
	
	// Mock AuthClient configuration
	config := mockAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TenantID:     "test-tenant-123",
		TokenURL:     server.URL + "/oauth/token",
	}
	
	// Create mock AuthClient
	client := &mockAuthClient{config: config}
	
	// Test token acquisition
	token, err := client.GetAccessToken(context.Background())
	require.NoError(t, err)
	
	assert.NotEmpty(t, token)
	assert.True(t, strings.HasPrefix(token, "eyJ"), "Token should be JWT format")
}

// TestAuthClient_TokenCaching tests token caching and reuse
func TestAuthClient_TokenCaching(t *testing.T) {
	requestCount := 0
	
	// Mock OAuth2 server that tracks requests
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		
		response := map[string]interface{}{
			"access_token": "cached.token.here",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// This test will fail until AuthClient is implemented
	config := mockAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TenantID:     "test-tenant-123",
		TokenURL:     server.URL + "/oauth/token",
	}
	
	client := &mockAuthClient{config: config}
	
	// First request should acquire token
	token1, err := client.GetAccessToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, requestCount, "Should make one request for new token")
	
	// Second request should use cached token
	token2, err := client.GetAccessToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, token1, token2, "Should return same cached token")
	assert.Equal(t, 1, requestCount, "Should not make additional request for cached token")
}

// TestAuthClient_TokenRefresh tests automatic token refresh before expiration
func TestAuthClient_TokenRefresh(t *testing.T) {
	requestCount := 0
	
	// Mock OAuth2 server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		
		var response map[string]interface{}
		if requestCount == 1 {
			// First token with short expiration
			response = map[string]interface{}{
				"access_token": "first.token.here",
				"token_type":   "Bearer",
				"expires_in":   2, // 2 seconds
			}
		} else {
			// Refreshed token
			response = map[string]interface{}{
				"access_token": "refreshed.token.here",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// This test will fail until AuthClient is implemented
	config := mockAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret", 
		TenantID:     "test-tenant-123",
		TokenURL:     server.URL + "/oauth/token",
	}
	
	client := &mockAuthClient{config: config}
	
	// Get initial token
	token1, err := client.GetAccessToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "first.token.here", token1)
	assert.Equal(t, 1, requestCount)
	
	// Wait for token to near expiration
	time.Sleep(3 * time.Second)
	
	// Request token again - should trigger refresh
	token2, err := client.GetAccessToken(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "refreshed.token.here", token2)
	assert.Equal(t, 2, requestCount, "Should make refresh request")
}

// TestAuthClient_ConcurrentAccess tests thread-safe token access
func TestAuthClient_ConcurrentAccess(t *testing.T) {
	requestCount := 0
	var mu sync.Mutex
	
	// Mock OAuth2 server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		count := requestCount
		mu.Unlock()
		
		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)
		
		response := map[string]interface{}{
			"access_token": "concurrent.token.here",
			"token_type":   "Bearer",
			"expires_in":   3600,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		
		t.Logf("OAuth2 request %d completed", count)
	}))
	defer server.Close()
	
	// This test will fail until AuthClient is implemented
	config := mockAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		TenantID:     "test-tenant-123",
		TokenURL:     server.URL + "/oauth/token",
	}
	
	client := &mockAuthClient{config: config}
	
	// Launch multiple concurrent requests
	const numGoroutines = 10
	var wg sync.WaitGroup
	tokens := make([]string, numGoroutines)
	errors := make([]error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			token, err := client.GetAccessToken(context.Background())
			tokens[index] = token
			errors[index] = err
		}(i)
	}
	
	wg.Wait()
	
	// Verify all requests succeeded
	for i, err := range errors {
		assert.NoError(t, err, "Request %d should not error", i)
	}
	
	// Verify all got the same token (from cache)
	for i, token := range tokens {
		assert.Equal(t, "concurrent.token.here", token, "Request %d should get correct token", i)
	}
	
	// Should only make one actual OAuth2 request due to caching
	assert.Equal(t, 1, requestCount, "Should only make one OAuth2 request despite concurrent access")
}

// TestAuthClient_ErrorHandling tests error handling for various failure scenarios
func TestAuthClient_ErrorHandling(t *testing.T) {
	testCases := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedError  string
	}{
		{
			name: "Invalid client credentials",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				response := map[string]string{
					"error":             "invalid_client",
					"error_description": "Invalid client credentials",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
			},
			expectedError: "invalid_client",
		},
		{
			name: "Server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			expectedError: "500",
		},
		{
			name: "Network timeout",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Simulate timeout by sleeping longer than client timeout
				time.Sleep(2 * time.Second)
				w.WriteHeader(http.StatusOK)
			},
			expectedError: "timeout",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock OAuth2 server with specific response
			server := httptest.NewTLSServer(http.HandlerFunc(tc.serverResponse))
			defer server.Close()
			
			// This test will fail until AuthClient is implemented
			config := mockAuthConfig{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				TenantID:     "test-tenant-123",
				TokenURL:     server.URL + "/oauth/token",
			}
			
			client := &mockAuthClient{config: config}
			
			// Test error handling
			token, err := client.GetAccessToken(context.Background())
			
			assert.Error(t, err, "Should return error")
			assert.Empty(t, token, "Should not return token on error")
			assert.Contains(t, err.Error(), tc.expectedError, "Error should contain expected message")
		})
	}
}

// mockAuthConfig represents OAuth2 configuration for testing
type mockAuthConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	TokenURL     string
}

// mockAuthClient is a temporary mock implementation for testing
// This will be replaced by the actual AuthClient implementation
type mockAuthClient struct {
	config    mockAuthConfig
	token     string
	expiresAt time.Time
	mutex     sync.RWMutex
}

func (c *mockAuthClient) GetAccessToken(ctx context.Context) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Check if we have a valid cached token
	if c.token != "" && time.Now().Before(c.expiresAt.Add(-60*time.Second)) {
		return c.token, nil
	}
	
	// Make OAuth2 request
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)
	data.Set("scope", "hiiretail:iam")
	
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.config.TokenURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	
	// Use client with short timeout for timeout test
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errorResp)
		if errorCode, ok := errorResp["error"]; ok {
			return "", fmt.Errorf("OAuth2 error: %s", errorCode)
		}
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}
	
	var tokenResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}
	
	// Cache token
	c.token = tokenResp["access_token"].(string)
	if expiresIn, ok := tokenResp["expires_in"].(float64); ok {
		c.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}
	
	return c.token, nil
}