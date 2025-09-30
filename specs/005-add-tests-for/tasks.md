# Tasks: IAM Role Binding Resource Implementation and Testing

**Input**: Design documents from `/specs/005-add-tests-for/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Go 1.21+ with terraform-plugin-framework v1.16.0
   → Structure: Single Terraform provider project structure
2. Load design documents:
   → data-model.md: IAM Role Binding entity with bindings sub-entity
   → contracts/: OpenAPI spec + contract tests for CRUD operations
   → quickstart.md: 5 user scenarios for validation
3. Generate tasks by category:
   → Setup: Provider structure, dependencies, mock server
   → Tests: Contract tests, unit tests, acceptance tests
   → Core: Resource implementation, provider registration
   → Integration: OAuth2 client, mock server integration
   → Polish: Documentation, performance validation
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T025)
6. Generate dependency graph
7. SUCCESS - ready for implementation
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions (Terraform Provider)
```
internal/provider/resource_iam_role_binding/
├── iam_role_binding_resource_gen.go (schema - exists)
├── iam_role_binding_resource.go (implementation)
└── iam_role_binding_resource_test.go (unit tests)

acceptance_tests/
└── iam_role_binding_resource_test.go (acceptance tests)
```

## Phase 3.1: Setup
- [x] T001 Verify existing provider structure and generated schema in `internal/provider/resource_iam_role_binding/iam_role_binding_resource_gen.go`
- [x] T002 Analyze existing OAuth2 client integration in `internal/client/oauth2_client.go` for authentication patterns
- [x] T003 [P] Review mock server setup in `acceptance_tests/mock_server_test.go` for test infrastructure

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T004 [P] Contract test CREATE role binding in `specs/005-add-tests-for/contracts/role_binding_contract_test.go` (enhance existing)
- [ ] T005 [P] Contract test READ role binding in `specs/005-add-tests-for/contracts/role_binding_contract_test.go` (enhance existing)  
- [ ] T006 [P] Contract test UPDATE role binding in `specs/005-add-tests-for/contracts/role_binding_contract_test.go` (enhance existing)
- [ ] T007 [P] Contract test DELETE role binding in `specs/005-add-tests-for/contracts/role_binding_contract_test.go` (enhance existing)
- [ ] T008 [P] Unit test for role binding model validation in `internal/provider/resource_iam_role_binding/iam_role_binding_resource_test.go`
- [ ] T009 [P] Unit test for max 10 bindings validation in `internal/provider/resource_iam_role_binding/iam_role_binding_resource_test.go`
- [ ] T010 [P] Unit test for tenant isolation logic in `internal/provider/resource_iam_role_binding/iam_role_binding_resource_test.go`
- [ ] T011 [P] Acceptance test for basic CRUD lifecycle in `acceptance_tests/iam_role_binding_resource_test.go`
- [ ] T012 [P] Acceptance test for import functionality in `acceptance_tests/iam_role_binding_resource_test.go`

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T013 [P] Create role binding data models and validation functions in `internal/models/role_binding_models.go`
- [ ] T014 Implement Create operation in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`
- [ ] T015 Implement Read operation in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`
- [ ] T016 Implement Update operation in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`
- [ ] T017 Implement Delete operation in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`
- [ ] T018 Implement Import functionality in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`
- [ ] T019 Add comprehensive input validation and error handling in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`

## Phase 3.4: Integration  
- [ ] T020 Register role binding resource in provider configuration in `internal/provider/provider_hiiretail_iam/hiiretail_iam_provider_gen.go`
- [ ] T021 Integrate OAuth2 authentication for all role binding operations in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`
- [ ] T022 Create mock server endpoints for role binding API in `acceptance_tests/mock_server_test.go`
- [ ] T023 Implement proper diagnostic messages and troubleshooting support in `internal/provider/resource_iam_role_binding/iam_role_binding_resource.go`

## Phase 3.5: Polish
- [ ] T024 [P] Performance validation - ensure <2s operations per quickstart requirements
- [ ] T025 [P] Run all quickstart scenarios from `specs/005-add-tests-for/quickstart.md` for end-to-end validation

## Dependencies
**Critical Path**:
- Setup (T001-T003) before tests
- All tests (T004-T012) before implementation (T013-T019)
- Core implementation (T013-T019) before integration (T020-T023)
- Integration before polish (T024-T025)

**Detailed Dependencies**:
- T001-T003 block T004-T012 (need existing code analysis)
- T004-T012 must all be complete and FAILING before T013
- T013 blocks T014-T019 (models needed for resource implementation)
- T020 requires T014-T019 (need resource implementation to register)
- T021 must integrate with T014-T019 (OAuth2 for all operations)
- T022 supports T011-T012 (mock server for acceptance tests)

## Parallel Execution Examples

### Phase 3.2 - Contract Tests (can run together)
```bash
# Launch T004-T007 together (same file, different test functions):
Task: "Enhance CREATE role binding contract test with proper API validation"
Task: "Enhance READ role binding contract test with response validation"  
Task: "Enhance UPDATE role binding contract test with atomic update validation"
Task: "Enhance DELETE role binding contract test with cleanup validation"
```

### Phase 3.2 - Unit Tests (can run together)
```bash
# Launch T008-T010 together (same file, different test functions):
Task: "Unit test role binding model validation with all field types"
Task: "Unit test max 10 bindings validation with boundary conditions"
Task: "Unit test tenant isolation logic with multiple tenant scenarios"
```

### Phase 3.2 - Acceptance Tests (can run together)
```bash  
# Launch T011-T012 together (same file, different test functions):
Task: "Acceptance test basic CRUD lifecycle with terraform-plugin-testing"
Task: "Acceptance test import functionality with proper state synchronization"
```

### Phase 3.5 - Polish (can run together)
```bash
# Launch T024-T025 together (independent validation):
Task: "Performance validation ensuring <2s operations"
Task: "Run all quickstart scenarios for end-to-end validation"
```

## Task Generation Rules Applied

### From Contracts (specs/005-add-tests-for/contracts/)
- `role_binding_api.yaml` → T004-T007 (one per CRUD operation)
- `role_binding_contract_test.go` → Enhanced existing contract tests

### From Data Model (specs/005-add-tests-for/data-model.md)
- IAM Role Binding entity → T013 (model creation)
- Binding sub-entity → Included in T013
- Validation rules → T008-T010, T019

### From User Stories (specs/005-add-tests-for/quickstart.md)
- 5 scenarios → T011-T012 (acceptance tests)
- End-to-end validation → T025

### From Technical Context (specs/005-add-tests-for/plan.md)
- terraform-plugin-framework → T014-T018 (resource implementation)
- OAuth2 integration → T021
- Mock server → T022
- Provider registration → T020

## Validation Checklist
*GATE: Verified before task execution*

- [x] All contracts have corresponding tests (T004-T007)
- [x] All entities have model tasks (T013)
- [x] All tests come before implementation (T004-T012 → T013-T019)
- [x] Parallel tasks truly independent ([P] tasks use different files or functions)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Provider registration included (T020)
- [x] Authentication integration included (T021)
- [x] Mock server integration included (T022)
- [x] Performance requirements addressed (T024)
- [x] End-to-end validation included (T025)

## Notes
- Generated schema already exists in `iam_role_binding_resource_gen.go`
- OAuth2 client patterns available from existing resources
- Mock server infrastructure established
- All tests must fail before implementation begins
- Follow terraform-plugin-framework patterns from iam_custom_role
- Commit after each task completion
- Use existing provider factory for testing