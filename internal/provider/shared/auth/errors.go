package auth

import (
	"fmt"
	"strconv"
	"time"
)

// AuthErrorType represents the type of authentication error
type AuthErrorType int

const (
	// AuthErrorUnknown represents an unclassified error
	AuthErrorUnknown AuthErrorType = iota

	// AuthErrorConfiguration represents configuration validation errors
	AuthErrorConfiguration

	// AuthErrorDiscovery represents OAuth2 discovery endpoint errors
	AuthErrorDiscovery

	// AuthErrorCredentials represents invalid client credentials
	AuthErrorCredentials

	// AuthErrorNetwork represents network connectivity issues
	AuthErrorNetwork

	// AuthErrorServerError represents OAuth2 server errors (5xx)
	AuthErrorServerError

	// AuthErrorRateLimit represents rate limiting errors (429)
	AuthErrorRateLimit

	// AuthErrorTokenExpired represents token expiration during operation
	AuthErrorTokenExpired
)

// String returns a human-readable string representation of the error type
func (t AuthErrorType) String() string {
	switch t {
	case AuthErrorConfiguration:
		return "Configuration Error"
	case AuthErrorDiscovery:
		return "Discovery Error"
	case AuthErrorCredentials:
		return "Credentials Error"
	case AuthErrorNetwork:
		return "Network Error"
	case AuthErrorServerError:
		return "Server Error"
	case AuthErrorRateLimit:
		return "Rate Limit Error"
	case AuthErrorTokenExpired:
		return "Token Expired Error"
	default:
		return "Unknown Error"
	}
}

// AuthError represents authentication-specific errors with enhanced context
type AuthError struct {
	// Type classifies the error for appropriate handling
	Type AuthErrorType

	// Message provides a user-friendly error description
	Message string

	// Underlying contains the original error that caused this AuthError
	Underlying error

	// Retryable indicates whether the operation should be retried
	Retryable bool

	// RetryAfter specifies the minimum delay before retrying (for rate limiting)
	RetryAfter time.Duration

	// Context provides additional debugging information
	Context map[string]interface{}
}

// Error implements the error interface
func (e *AuthError) Error() string {
	if e.Underlying != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type.String(), e.Message, e.Underlying)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

// Unwrap implements the errors.Unwrap interface for Go 1.13+ error handling
func (e *AuthError) Unwrap() error {
	return e.Underlying
}

// IsRetryable returns true if the error indicates a retryable condition
func (e *AuthError) IsRetryable() bool {
	return e.Retryable
}

// ShouldRetryAfter returns the delay before retrying, if applicable
func (e *AuthError) ShouldRetryAfter() time.Duration {
	return e.RetryAfter
}

// WithContext adds context information to the error
func (e *AuthError) WithContext(key string, value interface{}) *AuthError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// ConfigValidationError represents configuration validation errors with field-specific details
type ConfigValidationError struct {
	// Field is the name of the configuration field that failed validation
	Field string

	// Value is the invalid value (may be redacted for sensitive fields)
	Value interface{}

	// Constraint describes what validation rule was violated
	Constraint string

	// Suggestion provides guidance on how to fix the validation error
	Suggestion string
}

// Error implements the error interface
func (e *ConfigValidationError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("invalid %s: %s. Suggestion: %s", e.Field, e.Constraint, e.Suggestion)
	}
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Constraint)
}

// NewConfigValidationError creates a new configuration validation error
func NewConfigValidationError(field, constraint, suggestion string, value interface{}) *ConfigValidationError {
	// Redact sensitive values
	if isSensitiveField(field) {
		value = "[REDACTED]"
	}

	return &ConfigValidationError{
		Field:      field,
		Value:      value,
		Constraint: constraint,
		Suggestion: suggestion,
	}
}

// RetryConfig defines retry behavior for different error types
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int

	// BaseDelay is the initial delay between retries
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration

	// Multiplier is the exponential backoff multiplier
	Multiplier float64

	// Jitter indicates whether to add random jitter to retry delays
	Jitter bool

	// RetryableErrors maps error types to their retry eligibility
	RetryableErrors map[AuthErrorType]bool
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
		Jitter:      true,
		RetryableErrors: map[AuthErrorType]bool{
			AuthErrorNetwork:       true,
			AuthErrorServerError:   true,
			AuthErrorRateLimit:     true,
			AuthErrorTokenExpired:  true,
			AuthErrorCredentials:   false, // Never retry credential errors
			AuthErrorConfiguration: false, // Never retry configuration errors
			AuthErrorDiscovery:     false, // Never retry discovery errors (usually permanent)
		},
	}
}

// ShouldRetry determines if an error should be retried based on the retry configuration
func (rc *RetryConfig) ShouldRetry(err error, attempt int) bool {
	if attempt >= rc.MaxAttempts {
		return false
	}

	var authErr *AuthError
	if !isAuthError(err, &authErr) {
		return false // Only retry AuthError types
	}

	retryable, exists := rc.RetryableErrors[authErr.Type]
	return exists && retryable && authErr.Retryable
}

// GetDelay calculates the delay before the next retry attempt
func (rc *RetryConfig) GetDelay(attempt int, err error) time.Duration {
	var authErr *AuthError
	if isAuthError(err, &authErr) && authErr.RetryAfter > 0 {
		// Honor server-specified retry delay (e.g., from Retry-After header)
		return authErr.RetryAfter
	}

	// Calculate exponential backoff delay
	delay := time.Duration(float64(rc.BaseDelay) * pow(rc.Multiplier, float64(attempt)))

	if delay > rc.MaxDelay {
		delay = rc.MaxDelay
	}

	// Add jitter if enabled
	if rc.Jitter {
		jitter := time.Duration(float64(delay) * 0.1 * (2.0*random() - 1.0))
		delay += jitter
	}

	return delay
}

// Error factory functions for common authentication errors

// NewCredentialsError creates an error for invalid OAuth2 credentials
func NewCredentialsError(message string, underlying error) *AuthError {
	return &AuthError{
		Type:       AuthErrorCredentials,
		Message:    message,
		Underlying: underlying,
		Retryable:  false,
	}
}

// NewNetworkError creates an error for network connectivity issues
func NewNetworkError(message string, underlying error) *AuthError {
	return &AuthError{
		Type:       AuthErrorNetwork,
		Message:    message,
		Underlying: underlying,
		Retryable:  true,
	}
}

// NewServerError creates an error for OAuth2 server errors
func NewServerError(message string, underlying error) *AuthError {
	return &AuthError{
		Type:       AuthErrorServerError,
		Message:    message,
		Underlying: underlying,
		Retryable:  true,
	}
}

// NewRateLimitError creates an error for rate limiting with retry delay
func NewRateLimitError(message string, retryAfter time.Duration) *AuthError {
	return &AuthError{
		Type:       AuthErrorRateLimit,
		Message:    message,
		Underlying: fmt.Errorf("rate limited, retry after %v", retryAfter),
		Retryable:  true,
		RetryAfter: retryAfter,
	}
}

// NewDiscoveryError creates an error for OAuth2 discovery failures
func NewDiscoveryError(message string, underlying error) *AuthError {
	return &AuthError{
		Type:       AuthErrorDiscovery,
		Message:    message,
		Underlying: underlying,
		Retryable:  false, // Discovery errors are usually permanent
	}
}

// NewConfigurationError creates an error for configuration validation failures
func NewConfigurationError(message string, underlying error) *AuthError {
	return &AuthError{
		Type:       AuthErrorConfiguration,
		Message:    message,
		Underlying: underlying,
		Retryable:  false,
	}
}

// NewTokenExpiredError creates an error for token expiration during operations
func NewTokenExpiredError(message string) *AuthError {
	return &AuthError{
		Type:       AuthErrorTokenExpired,
		Message:    message,
		Underlying: fmt.Errorf("token expired during operation"),
		Retryable:  true, // Token refresh should be attempted
	}
}

// Helper functions

// isSensitiveField checks if a configuration field contains sensitive information
func isSensitiveField(field string) bool {
	sensitiveFields := []string{
		"client_secret",
		"access_token",
		"refresh_token",
		"bearer_token",
		"api_key",
		"private_key",
		"password",
		"secret",
		"token",
	}

	for _, sensitive := range sensitiveFields {
		if field == sensitive {
			return true
		}
	}
	return false
}

// isAuthError checks if an error is an AuthError and extracts it
func isAuthError(err error, authErr **AuthError) bool {
	if err == nil {
		return false
	}

	switch e := err.(type) {
	case *AuthError:
		*authErr = e
		return true
	default:
		return false
	}
}

// pow calculates base^exp for floating point numbers
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

// random returns a pseudo-random number between 0 and 1
func random() float64 {
	// Simple linear congruential generator for jitter
	// In production, you might want to use crypto/rand for better randomness
	seed := time.Now().UnixNano()
	return float64((seed*1103515245+12345)&0x7fffffff) / float64(0x7fffffff)
}

// ParseRetryAfterHeader parses the Retry-After header value from HTTP responses
func ParseRetryAfterHeader(retryAfter string) time.Duration {
	if retryAfter == "" {
		return 0
	}

	// Try to parse as seconds (integer)
	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		return time.Duration(seconds) * time.Second
	}

	// Try to parse as HTTP date (not commonly used for OAuth2)
	if parsedTime, err := time.Parse(time.RFC1123, retryAfter); err == nil {
		duration := time.Until(parsedTime)
		if duration > 0 {
			return duration
		}
	}

	// Default fallback
	return 60 * time.Second
}

// ErrorCode returns a standardized error code for API responses
func (e *AuthError) ErrorCode() string {
	switch e.Type {
	case AuthErrorCredentials:
		return "invalid_client"
	case AuthErrorConfiguration:
		return "invalid_request"
	case AuthErrorDiscovery:
		return "discovery_failed"
	case AuthErrorNetwork:
		return "network_error"
	case AuthErrorServerError:
		return "server_error"
	case AuthErrorRateLimit:
		return "rate_limited"
	case AuthErrorTokenExpired:
		return "token_expired"
	default:
		return "unknown_error"
	}
}

// TroubleshootingGuidance returns helpful troubleshooting information
func (e *AuthError) TroubleshootingGuidance() string {
	switch e.Type {
	case AuthErrorCredentials:
		return "Verify your client_id and client_secret are correct. Check that the OAuth2 client is properly configured in OCMS."
	case AuthErrorConfiguration:
		return "Review your provider configuration. Ensure all required parameters are set and URLs are valid HTTPS endpoints."
	case AuthErrorDiscovery:
		return "Check if the OCMS discovery endpoint is accessible. Verify the base URL and network connectivity. Consider providing explicit token_url."
	case AuthErrorNetwork:
		return "Check network connectivity to OCMS endpoints. Verify firewall settings and DNS resolution."
	case AuthErrorServerError:
		return "OCMS service may be temporarily unavailable. Check service status and try again later."
	case AuthErrorRateLimit:
		return "Too many requests to OCMS. Wait before retrying or check if your application is making excessive token requests."
	case AuthErrorTokenExpired:
		return "Token expired during operation. This should be handled automatically by token refresh."
	default:
		return "Check the error details and logs for more information."
	}
}
