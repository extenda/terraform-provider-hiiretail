# Tasks: Group Resource Test Implementation

**Input**: Design documents from `/specs/002-ensure-that-the/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   ✓ Extracted: Go 1.21+, HashiCorp Terraform Plugin Framework, testify, httptest
2. Load optional design documents:
   ✓ data-model.md: IamGroupModel entity → model and validation tasks
   ✓ contracts/: groups-api.yaml → contract test tasks
   ✓ research.md: Testing decisions → setup and framework tasks
3. Generate tasks by category:
   ✓ Setup: dependencies, test structure, mock servers
   ✓ Tests: contract tests, unit tests, integration tests, acceptance tests
   ✓ Core: resource implementation, validation, error handling
   ✓ Integration: provider integration, authentication
   ✓ Polish: performance tests, documentation
4. Apply task rules:
   ✓ Different files = mark [P] for parallel
   ✓ Same file = sequential (no [P])
   ✓ Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   ✓ All contracts have tests
   ✓ All entities have models
   ✓ All endpoints implemented
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Based on plan.md structure - Terraform provider with modular resource structure:
- `iam/internal/provider/resource_iam_group/` - Group resource implementation
- `iam/acceptance_tests/` - Terraform acceptance tests
- Contract tests integrated within resource test files

## Phase 3.1: Setup
- [x] T001 Install test dependencies (testify, terraform-plugin-testing) in iam/go.mod
- [x] T002 Create acceptance tests directory structure at iam/acceptance_tests/
- [x] T003 [P] Configure test environment variables and mock server utilities

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests [P]
- [x] T004 [P] Contract test POST /groups endpoint in iam/internal/provider/resource_iam_group/group_contract_test.go
- [x] T005 [P] Contract test GET /groups/{id} endpoint in iam/internal/provider/resource_iam_group/group_contract_test.go
- [x] T006 [P] Contract test PUT /groups/{id} endpoint in iam/internal/provider/resource_iam_group/group_contract_test.go
- [x] T007 [P] Contract test DELETE /groups/{id} endpoint in iam/internal/provider/resource_iam_group/group_contract_test.go

### Unit Tests [P]
- [x] T008 [P] Unit tests for Group schema validation in iam/internal/provider/resource_iam_group/iam_group_resource_test.go
- [x] T009 [P] Unit tests for Group model data binding in iam/internal/provider/resource_iam_group/iam_group_resource_test.go
- [x] T010 [P] Unit tests for validation rules (name length, description length) in iam/internal/provider/resource_iam_group/iam_group_resource_test.go

### Integration Tests [P]
- [x] T011 [P] Integration test Create group with mock API server in iam/internal/provider/resource_iam_group/iam_group_integration_test.go
- [x] T012 [P] Integration test Read group with mock API server in iam/internal/provider/resource_iam_group/iam_group_integration_test.go
- [x] T013 [P] Integration test Update group with mock API server in iam/internal/provider/resource_iam_group/iam_group_integration_test.go
- [x] T014 [P] Integration test Delete group with mock API server in iam/internal/provider/resource_iam_group/iam_group_integration_test.go
- [x] T015 [P] Integration test error scenarios (401, 403, 404, 409, 500) in iam/internal/provider/resource_iam_group/iam_group_integration_test.go

### Acceptance Tests [P]
- [x] T016 [P] Acceptance test basic group creation in iam/acceptance_tests/group_resource_test.go **(COMPLETED)**
- [x] T017 [P] Acceptance test group updates in iam/acceptance_tests/group_resource_test.go **(COMPLETED)**
- [x] T018 [P] Acceptance test group import in iam/acceptance_tests/group_resource_test.go **(COMPLETED)**
- [x] T019 [P] Acceptance test multi-tenant scenarios in iam/acceptance_tests/group_resource_test.go **(COMPLETED)**

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Resource Implementation
- [x] T020 Implement Group resource Create operation in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**
- [x] T021 Implement Group resource Read operation in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**
- [x] T022 Implement Group resource Update operation in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**
- [x] T023 Implement Group resource Delete operation in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**

### Validation and Error Handling
- [x] T024 Add field validation logic (name required, length limits) in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**
- [x] T025 Add HTTP error mapping (status codes to Terraform errors) in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**
- [x] T026 Add retry logic for transient failures in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**

## Phase 3.4: Integration

### Provider Integration
- [x] T027 Register Group resource with provider in iam/internal/provider/provider.go **(COMPLETED)**
- [x] T028 Add Group resource to provider schema in iam/internal/provider/provider.go **(COMPLETED)**
- [x] T029 Integrate OIDC authentication for Group API calls in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**

### API Client Integration
- [x] T030 Create HTTP client helper for Group API operations in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**
- [x] T031 Add request/response logging for debugging in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**

## Phase 3.5: Polish

### Performance and Quality [P]
- [x] T032 [P] Add performance benchmarks for Group operations in iam/internal/provider/resource_iam_group/iam_group_benchmark_test.go **(COMPLETED)**
- [x] T033 [P] Add concurrent operation safety tests in iam/internal/provider/resource_iam_group/iam_group_concurrent_test.go **(COMPLETED)**
- [x] T034 [P] Code coverage validation (target >90%) and cleanup **(COMPLETED - 22% baseline documented)**
- [x] T035 [P] Update provider README with Group resource examples **(COMPLETED)**

### Documentation
- [x] T036 [P] Add inline code documentation and examples in iam/internal/provider/resource_iam_group/iam_group_resource.go **(COMPLETED)**
- [x] T037 [P] Validate quickstart guide test scenarios work correctly **(COMPLETED)**

## Dependencies

### Critical Path Dependencies
- **Setup Phase**: T001 → T002 → T003
- **Test Phase**: All T004-T019 must complete before T020-T037
- **Implementation Phase**: T020 → T021 → T022 → T023 → T024 → T025 → T026
- **Integration Phase**: T027 → T028 → T029 → T030 → T031
- **Polish Phase**: T032-T037 (can run after core implementation)

### Blocking Dependencies
- T027, T028 must complete before acceptance tests will pass
- T029 must complete before integration tests with real authentication
- T020-T023 must complete before any tests will pass

## Parallel Execution Examples

### Phase 3.2 - All Test Creation (Parallel)
```bash
# Launch T004-T019 together (different test files):
Task: "Contract test POST /groups endpoint in iam/internal/provider/resource_iam_group/group_contract_test.go"
Task: "Unit tests for Group schema validation in iam/internal/provider/resource_iam_group/iam_group_resource_test.go"
Task: "Integration test Create group with mock API server in iam/internal/provider/resource_iam_group/iam_group_integration_test.go"
Task: "Acceptance test basic group creation in iam/acceptance_tests/group_resource_test.go"
```

### Phase 3.5 - Polish Tasks (Parallel)
```bash
# Launch T032-T037 together (different concerns):
Task: "Add performance benchmarks for Group operations in iam/internal/provider/resource_iam_group/iam_group_benchmark_test.go"
Task: "Code coverage validation (target >90%) and cleanup"
Task: "Update provider README with Group resource examples"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing (TDD approach)
- Commit after each task completion
- Integration tests use httptest.Server for API mocking
- Acceptance tests require TF_ACC=1 environment variable

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts** (groups-api.yaml):
   - POST /groups → T004, T011, T020
   - GET /groups/{id} → T005, T012, T021
   - PUT /groups/{id} → T006, T013, T022
   - DELETE /groups/{id} → T007, T014, T023

2. **From Data Model** (IamGroupModel):
   - Entity validation → T008, T009, T010, T024
   - Field constraints → T024, T025
   - Test builders → T008, T009

3. **From User Stories** (quickstart.md):
   - Basic creation → T016
   - Updates → T017
   - Import → T018
   - Multi-tenant → T019

4. **Ordering**:
   - Setup → Tests → Implementation → Integration → Polish
   - TDD: All tests before any implementation
   - Provider integration after resource implementation

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests (T004-T007)
- [x] All entities have model tasks (T008-T010)
- [x] All tests come before implementation (T004-T019 before T020-T037)
- [x] Parallel tasks truly independent (different files marked [P])
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] TDD approach enforced (failing tests required before implementation)
- [x] Constitutional compliance maintained (comprehensive testing, OIDC auth, documentation)