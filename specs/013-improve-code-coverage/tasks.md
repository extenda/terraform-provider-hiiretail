# Tasks: Improve Code Coverage (IAM Service)

**Input**: Design documents from `/specs/013-improve-code-coverage/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/

## Execution Flow (main)

### Phase 3.1: Setup
- [ ] T001 Ensure Go test and coverage tools are available
- [ ] T002 [P] Configure coverage reporting in CI pipeline (`.github/workflows/ci-sonarcloud.yml`)
- [ ] T003 [P] Review and update linting/formatting tools for test files

### Phase 3.2: Tests First (TDD)
- [ ] T004 [P] Create/expand unit test file for IAM service: `internal/provider/iam/service_test.go`
- [ ] T005 [P] Write contract tests for each method in `service.go` per `contracts/service_coverage_contract.md`
- [ ] T006 [P] Add table-driven tests for edge cases and error scenarios in `service_test.go`
- [ ] T007 [P] Mock API client and network calls for isolated unit testing
- [ ] T008 [P] Add integration tests for IAM service covering real API interactions (if possible)

### Phase 3.3: Core Implementation
- [ ] T009 [P] Refactor IAM service code to improve testability (inject dependencies, expose error paths)
- [ ] T010 [P] Document and justify any code excluded from coverage in `service.go` and coverage report

### Phase 3.4: Polish
- [ ] T011 [P] Run coverage analysis and verify 90%+ coverage for `service.go`
- [ ] T012 [P] Update documentation to describe coverage measurement and improvement process
- [ ] T013 [P] Commit all changes and validate CI pipeline enforces coverage threshold

## Parallel Execution Guidance
- Tasks T002, T003, T004, T005, T006, T007, T008 can be run in parallel (different files, no dependencies)
- T009 depends on completion of initial tests (T004-T008)
- T011, T012, T013 can be run in parallel after implementation and refactoring

## Dependency Notes
- Tests (T004-T008) must be written and run before refactoring (T009)
- Coverage analysis (T011) and documentation (T012) follow implementation
- Final commit and CI validation (T013) after all other tasks

## File Paths
- `internal/provider/iam/service.go` (target for coverage)
- `internal/provider/iam/service_test.go` (unit tests)
- `contracts/service_coverage_contract.md` (test contract)
- `.github/workflows/ci-sonarcloud.yml` (CI coverage enforcement)

## Validation Checklist
- [x] All contracts have corresponding tests
- [x] All entities have model tasks
- [x] All tests come before implementation
- [x] Parallel tasks truly independent
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
