
# Implementation Plan: IAM Role Binding Resource Implementation and Testing

**Branch**: `005-add-tests-for` | **Date**: September 30, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-add-tests-for/spec.md`

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
Implement complete IAM Role Binding resource with comprehensive testing infrastructure. Primary requirement: system administrators need to manage role assignments through Terraform with full CRUD operations, enforcing max 10 bindings per resource, tenant isolation, and OAuth2 authentication. Technical approach: Build on existing terraform-plugin-framework patterns from iam_custom_role and iam_group resources, integrating with established mock server infrastructure and acceptance testing framework.

## Technical Context
**Language/Version**: Go 1.21+ (terraform-plugin-framework v1.16.0)  
**Primary Dependencies**: HashiCorp terraform-plugin-framework, terraform-plugin-testing, golang.org/x/oauth2  
**Storage**: RESTful API backend with OAuth2 client credentials authentication  
**Testing**: Go testing framework with terraform-plugin-testing acceptance tests, mock HTTP server integration  
**Target Platform**: Terraform provider for HiiRetail IAM system (cross-platform)
**Project Type**: single (Terraform provider module)  
**Performance Goals**: <2s resource operations, stateless provider design for concurrent usage  
**Constraints**: Max 10 bindings per resource, tenant isolation required, OAuth2 token refresh handling  
**Scale/Scope**: Single resource implementation with comprehensive test coverage (unit, integration, acceptance)

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Provider Interface Standards**: ✅ PASS - Will implement complete CRUD operations using terraform-plugin-framework with proper schema, validation, and state management  
**II. Resource Implementation Completeness**: ✅ PASS - Full Create, Read, Update, Delete operations with API error handling, retry logic, and Terraform state management  
**III. Authentication & Security**: ✅ PASS - Using existing OAuth2 client credentials flow with secure token management and TLS encryption  
**IV. Testing & Validation**: ✅ PASS - Unit tests, integration tests with mock server, and acceptance tests following Terraform conventions  
**V. Documentation & Examples**: ✅ PASS - Will include comprehensive documentation with working examples and usage patterns

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
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->
```
internal/
├── provider/
│   ├── provider_hiiretail_iam/
│   │   └── hiiretail_iam_provider_gen.go
│   └── resource_iam_role_binding/
│       ├── iam_role_binding_resource_gen.go (schema)
│       ├── iam_role_binding_resource.go (implementation)
│       └── iam_role_binding_resource_test.go (unit tests)
├── client/
│   └── oauth2_client.go
└── models/
    └── role_binding_models.go

acceptance_tests/
├── iam_role_binding_resource_test.go
└── mock_server_test.go
```

**Structure Decision**: Single Terraform provider project following HashiCorp's terraform-plugin-framework conventions. Resource implementation goes in `internal/provider/resource_iam_role_binding/` alongside existing resources. Acceptance tests in dedicated `acceptance_tests/` directory with mock server integration.

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
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each contract → contract test task [P]
- Each entity → model creation task [P] 
- Each user story → integration test task
- Implementation tasks to make tests pass

**Ordering Strategy**:
- TDD order: Tests before implementation 
- Dependency order: Models before services before UI
- Mark [P] for parallel execution (independent files)

**Estimated Output**: 25-30 numbered, ordered tasks in tasks.md

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
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
