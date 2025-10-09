# Tasks: Support TDD for Authentication

**Input**: Design documents from `/specs/011-support-tdd-for/`
**Prerequisites**: plan.md (required)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: tech stack, libraries, structure
2. No optional design documents present (no data-model.md, contracts/, research.md, quickstart.md)
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: authentication contract and integration tests
   → Core: authentication models, services
   → Polish: unit tests, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness
```

## Phase 3.1: Setup
- [ ] T001 Ensure Go 1.21+ and Terraform Plugin Framework dependencies in go.mod
- [ ] T002 [P] Configure linting and formatting tools (gofmt, golangci-lint)

## Phase 3.2: Tests First (TDD)
- [ ] T003 [P] Refactor and fix authentication tests in `internal/provider/shared/auth/client_test.go`
- [ ] T004 [P] Refactor and fix authentication integration tests in `internal/provider/shared/auth/integration_test.go`

## Phase 3.3: Core Implementation
- [ ] T005 Update authentication model and client in `internal/provider/shared/auth/auth.go` as needed to support passing tests
- [ ] T006 Update supporting code in `internal/provider/shared/auth/client.go` for test reliability

## Phase 3.4: Polish
- [ ] T007 [P] Add/validate unit tests for edge cases in `internal/provider/shared/auth/client_test.go`
- [ ] T008 [P] Update documentation in `docs/guides/authentication.md` to reflect new authentication flow and test coverage

## Dependencies
- Setup (T001, T002) before all
- Tests (T003, T004) before implementation (T005, T006)
- Core (T005, T006) before polish (T007, T008)
- [P] tasks can run in parallel if in different files

## Parallel Execution Example
- T003 and T004 can be executed in parallel
- T007 and T008 can be executed in parallel

## Task Agent Commands
- agent run T003 & agent run T004
- agent run T007 & agent run T008
