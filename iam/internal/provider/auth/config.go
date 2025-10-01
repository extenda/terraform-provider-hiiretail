package auth

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// OAuth2Configuration holds the configuration for OAuth2 authentication
type OAuth2Configuration struct {
	// Required fields
	ClientID     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	TenantID     string `json:"tenant_id" validate:"required"`

	// Optional fields with defaults
	Environment string `json:"environment,omitempty"`
	AuthURL     string `json:"auth_url,omitempty"`
	APIURL      string `json:"api_url,omitempty"`

	// OAuth2 specific settings
	Scopes  []string `json:"scopes,omitempty"`
	Timeout int      `json:"timeout,omitempty"` // seconds

	// Advanced options
	MaxRetries int  `json:"max_retries,omitempty"`
	SkipTLS    bool `json:"skip_tls,omitempty"` // for testing only
}

// ValidationResult holds validation results
type ValidationResult struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

// NewOAuth2Configuration creates a new OAuth2Configuration with defaults
func NewOAuth2Configuration() *OAuth2Configuration {
	return &OAuth2Configuration{
		Environment: "production",
		Scopes:      []string{"hiiretail:iam"},
		Timeout:     30,
		MaxRetries:  3,
		SkipTLS:     false,
	}
}

// LoadFromEnvironment loads configuration from environment variables
func (c *OAuth2Configuration) LoadFromEnvironment() {
	if clientID := os.Getenv("HIIRETAIL_CLIENT_ID"); clientID != "" {
		c.ClientID = clientID
	}

	if clientSecret := os.Getenv("HIIRETAIL_CLIENT_SECRET"); clientSecret != "" {
		c.ClientSecret = clientSecret
	}

	if tenantID := os.Getenv("HIIRETAIL_TENANT_ID"); tenantID != "" {
		c.TenantID = tenantID
	}

	if environment := os.Getenv("HIIRETAIL_ENVIRONMENT"); environment != "" {
		c.Environment = environment
	}

	if authURL := os.Getenv("HIIRETAIL_AUTH_URL"); authURL != "" {
		c.AuthURL = authURL
	}

	if apiURL := os.Getenv("HIIRETAIL_API_URL"); apiURL != "" {
		c.APIURL = apiURL
	}
}

// Validate validates the OAuth2 configuration
func (c *OAuth2Configuration) Validate() *ValidationResult {
	result := &ValidationResult{Valid: true, Errors: []string{}}

	// Validate required fields
	if c.ClientID == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "client_id is required and cannot be empty")
	}

	if c.ClientSecret == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "client_secret is required and cannot be empty")
	}

	if c.TenantID == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "tenant_id is required and cannot be empty")
	}

	// Validate tenant ID format
	if c.TenantID != "" && !isValidTenantID(c.TenantID) {
		result.Valid = false
		result.Errors = append(result.Errors, "tenant_id must contain only alphanumeric characters and hyphens")
	}

	// Validate environment
	if c.Environment != "" && !isValidEnvironment(c.Environment) {
		result.Valid = false
		result.Errors = append(result.Errors, "environment must be one of: production, test, dev")
	}

	// Validate URLs if provided
	if c.AuthURL != "" && !isValidHTTPSURL(c.AuthURL) {
		result.Valid = false
		result.Errors = append(result.Errors, "auth_url must be a valid HTTPS URL")
	}

	if c.APIURL != "" && !isValidHTTPSURL(c.APIURL) {
		result.Valid = false
		result.Errors = append(result.Errors, "api_url must be a valid HTTPS URL")
	}

	// Validate timeout
	if c.Timeout <= 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "timeout must be greater than 0")
	}

	// Validate max retries
	if c.MaxRetries < 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "max_retries cannot be negative")
	}

	return result
}

// ResolveEndpoints resolves OAuth2 and API endpoints based on tenant and environment
func (c *OAuth2Configuration) ResolveEndpoints() error {
	resolver := NewEndpointResolver(c.TenantID, c.Environment)

	// Resolve auth URL if not provided
	if c.AuthURL == "" {
		var err error
		c.AuthURL, err = resolver.ResolveAuthURL()
		if err != nil {
			return fmt.Errorf("failed to resolve auth URL: %w", err)
		}
	}

	// Resolve API URL if not provided
	if c.APIURL == "" {
		var err error
		c.APIURL, err = resolver.ResolveAPIURL()
		if err != nil {
			return fmt.Errorf("failed to resolve API URL: %w", err)
		}
	}

	return nil
}

// GetTokenURL returns the OAuth2 token endpoint URL
func (c *OAuth2Configuration) GetTokenURL() string {
	if c.AuthURL == "" {
		return ""
	}

	// Ensure the auth URL has the correct token endpoint path
	if strings.HasSuffix(c.AuthURL, "/oauth/token") {
		return c.AuthURL
	}

	// Remove trailing slash and add token endpoint
	authURL := strings.TrimSuffix(c.AuthURL, "/")
	return authURL + "/oauth/token"
}

// GetScope returns the OAuth2 scope string
func (c *OAuth2Configuration) GetScope() string {
	if len(c.Scopes) == 0 {
		return "hiiretail:iam"
	}
	return strings.Join(c.Scopes, " ")
}

// IsTestEnvironment returns true if this is a test/dev environment
func (c *OAuth2Configuration) IsTestEnvironment() bool {
	env := strings.ToLower(c.Environment)
	tenantLower := strings.ToLower(c.TenantID)

	return env == "test" || env == "dev" ||
		strings.Contains(tenantLower, "test") ||
		strings.Contains(tenantLower, "dev")
}

// RedactSensitiveData returns a copy of the configuration with sensitive data redacted
func (c *OAuth2Configuration) RedactSensitiveData() *OAuth2Configuration {
	redacted := *c

	if redacted.ClientID != "" {
		redacted.ClientID = redactCredential(redacted.ClientID)
	}

	if redacted.ClientSecret != "" {
		redacted.ClientSecret = redactCredential(redacted.ClientSecret)
	}

	return &redacted
}

// String returns a string representation with sensitive data redacted
func (c *OAuth2Configuration) String() string {
	redacted := c.RedactSensitiveData()
	return fmt.Sprintf("OAuth2Configuration{ClientID: %s, TenantID: %s, Environment: %s, AuthURL: %s, APIURL: %s}",
		redacted.ClientID, redacted.TenantID, redacted.Environment, redacted.AuthURL, redacted.APIURL)
}

// Helper functions for validation

// isValidTenantID validates tenant ID format
func isValidTenantID(tenantID string) bool {
	if len(tenantID) == 0 || len(tenantID) > 100 {
		return false
	}

	// Allow alphanumeric characters, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9_-]*[a-zA-Z0-9]$`, tenantID)
	return matched
}

// isValidEnvironment validates environment value
func isValidEnvironment(env string) bool {
	validEnvironments := []string{"production", "test", "dev", "staging"}
	for _, validEnv := range validEnvironments {
		if strings.EqualFold(env, validEnv) {
			return true
		}
	}
	return false
}

// isValidHTTPSURL validates HTTPS URL format
func isValidHTTPSURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Scheme == "https" && parsedURL.Host != ""
}

// redactCredential redacts sensitive credential data
func redactCredential(credential string) string {
	if len(credential) <= 8 {
		return "***"
	}

	// Show first 4 and last 4 characters
	return credential[:4] + "***" + credential[len(credential)-4:]
}

// ConfigurationBuilder helps build OAuth2Configuration with validation
type ConfigurationBuilder struct {
	config *OAuth2Configuration
	errors []string
}

// NewConfigurationBuilder creates a new builder
func NewConfigurationBuilder() *ConfigurationBuilder {
	return &ConfigurationBuilder{
		config: NewOAuth2Configuration(),
		errors: []string{},
	}
}

// WithClientID sets the client ID
func (b *ConfigurationBuilder) WithClientID(clientID string) *ConfigurationBuilder {
	if clientID == "" {
		b.errors = append(b.errors, "client_id cannot be empty")
	}
	b.config.ClientID = clientID
	return b
}

// WithClientSecret sets the client secret
func (b *ConfigurationBuilder) WithClientSecret(clientSecret string) *ConfigurationBuilder {
	if clientSecret == "" {
		b.errors = append(b.errors, "client_secret cannot be empty")
	}
	b.config.ClientSecret = clientSecret
	return b
}

// WithTenantID sets the tenant ID
func (b *ConfigurationBuilder) WithTenantID(tenantID string) *ConfigurationBuilder {
	if tenantID == "" {
		b.errors = append(b.errors, "tenant_id cannot be empty")
	} else if !isValidTenantID(tenantID) {
		b.errors = append(b.errors, "tenant_id format is invalid")
	}
	b.config.TenantID = tenantID
	return b
}

// WithEnvironment sets the environment
func (b *ConfigurationBuilder) WithEnvironment(environment string) *ConfigurationBuilder {
	if environment != "" && !isValidEnvironment(environment) {
		b.errors = append(b.errors, "invalid environment")
	}
	b.config.Environment = environment
	return b
}

// WithAuthURL sets the auth URL
func (b *ConfigurationBuilder) WithAuthURL(authURL string) *ConfigurationBuilder {
	if authURL != "" && !isValidHTTPSURL(authURL) {
		b.errors = append(b.errors, "auth_url must be a valid HTTPS URL")
	}
	b.config.AuthURL = authURL
	return b
}

// WithAPIURL sets the API URL
func (b *ConfigurationBuilder) WithAPIURL(apiURL string) *ConfigurationBuilder {
	if apiURL != "" && !isValidHTTPSURL(apiURL) {
		b.errors = append(b.errors, "api_url must be a valid HTTPS URL")
	}
	b.config.APIURL = apiURL
	return b
}

// WithScopes sets the OAuth2 scopes
func (b *ConfigurationBuilder) WithScopes(scopes []string) *ConfigurationBuilder {
	b.config.Scopes = scopes
	return b
}

// WithTimeout sets the timeout in seconds
func (b *ConfigurationBuilder) WithTimeout(timeout int) *ConfigurationBuilder {
	if timeout <= 0 {
		b.errors = append(b.errors, "timeout must be greater than 0")
	}
	b.config.Timeout = timeout
	return b
}

// WithMaxRetries sets the maximum number of retries
func (b *ConfigurationBuilder) WithMaxRetries(maxRetries int) *ConfigurationBuilder {
	if maxRetries < 0 {
		b.errors = append(b.errors, "max_retries cannot be negative")
	}
	b.config.MaxRetries = maxRetries
	return b
}

// LoadFromEnvironment loads values from environment variables
func (b *ConfigurationBuilder) LoadFromEnvironment() *ConfigurationBuilder {
	b.config.LoadFromEnvironment()
	return b
}

// Build builds and validates the configuration
func (b *ConfigurationBuilder) Build() (*OAuth2Configuration, error) {
	// Add validation errors
	validation := b.config.Validate()
	if !validation.Valid {
		b.errors = append(b.errors, validation.Errors...)
	}

	if len(b.errors) > 0 {
		return nil, fmt.Errorf("configuration validation failed: %s", strings.Join(b.errors, "; "))
	}

	// Resolve endpoints
	if err := b.config.ResolveEndpoints(); err != nil {
		return nil, fmt.Errorf("endpoint resolution failed: %w", err)
	}

	return b.config, nil
}
