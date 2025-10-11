# Research: Achieving 90%+ Code Coverage for internal/provider/iam/service.go

## Decision
- Target file: `internal/provider/iam/service.go`
- Target coverage: 90% or higher (as measured by Go coverage tools)
- Approach: Add/expand unit tests, use mocks for external dependencies, cover all branches and error paths, document exclusions.

## Rationale
- High code coverage ensures reliability, maintainability, and easier refactoring.
- Coverage must include all CRUD operations, error handling, and edge cases in the IAM service client.
- External API calls and side effects will be mocked to isolate logic.
- Exclusions (e.g., generated code, unreachable branches) will be documented and justified.

## Alternatives Considered
- **Integration-only testing**: Not sufficient for coverage of all code paths, especially error handling and edge cases.
- **Partial coverage**: Would not meet reliability and maintainability goals.
- **Manual testing**: Not measurable or automatable for CI/CD enforcement.

## Best Practices
- Use Go's built-in testing and coverage tools (`go test -cover`, `go tool cover`)
- Mock external dependencies (API client, network calls) to isolate service logic
- Write table-driven tests for all methods, including error scenarios
- Use subtests for edge cases and boundary conditions
- Document and justify any code intentionally excluded from coverage
- Integrate coverage reporting into CI pipeline and enforce threshold

## Unknowns/Clarifications
- None: All requirements and constraints are clear from the feature spec and constitution.

## Next Steps
- Proceed to Phase 1: Design & Contracts
