
# Implementation Plan: Group Resource Test Implementation

**Branch**: `002-ensure-that-the` | **Date**: September 28, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-ensure-that-the/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Create comprehensive test coverage for the existing IAM Group resource in the HiiRetail Terraform Provider. The Group resource schema is already implemented with attributes for name (required), description (optional), id (computed), status (computed), and tenant_id (optional). This feature focuses on creating unit tests, integration tests, and validation tests to ensure the resource functions correctly with proper CRUD operations, validation logic, and error handling.

## Technical Context
**Language/Version**: Go 1.21+  
**Primary Dependencies**: HashiCorp Terraform Plugin Framework, golang.org/x/oauth2, terraform-plugin-framework-validators  
**Storage**: HiiRetail IAM API (RESTful service)  
**Testing**: Go testing framework, testify for assertions, httptest for mock servers  
**Target Platform**: Cross-platform (Linux, macOS, Windows) Terraform provider  
**Project Type**: Single project - Terraform provider with resource implementations  
**Performance Goals**: Resource operations complete within 30 seconds, efficient state management  
**Constraints**: Must follow HashiCorp Plugin Framework conventions, secure OIDC authentication required  
**Scale/Scope**: Multi-tenant IAM group management, enterprise-grade reliability and validation## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Provider Interface Standards
✅ **COMPLIANT**: Group resource schema follows HashiCorp Plugin Framework conventions with proper validation and computed fields

### II. Resource Implementation Completeness  
✅ **DESIGN COMPLIANT**: Comprehensive CRUD operations planned with proper API integration, error handling, and state management

### III. Authentication & Security (NON-NEGOTIABLE)
✅ **COMPLIANT**: Uses existing OIDC authentication implementation with secure credential handling

### IV. Testing & Validation
✅ **DESIGN COMPLIANT**: Comprehensive test suite planned including unit tests (>90% coverage), integration tests with mock servers, and acceptance tests following Terraform conventions

### V. Documentation & Examples
✅ **DESIGN COMPLIANT**: Complete documentation including quickstart guide, API contracts, and usage examples provided

**Post-Design Status**: All constitutional requirements are now satisfied through comprehensive design artifacts.

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
iam/
├── main.go                           # Provider entry point
├── go.mod                           # Go module dependencies
├── internal/provider/
│   ├── provider.go                  # Main provider implementation with OIDC auth
│   ├── provider_test.go             # Existing provider tests
│   ├── provider_integration_test.go # Existing integration tests
│   ├── resource_iam_group/
│   │   ├── iam_group_resource_gen.go    # Existing schema (needs full implementation)
│   │   ├── iam_group_resource.go        # NEW: Full resource implementation
│   │   ├── iam_group_resource_test.go   # NEW: Unit tests
│   │   └── iam_group_integration_test.go # NEW: Integration tests
│   ├── resource_iam_role/           # Existing role resource
│   ├── resource_iam_custom_role/    # Existing custom role resource
│   └── resource_iam_role_binding/   # Existing role binding resource
└── acceptance_tests/                # NEW: Acceptance tests directory
    └── group_resource_test.go       # NEW: Terraform acceptance tests
```

**Structure Decision**: Single Terraform provider project with modular resource structure. Following HashiCorp's conventions with separate directories for each resource type and comprehensive test coverage at unit, integration, and acceptance levels.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh copilot`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base structure
- Generate test-driven tasks from design artifacts:
  * Contract tests from `contracts/groups-api.yaml` API specification [P]
  * Unit tests from `data-model.md` entity specifications [P]
  * Integration tests from user scenarios in feature spec
  * Acceptance tests from quickstart.md Terraform configurations
- Implementation tasks to make failing tests pass
- Documentation tasks for code comments and usage examples

**Ordering Strategy**:
- **Phase 1**: Test scaffolding (all can run in parallel)
  1. Create unit test files with failing test cases [P]
  2. Create integration test files with mock server setup [P] 
  3. Create acceptance test files with Terraform configurations [P]
  4. Create contract test validation [P]
- **Phase 2**: Core implementation (sequential dependencies)
  5. Implement Group resource CRUD operations
  6. Add validation logic and error handling
  7. Integrate with provider authentication
- **Phase 3**: Test completion and validation
  8. Complete unit test implementations
  9. Complete integration test scenarios
  10. Complete acceptance test cases
  11. Performance benchmarks and optimization

**Estimated Output**: 15-20 numbered, ordered tasks focusing on test-driven development

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [x] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none required)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
