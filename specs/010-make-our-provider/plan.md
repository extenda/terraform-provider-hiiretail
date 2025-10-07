
# Implementation Plan: Provider Distribution and Availability

**Branch**: `010-make-our-provider` | **Date**: October 7, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/010-make-our-provider/spec.md`

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
Enable public distribution of the HiiRetail Terraform provider through HashiCorp's Terraform Registry using GitHub releases with automated publishing via GitHub Actions. This implements the complete provider distribution pipeline following HashiCorp's recommended practices for provider publishing.

## Technical Context
**Language/Version**: Go 1.21+ (existing Terraform provider)  
**Primary Dependencies**: HashiCorp Terraform Plugin Framework, GoReleaser, GitHub Actions  
**Storage**: N/A (distribution artifacts only)  
**Testing**: Go test, terraform-plugin-testing for acceptance tests  
**Target Platform**: Multi-platform (Linux, macOS, Windows - x86_64, ARM64)
**Project Type**: Single project (Terraform provider distribution)  
**Performance Goals**: Release build time <10 minutes, artifact size <50MB per platform  
**Constraints**: Must follow HashiCorp Registry requirements, GPG signing mandatory, semantic versioning required  
**Scale/Scope**: Public distribution via Terraform Registry, automated multi-platform releases

**Implementation Details from User**:
- Use GitHub Actions for automated publishing (not local GoReleaser)
- Follow HashiCorp documentation: https://developer.hashicorp.com/terraform/registry/providers/publishing#creating-a-github-release
- GPG key setup and GitHub secrets already configured
- Focus on GitHub Actions workflow: https://developer.hashicorp.com/terraform/registry/providers/publishing#github-actions-preferred

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Initial Check**:
**Provider Interface Standards**: ✅ PASS - Distribution feature does not modify provider interface
**Resource Implementation Completeness**: ✅ PASS - No new resources, existing resources remain complete  
**Authentication & Security**: ✅ PASS - GPG signing enforced, secure credential handling via GitHub secrets
**Testing & Validation**: ✅ PASS - Existing test suite maintained, release validation via automation
**Documentation & Examples**: ✅ PASS - Registry publication includes complete documentation and examples

**Post-Design Re-evaluation**:
**Provider Interface Standards**: ✅ PASS - GoReleaser config maintains proper binary naming and versioning
**Resource Implementation Completeness**: ✅ PASS - Distribution doesn't affect existing resource implementations
**Authentication & Security**: ✅ PASS - GPG signing workflow enforces signature validation, secure secret handling
**Testing & Validation**: ✅ PASS - Quickstart includes comprehensive validation steps, build verification
**Documentation & Examples**: ✅ PASS - Registry sync includes all docs/, examples/, and generated schemas

**Gate Status**: PASS - All constitutional principles satisfied, design meets standards

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
.github/
├── workflows/
│   └── release.yml           # GitHub Actions release workflow
├── ISSUE_TEMPLATE/
└── copilot-instructions.md

.goreleaser.yml               # GoReleaser configuration
docs/                         # Provider documentation
examples/                     # Usage examples
internal/
├── provider/                 # Provider implementation
└── validation/              # Input validation
main.go                      # Provider entry point
go.mod                       # Go module definition
go.sum                       # Dependency checksums
Makefile                     # Build automation
README.md                    # Installation and usage guide
```

**Structure Decision**: Single project structure selected. This is a Terraform provider distribution feature that adds release automation to the existing provider codebase without modifying the core provider structure.

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
- Generate implementation tasks from contracts and quickstart validation steps
- Configuration file creation tasks from GoReleaser and GitHub Actions contracts [P]
- Validation and testing tasks from quickstart scenarios
- Documentation update tasks for public distribution readiness

**Ordering Strategy**:
- Configuration first: GoReleaser config, GitHub Actions workflow
- Validation second: Local testing, workflow testing
- Documentation third: README updates, example improvements
- Release process fourth: Tag creation, registry verification
- Mark [P] for parallel execution (independent configuration files)

**Estimated Output**: 12-15 numbered, ordered tasks covering:
1. GoReleaser configuration creation
2. GitHub Actions workflow setup  
3. Local build validation
4. Documentation updates for public distribution
5. Release process testing and verification

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
- [x] Complexity deviations documented (N/A - no deviations)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
