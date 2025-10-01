# Tasks: OAuth2 Client Credentials Authentication

**Input**: Design documents from `/specs/006-correctly-handle-client/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Go 1.21+ with Terraform Plugin Framework v1.16.0
   → Libraries: golang.org/x/oauth2, HashiCorp Terraform Plugin Framework
   → Structure: Single binary Terraform provider with OAuth2 client integration
2. Load optional design documents:
   → data-model.md: Extract entities → OAuth2 models, auth client, error handling
   → contracts/: oauth2_authentication.md → OAuth2 contract tests  
   → research.md: Extract decisions → discovery protocol, token management
3. Generate tasks by category:
   → Setup: auth package structure, dependencies, OAuth2 configuration
   → Tests: OAuth2 discovery tests, authentication tests, provider tests
   → Core: discovery client, auth client, provider configuration updates
   → Integration: provider OAuth2 integration, error handling
   → Polish: unit tests, integration tests, documentation
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → OAuth2 contract has tests ✅
   → OAuth2 entities have models ✅ 
   → All endpoints implemented ✅
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Terraform Provider**: `internal/provider/` at repository root
- **Auth Package**: `internal/provider/auth/` for OAuth2 components
- **Tests**: Package-level tests alongside implementation files
- Paths shown below assume Terraform provider structure per plan.md

## Phase 3.1: Setup

- [x] **T001** Create OAuth2 authentication package structure
  - Create directory: `internal/provider/auth/`
  - Initialize package with proper Go module imports
  - Set up basic package documentation

- [x] **T002** [P] Update Go module dependencies for OAuth2 enhancement
  - Verify golang.org/x/oauth2 version compatibility
  - Update go.mod with required OAuth2 dependencies
  - Run `go mod tidy` to clean up dependencies

- [x] **T003** [P] Configure OAuth2-specific linting and validation rules
  - Add security-focused linting rules for credential handling
  - Configure static analysis for OAuth2 best practices
  - Set up pre-commit hooks for credential exposure detection

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

- [x] **T004** [P] Create OAuth2 discovery client contract tests
  - File: `internal/provider/auth/discovery_test.go`
  - Test OAuth2 discovery endpoint parsing per contract
  - Test discovery response validation and error handling
  - Test discovery result caching behavior

- [x] **T005** [P] Create OAuth2 authentication client contract tests  
  - File: `internal/provider/auth/client_test.go`
  - Test OAuth2 client credentials flow per contract
  - Test token acquisition, refresh, and expiration handling
  - Test authentication error scenarios and retry logic

- [x] **T006** [P] Create provider OAuth2 configuration tests
  - File: `internal/provider/provider_oauth2_test.go`
  - Test enhanced provider Configure() method with OAuth2
  - Test configuration validation for OAuth2 parameters
  - Test environment variable support for OAuth2 credentials

- [x] **T007** [P] Create OAuth2 integration test scenarios
  - File: `internal/provider/auth/integration_test.go`
  - Test end-to-end OAuth2 flow with mock OCMS endpoints
  - Test concurrent token acquisition and refresh scenarios
  - Test long-running operations with token lifecycle

## Phase 3.3: Core Implementation

- [x] **T008** Implement OAuth2 discovery client
  - File: `internal/provider/auth/discovery.go`
  - Implement OIDCDiscoveryResponse struct per data model
  - Implement DiscoveryClient with caching per technical decisions
  - Add discovery endpoint validation and fallback logic

- [x] **T009** Implement OAuth2 authentication error models
  - File: `internal/provider/auth/errors.go` 
  - Implement AuthError types per data model
  - Add error classification (credential, network, server errors)
  - Implement retry logic with exponential backoff

- [ ] **T010** Implement enhanced OAuth2 authentication client
  - File: `internal/provider/auth/client.go`
  - Implement AuthClient struct per data model
  - Integrate oauth2.ClientCredentials with discovery client
  - Add token lifecycle management and validation

- [ ] **T011** Implement OAuth2 configuration validation
  - File: `internal/provider/auth/validation.go`
  - Implement configuration validation rules per data model
  - Add URL validation for OAuth2 endpoints
  - Implement credential format validation

- [ ] **T012** Update provider configuration schema for OAuth2
  - File: `internal/provider/provider.go` (enhance existing)
  - Add OAuth2-specific configuration parameters to schema
  - Mark client_secret as sensitive field
  - Add timeout and retry configuration options

- [ ] **T013** Enhance provider Configure() method with OAuth2
  - File: `internal/provider/provider.go` (enhance existing)
  - Integrate OAuth2 discovery and authentication client
  - Add comprehensive OAuth2 credential validation
  - Implement OAuth2 error handling with clear user messages

## Phase 3.4: Integration

- [ ] **T014** [P] Integrate OAuth2 authentication with existing APIClient
  - File: `internal/provider/provider.go` (enhance existing)
  - Update APIClient struct to include OAuth2 AuthClient
  - Ensure backward compatibility with existing resource implementations
  - Add OAuth2 authentication to HTTP client factory

- [ ] **T015** [P] Add OAuth2 retry logic for API operations
  - File: `internal/provider/auth/retry.go`
  - Implement retry wrapper for API calls with token refresh
  - Add automatic token refresh on 401 authentication errors
  - Implement exponential backoff for retryable errors

- [ ] **T016** [P] Update provider error handling for OAuth2 scenarios
  - File: `internal/provider/auth/error_handler.go`
  - Add OAuth2-specific error mapping and user-friendly messages
  - Implement troubleshooting guidance for common OAuth2 issues
  - Add security audit for error message content (no credential exposure)

## Phase 3.5: Polish

- [ ] **T017** [P] Add comprehensive unit tests for OAuth2 components
  - Files: `internal/provider/auth/*_test.go`
  - Achieve >90% test coverage for OAuth2 authentication logic
  - Add edge case testing for token expiration and refresh
  - Test concurrent authentication scenarios

- [ ] **T018** [P] Create OAuth2 integration tests with real OCMS endpoints
  - File: `internal/provider/auth/ocms_integration_test.go`
  - Test against real https://auth.retailsvc.com endpoints
  - Validate OAuth2 discovery and token acquisition flows
  - Add integration test environment setup documentation

- [ ] **T019** [P] Add provider acceptance tests for OAuth2 configuration
  - File: `internal/provider/provider_acceptance_test.go`
  - Test provider configuration with various OAuth2 parameter combinations
  - Test environment variable configuration scenarios
  - Validate OAuth2 authentication in full Terraform workflow

- [ ] **T020** [P] Update provider documentation for OAuth2 authentication
  - File: `README.md` and docs directory
  - Add OAuth2 configuration examples and best practices
  - Document environment variable support for OAuth2 credentials
  - Create troubleshooting guide for OAuth2 authentication issues

- [ ] **T021** [P] Add OAuth2 performance benchmarks and monitoring
  - File: `internal/provider/auth/benchmark_test.go`
  - Benchmark token acquisition time (<500ms requirement)
  - Test concurrent authentication performance
  - Add metrics for token cache hit rates

- [ ] **T022** [P] Security audit for OAuth2 implementation
  - Review all OAuth2 code for credential exposure risks
  - Validate TLS-only communication enforcement
  - Audit logging statements for sensitive information
  - Add security testing scenarios

## Dependency Graph

```
T001 (setup) → T004,T005,T006,T007 (tests) → T008,T009,T010,T011 (core)
                                                ↓
T002,T003 (deps) → T012,T013 (provider) → T014,T015,T016 (integration)
                                                ↓  
                                         T017,T018,T019,T020,T021,T022 (polish)
```

## Parallel Execution Examples

### Phase 3.2 Tests (All Parallel)
```bash
# Run all test creation tasks in parallel
Task T004 & Task T005 & Task T006 & Task T007 & wait
```

### Phase 3.3 Core (Sequential + Some Parallel)
```bash
# T008 must complete first (discovery client)
Task T008
# T009,T010,T011 can run in parallel after T008
Task T009 & Task T010 & Task T011 & wait
# T012,T013 must be sequential (same file)
Task T012
Task T013
```

### Phase 3.5 Polish (All Parallel)
```bash
# All polish tasks can run in parallel
Task T017 & Task T018 & Task T019 & Task T020 & Task T021 & Task T022 & wait
```

## File Impact Summary

### New Files Created
- `internal/provider/auth/discovery.go` - OAuth2 discovery client
- `internal/provider/auth/client.go` - Enhanced authentication client  
- `internal/provider/auth/errors.go` - OAuth2 error handling
- `internal/provider/auth/validation.go` - Configuration validation
- `internal/provider/auth/retry.go` - Retry logic for API operations
- `internal/provider/auth/error_handler.go` - Error mapping and messages

### Existing Files Enhanced  
- `internal/provider/provider.go` - OAuth2 configuration and integration
- Multiple test files - Comprehensive OAuth2 testing

### Testing Files
- `internal/provider/auth/*_test.go` - Unit tests for OAuth2 components
- `internal/provider/auth/integration_test.go` - OCMS integration tests
- `internal/provider/provider_oauth2_test.go` - Provider OAuth2 tests
- `internal/provider/auth/benchmark_test.go` - Performance benchmarks

## Success Criteria

- [ ] OAuth2 discovery working with https://auth.retailsvc.com/.well-known/openid-configuration
- [ ] OAuth2 client credentials flow integrated with existing provider
- [ ] Token acquisition <500ms, validation <100ms performance targets met
- [ ] Comprehensive error handling with user-friendly messages
- [ ] Security audit passed (no credential exposure)
- [ ] >90% test coverage for OAuth2 components
- [ ] Backward compatibility maintained with existing configurations
- [ ] Documentation updated with OAuth2 examples and troubleshooting

## Validation Commands

```bash
# Run all OAuth2 tests
go test ./internal/provider/auth -v

# Test OAuth2 provider configuration
go test ./internal/provider -v -run TestProvider.*OAuth2

# Integration test with real OCMS
INTEGRATION_TEST=true go test ./internal/provider/auth -v -run TestOCMS

# Build and test provider
make build && terraform init && terraform plan

# Security validation
make security-audit

# Performance benchmarks  
go test -bench=. ./internal/provider/auth
```

---

**Tasks Status**: ✅ Ready for execution  
**Total Tasks**: 22 tasks across 5 phases  
**Parallel Tasks**: 12 tasks can run in parallel (marked with [P])  
**Implementation Time**: Estimated 6-8 hours for core implementation + testing