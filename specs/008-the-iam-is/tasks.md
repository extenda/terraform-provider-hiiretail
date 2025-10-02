# Tasks: Multi-API Provider User Experience Enhancement

**Input**: Design documents from `/specs/008-the-iam-is/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Based on plan.md structure: Multi-API Terraform provider with service modules under `internal/provider/`

## Phase 3.1: Setup & Project Restructuring
- [x] T001 Rename provider module from hiiretail-iam to hiiretail in go.mod
- [x] T002 Update main.go to reference new provider structure and remove hiiretail-iam references
- [x] T003 [P] Create internal/provider/shared/auth/ directory structure for shared OAuth2 implementation
- [x] T004 [P] Create internal/provider/shared/client/ directory structure for common HTTP client
- [x] T005 [P] Create internal/provider/shared/validators/ directory structure for common validation logic
- [x] T006 [P] Create internal/provider/iam/ directory structure with resources/ and data_sources/ subdirs
- [x] T007 [P] Create internal/provider/ccc/ directory structure for future API (resources/ and data_sources/)
- [x] T008 [P] Create enhanced docs/ directory structure (guides/, resources/, examples/)
- [x] T009 [P] Create tests/ directory structure (acceptance/iam/, integration/, unit/)

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T010 [P] Contract test for provider configuration schema in tests/unit/provider_config_test.go
- [ ] T011 [P] Contract test for resource naming validation in tests/unit/resource_naming_test.go
- [ ] T012 [P] Contract test for service module interface in tests/unit/service_module_test.go
- [ ] T013 [P] Contract test for documentation structure validation in tests/unit/documentation_test.go
- [ ] T014 [P] Integration test for provider initialization with new name in tests/integration/provider_init_test.go
- [ ] T015 [P] Integration test for IAM service module registration in tests/integration/iam_service_test.go
- [ ] T016 [P] Acceptance test for hiiretail_iam_group resource in tests/acceptance/iam/group_test.go
- [ ] T017 [P] Acceptance test for hiiretail_iam_custom_role resource in tests/acceptance/iam/custom_role_test.go
- [ ] T018 [P] Acceptance test for hiiretail_iam_role_binding resource in tests/acceptance/iam/role_binding_test.go
- [ ] T019 [P] Integration test for quickstart scenario in tests/integration/quickstart_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T020 [P] Move existing auth package to internal/provider/shared/auth/ and update imports
- [ ] T021 [P] Create shared HTTP client in internal/provider/shared/client/client.go
- [ ] T022 [P] Create common validators in internal/provider/shared/validators/validators.go
- [ ] T023 Create ServiceModule interface in internal/provider/service_module.go
- [ ] T024 Create IAM service module in internal/provider/iam/service.go implementing ServiceModule interface
- [ ] T025 Update internal/provider/provider.go to register all service modules with new hiiretail provider
- [ ] T026 Rename hiiretail_iam_group resource: move from resource_iam_group/ to internal/provider/iam/resources/group.go
- [ ] T027 Rename hiiretail_iam_custom_role resource: move from resource_iam_custom_role/ to internal/provider/iam/resources/custom_role.go
- [ ] T028 Rename hiiretail_iam_role_binding resource: move from resource_iam_role_binding/ to internal/provider/iam/resources/role_binding.go
- [ ] T029 Create hiiretail_iam_groups data source in internal/provider/iam/data_sources/groups.go
- [ ] T030 Update resource TypeName() methods to return new naming convention (hiiretail_iam_*)
- [ ] T031 Update provider resource registration to use service module pattern
- [ ] T032 Create placeholder CCC service module in internal/provider/ccc/service.go for future extensibility

## Phase 3.4: Integration & Migration Support
- [ ] T033 Update generator_config.yaml to reflect new provider name and resource structure
- [ ] T034 Create backward compatibility aliases for old resource names during transition
- [ ] T035 Update OAuth2 configuration to support service-specific endpoints
- [ ] T036 Implement service health validation in provider configuration
- [ ] T037 Add provider configuration validation for new schema
- [ ] T038 Update comprehensive_test.tf to use new provider and resource names

## Phase 3.5: Documentation & Examples
- [ ] T039 [P] Create docs/index.md with provider overview and multi-API introduction
- [ ] T040 [P] Create docs/guides/getting-started.md based on quickstart.md content
- [ ] T041 [P] Create docs/guides/authentication.md with OAuth2 configuration details
- [ ] T042 [P] Create docs/guides/migration-guides/from-hiiretail-iam.md with migration instructions
- [ ] T043 [P] Create docs/resources/iam/overview.md with IAM service documentation
- [ ] T044 [P] Create docs/resources/iam/group.md for hiiretail_iam_group resource
- [ ] T045 [P] Create docs/resources/iam/custom_role.md for hiiretail_iam_custom_role resource
- [ ] T046 [P] Create docs/resources/iam/role_binding.md for hiiretail_iam_role_binding resource
- [ ] T047 [P] Create docs/examples/basic-iam-setup/ with complete working example
- [ ] T048 [P] Create docs/examples/multi-service-deployment/ showing future multi-API usage
- [ ] T049 [P] Update README.md to reflect new provider name and capabilities

## Phase 3.6: Polish & Validation
- [ ] T050 [P] Update all existing unit tests to work with new structure in tests/unit/
- [ ] T051 Run full acceptance test suite and fix any failures
- [ ] T052 [P] Add performance benchmarks for provider initialization in tests/performance/
- [ ] T053 [P] Create migration automation script for users upgrading from hiiretail-iam
- [ ] T054 Update build and release scripts to use new provider name
- [ ] T055 Validate quickstart guide with real deployment
- [ ] T056 Update comprehensive test configuration and run end-to-end validation
- [ ] T057 Clean up old directory structure and remove deprecated files

## Dependencies
- Setup (T001-T009) before everything else
- Tests (T010-T019) before implementation (T020-T032)
- T020 (auth move) blocks T021, T022 (shared components)
- T023 (ServiceModule interface) blocks T024, T032 (service implementations)
- T024 (IAM service) blocks T026-T028 (resource moves)
- T025 (provider update) requires T023, T024 (service modules)
- T026-T030 (resource renames) before T031 (registration update)
- Implementation (T020-T032) before integration (T033-T038)
- Integration before documentation (T039-T049)
- Documentation before polish (T050-T057)

## Parallel Example
```
# Launch setup tasks together:
Task: "Create internal/provider/shared/auth/ directory structure for shared OAuth2 implementation"
Task: "Create internal/provider/shared/client/ directory structure for common HTTP client"
Task: "Create internal/provider/shared/validators/ directory structure for common validation logic"
Task: "Create internal/provider/iam/ directory structure with resources/ and data_sources/ subdirs"

# Launch contract tests together:
Task: "Contract test for provider configuration schema in tests/unit/provider_config_test.go"
Task: "Contract test for resource naming validation in tests/unit/resource_naming_test.go"
Task: "Contract test for service module interface in tests/unit/service_module_test.go"
Task: "Contract test for documentation structure validation in tests/unit/documentation_test.go"

# Launch documentation tasks together:
Task: "Create docs/index.md with provider overview and multi-API introduction"
Task: "Create docs/guides/getting-started.md based on quickstart.md content"
Task: "Create docs/guides/authentication.md with OAuth2 configuration details"
Task: "Create docs/resources/iam/overview.md with IAM service documentation"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each major phase
- Maintain backward compatibility during transition
- Test all examples and documentation before completion

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - provider-config.md → T010 (provider configuration test)
   - resource-naming.md → T011 (naming validation test)  
   - service-module.md → T012 (service module interface test)
   - documentation-structure.md → T013 (documentation validation test)

2. **From Data Model**:
   - Provider Configuration Entity → T020, T035, T037 (auth and config tasks)
   - Service Module Entity → T023, T024, T032 (service module implementation)
   - Resource Naming Entity → T026-T030 (resource renaming)
   - Documentation Entity → T039-T049 (documentation generation)

3. **From Quickstart Scenarios**:
   - Provider setup → T014, T019 (integration tests)
   - Resource creation → T016-T018 (acceptance tests)
   - End-to-end validation → T055, T056 (validation tasks)

4. **Ordering**:
   - Setup → Tests → Core → Integration → Documentation → Polish
   - Service architecture before resource migration
   - Tests before implementation (TDD)

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All contracts have corresponding tests (T010-T013)
- [x] All entities have implementation tasks (T020-T032)
- [x] All tests come before implementation (T010-T019 before T020-T032)
- [x] Parallel tasks truly independent (different files/directories)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Migration and backward compatibility addressed (T034, T038, T053)
- [x] Documentation covers all aspects (T039-T049)
- [x] Validation and testing comprehensive (T050-T057)