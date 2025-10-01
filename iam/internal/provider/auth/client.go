package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// AuthClientConfig contains configuration for OAuth2 authentication client
type AuthClientConfig struct {
	// Required authentication parameters
	TenantID     string
	ClientID     string
	ClientSecret string

	// Optional configuration
	BaseURL  string
	TokenURL string
	Scopes   []string

	// Timeout and retry configuration
	Timeout    time.Duration
	MaxRetries int

	// Advanced configuration
	DisableDiscovery bool
	CustomHeaders    map[string]string
}

// AuthClient manages OAuth2 authentication and token lifecycle
type AuthClient struct {
	// Configuration
	config       *AuthClientConfig
	oauth2Config *clientcredentials.Config
	tokenSource  oauth2.TokenSource
	httpClient   *http.Client

	// Discovery integration
	discoveryClient *DiscoveryClient

	// Token management
	tokenCache *TokenCache

	// Retry configuration
	retryConfig *RetryConfig

	// Thread safety
	mutex sync.RWMutex
}

// TokenCache manages token caching and validation
type TokenCache struct {
	token        *oauth2.Token
	hash         string
	createdAt    time.Time
	lastUsed     time.Time
	refreshCount int
	mutex        sync.RWMutex
}

// NewAuthClient creates a new OAuth2 authentication client
func NewAuthClient(config *AuthClientConfig) (*AuthClient, error) {
	if err := validateAuthConfig(config); err != nil {
		return nil, NewConfigurationError("invalid auth client configuration", err)
	}

	client := &AuthClient{
		config:      config,
		tokenCache:  &TokenCache{},
		retryConfig: DefaultRetryConfig(),
	}

	// Initialize discovery client if not disabled
	if !config.DisableDiscovery && config.BaseURL != "" {
		client.discoveryClient = NewDiscoveryClient(config.BaseURL, config.Timeout)
	}

	// Set up HTTP client with timeout
	client.httpClient = &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: config.Timeout,
		},
	}

	// Initialize OAuth2 configuration
	if err := client.initializeOAuth2Config(); err != nil {
		return nil, err
	}

	return client, nil
}

// initializeOAuth2Config sets up the OAuth2 client credentials configuration
func (c *AuthClient) initializeOAuth2Config() error {
	tokenURL := c.config.TokenURL

	// Use discovery to get token endpoint if not explicitly configured
	if tokenURL == "" && c.discoveryClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
		defer cancel()

		discoveredURL, err := c.discoveryClient.GetTokenEndpoint(ctx, "")
		if err != nil {
			return NewDiscoveryError("failed to discover token endpoint", err)
		}
		tokenURL = discoveredURL
	}

	if tokenURL == "" {
		return NewConfigurationError("no token URL available - discovery failed and no explicit token_url configured", nil)
	}

	// Validate scopes if discovery is available
	if c.discoveryClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
		defer cancel()

		if err := c.discoveryClient.ValidateScopes(ctx, c.config.Scopes); err != nil {
			// Log warning but don't fail - scope validation is optional
		}
	}

	// Create OAuth2 client credentials configuration
	c.oauth2Config = &clientcredentials.Config{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		TokenURL:     tokenURL,
		Scopes:       c.config.Scopes,
		AuthStyle:    oauth2.AuthStyleInParams, // Default to form parameters
	}

	// Create token source with custom HTTP client
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c.httpClient)
	c.tokenSource = c.oauth2Config.TokenSource(ctx)

	return nil
}

// GetToken acquires or returns a cached OAuth2 access token
func (c *AuthClient) GetToken(ctx context.Context) (*oauth2.Token, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if we have a valid cached token
	if token := c.tokenCache.getValidToken(); token != nil {
		c.tokenCache.updateLastUsed()
		return token, nil
	}

	// Acquire new token with retry logic
	token, err := c.acquireTokenWithRetry(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the new token
	c.tokenCache.setToken(token)

	return token, nil
}

// acquireTokenWithRetry attempts to acquire a token with exponential backoff retry
func (c *AuthClient) acquireTokenWithRetry(ctx context.Context) (*oauth2.Token, error) {
	var lastErr error

	for attempt := 0; attempt < c.retryConfig.MaxAttempts; attempt++ {
		token, err := c.acquireToken(ctx)
		if err == nil {
			return token, nil
		}

		lastErr = err

		// Check if we should retry
		if !c.retryConfig.ShouldRetry(err, attempt) {
			break
		}

		// Calculate delay and wait
		delay := c.retryConfig.GetDelay(attempt, err)
		select {
		case <-ctx.Done():
			return nil, NewNetworkError("context canceled during token acquisition retry", ctx.Err())
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return nil, lastErr
}

// acquireToken performs a single token acquisition attempt
func (c *AuthClient) acquireToken(ctx context.Context) (*oauth2.Token, error) {
	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, c.mapOAuth2Error(err)
	}

	if !token.Valid() {
		return nil, NewCredentialsError("received invalid token from OAuth2 server", nil)
	}

	return token, nil
}

// ValidateToken checks if a token is valid and not expired
func (c *AuthClient) ValidateToken(ctx context.Context, token *oauth2.Token) (bool, error) {
	if token == nil {
		return false, nil
	}

	// Basic validation
	if !token.Valid() {
		return false, nil
	}

	// Optional: Perform introspection if supported
	// This would require additional OCMS endpoint support

	return true, nil
}

// HTTPClient returns an HTTP client that automatically includes OAuth2 authentication
func (c *AuthClient) HTTPClient(ctx context.Context) (*http.Client, error) {
	// Get current token
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication token: %w", err)
	}

	// Create transport that adds authentication headers
	transport := &AuthenticatedTransport{
		Base:     c.httpClient.Transport,
		Token:    token,
		TenantID: c.config.TenantID,
		Headers:  c.config.CustomHeaders,
	}

	// Return client with authenticated transport
	return &http.Client{
		Transport: transport,
		Timeout:   c.config.Timeout,
	}, nil
}

// HTTPClientWithRetry returns an HTTP client with automatic token refresh on authentication errors
func (c *AuthClient) HTTPClientWithRetry(ctx context.Context) (*http.Client, error) {
	// Get current token
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication token: %w", err)
	}

	// Create transport with retry capability
	transport := &RetryTransport{
		Base:       c.httpClient.Transport,
		AuthClient: c,
		Token:      token,
		TenantID:   c.config.TenantID,
		Headers:    c.config.CustomHeaders,
		MaxRetries: c.config.MaxRetries,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   c.config.Timeout,
	}, nil
}

// RefreshToken forces a token refresh
func (c *AuthClient) RefreshToken(ctx context.Context) (*oauth2.Token, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Clear cached token to force refresh
	c.tokenCache.clearToken()

	// Acquire new token
	token, err := c.acquireTokenWithRetry(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the new token
	c.tokenCache.setToken(token)

	return token, nil
}

// Close cleans up resources and clears sensitive data
func (c *AuthClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Clear cached token
	c.tokenCache.clearToken()

	// Clear sensitive configuration
	c.config.ClientSecret = ""

	return nil
}

// TokenCache methods

// getValidToken returns a cached token if it's valid, nil otherwise
func (tc *TokenCache) getValidToken() *oauth2.Token {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	if tc.token == nil {
		return nil
	}

	// Check token validity
	if !tc.token.Valid() {
		return nil
	}

	// Validate token hash to detect tampering
	if tc.hash != "" && computeTokenHash(tc.token) != tc.hash {
		return nil
	}

	return tc.token
}

// setToken stores a token in the cache with security measures
func (tc *TokenCache) setToken(token *oauth2.Token) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	tc.token = token
	tc.hash = computeTokenHash(token)
	tc.createdAt = time.Now()
	tc.refreshCount++
}

// updateLastUsed updates the last used timestamp
func (tc *TokenCache) updateLastUsed() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	tc.lastUsed = time.Now()
}

// clearToken removes the cached token and clears sensitive data
func (tc *TokenCache) clearToken() {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	tc.token = nil
	tc.hash = ""
}

// AuthenticatedTransport adds OAuth2 authentication headers to HTTP requests
type AuthenticatedTransport struct {
	Base     http.RoundTripper
	Token    *oauth2.Token
	TenantID string
	Headers  map[string]string
}

// RoundTrip implements the http.RoundTripper interface
func (t *AuthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	newReq := req.Clone(req.Context())

	// Add OAuth2 authorization header
	if t.Token != nil && t.Token.AccessToken != "" {
		newReq.Header.Set("Authorization", "Bearer "+t.Token.AccessToken)
	}

	// Add tenant ID header
	if t.TenantID != "" {
		newReq.Header.Set("X-Tenant-ID", t.TenantID)
	}

	// Add custom headers
	for key, value := range t.Headers {
		newReq.Header.Set(key, value)
	}

	// Use base transport or default
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	return base.RoundTrip(newReq)
}

// RetryTransport provides automatic token refresh on authentication errors
type RetryTransport struct {
	Base       http.RoundTripper
	AuthClient *AuthClient
	Token      *oauth2.Token
	TenantID   string
	Headers    map[string]string
	MaxRetries int
}

// RoundTrip implements the http.RoundTripper interface with retry logic
func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for attempt := 0; attempt <= t.MaxRetries; attempt++ {
		// Clone the request
		newReq := req.Clone(req.Context())

		// Add authentication headers
		if t.Token != nil && t.Token.AccessToken != "" {
			newReq.Header.Set("Authorization", "Bearer "+t.Token.AccessToken)
		}

		if t.TenantID != "" {
			newReq.Header.Set("X-Tenant-ID", t.TenantID)
		}

		for key, value := range t.Headers {
			newReq.Header.Set(key, value)
		}

		// Execute request
		base := t.Base
		if base == nil {
			base = http.DefaultTransport
		}

		resp, err := base.RoundTrip(newReq)
		if err != nil {
			return nil, err
		}

		// Check for authentication errors
		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()

			// Parse error response
			if attempt < t.MaxRetries && t.isTokenExpiredError(resp) {
				// Refresh token and retry
				newToken, refreshErr := t.AuthClient.RefreshToken(req.Context())
				if refreshErr != nil {
					return nil, NewTokenExpiredError("failed to refresh expired token")
				}

				t.Token = newToken
				continue
			}
		}

		return resp, nil
	}

	return nil, NewCredentialsError("authentication failed after maximum retries", nil)
}

// isTokenExpiredError checks if the response indicates an expired token
func (t *RetryTransport) isTokenExpiredError(resp *http.Response) bool {
	// Check WWW-Authenticate header for token expiration
	wwwAuth := resp.Header.Get("WWW-Authenticate")
	if wwwAuth != "" {
		return contains([]string{wwwAuth}, "invalid_token") || contains([]string{wwwAuth}, "token_expired")
	}

	// Could also parse response body for OAuth2 error format
	return true // Assume 401 is token expiration for retry purposes
}

// Helper functions

// validateAuthConfig validates the authentication client configuration
func validateAuthConfig(config *AuthClientConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if config.TenantID == "" {
		return NewConfigValidationError("tenant_id", "required", "provide a valid tenant identifier", "")
	}

	if config.ClientID == "" {
		return NewConfigValidationError("client_id", "required", "provide a valid OAuth2 client ID", "")
	}

	if config.ClientSecret == "" {
		return NewConfigValidationError("client_secret", "required", "provide a valid OAuth2 client secret", "[REDACTED]")
	}

	if len(config.ClientSecret) < 8 {
		return NewConfigValidationError("client_secret", "minimum length 8 characters", "use a stronger client secret", "[REDACTED]")
	}

	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second // Default timeout
	}

	if config.MaxRetries < 0 {
		config.MaxRetries = 3 // Default max retries
	}

	if len(config.Scopes) == 0 {
		config.Scopes = []string{"iam:read", "iam:write"} // Default scopes
	}

	return nil
}

// mapOAuth2Error converts golang.org/x/oauth2 errors to AuthError types
func (c *AuthClient) mapOAuth2Error(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Check for specific OAuth2 error patterns
	if contains([]string{errStr}, "invalid_client") {
		return NewCredentialsError("OAuth2 client authentication failed - check client_id and client_secret", err)
	}

	if contains([]string{errStr}, "invalid_grant") {
		return NewCredentialsError("OAuth2 grant type not supported or invalid", err)
	}

	if contains([]string{errStr}, "invalid_scope") {
		return NewCredentialsError("requested OAuth2 scopes are invalid or not authorized", err)
	}

	if contains([]string{errStr}, "temporarily_unavailable") || contains([]string{errStr}, "server_error") {
		return NewServerError("OAuth2 server temporarily unavailable", err)
	}

	if contains([]string{errStr}, "timeout") || contains([]string{errStr}, "connection") {
		return NewNetworkError("network error during OAuth2 token acquisition", err)
	}

	// Default to credentials error for unknown OAuth2 errors
	return NewCredentialsError("OAuth2 authentication failed", err)
}

// computeTokenHash creates a hash of the token for integrity validation
func computeTokenHash(token *oauth2.Token) string {
	if token == nil {
		return ""
	}

	// Simple hash based on token content
	// In production, you might want to use a proper cryptographic hash
	data := fmt.Sprintf("%s:%s:%v", token.AccessToken[:min(len(token.AccessToken), 10)], token.TokenType, token.Expiry.Unix())
	hash := 0
	for _, char := range data {
		hash = hash*31 + int(char)
	}
	return fmt.Sprintf("%x", hash)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
