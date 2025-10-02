package auth

import (
	"fmt"
	"strings"
)

// EndpointResolver resolves OAuth2 and API endpoints based on tenant and environment
type EndpointResolver struct {
	TenantID    string
	Environment string
}

// EndpointMapping defines the mapping of environments to endpoints
type EndpointMapping struct {
	AuthBaseURL string
	APIBaseURL  string
}

// NewEndpointResolver creates a new endpoint resolver
func NewEndpointResolver(tenantID, environment string) *EndpointResolver {
	return &EndpointResolver{
		TenantID:    tenantID,
		Environment: environment,
	}
}

// ResolveAuthURL resolves the OAuth2 authentication endpoint URL
func (r *EndpointResolver) ResolveAuthURL() (string, error) {
	mapping, err := r.getEndpointMapping()
	if err != nil {
		return "", fmt.Errorf("failed to get endpoint mapping: %w", err)
	}

	authURL := mapping.AuthBaseURL + "/oauth2/token"
	return authURL, nil
} // ResolveAPIURL resolves the IAM API base URL
func (r *EndpointResolver) ResolveAPIURL() (string, error) {
	mapping, err := r.getEndpointMapping()
	if err != nil {
		return "", fmt.Errorf("failed to get endpoint mapping: %w", err)
	}

	return mapping.APIBaseURL, nil
}

// getEndpointMapping returns the endpoint mapping based on tenant and environment
func (r *EndpointResolver) getEndpointMapping() (*EndpointMapping, error) {
	// Determine environment from tenant ID if not explicitly set
	effectiveEnv := r.determineEffectiveEnvironment()

	switch effectiveEnv {
	case "production":
		return &EndpointMapping{
			AuthBaseURL: "https://auth.retailsvc.com",
			APIBaseURL:  "https://iam-api.retailsvc.com",
		}, nil

	case "test":
		return &EndpointMapping{
			AuthBaseURL: "https://auth.retailsvc-test.com",
			APIBaseURL:  "https://iam-api.retailsvc-test.com",
		}, nil

	case "dev", "development":
		return &EndpointMapping{
			AuthBaseURL: "https://auth.retailsvc-dev.com",
			APIBaseURL:  "https://iam-api.retailsvc-dev.com",
		}, nil

	case "staging":
		return &EndpointMapping{
			AuthBaseURL: "https://auth.retailsvc-staging.com",
			APIBaseURL:  "https://iam-api.retailsvc-staging.com",
		}, nil

	default:
		return nil, fmt.Errorf("unsupported environment: %s", effectiveEnv)
	}
}

// determineEffectiveEnvironment determines the effective environment based on tenant ID and explicit environment
func (r *EndpointResolver) determineEffectiveEnvironment() string {
	// If environment is explicitly set, use it
	if r.Environment != "" {
		return strings.ToLower(r.Environment)
	}

	// Determine environment from tenant ID patterns
	tenantLower := strings.ToLower(r.TenantID)

	// Check for test patterns
	if strings.Contains(tenantLower, "test") ||
		strings.Contains(tenantLower, "tst") ||
		strings.HasPrefix(tenantLower, "test-") ||
		strings.HasSuffix(tenantLower, "-test") {
		return "test"
	}

	// Check for dev patterns
	if strings.Contains(tenantLower, "dev") ||
		strings.Contains(tenantLower, "development") ||
		strings.HasPrefix(tenantLower, "dev-") ||
		strings.HasSuffix(tenantLower, "-dev") {
		return "dev"
	}

	// Check for staging patterns
	if strings.Contains(tenantLower, "staging") ||
		strings.Contains(tenantLower, "stage") ||
		strings.Contains(tenantLower, "stg") ||
		strings.HasPrefix(tenantLower, "staging-") ||
		strings.HasSuffix(tenantLower, "-staging") {
		return "staging"
	}

	// Default to production
	return "production"
}

// IsTestEnvironment returns true if the resolved environment is for testing
func (r *EndpointResolver) IsTestEnvironment() bool {
	env := r.determineEffectiveEnvironment()
	return env == "test" || env == "dev" || env == "development" || env == "staging"
}

// ValidateEndpoints validates that the resolved endpoints are accessible
func (r *EndpointResolver) ValidateEndpoints() error {
	authURL, err := r.ResolveAuthURL()
	if err != nil {
		return fmt.Errorf("failed to resolve auth URL: %w", err)
	}

	apiURL, err := r.ResolveAPIURL()
	if err != nil {
		return fmt.Errorf("failed to resolve API URL: %w", err)
	}

	// Validate URL formats
	if !isValidHTTPSURL(authURL) {
		return fmt.Errorf("resolved auth URL is not a valid HTTPS URL: %s", authURL)
	}

	if !isValidHTTPSURL(apiURL) {
		return fmt.Errorf("resolved API URL is not a valid HTTPS URL: %s", apiURL)
	}

	return nil
}

// GetEndpointInfo returns detailed information about resolved endpoints
func (r *EndpointResolver) GetEndpointInfo() (*EndpointInfo, error) {
	authURL, err := r.ResolveAuthURL()
	if err != nil {
		return nil, err
	}

	apiURL, err := r.ResolveAPIURL()
	if err != nil {
		return nil, err
	}

	effectiveEnv := r.determineEffectiveEnvironment()

	return &EndpointInfo{
		TenantID:             r.TenantID,
		Environment:          r.Environment,
		EffectiveEnvironment: effectiveEnv,
		AuthURL:              authURL,
		APIURL:               apiURL,
		IsTestEnvironment:    r.IsTestEnvironment(),
	}, nil
}

// EndpointInfo contains detailed endpoint resolution information
type EndpointInfo struct {
	TenantID             string `json:"tenant_id"`
	Environment          string `json:"environment"`
	EffectiveEnvironment string `json:"effective_environment"`
	AuthURL              string `json:"auth_url"`
	APIURL               string `json:"api_url"`
	IsTestEnvironment    bool   `json:"is_test_environment"`
}

// String returns a string representation of endpoint info
func (e *EndpointInfo) String() string {
	return fmt.Sprintf("EndpointInfo{TenantID: %s, Environment: %s -> %s, AuthURL: %s, APIURL: %s, IsTest: %t}",
		e.TenantID, e.Environment, e.EffectiveEnvironment, e.AuthURL, e.APIURL, e.IsTestEnvironment)
}

// CustomEndpointResolver allows for custom endpoint resolution logic
type CustomEndpointResolver struct {
	*EndpointResolver
	CustomAuthURL string
	CustomAPIURL  string
}

// NewCustomEndpointResolver creates a resolver with custom URLs
func NewCustomEndpointResolver(tenantID, environment, customAuthURL, customAPIURL string) *CustomEndpointResolver {
	return &CustomEndpointResolver{
		EndpointResolver: NewEndpointResolver(tenantID, environment),
		CustomAuthURL:    customAuthURL,
		CustomAPIURL:     customAPIURL,
	}
}

// ResolveAuthURL resolves auth URL, preferring custom URL if provided
func (r *CustomEndpointResolver) ResolveAuthURL() (string, error) {
	if r.CustomAuthURL != "" {
		if !isValidHTTPSURL(r.CustomAuthURL) {
			return "", fmt.Errorf("custom auth URL is not a valid HTTPS URL: %s", r.CustomAuthURL)
		}

		// Ensure it has the correct token endpoint path
		if strings.HasSuffix(r.CustomAuthURL, "/oauth/token") {
			return r.CustomAuthURL, nil
		}

		// Add token endpoint path
		customURL := strings.TrimSuffix(r.CustomAuthURL, "/")
		return customURL + "/oauth/token", nil
	}

	return r.EndpointResolver.ResolveAuthURL()
}

// ResolveAPIURL resolves API URL, preferring custom URL if provided
func (r *CustomEndpointResolver) ResolveAPIURL() (string, error) {
	if r.CustomAPIURL != "" {
		if !isValidHTTPSURL(r.CustomAPIURL) {
			return "", fmt.Errorf("custom API URL is not a valid HTTPS URL: %s", r.CustomAPIURL)
		}
		return r.CustomAPIURL, nil
	}

	return r.EndpointResolver.ResolveAPIURL()
}

// ResolverInterface defines the interface for endpoint resolution
type ResolverInterface interface {
	ResolveAuthURL() (string, error)
	ResolveAPIURL() (string, error)
	IsTestEnvironment() bool
	ValidateEndpoints() error
}

// Ensure our types implement the interface
var (
	_ ResolverInterface = (*EndpointResolver)(nil)
	_ ResolverInterface = (*CustomEndpointResolver)(nil)
)
