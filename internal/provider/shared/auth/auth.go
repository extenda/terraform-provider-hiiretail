// Package auth provides OAuth2 authentication for the HiiRetail IAM Terraform provider.
//
// This package implements OAuth2 client credentials flow for authenticating with
// HiiRetail IAM APIs, including automatic token management, caching, and retry logic.
//
// # Basic Usage
//
//	// Create configuration
//	config := &auth.Config{
//		ClientID:     "your-client-id",
//		ClientSecret: "your-client-secret",
//		TenantID:     "your-tenant-id",
//	}
//
//	// Create authenticated client
//	client, err := auth.NewClient(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Get authenticated HTTP client
//	httpClient, err := client.HTTPClient(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Use httpClient for API requests
//	resp, err := httpClient.Get("https://iam-api.retailsvc.com/roles")
//
// # Environment Variables
//
// The auth package supports configuration via environment variables:
//
//	HIIRETAIL_CLIENT_ID      - OAuth2 client ID
//	HIIRETAIL_CLIENT_SECRET  - OAuth2 client secret
//	HIIRETAIL_TENANT_ID      - HiiRetail tenant identifier
//	HIIRETAIL_ENVIRONMENT    - Target environment (production, test, dev)
//	HIIRETAIL_AUTH_URL       - Custom OAuth2 token endpoint URL
//	HIIRETAIL_API_URL        - Custom IAM API base URL
//
// # Endpoint Resolution
//
// The package automatically resolves OAuth2 and API endpoints based on the tenant ID
// and environment:
//
//   - Production tenants → https://auth.retailsvc.com, https://iam-api.retailsvc.com
//   - Test/dev tenants  → https://auth.retailsvc-test.com, https://iam-api.retailsvc-test.com
//
// Custom endpoints can be provided via configuration or environment variables.
//
// # Error Handling
//
// All functions return typed errors that can be inspected for retry logic:
//
//	client, err := auth.NewClient(config)
//	if err != nil {
//		if authErr, ok := err.(*auth.AuthError); ok {
//			switch authErr.Type {
//			case auth.AuthErrorCredentials:
//				// Handle invalid credentials
//			case auth.AuthErrorNetwork:
//				// Handle network errors
//			}
//		}
//	}
//
// # Security
//
// The package implements security best practices:
//
//   - Automatic credential redaction in logs and error messages
//   - Secure token caching with integrity validation
//   - HTTPS enforcement for all OAuth2 and API endpoints
//   - Automatic token refresh and retry logic (rate limiting logic removed)
package auth

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// Version of the auth package
const Version = "1.0.0"

// Config holds OAuth2 authentication configuration
type Config struct {
	// Required authentication parameters
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TenantID     string `json:"tenant_id"`

	// Optional configuration
	Environment string `json:"environment,omitempty"`
	AuthURL     string `json:"auth_url,omitempty"`
	APIURL      string `json:"api_url,omitempty"`

	// OAuth2 settings
	Scopes  []string      `json:"scopes,omitempty"`
	Timeout time.Duration `json:"timeout,omitempty"`

	// Advanced options
	MaxRetries       int  `json:"max_retries,omitempty"`
	DisableDiscovery bool `json:"disable_discovery,omitempty"`
	SkipTLS          bool `json:"skip_tls,omitempty"` // For testing only
}

// Client provides OAuth2 authentication for HiiRetail IAM APIs
type Client interface {
	// GetToken returns a valid OAuth2 access token, refreshing if necessary
	GetToken(ctx context.Context) (*oauth2.Token, error)

	// HTTPClient returns an HTTP client with automatic OAuth2 authentication
	HTTPClient(ctx context.Context) (*http.Client, error)

	// HTTPClientWithRetry returns an HTTP client with automatic token refresh on auth errors
	HTTPClientWithRetry(ctx context.Context) (*http.Client, error)

	// RefreshToken forces a token refresh
	RefreshToken(ctx context.Context) (*oauth2.Token, error)

	// ValidateToken checks if a token is valid
	ValidateToken(ctx context.Context, token *oauth2.Token) (bool, error)

	// Close cleans up resources and clears sensitive data
	Close() error
}

// TokenProvider defines the interface for OAuth2 token acquisition
type TokenProvider interface {
	// Token returns a valid OAuth2 token
	Token(ctx context.Context) (*oauth2.Token, error)
}

// ConfigValidator validates OAuth2 configuration
type ConfigValidator interface {
	// Validate checks configuration validity
	Validate(config *Config) error
}

// Public API functions

// New creates a new OAuth2 authentication client with the provided configuration.
//
// The client automatically resolves endpoints, validates configuration, and sets up
// token caching and retry logic.
//
// Example:
//
//	config := &auth.Config{
//		ClientID:     "your-client-id",
//		ClientSecret: "your-client-secret",
//		TenantID:     "your-tenant-id",
//	}
//
//	client, err := auth.New(config)
//	if err != nil {
//		return err
//	}
//	defer client.Close()
func New(config *Config) (Client, error) {
	return NewClient(config)
}

// NewFromEnvironment creates a new client using environment variables for configuration.
//
// This is a convenience function that loads configuration from standard environment
// variables and creates a client.
//
// Example:
//
//	// Set environment variables:
//	// export HIIRETAIL_CLIENT_ID="your-client-id"
//	// export HIIRETAIL_CLIENT_SECRET="your-client-secret"
//	// export HIIRETAIL_TENANT_ID="your-tenant-id"
//
//	client, err := auth.NewFromEnvironment()
//	if err != nil {
//		return err
//	}
//	defer client.Close()
func NewFromEnvironment() (Client, error) {
	config := NewOAuth2Configuration()
	config.LoadFromEnvironment()

	// Convert OAuth2Configuration to Config for compatibility
	clientConfig := &Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TenantID:     config.TenantID,
		Environment:  config.Environment,
		AuthURL:      config.AuthURL,
		APIURL:       config.APIURL,
		Scopes:       config.Scopes,
		Timeout:      time.Duration(config.Timeout) * time.Second,
		MaxRetries:   config.MaxRetries,
		SkipTLS:      config.SkipTLS,
	}

	return NewClient(clientConfig)
}

// NewClient creates a new OAuth2 authentication client.
//
// This is the primary constructor that converts the public Config to the internal
// AuthClientConfig and creates an AuthClient.
func NewClient(config *Config) (Client, error) {
	if config == nil {
		return nil, NewConfigurationError("configuration cannot be nil", nil)
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	if len(config.Scopes) == 0 {
		config.Scopes = []string{"hiiretail:iam"}
	}

	// Convert to internal configuration format
	authConfig := &AuthClientConfig{
		TenantID:         config.TenantID,
		ClientID:         config.ClientID,
		ClientSecret:     config.ClientSecret,
		TokenURL:         config.AuthURL,
		Scopes:           config.Scopes,
		Timeout:          config.Timeout,
		MaxRetries:       config.MaxRetries,
		DisableDiscovery: config.DisableDiscovery,
	}

	// Resolve endpoints if not provided
	if config.AuthURL == "" || config.APIURL == "" {
		resolver := NewEndpointResolver(config.TenantID, config.Environment)

		if config.AuthURL == "" {
			authURL, err := resolver.ResolveAuthURL()
			if err != nil {
				return nil, NewConfigurationError("failed to resolve auth URL", err)
			}
			authConfig.TokenURL = authURL
		}

		if config.APIURL == "" {
			apiURL, err := resolver.ResolveAPIURL()
			if err != nil {
				return nil, NewConfigurationError("failed to resolve API URL", err)
			}
			authConfig.BaseURL = apiURL
		}
	} else {
		authConfig.BaseURL = config.APIURL
	}

	return NewAuthClient(authConfig)
}

// ValidateConfig validates OAuth2 configuration without creating a client.
//
// This is useful for configuration validation in Terraform providers or
// other scenarios where you want to validate config before client creation.
//
// Example:
//
//	if err := auth.ValidateConfig(config); err != nil {
//		return fmt.Errorf("invalid auth config: %w", err)
//	}
func ValidateConfig(config *Config) error {
	if config == nil {
		return NewConfigurationError("configuration cannot be nil", nil)
	}

	// Use the OAuth2Configuration validator
	oauth2Config := &OAuth2Configuration{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TenantID:     config.TenantID,
		Environment:  config.Environment,
		AuthURL:      config.AuthURL,
		APIURL:       config.APIURL,
		Scopes:       config.Scopes,
		Timeout:      int(config.Timeout.Seconds()),
		MaxRetries:   config.MaxRetries,
		SkipTLS:      config.SkipTLS,
	}

	result := oauth2Config.Validate()
	if !result.Valid {
		return NewConfigurationError("configuration validation failed: "+result.Errors[0], nil)
	}

	return nil
}

// ResolveEndpoints resolves OAuth2 and API endpoints for a given tenant and environment.
//
// This is a utility function that can be used to preview endpoint resolution
// without creating a full client.
//
// Example:
//
//	authURL, apiURL, err := auth.ResolveEndpoints("my-tenant-123", "production")
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("Auth URL: %s\n", authURL)
//	fmt.Printf("API URL: %s\n", apiURL)
func ResolveEndpoints(tenantID, environment string) (authURL, apiURL string, err error) {
	resolver := NewEndpointResolver(tenantID, environment)

	authURL, err = resolver.ResolveAuthURL()
	if err != nil {
		return "", "", err
	}

	apiURL, err = resolver.ResolveAPIURL()
	if err != nil {
		return "", "", err
	}

	return authURL, apiURL, nil
}

// IsTestEnvironment checks if a tenant ID indicates a test environment.
//
// This is a utility function for determining if special test/dev behavior
// should be enabled based on tenant naming patterns.
//
// Example:
//
//	if auth.IsTestEnvironment("test-tenant-123") {
//		// Enable debug logging, relaxed SSL, etc.
//	}
func IsTestEnvironment(tenantID string) bool {
	resolver := NewEndpointResolver(tenantID, "")
	return resolver.IsTestEnvironment()
}

// HTTPClientOption defines options for HTTP client configuration
type HTTPClientOption func(*http.Client)

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) HTTPClientOption {
	return func(client *http.Client) {
		client.Timeout = timeout
	}
}

// WithTransport sets a custom HTTP transport
func WithTransport(transport http.RoundTripper) HTTPClientOption {
	return func(client *http.Client) {
		client.Transport = transport
	}
}

// NewHTTPClient creates an HTTP client with OAuth2 authentication.
//
// This is a convenience function that creates a client and returns an
// authenticated HTTP client in one step.
//
// Example:
//
//	httpClient, err := auth.NewHTTPClient(ctx, config)
//	if err != nil {
//		return err
//	}
//
//	resp, err := httpClient.Get("https://iam-api.retailsvc.com/roles")
func NewHTTPClient(ctx context.Context, config *Config, opts ...HTTPClientOption) (*http.Client, error) {
	client, err := New(config)
	if err != nil {
		return nil, err
	}

	httpClient, err := client.HTTPClient(ctx)
	if err != nil {
		client.Close()
		return nil, err
	}

	// Apply options
	for _, opt := range opts {
		opt(httpClient)
	}

	return httpClient, nil
}

// MustNew creates a new client and panics if there's an error.
//
// This is useful for package-level initialization where configuration
// is known to be valid.
//
// Example:
//
//	var authClient = auth.MustNew(&auth.Config{
//		ClientID:     "known-valid-client-id",
//		ClientSecret: "known-valid-client-secret",
//		TenantID:     "known-valid-tenant-id",
//	})
func MustNew(config *Config) Client {
	client, err := New(config)
	if err != nil {
		panic("auth.MustNew: " + err.Error())
	}
	return client
}

// Ensure our implementation satisfies the interface
var _ Client = (*AuthClient)(nil)
