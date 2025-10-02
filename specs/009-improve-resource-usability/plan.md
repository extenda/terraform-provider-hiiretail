
# Implementation Plan: Improve Resource Usability

**Branch**: `009-improve-resource-usability` | **Date**: October 2, 2025 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/009-improve-resource-usability/spec.md`

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
Improve the usability of HiiRetail Terraform provider IAM resources through enhanced validation, error messages, documentation, and user experience. The primary requirement is to provide clear, actionable feedback for DevOps engineers configuring IAM resources (groups, custom roles, role bindings) with comprehensive examples and troubleshooting guidance. Technical approach focuses on provider-level improvements to schema validation, error handling, and documentation based on the concrete resource implementations defined in simple_test.tf.

## Technical Context
**Language/Version**: Go 1.19+ (Terraform Plugin Framework)  
**Primary Dependencies**: HashiCorp Terraform Plugin Framework, OAuth2 client libraries, HiiRetail IAM API clients  
**Storage**: N/A (provider manages remote API resources)  
**Testing**: Go testing framework, Terraform acceptance testing framework  
**Target Platform**: Cross-platform (provider runs where Terraform runs)
**Project Type**: single (Terraform provider with multiple resource types)  
**Performance Goals**: <5s resource CRUD operations, efficient OAuth2 token management  
**Constraints**: Backward compatibility with existing configurations, secure credential handling, proper Terraform state management  
**Scale/Scope**: Support for enterprise IAM scenarios with hundreds of groups/roles, clear validation for complex role binding configurations  

**Resource Analysis from simple_test.tf**:
- **hiiretail_iam_group**: Basic resource with name and description fields
- **hiiretail_iam_custom_role**: Complex resource with permissions array, stage field, title/description
- **hiiretail_iam_role_binding**: Most complex with role references, members array, optional conditions
- **Data sources**: hiiretail_iam_groups and hiiretail_iam_roles for querying existing resources
- **Provider config**: OAuth2 with client credentials, configurable endpoints, scopes, timeouts, retries

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Core Principle Alignment**:
- ✅ **Provider Interface Standards**: Improving existing CRUD operations with better validation and error handling
- ✅ **Resource Implementation Completeness**: Enhancing existing complete implementations with usability improvements
- ✅ **Authentication & Security**: No changes to OAuth2 implementation, maintaining security standards
- ✅ **Testing & Validation**: Will add comprehensive validation tests for improved user experience
- ✅ **Documentation & Examples**: Primary focus on improving documentation and providing working examples

**Technical Standards Compliance**:
- ✅ **Go Best Practices**: Improvements follow Terraform Plugin Framework conventions
- ✅ **Error Handling**: Core focus on improving error messages and HTTP status mapping
- ✅ **OpenAPI Alignment**: Validation improvements based on existing API specifications

**Quality Assurance**:
- ✅ **Non-Breaking Changes**: All improvements maintain backward compatibility
- ✅ **Test Coverage**: Enhanced validation requires expanded test coverage
- ✅ **Documentation Updates**: Primary deliverable includes comprehensive documentation

**Status**: PASS - No constitutional violations identified. All improvements align with provider development standards.

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
iam/
├── internal/
│   └── provider/
│       ├── provider_hiiretail_iam/
│       │   └── hiiretail_iam_provider_gen.go
│       ├── resource_iam_custom_role/
│       │   └── iam_custom_role_resource_gen.go
│       ├── resource_iam_group/
│       │   └── iam_group_resource_gen.co
│       ├── resource_iam_role/
│       │   └── iam_role_resource_gen.go
│       └── resource_iam_role_binding/
│           └── iam_role_binding_resource_gen.go
├── examples/
│   ├── basic/
│   ├── enterprise/
│   └── troubleshooting/
├── docs/
│   ├── resources/
│   └── guides/
└── tests/
    ├── acceptance/
    ├── unit/
    └── validation/
```

**Structure Decision**: Single project structure based on existing Terraform provider codebase. The improvements will enhance the existing resource implementations in `internal/provider/` with better validation, error handling, and documentation. New directories for examples and enhanced documentation will be added to support the usability improvements.

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
- **Validation Enhancement Tasks**: Implement custom validators for each resource type [P]
- **Error Message Tasks**: Create enhanced error message system [P]  
- **Schema Enhancement Tasks**: Update existing resource schemas with improved validation
- **Documentation Tasks**: Create comprehensive examples and troubleshooting guides [P]
- **Testing Tasks**: Contract tests for validation, integration tests for user scenarios

**Resource-Specific Task Breakdown**:
- **Group Resource**: Name validation, uniqueness checks, description validation
- **Custom Role Resource**: Permission format validation, stage validation, title requirements
- **Role Binding Resource**: Reference validation, member format validation, condition parsing
- **Cross-Resource**: Reference resolution, dependency validation, circular dependency detection

**Ordering Strategy**:
- **Phase A**: Enhanced error system and base validators (foundation)
- **Phase B**: Resource-specific schema enhancements [P] (parallel per resource)
- **Phase C**: Cross-resource validation and reference resolution
- **Phase D**: Documentation and examples [P] (parallel with implementation)
- **Phase E**: Comprehensive testing and validation

**Estimated Output**: 35-40 numbered, ordered tasks focusing on provider usability improvements

**Key Dependencies**:
- Enhanced error system must be implemented before resource-specific validators
- Schema enhancements can be done in parallel per resource type
- Cross-resource validation requires completed individual resource validators
- Documentation tasks can run parallel with implementation tasks

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
- [x] Phase 0: Research complete (/plan command) ✅
- [x] Phase 1: Design complete (/plan command) ✅  
- [x] Phase 2: Task planning approach described (/plan command) ✅
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS ✅
- [x] Post-Design Constitution Check: PASS ✅
- [x] All NEEDS CLARIFICATION resolved ✅
- [x] Complexity deviations documented: N/A ✅

**Constitution Re-check After Phase 1**:
- ✅ **Provider Interface Standards**: Design maintains CRUD operations, enhances validation
- ✅ **Resource Implementation Completeness**: Builds on existing complete implementations
- ✅ **Authentication & Security**: No changes to security model, maintains OAuth2 standards
- ✅ **Testing & Validation**: Enhanced testing framework for validation improvements
- ✅ **Documentation & Examples**: Core deliverable addresses documentation requirements

**Status**: All constitutional requirements satisfied. Ready for Phase 3 task generation.

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
