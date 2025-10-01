# Research: OAuth2 Client Credentials Authentication

**Feature**: Correctly Handle Client Credentials  
**Date**: October 1, 2025  
**Status**: Phase 0 Research Complete

## Research Questions

### 1. Hii Retail OAuth Client Management Service (OCMS) Integration

**Question**: How does the OCMS OAuth2 flow work and what are the endpoint requirements?

**Findings**:
- **Discovery Endpoint**: https://auth.retailsvc.com/.well-known/openid-configuration
- **Service**: Hii Retail OAuth Client Management Service (OCMS)
- **Flow**: OAuth2 Client Credentials Grant (RFC 6749)
- **Documentation**: https://developer.hiiretail.com/docs/ocms/public/concepts/oauth2-authentication/

**Expected Endpoints from Discovery**:
- `token_endpoint`: For token acquisition and refresh
- `authorization_endpoint`: For OAuth2 flows (not needed for client credentials)
- `jwks_uri`: For token validation (if JWT tokens are used)
- `issuer`: Token issuer identification
- `supported_grant_types`: Should include "client_credentials"

### 2. Current Provider Authentication Issues

**Question**: What are the current authentication problems in the provider?

**Findings from Codebase Analysis**:
- Current implementation uses golang.org/x/oauth2 clientcredentials package
- Token endpoint is hardcoded as `fmt.Sprintf("%s/oauth2/token", baseUrl)`
- No discovery endpoint usage - relies on manual endpoint construction
- Limited error handling for token acquisition failures
- No retry logic for transient authentication failures
- Token caching handled by oauth2 library but no custom logic

**Current Code Location**: `internal/provider/provider.go` Configure() method

### 3. OAuth2 Best Practices for Terraform Providers

**Question**: What are the HashiCorp recommendations for OAuth2 in providers?

**Findings**:
- Use `context.Background()` for OAuth2 operations to avoid cancellation
- Set reasonable timeouts for token acquisition (10s recommended)
- Cache HTTP client with OAuth2 configuration
- Handle token refresh automatically via oauth2.TokenSource
- Mark client_secret as sensitive in schema
- Support both configuration block and environment variables
- Provide clear error messages for authentication failures

### 4. Token Lifecycle Management

**Question**: How should token acquisition, caching, and refresh be handled?

**Findings**:
- **Acquisition**: Use oauth2/clientcredentials.Config for initial token
- **Caching**: golang.org/x/oauth2 handles automatic caching and refresh
- **Refresh**: Automatic when token expires during API calls
- **Error Handling**: Distinguish between network errors and credential errors
- **Concurrency**: oauth2.TokenSource is thread-safe for concurrent operations
- **Context**: Use background context to prevent premature cancellation

### 5. Error Scenarios and Recovery

**Question**: What authentication error scenarios need handling?

**Critical Scenarios**:
1. **Invalid Credentials**: 401/403 responses from token endpoint
2. **Network Failures**: Connection timeouts, DNS resolution issues
3. **Service Unavailable**: 5xx responses from OCMS
4. **Token Revocation**: External token revocation during operations
5. **Malformed Configuration**: Invalid URLs, missing parameters
6. **Rate Limiting**: 429 responses from token endpoint

**Recovery Strategies**:
- Exponential backoff for retryable errors (5xx, network)
- Immediate failure for credential errors (4xx except 429)
- Clear error messages with troubleshooting guidance
- Validation of configuration before first API call

## Technical Decisions

### 1. Discovery Endpoint Usage
**Decision**: Implement OAuth2 discovery protocol for dynamic endpoint configuration
**Rationale**: 
- Follows OAuth2 best practices (RFC 8414)
- Provides resilience against endpoint changes
- Enables proper endpoint validation
- Supports future OAuth2 extensions

### 2. Token Management Strategy
**Decision**: Enhance existing oauth2.ClientCredentials with discovery and error handling
**Rationale**:
- Maintains compatibility with existing provider structure
- Leverages proven oauth2 library for token lifecycle
- Adds necessary discovery and error handling layers
- Preserves thread-safety for concurrent operations

### 3. Configuration Enhancement
**Decision**: Maintain current configuration schema with enhanced validation
**Rationale**:
- No breaking changes to existing configurations
- Validates token endpoint URL against discovery
- Provides clear error messages for misconfigurations
- Supports both explicit token_url and discovery-based configuration

### 4. Error Handling Approach
**Decision**: Implement layered error handling with clear user guidance
**Rationale**:
- Distinguishes between configuration, network, and credential errors
- Provides actionable error messages with troubleshooting steps
- Implements appropriate retry strategies for different error types
- Maintains security by not exposing sensitive information

## Implementation Strategy

### Phase 1: Discovery Integration
1. Add OAuth2 discovery client to fetch endpoint configuration
2. Validate discovered endpoints against expected OAuth2 flows
3. Cache discovery results to minimize network calls
4. Implement fallback to manual configuration if discovery fails

### Phase 2: Enhanced Token Management
1. Integrate discovery endpoints with existing oauth2.ClientCredentials
2. Add comprehensive error handling for token acquisition
3. Implement retry logic with exponential backoff
4. Enhance logging while preserving credential security

### Phase 3: Validation and Testing
1. Add configuration validation for OAuth2 parameters
2. Implement unit tests for discovery and token management
3. Add integration tests with actual OCMS endpoints
4. Create acceptance tests for various configuration scenarios

### Phase 4: Documentation and Examples
1. Update provider documentation with OAuth2 configuration examples
2. Add troubleshooting guide for authentication issues
3. Provide migration examples for enhanced credential handling
4. Document security best practices for credential management

## Risks and Mitigations

### Risk 1: Discovery Endpoint Availability
**Impact**: Provider fails if discovery endpoint is unreachable
**Mitigation**: Implement fallback to manual endpoint configuration

### Risk 2: Breaking Changes
**Impact**: Existing provider configurations stop working
**Mitigation**: Maintain backward compatibility with current configuration schema

### Risk 3: Token Refresh During Operations
**Impact**: Long-running Terraform operations fail if token expires
**Mitigation**: oauth2 library handles automatic refresh; add error recovery for edge cases

### Risk 4: Credential Exposure
**Impact**: Client secrets exposed in logs or error messages
**Mitigation**: Implement secure error handling and audit all logging statements

## Success Criteria

1. **Discovery Integration**: Provider successfully discovers OAuth2 endpoints from OCMS
2. **Token Management**: Automatic token acquisition, refresh, and error handling
3. **Error Handling**: Clear, actionable error messages for all authentication scenarios
4. **Security Compliance**: No credential exposure in logs or debug output
5. **Performance**: <500ms token acquisition, efficient token reuse
6. **Compatibility**: No breaking changes to existing provider configurations
7. **Testing**: Comprehensive test coverage for all authentication scenarios

## References

- [Hii Retail OCMS Documentation](https://developer.hiiretail.com/docs/ocms/public/concepts/oauth2-authentication/)
- [OAuth2 Discovery Specification (RFC 8414)](https://tools.ietf.org/html/rfc8414)
- [OAuth2 Client Credentials Grant (RFC 6749)](https://tools.ietf.org/html/rfc6749#section-4.4)
- [HashiCorp Provider Development](https://developer.hashicorp.com/terraform/plugin)
- [golang.org/x/oauth2 Documentation](https://pkg.go.dev/golang.org/x/oauth2)

---

**Research Status**: ✅ Complete  
**Next Phase**: Phase 1 - Design (contracts, data-model, quickstart)