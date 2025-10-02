# Phase 0: Research & Technical Decisions
**Feature**: OAuth2 Authentication with Environment-Specific Endpoints  
**Date**: October 1, 2025

## Research Tasks Completed

### 1. OAuth2 Client Credentials Flow for Terraform Providers

**Decision**: Use golang.org/x/oauth2 with client credentials configuration  
**Rationale**: 
- Official Go OAuth2 library with Terraform provider ecosystem adoption
- Built-in token refresh and expiration handling
- Thread-safe for concurrent Terraform operations
- Supports custom token endpoints and HTTP clients

**Alternatives considered**:
- Custom OAuth2 implementation: Rejected due to security complexity and maintenance burden
- Third-party OAuth2 libraries: Rejected due to additional dependencies and unknown security posture

### 2. Tenant ID Parsing Strategy

**Decision**: Pattern-based parsing with fallback to Live environment  
**Rationale**:
- Test tenant IDs typically contain patterns like "test", "dev", "staging" in subdomain or prefix
- Conservative approach defaults to Live (production) to prevent accidental cross-environment operations
- Configurable override for edge cases via environment variables

**Pattern Detection Logic**:
```
Test Tenant Indicators:
- Contains "test" (case-insensitive)
- Contains "dev" (case-insensitive) 
- Contains "staging" (case-insensitive)
- Matches pattern: *-test-*, *-dev-*, *-staging-*
- Environment variable override: HIIRETAIL_FORCE_TEST_ENV=true

Default: Live Tenant (iam-api.retailsvc.com)
```

**Alternatives considered**:
- Explicit environment configuration: Rejected to maintain simplicity and reduce configuration burden
- API-based environment detection: Rejected due to additional API calls and potential circular dependencies

### 3. Mock Server Override for Testing

**Decision**: Environment variable based URL override with validation  
**Rationale**:
- Allows seamless testing without code changes
- Validates URLs to prevent accidental production usage in tests
- Supports both authentication and API endpoint overrides

**Override Mechanism**:
```
Environment Variables:
- HIIRETAIL_AUTH_URL: Override auth.retailsvc.com (testing only)
- HIIRETAIL_API_URL: Override iam-api.retailsvc.com/iam-api.retailsvc-test.com (testing only)
- HIIRETAIL_MOCK_MODE: Enable mock server mode with validation bypass
```

**Alternatives considered**:
- Configuration file based overrides: Rejected due to additional file management complexity
- Command-line flags: Rejected as Terraform providers don't support custom CLI arguments

### 4. Token Management and Caching Strategy

**Decision**: In-memory token caching with automatic refresh  
**Rationale**:
- Reduces authentication overhead for multiple resource operations
- Thread-safe caching for concurrent Terraform operations
- Automatic refresh before expiration prevents authentication failures
- No persistent storage reduces security attack surface

**Implementation Details**:
- Token cached in memory per provider instance
- Refresh 60 seconds before expiration
- Fallback to new authentication on refresh failure
- Thread-safe access using sync.RWMutex

**Alternatives considered**:
- File-based token caching: Rejected due to security risks and cleanup complexity
- No caching: Rejected due to performance impact on multi-resource operations

### 5. Error Handling and Retry Logic

**Decision**: Exponential backoff with authentication-specific error classification  
**Rationale**:
- Distinguishes between transient network errors and permanent authentication failures
- Avoids overwhelming authentication servers during outages
- Provides clear error messages for debugging

**Error Categories**:
- **Retryable**: Network timeouts, 5xx server errors, rate limiting (429)
- **Non-retryable**: Invalid credentials (401), malformed requests (400), not found (404)
- **Authentication-specific**: Token expired, invalid client credentials, malformed tenant ID

**Alternatives considered**:
- Fixed retry intervals: Rejected due to potential server overload during outages
- No retry logic: Rejected due to poor user experience with transient failures

### 6. Security Considerations

**Decision**: Comprehensive credential protection and TLS enforcement  
**Rationale**:
- Prevents credential exposure in logs, debug output, or error messages
- Enforces TLS for all authentication communications
- Validates TLS certificates to prevent man-in-the-middle attacks

**Security Measures**:
- Mark all credential fields as sensitive in Terraform schema
- Redact credentials from all log outputs
- Enforce TLS 1.2+ for OAuth2 communications
- Validate server certificates
- Clear credentials from memory on provider disposal

**Alternatives considered**:
- Basic HTTP authentication: Rejected due to security vulnerabilities
- Custom encryption: Rejected in favor of proven TLS standards

## Technical Decisions Summary

| Component | Decision | Primary Rationale |
|-----------|----------|-------------------|
| OAuth2 Library | golang.org/x/oauth2 | Official library, proven in Terraform ecosystem |
| Tenant Detection | Pattern-based parsing | Automatic environment detection without configuration |
| Mock Testing | Environment variable override | Simple testing without code changes |
| Token Caching | In-memory with auto-refresh | Performance with security |
| Error Handling | Exponential backoff | Reliability with server protection |
| Security | TLS enforcement + credential protection | Comprehensive security posture |

## Dependencies Resolved

- **golang.org/x/oauth2 v0.30.0**: OAuth2 client implementation
- **sync package**: Thread-safe token caching
- **net/http**: HTTP client with TLS configuration
- **time package**: Token expiration and retry timing
- **context package**: Request cancellation and timeouts

All dependencies are part of Go standard library or well-established packages already used in the Terraform ecosystem.

## Implementation Ready

All research tasks completed. No remaining NEEDS CLARIFICATION items. Ready to proceed to Phase 1 design.