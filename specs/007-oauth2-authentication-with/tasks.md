# Tasks: OAuth2 Authentication with Environment-Specific Endpoints

**Input**: Design documents from `/specs/007-oauth2-authentication-with/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Found: OAuth2 authentication for HiiRetail IAM Terraform provider
   → Extract: Go 1.21+, Terraform Plugin Framework v1.16.0, golang.org/x/oauth2 v0.30.0
2. Load optional design documents:
   → data-model.md: 5 entities (OAuth2Configuration, AuthClient, TokenCache, EndpointResolver, AuthError)
   → contracts/: 2 files (oauth2-token.md, iam-api-auth.md)
   → research.md: OAuth2 implementation decisions, security requirements
3. Generate tasks by category:
   → Setup: Go dependencies, auth package structure, linting
   → Tests: OAuth2 contract tests, IAM API tests, integration tests
   → Core: Auth entities, OAuth2 client, endpoint resolver
   → Integration: Provider integration, HTTP client, error handling
   → Polish: Demo program, documentation, security audit
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph for OAuth2 authentication system
7. Create parallel execution examples for independent components
8. Validate task completeness:
   → OAuth2 token contract has tests ✓
   → IAM API authentication contract has tests ✓
   → All 5 entities have implementation tasks ✓
   → Provider integration covered ✓
9. Return: SUCCESS (OAuth2 tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Terraform provider extension**: `internal/provider/auth/` for new authentication package
- **Tests**: `tests/auth/` for authentication-specific tests
- **Demo**: `demo/` for usage examples
- Paths follow existing HiiRetail IAM provider structure

## Phase 3.1: Setup
- [x] T001 Create auth package structure in internal/provider/auth/
- [x] T002 Update go.mod with golang.org/x/oauth2 v0.30.0 dependency  
- [x] T003 [P] Configure golangci-lint rules for OAuth2 credential exposure detection

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T004 [P] Contract test OAuth2 token endpoint in tests/auth/oauth2_token_test.go
- [ ] T005 [P] Contract test IAM API authentication in tests/auth/iam_api_auth_test.go
- [ ] T006 [P] Unit test EndpointResolver environment detection in tests/auth/discovery_test.go
- [ ] T007 [P] Unit test AuthClient token lifecycle in tests/auth/client_test.go
- [ ] T008 [P] Unit test credential validation logic in tests/auth/validation_test.go
- [ ] T009 [P] Integration test mock OAuth2 server flow in tests/auth/integration_test.go
- [ ] T010 [P] Provider test OAuth2 configuration schema in tests/provider_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T011 [P] OAuth2Configuration struct and validation in internal/provider/auth/validation.go
- [ ] T012 [P] AuthError types and classification in internal/provider/auth/errors.go
- [ ] T013 [P] EndpointResolver with tenant ID parsing in internal/provider/auth/discovery.go
- [ ] T014 TokenCache with thread-safe operations in internal/provider/auth/client.go
- [ ] T015 AuthClient OAuth2 flow implementation in internal/provider/auth/client.go  
- [ ] T016 Token refresh and expiration handling in internal/provider/auth/client.go
- [ ] T017 TLS enforcement and HTTP client configuration in internal/provider/auth/client.go

## Phase 3.4: Integration
- [ ] T018 Enhanced provider schema with OAuth2 fields in internal/provider/provider.go
- [ ] T019 buildAuthConfig function with environment variables in internal/provider/provider.go
- [ ] T020 OAuth2 client integration in provider Configure method in internal/provider/provider.go
- [ ] T021 Authenticated HTTP client creation in internal/provider/provider.go
- [ ] T022 [P] Update all resource files to use authenticated client

## Phase 3.5: Polish  
- [ ] T023 [P] OAuth2 demo program in demo/oauth2_demo.go
- [ ] T024 [P] Performance benchmarks for token acquisition in tests/auth/benchmark_test.go
- [ ] T025 [P] Security audit for credential exposure in all files
- [ ] T026 [P] Update README.md with OAuth2 configuration examples
- [ ] T027 Validate quickstart.md scenarios with real implementation

## Dependencies
- Setup (T001-T003) before everything
- Tests (T004-T010) before implementation (T011-T027)  
- T011 (OAuth2Configuration) blocks T013, T015
- T012 (AuthError) blocks T015, T016
- T013 (EndpointResolver) blocks T015
- T014-T017 (AuthClient components) must be done in order
- T018-T021 (Provider integration) must be done in order
- T022 (Resource updates) requires T021 completion
- Implementation (T011-T022) before polish (T023-T027)

## Parallel Example
```bash
# Phase 3.2: Launch all test files together (different files, independent):
Task: "Contract test OAuth2 token endpoint in tests/auth/oauth2_token_test.go"
Task: "Contract test IAM API authentication in tests/auth/iam_api_auth_test.go"  
Task: "Unit test EndpointResolver environment detection in tests/auth/discovery_test.go"
Task: "Unit test AuthClient token lifecycle in tests/auth/client_test.go"
Task: "Unit test credential validation logic in tests/auth/validation_test.go"
Task: "Integration test mock OAuth2 server flow in tests/auth/integration_test.go"
Task: "Provider test OAuth2 configuration schema in tests/provider_test.go"

# Phase 3.3: Launch independent entity files together:
Task: "OAuth2Configuration struct and validation in internal/provider/auth/validation.go"
Task: "AuthError types and classification in internal/provider/auth/errors.go"  
Task: "EndpointResolver with tenant ID parsing in internal/provider/auth/discovery.go"

# Phase 3.5: Launch polish tasks together:
Task: "OAuth2 demo program in demo/oauth2_demo.go"
Task: "Performance benchmarks for token acquisition in tests/auth/benchmark_test.go"
Task: "Security audit for credential exposure in all files"
Task: "Update README.md with OAuth2 configuration examples"
```

## OAuth2-Specific Implementation Notes

### Security Requirements (ALL TASKS)
- Mark client_secret as sensitive in Terraform schema
- Redact credentials from all log outputs  
- Enforce TLS 1.2+ for all OAuth2 communications
- Validate server certificates
- Clear credentials from memory on disposal

### Thread Safety Requirements (AuthClient tasks)
- Use sync.RWMutex for token cache access
- Ensure concurrent token refresh safety
- Protect configuration changes during authentication

### Environment Detection Logic (T013)
- Pattern matching: test|dev|staging (case-insensitive)
- Environment variable override: HIIRETAIL_FORCE_TEST_ENV
- Default to Live environment for security

### Mock Server Support (T009, T023)
- Environment variable URL overrides
- Validation bypass in mock mode
- Test-only URL acceptance

## Contract Test Details

### T004: OAuth2 Token Endpoint Tests
```go
// tests/auth/oauth2_token_test.go
func TestOAuth2TokenEndpoint_ValidCredentials(t *testing.T) // → 200 OK with token
func TestOAuth2TokenEndpoint_InvalidCredentials(t *testing.T) // → 401 Unauthorized  
func TestOAuth2TokenEndpoint_MissingClientID(t *testing.T) // → 400 Bad Request
func TestOAuth2TokenEndpoint_WrongGrantType(t *testing.T) // → 400 Bad Request
```

### T005: IAM API Authentication Tests  
```go
// tests/auth/iam_api_auth_test.go
func TestIAMAPI_ValidBearerToken(t *testing.T) // → 200 OK with resource data
func TestIAMAPI_ExpiredToken(t *testing.T) // → 401 Unauthorized
func TestIAMAPI_InvalidTenant(t *testing.T) // → 404 Not Found
func TestIAMAPI_EndpointResolution(t *testing.T) // → Correct URL based on tenant
```

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests (T004-T005)
- [x] All entities have model tasks (T011-T016) 
- [x] All tests come before implementation (T004-T010 → T011-T027)
- [x] Parallel tasks truly independent (different files)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] OAuth2 security requirements specified for all tasks
- [x] Environment detection and mock server support covered
- [x] Provider integration maintains existing functionality

## Task Execution Success Criteria

**Phase 3.2 Complete**: All tests written and failing, no implementation exists
**Phase 3.3 Complete**: Core OAuth2 authentication working, tests passing
**Phase 3.4 Complete**: Provider integrated with OAuth2, existing resources unaffected  
**Phase 3.5 Complete**: Demo working, documentation updated, security validated

**Overall Success**: OAuth2 client credentials flow working with automatic environment detection, secure credential handling, and comprehensive test coverage.