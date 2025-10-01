package auth

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ConfigValidationResult contains the result of configuration validation
type ConfigValidationResult struct {
	Valid    bool                 `json:"valid"`
	Errors   []ConfigFieldError   `json:"errors,omitempty"`
	Warnings []ConfigFieldWarning `json:"warnings,omitempty"`
}

// ConfigFieldError represents a configuration field validation error
type ConfigFieldError struct {
	Field      string `json:"field"`
	Value      string `json:"value,omitempty"` // Sanitized value for logging
	Rule       string `json:"rule"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

// ConfigFieldWarning represents a configuration field validation warning
type ConfigFieldWarning struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

// ValidationRules contains validation rules for configuration fields
type ValidationRules struct {
	TenantIDPattern    *regexp.Regexp
	ClientIDPattern    *regexp.Regexp
	MinClientSecretLen int
	MaxClientSecretLen int
	MinTimeoutSeconds  int
	MaxTimeoutSeconds  int
	MaxRetries         int
	RequiredURLSchemes []string
	AllowedScopes      []string
}

// DefaultValidationRules returns the default validation rules for HiiRetail IAM
func DefaultValidationRules() *ValidationRules {
	return &ValidationRules{
		// Tenant ID: alphanumeric with hyphens, 3-64 characters
		TenantIDPattern: regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-]{1,62}[a-zA-Z0-9]$`),

		// Client ID: UUID format or alphanumeric, 8-128 characters
		ClientIDPattern: regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_]{6,126}[a-zA-Z0-9]$`),

		// Client secret: minimum 8 characters, maximum 256 characters
		MinClientSecretLen: 8,
		MaxClientSecretLen: 256,

		// Timeout: 5 seconds to 300 seconds (5 minutes)
		MinTimeoutSeconds: 5,
		MaxTimeoutSeconds: 300,

		// Max retries: 0 to 10
		MaxRetries: 10,

		// Required URL schemes
		RequiredURLSchemes: []string{"https"},

		// Allowed OAuth2 scopes for HiiRetail IAM
		AllowedScopes: []string{
			"iam:read",
			"iam:write",
			"iam:admin",
			"roles:read",
			"roles:write",
			"groups:read",
			"groups:write",
			"bindings:read",
			"bindings:write",
		},
	}
}

// ValidateAuthConfig validates the complete OAuth2 authentication configuration
func ValidateAuthConfig(config *AuthClientConfig, rules *ValidationRules) *ConfigValidationResult {
	if config == nil {
		return &ConfigValidationResult{
			Valid: false,
			Errors: []ConfigFieldError{{
				Field:   "config",
				Rule:    "required",
				Message: "Configuration cannot be nil",
			}},
		}
	}

	if rules == nil {
		rules = DefaultValidationRules()
	}

	result := &ConfigValidationResult{
		Valid:    true,
		Errors:   []ConfigFieldError{},
		Warnings: []ConfigFieldWarning{},
	}

	// Validate each configuration field
	validateTenantID(config.TenantID, rules, result)
	validateClientID(config.ClientID, rules, result)
	validateClientSecret(config.ClientSecret, rules, result)
	validateBaseURL(config.BaseURL, rules, result)
	validateTokenURL(config.TokenURL, rules, result)
	validateScopes(config.Scopes, rules, result)
	validateTimeout(config.Timeout, rules, result)
	validateMaxRetries(config.MaxRetries, rules, result)
	validateCustomHeaders(config.CustomHeaders, rules, result)

	// Additional cross-field validation
	validateConfiguration(config, rules, result)

	// Set overall validity
	result.Valid = len(result.Errors) == 0

	return result
}

// validateTenantID validates the tenant ID field
func validateTenantID(tenantID string, rules *ValidationRules, result *ConfigValidationResult) {
	if tenantID == "" {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "tenant_id",
			Rule:       "required",
			Message:    "Tenant ID is required",
			Suggestion: "Provide a valid tenant identifier from your HiiRetail IAM configuration",
		})
		return
	}

	if !rules.TenantIDPattern.MatchString(tenantID) {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "tenant_id",
			Value:      sanitizeValue(tenantID),
			Rule:       "pattern",
			Message:    "Tenant ID format is invalid",
			Suggestion: "Use alphanumeric characters with hyphens, 3-64 characters long",
		})
	}
}

// validateClientID validates the OAuth2 client ID field
func validateClientID(clientID string, rules *ValidationRules, result *ConfigValidationResult) {
	if clientID == "" {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "client_id",
			Rule:       "required",
			Message:    "OAuth2 client ID is required",
			Suggestion: "Obtain a client ID from your OAuth2 provider configuration",
		})
		return
	}

	if !rules.ClientIDPattern.MatchString(clientID) {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "client_id",
			Value:      sanitizeValue(clientID),
			Rule:       "pattern",
			Message:    "Client ID format is invalid",
			Suggestion: "Use alphanumeric characters with hyphens or underscores, 8-128 characters long",
		})
	}
}

// validateClientSecret validates the OAuth2 client secret field
func validateClientSecret(clientSecret string, rules *ValidationRules, result *ConfigValidationResult) {
	if clientSecret == "" {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "client_secret",
			Value:      "[REDACTED]",
			Rule:       "required",
			Message:    "OAuth2 client secret is required",
			Suggestion: "Obtain a client secret from your OAuth2 provider configuration",
		})
		return
	}

	if len(clientSecret) < rules.MinClientSecretLen {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "client_secret",
			Value:      "[REDACTED]",
			Rule:       "min_length",
			Message:    fmt.Sprintf("Client secret must be at least %d characters long", rules.MinClientSecretLen),
			Suggestion: "Use a longer, more secure client secret",
		})
	}

	if len(clientSecret) > rules.MaxClientSecretLen {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:   "client_secret",
			Value:   "[REDACTED]",
			Rule:    "max_length",
			Message: fmt.Sprintf("Client secret must be no more than %d characters long", rules.MaxClientSecretLen),
		})
	}

	// Check for common weak patterns
	if isWeakSecret(clientSecret) {
		result.Warnings = append(result.Warnings, ConfigFieldWarning{
			Field:      "client_secret",
			Message:    "Client secret appears to be weak or follows a common pattern",
			Suggestion: "Use a cryptographically strong, randomly generated client secret",
		})
	}
}

// validateBaseURL validates the base URL field
func validateBaseURL(baseURL string, rules *ValidationRules, result *ConfigValidationResult) {
	if baseURL == "" {
		// Base URL is optional when token_url is provided
		return
	}

	if err := validateURL(baseURL, "base_url", rules.RequiredURLSchemes, result); err != nil {
		return
	}

	// Additional validation for base URL format
	if !strings.HasSuffix(baseURL, "/") {
		result.Warnings = append(result.Warnings, ConfigFieldWarning{
			Field:      "base_url",
			Message:    "Base URL should end with a trailing slash",
			Suggestion: fmt.Sprintf("Consider using '%s/' instead of '%s'", baseURL, baseURL),
		})
	}
}

// validateTokenURL validates the token URL field
func validateTokenURL(tokenURL string, rules *ValidationRules, result *ConfigValidationResult) {
	if tokenURL == "" {
		// Token URL is optional when base_url is provided for discovery
		return
	}

	validateURL(tokenURL, "token_url", rules.RequiredURLSchemes, result)
}

// validateURL validates a URL field against the validation rules
func validateURL(urlStr, fieldName string, requiredSchemes []string, result *ConfigValidationResult) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      fieldName,
			Value:      sanitizeValue(urlStr),
			Rule:       "url_format",
			Message:    "URL format is invalid",
			Suggestion: "Provide a valid URL with scheme, host, and optional path",
		})
		return err
	}

	if parsedURL.Scheme == "" {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      fieldName,
			Value:      sanitizeValue(urlStr),
			Rule:       "url_scheme",
			Message:    "URL must include a scheme (e.g., https://)",
			Suggestion: "Add https:// prefix to the URL",
		})
		return fmt.Errorf("missing URL scheme")
	}

	if parsedURL.Host == "" {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      fieldName,
			Value:      sanitizeValue(urlStr),
			Rule:       "url_host",
			Message:    "URL must include a host",
			Suggestion: "Provide a valid hostname or IP address",
		})
		return fmt.Errorf("missing URL host")
	}

	// Check scheme requirements
	schemeValid := false
	for _, scheme := range requiredSchemes {
		if parsedURL.Scheme == scheme {
			schemeValid = true
			break
		}
	}

	if !schemeValid {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      fieldName,
			Value:      sanitizeValue(urlStr),
			Rule:       "url_scheme_whitelist",
			Message:    fmt.Sprintf("URL scheme must be one of: %s", strings.Join(requiredSchemes, ", ")),
			Suggestion: "Use HTTPS for secure communication",
		})
		return fmt.Errorf("invalid URL scheme")
	}

	return nil
}

// validateScopes validates the OAuth2 scopes field
func validateScopes(scopes []string, rules *ValidationRules, result *ConfigValidationResult) {
	if len(scopes) == 0 {
		result.Warnings = append(result.Warnings, ConfigFieldWarning{
			Field:      "scopes",
			Message:    "No OAuth2 scopes specified, default scopes will be used",
			Suggestion: "Consider specifying explicit scopes for better security",
		})
		return
	}

	// Check for invalid scopes
	for _, scope := range scopes {
		if scope == "" {
			result.Errors = append(result.Errors, ConfigFieldError{
				Field:      "scopes",
				Rule:       "empty_scope",
				Message:    "Empty scope is not allowed",
				Suggestion: "Remove empty scopes or provide valid scope names",
			})
			continue
		}

		// Check if scope is in allowed list (optional validation)
		if len(rules.AllowedScopes) > 0 {
			scopeValid := false
			for _, allowedScope := range rules.AllowedScopes {
				if scope == allowedScope {
					scopeValid = true
					break
				}
			}

			if !scopeValid {
				result.Warnings = append(result.Warnings, ConfigFieldWarning{
					Field:      "scopes",
					Message:    fmt.Sprintf("Scope '%s' is not in the standard list", scope),
					Suggestion: fmt.Sprintf("Consider using standard scopes: %s", strings.Join(rules.AllowedScopes, ", ")),
				})
			}
		}
	}
}

// validateTimeout validates the timeout configuration
func validateTimeout(timeout time.Duration, rules *ValidationRules, result *ConfigValidationResult) {
	if timeout <= 0 {
		result.Warnings = append(result.Warnings, ConfigFieldWarning{
			Field:      "timeout",
			Message:    "Timeout not specified, default timeout will be used",
			Suggestion: "Consider specifying an explicit timeout for better control",
		})
		return
	}

	timeoutSeconds := int(timeout.Seconds())

	if timeoutSeconds < rules.MinTimeoutSeconds {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "timeout",
			Value:      timeout.String(),
			Rule:       "min_timeout",
			Message:    fmt.Sprintf("Timeout must be at least %d seconds", rules.MinTimeoutSeconds),
			Suggestion: fmt.Sprintf("Use a timeout of at least %ds for reliable operation", rules.MinTimeoutSeconds),
		})
	}

	if timeoutSeconds > rules.MaxTimeoutSeconds {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "timeout",
			Value:      timeout.String(),
			Rule:       "max_timeout",
			Message:    fmt.Sprintf("Timeout must be no more than %d seconds", rules.MaxTimeoutSeconds),
			Suggestion: fmt.Sprintf("Use a timeout of at most %ds to avoid hanging operations", rules.MaxTimeoutSeconds),
		})
	}
}

// validateMaxRetries validates the maximum retry configuration
func validateMaxRetries(maxRetries int, rules *ValidationRules, result *ConfigValidationResult) {
	if maxRetries < 0 {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "max_retries",
			Value:      fmt.Sprintf("%d", maxRetries),
			Rule:       "min_retries",
			Message:    "Maximum retries cannot be negative",
			Suggestion: "Use 0 for no retries or a positive number for retry attempts",
		})
	}

	if maxRetries > rules.MaxRetries {
		result.Warnings = append(result.Warnings, ConfigFieldWarning{
			Field:      "max_retries",
			Message:    fmt.Sprintf("High retry count (%d) may cause long delays", maxRetries),
			Suggestion: fmt.Sprintf("Consider using a lower retry count (max recommended: %d)", rules.MaxRetries),
		})
	}
}

// validateCustomHeaders validates the custom headers configuration
func validateCustomHeaders(headers map[string]string, rules *ValidationRules, result *ConfigValidationResult) {
	if len(headers) == 0 {
		return
	}

	// Check for problematic headers
	problematicHeaders := []string{
		"authorization",
		"content-type",
		"user-agent",
		"host",
		"content-length",
	}

	for headerName := range headers {
		lowerName := strings.ToLower(headerName)

		for _, problematic := range problematicHeaders {
			if lowerName == problematic {
				result.Warnings = append(result.Warnings, ConfigFieldWarning{
					Field:      "custom_headers",
					Message:    fmt.Sprintf("Custom header '%s' may conflict with automatically set headers", headerName),
					Suggestion: "Avoid setting headers that are automatically managed by the OAuth2 client",
				})
			}
		}

		// Check for empty header values
		if headers[headerName] == "" {
			result.Warnings = append(result.Warnings, ConfigFieldWarning{
				Field:      "custom_headers",
				Message:    fmt.Sprintf("Custom header '%s' has empty value", headerName),
				Suggestion: "Remove headers with empty values or provide meaningful values",
			})
		}
	}
}

// validateConfiguration performs cross-field validation
func validateConfiguration(config *AuthClientConfig, rules *ValidationRules, result *ConfigValidationResult) {
	// Check if either base_url or token_url is provided
	if config.BaseURL == "" && config.TokenURL == "" {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "base_url,token_url",
			Rule:       "one_required",
			Message:    "Either base_url (for discovery) or token_url (explicit) must be provided",
			Suggestion: "Provide base_url for automatic discovery or token_url for explicit configuration",
		})
	}

	// Warn if both are provided
	if config.BaseURL != "" && config.TokenURL != "" {
		result.Warnings = append(result.Warnings, ConfigFieldWarning{
			Field:      "base_url,token_url",
			Message:    "Both base_url and token_url are provided, token_url will take precedence",
			Suggestion: "Consider using only one method for clarity",
		})
	}

	// Check discovery configuration
	if config.DisableDiscovery && config.BaseURL != "" && config.TokenURL == "" {
		result.Errors = append(result.Errors, ConfigFieldError{
			Field:      "disable_discovery",
			Rule:       "logical_consistency",
			Message:    "Discovery disabled but no explicit token_url provided",
			Suggestion: "Enable discovery or provide an explicit token_url",
		})
	}
}

// Helper functions

// sanitizeValue creates a safe representation of a value for logging
func sanitizeValue(value string) string {
	if len(value) <= 8 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

// isWeakSecret checks if a client secret follows common weak patterns
func isWeakSecret(secret string) bool {
	// Check for common weak patterns
	weakPatterns := []string{
		"password",
		"secret",
		"123456",
		"qwerty",
		"admin",
		"test",
		"demo",
	}

	lowerSecret := strings.ToLower(secret)

	for _, pattern := range weakPatterns {
		if strings.Contains(lowerSecret, pattern) {
			return true
		}
	}

	// Check for sequential characters
	if hasSequentialChars(secret) {
		return true
	}

	// Check for repeated characters
	if hasRepeatedChars(secret, 4) {
		return true
	}

	return false
}

// hasSequentialChars checks if the string contains sequential characters
func hasSequentialChars(s string) bool {
	if len(s) < 3 {
		return false
	}

	for i := 0; i < len(s)-2; i++ {
		if s[i+1] == s[i]+1 && s[i+2] == s[i]+2 {
			return true
		}
	}

	return false
}

// hasRepeatedChars checks if the string has repeated characters
func hasRepeatedChars(s string, threshold int) bool {
	charCount := make(map[rune]int)

	for _, char := range s {
		charCount[char]++
		if charCount[char] >= threshold {
			return true
		}
	}

	return false
}

// ValidateProviderConfig validates the provider-level OAuth2 configuration
func ValidateProviderConfig(tenantID, clientID, clientSecret, baseURL, tokenURL string, scopes []string, timeout time.Duration, maxRetries int, disableDiscovery bool, customHeaders map[string]string) *ConfigValidationResult {
	config := &AuthClientConfig{
		TenantID:         tenantID,
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		BaseURL:          baseURL,
		TokenURL:         tokenURL,
		Scopes:           scopes,
		Timeout:          timeout,
		MaxRetries:       maxRetries,
		DisableDiscovery: disableDiscovery,
		CustomHeaders:    customHeaders,
	}

	return ValidateAuthConfig(config, DefaultValidationRules())
}
