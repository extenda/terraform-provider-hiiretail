# Tasks: Provider Distribution and Availability

**Input**: Design documents from `/specs/010-make-our-provider/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Go 1.21+, HashiCorp Terraform Plugin Framework, GoReleaser, GitHub Actions
   → Structure: Single project (Terraform provider distribution)
2. Load design documents:
   → data-model.md: Release artifacts, build configuration, documentation bundle
   → contracts/: GitHub Actions, GoReleaser, Terraform Registry contracts
   → quickstart.md: Complete validation pipeline with test scenarios
3. Generate tasks by category:
   → Setup: GoReleaser config, GitHub Actions workflow
   → Tests: Configuration validation, build testing
   → Core: Release automation, artifact generation
   → Integration: Registry publishing, documentation sync
   → Polish: Process documentation, monitoring
4. Apply task rules:
   → Configuration files = [P] for parallel creation
   → Sequential validation and testing workflow
   → Tests before live releases (TDD approach)
5. Number tasks sequentially (T001, T002...)
6. SUCCESS: Tasks ready for provider distribution implementation
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
Single project structure at repository root:
- `.goreleaser.yml` - GoReleaser configuration
- `.github/workflows/release.yml` - GitHub Actions workflow
- `README.md` - Installation and usage documentation
- `docs/` - Provider documentation for registry
- `examples/` - Usage examples for registry

## Phase 3.1: Setup
- [ ] T001 Verify prerequisites: GitHub secrets GPG_PRIVATE_KEY and PASSPHRASE are configured
- [ ] T002 Ensure clean repository state with no uncommitted changes
- [ ] T003 Validate existing Go module structure and dependencies

## Phase 3.2: Configuration Creation (TDD - Tests First)
**CRITICAL: Validation tests MUST be written and MUST FAIL before configuration implementation**
- [ ] T004 [P] Create GoReleaser validation test in `test_goreleaser_config.sh` to validate .goreleaser.yml syntax and required fields
- [ ] T005 [P] Create GitHub Actions workflow validation test in `test_github_actions.sh` to validate .github/workflows/release.yml syntax
- [ ] T006 [P] Create build artifacts validation test in `test_build_artifacts.sh` to verify expected output structure
- [ ] T007 Create release process integration test in `test_release_process.sh` to validate end-to-end workflow

## Phase 3.3: Core Configuration Implementation (ONLY after tests are failing)
- [ ] T008 [P] Create GoReleaser configuration file `.goreleaser.yml` with project metadata, build matrix, and signing setup
- [ ] T009 [P] Create GitHub Actions release workflow `.github/workflows/release.yml` with Go setup, GPG import, and GoReleaser execution
- [ ] T010 Update Makefile to include release validation targets and local build testing commands
- [ ] T011 Configure release validation pipeline with pre-release checks and artifact verification

## Phase 3.4: Documentation and Registry Preparation
- [ ] T012 [P] Update README.md with Terraform Registry installation instructions and provider discovery information
- [ ] T013 [P] Enhance docs/ directory with comprehensive provider documentation for registry publication
- [ ] T014 [P] Improve examples/ directory with working Terraform configurations demonstrating provider usage
- [ ] T015 Create CHANGELOG.md template for release notes and version history tracking

## Phase 3.5: Validation and Testing
- [ ] T016 Execute local GoReleaser build test using `goreleaser build --snapshot --clean` to verify configuration
- [ ] T017 Validate multi-platform binary generation and verify artifact naming conventions
- [ ] T018 Test GPG signing setup with test artifacts to ensure proper signature generation
- [ ] T019 Create beta release (v0.1.0-beta.1) to test complete GitHub Actions workflow

## Phase 3.6: Registry Integration and Polish
- [ ] T020 Validate GitHub release creation with proper artifacts, checksums, and signatures
- [ ] T021 Monitor Terraform Registry detection and provider publication process
- [ ] T022 Test provider installation from registry using `terraform init` with registry source
- [ ] T023 [P] Create release process documentation in docs/RELEASE_PROCESS.md for team reference
- [ ] T024 [P] Set up release monitoring and notification system for successful publications

## Dependencies
- Prerequisites (T001-T003) before configuration tests (T004-T007)
- Configuration tests (T004-T007) before implementation (T008-T011)
- T008 (.goreleaser.yml) blocks T016 (local build test)
- T009 (GitHub Actions) blocks T019 (beta release test)
- Documentation tasks (T012-T015) can run parallel with validation (T016-T018)
- T019 (beta release) blocks registry integration (T020-T022)
- Process documentation (T023-T024) after successful validation

## Parallel Execution Examples

### Configuration Creation Phase
```bash
# Launch T004-T006 together (different test files):
Task: "Create GoReleaser validation test in test_goreleaser_config.sh"
Task: "Create GitHub Actions workflow validation test in test_github_actions.sh"  
Task: "Create build artifacts validation test in test_build_artifacts.sh"
```

### Implementation Phase
```bash
# Launch T008-T009 together (different config files):
Task: "Create GoReleaser configuration file .goreleaser.yml"
Task: "Create GitHub Actions release workflow .github/workflows/release.yml"
```

### Documentation Phase
```bash
# Launch T012-T014 together (different documentation files):
Task: "Update README.md with Terraform Registry installation instructions"
Task: "Enhance docs/ directory with comprehensive provider documentation"
Task: "Improve examples/ directory with working Terraform configurations"
```

## Contract Validation Mapping

### GitHub Actions Contract → Tasks
- Workflow trigger contract → T005 (workflow validation test), T009 (workflow implementation)
- Job execution contract → T009 (workflow implementation), T019 (beta release test)
- Secrets contract → T001 (prerequisites verification)
- Output contract → T006 (artifacts validation test), T020 (release validation)

### GoReleaser Contract → Tasks
- Build configuration contract → T004 (config validation test), T008 (config implementation)
- Archive configuration contract → T008 (config implementation), T017 (binary validation)
- Checksum and signing contract → T008 (config implementation), T018 (GPG signing test)
- Release configuration contract → T008 (config implementation), T019 (release test)

### Terraform Registry Contract → Tasks
- Provider registration contract → T021 (registry monitoring), T022 (installation test)
- Documentation sync contract → T013-T014 (documentation preparation)
- Version management contract → T015 (changelog), T019-T020 (release process)

## Success Criteria
- [ ] All configuration validation tests pass
- [ ] Multi-platform binaries build successfully with proper naming
- [ ] GPG signing works for all artifacts
- [ ] GitHub releases created automatically on version tags
- [ ] Terraform Registry detects and publishes provider
- [ ] Users can install provider via `terraform init`
- [ ] Complete documentation available on registry
- [ ] Release process documented for team use

## Notes
- [P] tasks operate on different files with no dependencies
- Validation tests must fail initially to confirm proper TDD approach
- Beta release (T019) serves as integration test before production
- Registry detection may take several hours - plan accordingly
- GPG secrets must be properly configured before starting implementation
- Each task should be committed separately for audit trail