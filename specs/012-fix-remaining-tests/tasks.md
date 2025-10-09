# Tasks: Fix Remaining Tests

## Feature: Fix Remaining Tests

---

### Setup Tasks
- T001: Ensure Go 1.21+ and all dependencies are installed (go mod tidy, testify, Terraform Plugin Framework)
- T002: Clean and rebuild the project to clear any stale build artifacts

---

### Model & Schema Alignment [P]
- T003: Update `RoleBindingResourceModel` to include all required fields (id, tenant_id, roles, bindings, is_custom) in `/internal/provider/resource_iam_role_binding/`
- T004: Update `IamCustomRoleModel` to include all required fields (id, tenant_id, permissions, description) in `/internal/provider/resource_iam_custom_role/`
- T005: Update `IamResourceModel` to include all required fields (id, tenant_id, name, type, attributes) in `/internal/provider/resource_iam_resource/`
- T006: Update `HiiRetailProviderModel` to include tenant_id and ensure schema matches struct in `/internal/provider/provider_hiiretail_iam/`

---

### Test & Contract Fixes [P]
- T007: Fix argument and field mismatches in `ValidateRoleBindingModel` and related tests in `/internal/provider/resource_iam_role_binding/`
- T008: Fix build errors and type mismatches in `/internal/provider/resource_iam_custom_role/create_timeout_test.go`
- T009: Fix panic in `TestSetResourceContract/create_new_resource_success` by initializing IAMClient in `/internal/provider/resource_iam_resource/resource_test.go`
- T010: Fix provider integration tests for struct/object mismatches (tenant_id) in `/internal/provider/provider_integration_test.go`
- T011: Fix quickstart validation test for missing directory in `/internal/provider/resource_iam_group/quickstart_validation_test.go` (create or skip `acceptance_tests` directory)
- T012: Fix contract and quickstart tests to ensure they are either implemented or skipped, not failed

---

### Utility & Validation Tests [P]
- T013: Ensure all tests in `/internal/provider/testutils/` pass
- T014: Ensure all tests in `/internal/validation/` pass
- T015: Ensure all tests in `/internal/provider/shared/validators/` pass

---

### Polish & Final Validation [P]
- T016: Run all specified test commands and confirm all pass:
    - `go test ./internal/provider/resource_iam_role`
    - `go test ./internal/provider/resource_iam_role_binding`
    - `go test ./internal/provider/testutils`
    - `go test ./internal/validation`
    - `go test ./internal/provider/shared/validators`
    - `go test ./internal/provider/resource_iam_group`
    - `go test ./internal/provider/resource_iam_resource`
- T017: Document all fixes and update quickstart.md with final validation steps
- T018: Review for constitutional compliance and maintainability

---

## Parallel Execution Guidance
- Tasks marked [P] can be executed in parallel (T003-T006, T007-T015)
- Final validation (T016-T018) should be executed sequentially after all parallel tasks complete

---

## Dependency Notes
- Model/schema alignment tasks (T003-T006) must be completed before test/contract fixes (T007-T012)
- Utility and validation tests (T013-T015) depend on model/schema fixes
- Final validation (T016) depends on all previous tasks

---
