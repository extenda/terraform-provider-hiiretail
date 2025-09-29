# Tasks: IAM Custom Role Resource Testing Implementation

**Feature**: Add Comprehensive Tests for IAM Custom Role Resource  
**Branch**: `004-ensure-that-the`  
**Date**: September 28, 2025

## Task Execution Order

**Sequential Tasks**: Execute in numerical order
**Parallel Tasks [P]**: Can execute simultaneously with same-numbered tasks
**Dependencies**: Must complete previous phase before starting next

---

## Phase 1: Setup & Registration (Tasks 1-5)

### T001: Register Custom Role Resource with Provider ✅
**File**: `internal/provider/provider.go`
**Description**: Add iam_custom_role resource to provider Resources() method
**Dependencies**: None
**Priority**: Critical Path
**Status**: COMPLETED - Resource registered and provider builds successfully

### T002: Create Resource Package Structure [P] ✅
**Files**: `internal/provider/resource_iam_custom_role/`
**Description**: Create enhanced resource implementation file alongside generated schema
**Dependencies**: None
**Status**: COMPLETED - Resource and test files created, basic tests passing

### T003: Extend Mock Server for Custom Roles [P] ✅
**File**: `internal/provider/testutils/mock_server.go`
**Description**: Add custom role CRUD endpoints to existing mock server
**Dependencies**: None
**Status**: COMPLETED - All CRUD endpoints implemented with validation and error scenarios

### T004: Create Contract Test Structure [P]
**File**: `internal/provider/resource_iam_custom_role/custom_role_contract_test.go`
**Description**: Create failing contract tests based on test-contracts.md
**Dependencies**: None
**Test Categories**: Schema validation, CRUD operations, permission validation

### T005: Create Acceptance Test Structure [P]
**File**: `acceptance_tests/custom_role_resource_test.go`
**Description**: Create Terraform acceptance test file
**Dependencies**: None
**Test Scenarios**: Basic CRUD, permission limits, attribute validation

---

## Phase 2: Core Resource Implementation (Tasks 6-15)

### T006: Implement Resource Metadata and Schema
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Implement Metadata() and Schema() methods using generated schema
**Dependencies**: T001, T002
**Implementation**: Wrap generated schema with business logic

### T007: Implement Resource Configure Method
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Add Configure() method for OAuth2 client access
**Dependencies**: T006
**Pattern**: Follow iam_group resource OAuth2 integration pattern

### T008: Implement Create Operation
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Implement Create() method with API integration
**Dependencies**: T007
**Features**: Request validation, API call, error handling, state setting

### T009: Implement Read Operation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Implement Read() method for state refresh
**Dependencies**: T007
**Features**: API call, state mapping, error handling

### T010: Implement Update Operation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Implement Update() method for resource changes
**Dependencies**: T007
**Features**: Diff detection, API call, partial updates, state refresh

### T011: Implement Delete Operation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Implement Delete() method for resource removal
**Dependencies**: T007
**Features**: API call, error handling, state cleanup

### T012: Implement Import State Method [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Add ImportState() method for existing resource import
**Dependencies**: T009
**Features**: ID-based import, state population

### T013: Add Permission Validation Logic
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Implement permission pattern validation and limits
**Dependencies**: T008, T010
**Validation Rules**:
- Pattern: `^[a-z][-a-z]{2}\\.[a-z][-a-z]{1,15}\\.[a-z][-a-z]{1,15}$`
- Limits: 100 general, 500 POS permissions

### T014: Add Attribute Constraints Validation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Implement attribute object validation
**Dependencies**: T008, T010
**Constraints**: Max 10 properties, 40 char keys, 256 char values

### T015: Implement Error Handling and Retry Logic
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Add comprehensive error mapping and retry logic
**Dependencies**: T008-T011
**Features**: HTTP status mapping, retry with backoff, diagnostics

---

## Phase 3: Unit Testing (Tasks 16-25)

### T016: Unit Tests - Schema Validation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test resource schema and field validation
**Dependencies**: T006
**Test Cases**: Required fields, optional fields, validation rules

### T017: Unit Tests - Create Operation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test Create method with various scenarios
**Dependencies**: T008
**Test Cases**: Valid role, invalid permissions, API errors

### T018: Unit Tests - Read Operation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test Read method state management
**Dependencies**: T009
**Test Cases**: Existing role, not found, API errors

### T019: Unit Tests - Update Operation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test Update method diff handling
**Dependencies**: T010
**Test Cases**: Permission changes, attribute updates, partial updates

### T020: Unit Tests - Delete Operation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test Delete method cleanup
**Dependencies**: T011
**Test Cases**: Successful deletion, not found, API errors

### T021: Unit Tests - Permission Validation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test permission pattern and limit validation
**Dependencies**: T013
**Test Cases**: Valid patterns, invalid patterns, limit enforcement

### T022: Unit Tests - Attribute Validation [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test attribute constraint validation
**Dependencies**: T014
**Test Cases**: Size limits, key/value constraints

### T023: Unit Tests - Error Handling [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test error scenarios and retry logic
**Dependencies**: T015
**Test Cases**: Network errors, API errors, retry behavior

### T024: Unit Tests - OAuth2 Integration [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test authentication and client configuration
**Dependencies**: T007
**Test Cases**: Token handling, client setup, auth errors

### T025: Unit Tests - Import Functionality [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Test state import scenarios
**Dependencies**: T012
**Test Cases**: Valid import, invalid ID, state mapping

---

## Phase 4: Integration Testing (Tasks 26-30)

### T026: Integration Tests - API Client
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_integration_test.go`
**Description**: Test API client integration with mock server
**Dependencies**: T003, T008-T011
**Test Scenarios**: Full CRUD cycle, authentication flow

### T027: Integration Tests - Provider Integration [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_integration_test.go`
**Description**: Test resource registration and provider coupling
**Dependencies**: T001, T007
**Test Scenarios**: Resource discovery, configuration inheritance

### T028: Integration Tests - Concurrent Operations [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_concurrent_test.go`
**Description**: Test concurrent access and race conditions
**Dependencies**: T008-T011
**Test Scenarios**: Parallel operations, state consistency

### T029: Contract Tests - API Compliance
**File**: `internal/provider/resource_iam_custom_role/custom_role_contract_test.go`
**Description**: Implement contract tests from contracts/test-contracts.md
**Dependencies**: T004, T008-T011
**Test Coverage**: All contract scenarios must pass

### T030: Acceptance Tests - Terraform Lifecycle
**File**: `acceptance_tests/custom_role_resource_test.go`
**Description**: Implement Terraform acceptance tests
**Dependencies**: T005, T001, T008-T012
**Test Scenarios**: Plan/apply/destroy cycles, import, updates

---

## Phase 5: Performance & Benchmarks (Tasks 31-35)

### T031: Benchmark Tests - CRUD Operations [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_benchmark_test.go`
**Description**: Create performance benchmarks for all operations
**Dependencies**: T008-T011
**Benchmarks**: Create, Read, Update, Delete performance

### T032: Benchmark Tests - Large Permission Sets [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_benchmark_test.go`
**Description**: Test performance with maximum permissions (500)
**Dependencies**: T013, T031
**Scenarios**: 500 POS permissions, 100 general permissions

### T033: Benchmark Tests - Validation Performance [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_benchmark_test.go`
**Description**: Benchmark validation logic performance
**Dependencies**: T013, T014, T031
**Tests**: Pattern validation, attribute validation, limit checks

### T034: Memory Usage Optimization
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Optimize memory usage for large roles
**Dependencies**: T032
**Optimizations**: Efficient data structures, garbage collection

### T035: Performance Validation Tests [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource_test.go`
**Description**: Add performance assertion tests
**Dependencies**: T031-T034
**Assertions**: Operation times, memory limits, throughput

---

## Phase 6: Documentation & Polish (Tasks 36-40)

### T036: Coverage Report Generation [P]
**File**: `internal/provider/resource_iam_custom_role/COVERAGE_REPORT.md`
**Description**: Generate test coverage report
**Dependencies**: T016-T030
**Target**: >95% test coverage

### T037: Quickstart Validation [P]
**File**: Test against `quickstart.md` scenarios
**Description**: Validate all quickstart scenarios work
**Dependencies**: All previous tasks
**Validation**: Execute all quickstart steps successfully

### T038: Provider Tests Integration
**File**: `internal/provider/provider_test.go`
**Description**: Add custom role resource to provider tests
**Dependencies**: T001, T027
**Tests**: Resource registration, configuration validation

### T039: Error Message Improvement [P]
**File**: `internal/provider/resource_iam_custom_role/iam_custom_role_resource.go`
**Description**: Enhance error messages for better UX
**Dependencies**: T015
**Improvements**: Clear validation errors, actionable messages

### T040: Final Integration Validation
**File**: Run complete test suite
**Description**: Execute full test suite and validate all requirements
**Dependencies**: All previous tasks
**Validation**: All tests pass, coverage meets target, performance acceptable

---

## Task Dependencies Summary

**Critical Path**: T001 → T006 → T007 → T008 → T016 → T026 → T029 → T030 → T040
**Parallel Opportunities**: Most testing tasks (T016-T025, T031-T033, T036-T039)
**Prerequisites**: T001 must complete before any resource implementation
**Validation Gates**: T029 (contracts), T030 (acceptance), T040 (final validation)

## Success Criteria

- [ ] All 40 tasks completed successfully
- [ ] >95% test coverage for custom role resource package
- [ ] All contract tests pass (T029)
- [ ] All acceptance tests pass (T030)
- [ ] Performance benchmarks meet targets (T031-T035)
- [ ] All quickstart scenarios work (T037)
- [ ] Full test suite passes (T040)

**Estimated Completion**: 35-40 tasks covering comprehensive test implementation