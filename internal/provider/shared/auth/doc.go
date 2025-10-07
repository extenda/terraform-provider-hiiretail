// Package auth provides OAuth2 authentication capabilities for the HiiRetail IAM Terraform provider.
//
// This package implements OAuth2 client credentials flow with the Hii Retail OAuth Client Management Service (OCMS).
// It provides:
//
//   - OAuth2 discovery client with endpoint caching
//   - Enhanced authentication client with token lifecycle management
//   - Comprehensive error handling and retry logic
//   - Secure credential handling with no credential exposure in logs
//   - Thread-safe token management for concurrent operations
//
// The package integrates with the golang.org/x/oauth2 library and follows OAuth2 best practices
// for client credentials authentication with enterprise-grade error handling and security.
//
// Key components:
//   - DiscoveryClient: Handles OAuth2 endpoint discovery and caching
//   - AuthClient: Manages OAuth2 authentication and token lifecycle
//   - Error handling: Comprehensive error classification and retry logic
//   - Validation: Configuration and credential validation
//
// Example usage:
//
//	authClient, err := NewAuthClient(config)
//	if err != nil {
//		return fmt.Errorf("failed to create auth client: %w", err)
//	}
//
//	token, err := authClient.GetToken(ctx)
//	if err != nil {
//		return fmt.Errorf("failed to acquire token: %w", err)
//	}
//
//	// Use token for API calls
//	client := authClient.HTTPClient(ctx)
//
// Security considerations:
//   - All credentials are handled securely with no logging exposure
//   - TLS-only communication is enforced
//   - Token validation includes tampering detection
//   - Memory is cleared on client destruction
package auth
