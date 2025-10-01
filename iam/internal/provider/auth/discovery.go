package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// OIDCDiscoveryResponse represents the OpenID Connect discovery response
// as defined in RFC 8414: OAuth 2.0 Authorization Server Metadata
type OIDCDiscoveryResponse struct {
	Issuer                   string   `json:"issuer"`
	TokenEndpoint            string   `json:"token_endpoint"`
	AuthorizationEndpoint    string   `json:"authorization_endpoint"`
	JwksUri                  string   `json:"jwks_uri"`
	GrantTypesSupported      []string `json:"grant_types_supported"`
	TokenEndpointAuthMethods []string `json:"token_endpoint_auth_methods_supported"`
	ResponseTypesSupported   []string `json:"response_types_supported"`
	ScopesSupported          []string `json:"scopes_supported"`
}

// Validate validates the OIDC discovery response according to OAuth2 specifications
func (d *OIDCDiscoveryResponse) Validate() error {
	if d.Issuer == "" {
		return fmt.Errorf("missing required field: issuer")
	}

	if d.TokenEndpoint == "" {
		return fmt.Errorf("missing required field: token_endpoint")
	}

	if !contains(d.GrantTypesSupported, "client_credentials") {
		return fmt.Errorf("client_credentials grant type not supported")
	}

	if !isValidURL(d.TokenEndpoint) {
		return fmt.Errorf("invalid token_endpoint URL")
	}

	return nil
}

// DiscoveryCache provides thread-safe caching for OAuth2 discovery responses
type DiscoveryCache struct {
	cache map[string]*cacheEntry
	mutex sync.RWMutex
}

type cacheEntry struct {
	response  *OIDCDiscoveryResponse
	expiresAt time.Time
}

// Get retrieves a cached discovery response if it exists and hasn't expired
func (c *DiscoveryCache) Get(url string) (*OIDCDiscoveryResponse, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.cache == nil {
		return nil, false
	}

	entry, exists := c.cache[url]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.response, true
}

// Set stores a discovery response in the cache with the specified TTL
func (c *DiscoveryCache) Set(url string, response *OIDCDiscoveryResponse, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cache == nil {
		c.cache = make(map[string]*cacheEntry)
	}

	c.cache[url] = &cacheEntry{
		response:  response,
		expiresAt: time.Now().Add(ttl),
	}
}

// Clear removes all cached entries
func (c *DiscoveryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*cacheEntry)
}

// DiscoveryClient handles OAuth2 endpoint discovery with caching
type DiscoveryClient struct {
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
	cache      *DiscoveryCache
}

const (
	// DiscoveryCacheTTL defines how long discovery responses are cached
	DiscoveryCacheTTL = 1 * time.Hour

	// DiscoveryMaxAge defines the maximum age for cached discovery responses
	DiscoveryMaxAge = 24 * time.Hour

	// DiscoveryEndpointPath is the standard OAuth2 discovery endpoint path
	DiscoveryEndpointPath = "/.well-known/openid-configuration"
)

// NewDiscoveryClient creates a new OAuth2 discovery client
func NewDiscoveryClient(baseURL string, timeout time.Duration) *DiscoveryClient {
	return &DiscoveryClient{
		baseURL: baseURL,
		timeout: timeout,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		cache: &DiscoveryCache{},
	}
}

// FetchDiscovery retrieves OAuth2 discovery information from the OCMS discovery endpoint
func (c *DiscoveryClient) FetchDiscovery(ctx context.Context) (*OIDCDiscoveryResponse, error) {
	discoveryURL := c.buildDiscoveryURL()

	// Check cache first
	if cached, found := c.cache.Get(discoveryURL); found {
		return cached, nil
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return nil, &AuthError{
			Type:       AuthErrorConfiguration,
			Message:    fmt.Sprintf("failed to create discovery request: %v", err),
			Underlying: err,
			Retryable:  false,
		}
	}

	// Set appropriate headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "HiiRetail-IAM-Terraform-Provider/1.0")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &AuthError{
			Type:       AuthErrorNetwork,
			Message:    fmt.Sprintf("failed to fetch discovery endpoint: %v", err),
			Underlying: err,
			Retryable:  true,
		}
	}
	defer resp.Body.Close()

	// Handle HTTP errors
	switch resp.StatusCode {
	case http.StatusOK:
		// Continue processing
	case http.StatusNotFound:
		return nil, &AuthError{
			Type:       AuthErrorDiscovery,
			Message:    "discovery endpoint not available - OCMS service may not support OAuth2 discovery",
			Underlying: fmt.Errorf("discovery endpoint returned 404"),
			Retryable:  false,
		}
	case http.StatusTooManyRequests:
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		return nil, &AuthError{
			Type:       AuthErrorRateLimit,
			Message:    "discovery endpoint rate limited",
			Underlying: fmt.Errorf("discovery rate limited, retry after %v", retryAfter),
			Retryable:  true,
			RetryAfter: retryAfter,
		}
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return nil, &AuthError{
			Type:       AuthErrorServerError,
			Message:    "OCMS service unavailable",
			Underlying: fmt.Errorf("discovery endpoint returned status %d", resp.StatusCode),
			Retryable:  true,
		}
	default:
		return nil, &AuthError{
			Type:       AuthErrorDiscovery,
			Message:    fmt.Sprintf("discovery endpoint returned unexpected status: %d", resp.StatusCode),
			Underlying: fmt.Errorf("unexpected HTTP status %d", resp.StatusCode),
			Retryable:  false,
		}
	}

	// Parse JSON response
	var discoveryResponse OIDCDiscoveryResponse
	if err := json.NewDecoder(resp.Body).Decode(&discoveryResponse); err != nil {
		return nil, &AuthError{
			Type:       AuthErrorDiscovery,
			Message:    "failed to parse discovery response - invalid JSON format",
			Underlying: fmt.Errorf("JSON decode error: %w", err),
			Retryable:  false,
		}
	}

	// Validate response
	if err := discoveryResponse.Validate(); err != nil {
		return nil, &AuthError{
			Type:       AuthErrorDiscovery,
			Message:    fmt.Sprintf("invalid discovery response: %v", err),
			Underlying: err,
			Retryable:  false,
		}
	}

	// Cache the response
	c.cache.Set(discoveryURL, &discoveryResponse, DiscoveryCacheTTL)

	return &discoveryResponse, nil
}

// GetTokenEndpoint retrieves the OAuth2 token endpoint URL with fallback support
func (c *DiscoveryClient) GetTokenEndpoint(ctx context.Context, fallbackURL string) (string, error) {
	// Try discovery first
	discovery, err := c.FetchDiscovery(ctx)
	if err == nil {
		return discovery.TokenEndpoint, nil
	}

	// If discovery fails and we have a fallback URL, use it
	if fallbackURL != "" {
		if !isValidURL(fallbackURL) {
			return "", &AuthError{
				Type:       AuthErrorConfiguration,
				Message:    "fallback token URL is not a valid URL",
				Underlying: fmt.Errorf("invalid fallback URL: %s", fallbackURL),
				Retryable:  false,
			}
		}
		return fallbackURL, nil
	}

	// No fallback available
	return "", &AuthError{
		Type:       AuthErrorDiscovery,
		Message:    "discovery failed and no fallback token URL provided",
		Underlying: err,
		Retryable:  false,
	}
}

// GetSupportedScopes retrieves the list of supported OAuth2 scopes
func (c *DiscoveryClient) GetSupportedScopes(ctx context.Context) ([]string, error) {
	discovery, err := c.FetchDiscovery(ctx)
	if err != nil {
		return nil, err
	}

	return discovery.ScopesSupported, nil
}

// ValidateScopes checks if the requested scopes are supported by the OCMS
func (c *DiscoveryClient) ValidateScopes(ctx context.Context, requestedScopes []string) error {
	supportedScopes, err := c.GetSupportedScopes(ctx)
	if err != nil {
		// If we can't get supported scopes, skip validation
		return nil
	}

	for _, requested := range requestedScopes {
		if !contains(supportedScopes, requested) {
			return &AuthError{
				Type:       AuthErrorConfiguration,
				Message:    fmt.Sprintf("scope '%s' is not supported by OCMS", requested),
				Underlying: fmt.Errorf("unsupported scope: %s", requested),
				Retryable:  false,
			}
		}
	}

	return nil
}

// ClearCache removes all cached discovery responses
func (c *DiscoveryClient) ClearCache() {
	c.cache.Clear()
}

// buildDiscoveryURL constructs the full discovery endpoint URL
func (c *DiscoveryClient) buildDiscoveryURL() string {
	baseURL := c.baseURL
	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return baseURL + DiscoveryEndpointPath
}

// Helper functions

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isValidURL validates if a string is a valid HTTP/HTTPS URL
func isValidURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

// parseRetryAfter parses the Retry-After header value
func parseRetryAfter(retryAfter string) time.Duration {
	if retryAfter == "" {
		return 60 * time.Second // Default retry after 60 seconds
	}

	// Try to parse as seconds
	if duration, err := time.ParseDuration(retryAfter + "s"); err == nil {
		return duration
	}

	// Default fallback
	return 60 * time.Second
}
