
# Implementation Plan: OAuth2 Authentication with Environment-Specific Endpoints

**Branch**: `007-oauth2-authentication-with` | **Date**: October 1, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-oauth2-authentication-with/spec.md`

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
Implement OAuth2 client credentials authentication for the HiiRetail IAM Terraform provider with automatic environment detection and endpoint routing. The system authenticates against auth.retailsvc.com and automatically routes to iam-api.retailsvc.com for Live Tenants or iam-api.retailsvc-test.com for Test Tenants based on tenant ID parsing, with mock server override support for testing.

## Technical Context
**Language/Version**: Go 1.21+ (matches existing Terraform provider)  
**Primary Dependencies**: Terraform Plugin Framework v1.16.0, golang.org/x/oauth2 v0.30.0  
**Storage**: N/A (stateless authentication service)  
**Testing**: Go testing framework with testify, mock OAuth2 server for integration tests  
**Target Platform**: Cross-platform (Linux, macOS, Windows) - Terraform provider binary
**Project Type**: Single project (Terraform provider extension)  
**Performance Goals**: <100ms token acquisition, <10ms endpoint resolution  
**Constraints**: Must integrate with existing provider without breaking changes, secure credential handling  
**Scale/Scope**: Enterprise-scale multi-tenant IAM operations, thousands of concurrent provider instances

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Provider Interface Standards**: ✅ PASS - OAuth2 authentication extends existing provider without changing CRUD operations
**II. Resource Implementation Completeness**: ✅ PASS - Authentication layer supports all existing resources  
**III. Authentication & Security (NON-NEGOTIABLE)**: ✅ PASS - OAuth2 client credentials with proper token management, TLS enforcement
**IV. Testing & Validation**: ✅ PASS - Unit tests, integration tests with mock OAuth2 server, acceptance tests planned
**V. Documentation & Examples**: ✅ PASS - Complete documentation and usage examples to be provided

**Gate Status**: PASS - All constitutional requirements satisfied

**Post-Design Re-evaluation**: 
- ✅ OAuth2 implementation maintains all CRUD operations for existing resources
- ✅ Security requirements exceeded with comprehensive credential protection
- ✅ Testing strategy includes unit, integration, and acceptance tests with mock server
- ✅ Documentation includes contracts, quickstart, and troubleshooting guides
- ✅ Code quality enforced through TLS requirements and error handling standards

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
internal/provider/
├── auth/                    # NEW: OAuth2 authentication package
│   ├── client.go           # OAuth2 client implementation
│   ├── discovery.go        # Endpoint discovery and routing
│   ├── errors.go           # Authentication error types
│   └── validation.go       # Credential validation
├── provider.go             # MODIFIED: Enhanced with OAuth2 integration
├── resource_iam_custom_role/
├── resource_iam_group/
├── resource_iam_role/
└── resource_iam_role_binding/

tests/
├── auth/                   # NEW: OAuth2 authentication tests
│   ├── client_test.go      # OAuth2 client unit tests
│   ├── discovery_test.go   # Endpoint discovery tests
│   ├── integration_test.go # Mock server integration tests
│   └── validation_test.go  # Credential validation tests
└── provider_test.go        # MODIFIED: Enhanced provider tests

demo/                       # NEW: Usage examples
└── oauth2_demo.go          # OAuth2 configuration examples
```

**Structure Decision**: Single project extension - OAuth2 authentication is added as a new internal package that integrates with the existing Terraform provider structure. This maintains backward compatibility while adding secure authentication capabilities.

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
- Generate OAuth2-specific tasks from Phase 1 design docs
- **Contract Tests**: OAuth2 token endpoint tests, IAM API authentication tests [P]
- **Entity Implementation**: AuthClient, TokenCache, EndpointResolver, AuthError types [P]
- **Integration Tests**: Mock server tests, environment detection tests, error handling tests
- **Provider Integration**: Enhanced provider.go with OAuth2 configuration and buildAuthConfig()
- **Security Implementation**: TLS enforcement, credential protection, token refresh logic
- **Demo and Documentation**: OAuth2 demo program, usage examples, troubleshooting guides

**Ordering Strategy**:
- **Phase 1 (Setup)**: Dependencies, package structure, configuration types [P]
- **Phase 2 (TDD Tests)**: All test files before any implementation [P]
- **Phase 3 (Core Implementation)**: AuthClient, TokenCache, EndpointResolver in dependency order
- **Phase 4 (Provider Integration)**: Enhanced provider with OAuth2, HTTP client integration
- **Phase 5 (Polish)**: Demo program, documentation, security audit, performance testing

**OAuth2-Specific Considerations**:
- Security tasks prioritized (credential protection, TLS enforcement)
- Token management tasks with thread-safety requirements
- Environment detection with comprehensive test coverage
- Mock server integration for testing without external dependencies

**Estimated Output**: 22-25 numbered, ordered tasks with OAuth2 security and testing focus

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
