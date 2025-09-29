# Research: IAM Custom Role Resource Testing Implementation

**Date**: September 28, 2025  
**Feature**: Add Comprehensive Tests for IAM Custom Role Resource

## Research Findings

### Testing Framework Architecture

**Decision**: Use Terraform Plugin Framework Testing with testify assertions
**Rationale**: 
- Aligns with existing codebase patterns (seen in iam_group resource tests)
- Provides comprehensive testing capabilities for Terraform resources
- Supports unit, integration, and acceptance testing patterns
- Compatible with Go testing conventions and CI/CD pipelines

**Alternatives considered**:
- Plain Go testing: Insufficient for Terraform-specific resource testing
- Ginkgo/Gomega: Overcomplicated for provider testing needs
- Custom testing framework: Unnecessary reinvention

### Resource Implementation Patterns

**Decision**: Follow established iam_group resource implementation pattern
**Rationale**:
- Proven working implementation with 22% baseline coverage
- Consistent with provider architecture and authentication flow
- Reuses existing mock server infrastructure and test utilities
- Maintains code consistency across resource implementations

**Alternatives considered**:
- Generated-only implementation: Lacks business logic and proper error handling
- Different testing approach: Would break consistency with existing resources
- Custom implementation from scratch: Unnecessary when proven pattern exists

### Test Coverage Strategy

**Decision**: Implement comprehensive test suite with multiple test types
**Rationale**:
- Constitutional requirement for thorough testing and validation
- Critical IAM resource requires extensive validation
- Performance characteristics important for enterprise usage
- Concurrent access patterns need validation for Terraform parallelism

**Test Types Required**:
1. **Unit Tests**: Individual method validation, schema testing, validation logic
2. **Integration Tests**: API client integration, authentication flow, error handling
3. **Contract Tests**: API contract validation, request/response schemas
4. **Acceptance Tests**: End-to-end Terraform lifecycle operations  
5. **Benchmark Tests**: Performance characteristics with maximum permissions
6. **Concurrent Tests**: Race condition handling, parallel resource operations

### API Client Integration

**Decision**: Leverage existing OAuth2 client infrastructure from provider
**Rationale**:
- Existing authentication implementation is secure and functional
- Consistent with provider authentication approach
- Reuses configured HTTP client with proper token management
- Maintains security standards required by constitution

**Integration Points**:
- Provider APIClient structure for HTTP client access
- OAuth2 client credentials flow for authentication
- Base URL and tenant ID configuration from provider
- Error handling and retry logic patterns

### Mock Server Infrastructure

**Decision**: Extend existing mock server from testutils package
**Rationale**:
- Proven working mock server implementation exists
- Supports isolated testing without external dependencies
- Consistent with existing test infrastructure
- Enables fast test execution and reliable CI/CD

**Mock Server Capabilities**:
- Custom role CRUD operations simulation
- Permission validation and limit enforcement
- Attribute constraint validation
- Error scenario simulation (API failures, timeouts, etc.)

### Permission Validation Logic

**Decision**: Implement comprehensive permission pattern validation
**Rationale**:
- Critical business rule: {systemPrefix}.{resource}.{action} pattern
- Performance implications: up to 500 POS permissions allowed
- Security implications: proper permission structure validation
- User experience: clear error messages for malformed permissions

**Validation Requirements**:
- Pattern matching: `^[a-z][-a-z]{2}\\.[a-z][-a-z]{1,15}\\.[a-z][-a-z]{1,15}$`
- Limits: 100 general permissions, 500 POS permissions
- Attribute constraints: max 10 properties, 40 char keys, 256 char values
- Real-time validation during Terraform operations

### Error Handling Patterns

**Decision**: Implement comprehensive error mapping following Terraform conventions
**Rationale**:
- Terraform users need clear, actionable error messages
- API errors must be properly translated to Terraform diagnostics
- Retry logic needed for transient failures
- Graceful degradation for partial failures

**Error Categories**:
- Validation errors: Client-side validation with immediate feedback
- API errors: HTTP status code mapping to Terraform diagnostics
- Authentication errors: OAuth2 token refresh and credential validation
- Network errors: Retry logic with exponential backoff
- Business logic errors: Permission limits, attribute constraints, etc.

## Technical Decisions Summary

| Area | Decision | Rationale |
|------|----------|-----------|
| Testing Framework | Terraform Plugin Framework Testing + testify | Proven, comprehensive, consistent |
| Implementation Pattern | Follow iam_group resource pattern | Working implementation, consistency |
| Test Coverage | Multi-type comprehensive suite | Constitutional requirements, critical resource |
| Authentication | Reuse existing OAuth2 infrastructure | Security, consistency, proven |
| Mock Infrastructure | Extend existing testutils mock server | Fast, reliable, consistent |
| Validation | Comprehensive permission pattern validation | Business rules, security, UX |
| Error Handling | Terraform-native error mapping | User experience, debugging |

## Implementation Readiness

✅ **Technical Foundation**: All dependencies and patterns identified  
✅ **Architecture Decisions**: Clear implementation approach defined  
✅ **Testing Strategy**: Comprehensive coverage plan established  
✅ **Integration Points**: Authentication and client infrastructure mapped  
✅ **Validation Logic**: Permission and constraint validation defined  
✅ **Error Handling**: Comprehensive error mapping strategy planned  

**Status**: Ready for Phase 1 Design & Contracts