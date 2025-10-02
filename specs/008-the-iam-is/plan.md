
# Implementation Plan: Multi-API Provider User Experience Enhancement

**Branch**: `008-the-iam-is` | **Date**: October 2, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-the-iam-is/spec.md`

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
Transform the HiiRetail Terraform provider to feel familiar and intuitive like major cloud providers (e.g., GCP) when managing multiple APIs. Key changes include renaming the provider from "hiiretail-iam" to "hiiretail" and implementing consistent naming conventions (api_resource format: hiiretail_iam_role_binding, hiiretail_ccc_kind, etc.). The provider must organize resources in a logical, service-based structure with comprehensive documentation, consistent authentication patterns, and getting started guides that reduce time-to-first-success for DevOps engineers and infrastructure developers.

## Technical Context
**Language/Version**: Go 1.19+ (Terraform Plugin Framework requirement)
**Primary Dependencies**: HashiCorp Terraform Plugin Framework, golang.org/x/oauth2, existing HiiRetail APIs
**Storage**: N/A (provider manages external API resources, not local storage)
**Testing**: Go standard testing, Terraform acceptance tests, integration tests with real APIs
**Target Platform**: Cross-platform (Windows, macOS, Linux) - Terraform provider binary
**Project Type**: single (Terraform provider with multiple service modules)
**Performance Goals**: Handle 100+ concurrent API requests, <500ms response time for plan operations
**Constraints**: Must maintain backward compatibility, follow Terraform naming conventions, secure OAuth2 implementation
**Scale/Scope**: Support 10+ HiiRetail APIs, 50+ resource types, enterprise-scale infrastructure management

**User-Provided Implementation Details**: 
- Provider name changes from "hiiretail-iam" to "hiiretail"
- Resource naming convention: api_resource format (hiiretail_iam_role_binding, hiiretail_ccc_kind)
- Folder and naming structures need reorganization to support multi-API architecture
- Must follow conventions similar to GCP Terraform provider for discoverability## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Provider Interface Standards**: ✅ PASS - Refactoring maintains complete CRUD operations, enhances schema organization, and improves error handling
**II. Resource Implementation Completeness**: ✅ PASS - All existing resources will be maintained with full implementations, organizational changes don't affect functionality  
**III. Authentication & Security**: ✅ PASS - OAuth2 implementation remains secure, naming changes don't affect credential handling
**IV. Testing & Validation**: ✅ PASS - Comprehensive test suite will be updated for new naming conventions, acceptance tests maintained
**V. Documentation & Examples**: ✅ PASS - Enhanced documentation is core requirement, working examples will be updated for new structure

**Gate Status**: PASS - All constitutional principles are maintained and enhanced by this refactoring

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
# Multi-API Terraform Provider Structure
main.go                           # Provider entry point - renamed from hiiretail-iam
go.mod                            # Module renamed to hiiretail provider
go.sum

internal/
├── provider/
│   ├── provider.go               # Main provider configuration - supports all APIs
│   ├── iam/                      # IAM service module  
│   │   ├── resources/
│   │   │   ├── group.go          # hiiretail_iam_group
│   │   │   ├── custom_role.go    # hiiretail_iam_custom_role
│   │   │   └── role_binding.go   # hiiretail_iam_role_binding
│   │   └── data_sources/
│   │       └── groups.go         # hiiretail_iam_groups
│   ├── ccc/                      # Future API service module
│   │   ├── resources/
│   │   │   └── kind.go           # hiiretail_ccc_kind  
│   │   └── data_sources/
│   └── shared/
│       ├── auth/                 # Shared OAuth2 implementation
│       ├── client/               # Common HTTP client
│       └── validators/           # Common validation logic

docs/                             # Enhanced documentation structure
├── guides/
│   ├── getting-started.md
│   ├── authentication.md
│   └── migration-from-gcp.md
├── resources/
│   ├── iam/
│   └── ccc/
└── examples/
    ├── basic/
    ├── multi-service/
    └── migration/

tests/
├── acceptance/                   # Terraform acceptance tests
│   ├── iam/
│   └── ccc/
├── integration/                  # API integration tests  
└── unit/                        # Unit tests for provider logic
```

**Structure Decision**: Single multi-API provider project with service-based module organization. Each HiiRetail API (IAM, CCC, etc.) has its own module under `internal/provider/` with resources following the `hiiretail_{api}_{resource}` naming convention. Shared authentication and client logic supports all APIs from a central location.

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
- Provider refactoring tasks: Rename provider, update resource names, maintain compatibility
- Documentation tasks: Generate new docs structure, migration guides, examples
- Testing tasks: Update acceptance tests, create service module tests
- Service module tasks: Create IAM service module, implement extensibility patterns

**Ordering Strategy**:
1. **Foundation Tasks**: Provider rename, core structure changes [P]
2. **Service Module Tasks**: Create IAM service module, resource refactoring [Sequential]
3. **Testing Tasks**: Update test suite, create new acceptance tests [P after foundation]
4. **Documentation Tasks**: Generate docs, examples, migration guides [P after implementation]
5. **Validation Tasks**: End-to-end testing, performance validation [Sequential, final]

**Specific Task Categories**:
- **Provider Tasks**: Update provider name, configuration schema, resource registration
- **Resource Tasks**: Rename resources, update schemas, maintain backward compatibility  
- **Service Tasks**: Create service module architecture, implement IAM service
- **Test Tasks**: Update acceptance tests, create service tests, integration tests
- **Documentation Tasks**: Generate docs, create examples, write migration guides
- **Migration Tasks**: Create migration tooling, backward compatibility layers

**Estimated Output**: 35-40 numbered, ordered tasks in tasks.md

**Critical Dependencies**:
- Provider rename must complete before resource updates
- Service modules must be created before resource migration
- Tests must be updated alongside implementation changes
- Documentation should be generated from final implementation

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
- [x] Complexity deviations documented (none required)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
