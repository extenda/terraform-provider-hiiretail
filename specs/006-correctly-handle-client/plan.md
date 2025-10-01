
# Implementation Plan: Correctly Handle Client Credentials

**Branch**: `006-correctly-handle-client` | **Date**: October 1, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-correctly-handle-client/spec.md`

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

**IMPORTANT**: The /plan command STOPS at step 9. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Progress Tracking

### ✅ Phase 0: Research Complete
- OAuth2 discovery protocol research complete
- OCMS integration requirements documented
- Current provider authentication issues identified
- Technical decisions and implementation strategy defined

### ✅ Phase 1: Design Complete  
- OAuth2 authentication contract defined
- Data models for discovery, authentication, and error handling created
- Quickstart guide with implementation steps provided
- GitHub Copilot instructions updated with OAuth2 guidelines

### 🔄 Phase 2: Task Generation
- Ready for /tasks command to generate detailed implementation tasks
- All design artifacts completed and validated

**GATE STATUS**: ✅ Ready for Phase 2 Task Generation

## Phase 2 Task Generation Approach

The task generation will create detailed implementation tasks based on the research and design artifacts:

### Task Categories
1. **OAuth2 Discovery Implementation**
   - Create discovery client with caching
   - Implement endpoint validation
   - Add fallback configuration support

2. **Authentication Client Enhancement**  
   - Enhance existing OAuth2 client credentials flow
   - Add comprehensive error handling
   - Implement retry logic with exponential backoff

3. **Provider Configuration Updates**
   - Update provider schema with new OAuth2 parameters
   - Enhance Configure method with discovery integration
   - Add configuration validation

4. **Testing Implementation**
   - Unit tests for OAuth2 components
   - Integration tests with OCMS endpoints
   - Provider configuration test scenarios

5. **Documentation Updates**
   - Update provider documentation with OAuth2 examples
   - Add troubleshooting guide
   - Create migration guide for enhanced authentication

### Task Dependencies
- Discovery client → Authentication client → Provider integration
- Base implementation → Error handling → Testing
- Core functionality → Documentation → Final validation

### Implementation Order
1. OAuth2 discovery foundation
2. Authentication client enhancement  
3. Provider integration
4. Comprehensive testing
5. Documentation and examples

## Summary
Enhance the HiiRetail IAM Terraform provider to properly handle OAuth2 client credentials authentication using the Hii Retail OAuth Client Management Service (OCMS). The provider must securely authenticate with the IAM API, handle token lifecycle management (acquisition, refresh, expiration), and provide robust error handling for authentication failures. This ensures reliable Terraform operations with proper credential validation and automatic token management.

## Technical Context
**Language/Version**: Go 1.21+ (Terraform Plugin Framework v1.16.0)
**Primary Dependencies**: HashiCorp Terraform Plugin Framework, golang.org/x/oauth2, HTTP client with OAuth2 support
**Storage**: Token caching in memory (no persistent storage required)
**Testing**: Go testing framework with unit tests, integration tests, and Terraform acceptance tests
**Target Platform**: Cross-platform (Linux, macOS, Windows) - Terraform provider binary
**Project Type**: Single binary Terraform provider with OAuth2 client integration
**Performance Goals**: <500ms token acquisition, <100ms token validation, efficient token reuse
**Constraints**: Secure credential handling (no logging of secrets), TLS-only communication, thread-safe token management
**Scale/Scope**: Support concurrent Terraform operations, handle token refresh during long-running operations

**OAuth2 Endpoint**: https://auth.retailsvc.com/.well-known/openid-configuration
**Documentation**: https://developer.hiiretail.com/docs/ocms/public/concepts/oauth2-authentication/
**Service**: Hii Retail OAuth Client Management Service (OCMS) for token acquisition and management## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ I. Provider Interface Standards
- OAuth2 authentication enhancement maintains existing provider interface
- No breaking changes to resource CRUD operations
- Error handling improvements align with HashiCorp conventions
- Schema definitions remain consistent with Plugin Framework standards

### ✅ III. Authentication & Security (NON-NEGOTIABLE)
- **CRITICAL**: OAuth2 client credentials flow implementation required
- Secure token management with proper refresh handling
- Credential validation before API operations
- Sensitive configuration values properly marked in schemas
- No credential exposure in logs or debug output
- TLS-only communication enforced

### ✅ IV. Testing & Validation
- Unit tests required for OAuth2 authentication logic
- Integration tests with actual OCMS endpoints
- Acceptance tests for provider configuration scenarios
- Mock tests insufficient - real API validation required

### ✅ V. Documentation & Examples
- Provider configuration examples with OAuth2 setup
- Authentication troubleshooting documentation
- Migration guide for credential handling changes
- Usage patterns for different credential sources

**GATE STATUS**: ✅ PASS - No constitutional violations detected. Authentication security requirements align with NON-NEGOTIABLE principle III.

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
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

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
- [ ] Phase 0: Research complete (/plan command)
- [ ] Phase 1: Design complete (/plan command)
- [ ] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [ ] Initial Constitution Check: PASS
- [ ] Post-Design Constitution Check: PASS
- [ ] All NEEDS CLARIFICATION resolved
- [ ] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
