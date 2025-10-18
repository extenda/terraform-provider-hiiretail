# Quickstart: Improving Code Coverage for IAM Service

## Steps

1. Run Go coverage analysis:
   ```bash
   go test -coverprofile=coverage.out ./internal/provider/iam
   go tool cover -func=coverage.out
   go tool cover -html=coverage.out
   ```
2. Identify untested code paths in `service.go` (functions, branches, error handling).
3. Add/expand unit tests in `internal/provider/iam/service_test.go`:
   - Cover all CRUD methods and error scenarios
   - Use mocks for API client and network calls
   - Table-driven tests for edge cases
4. Run tests and verify coverage meets/exceeds 90%:
   ```bash
   go test -cover ./internal/provider/iam
   ```
5. Document any exclusions and justify in code comments and coverage report.
6. Integrate coverage check into CI pipeline (update workflow if needed).

## Validation
- Coverage report shows 90%+ for `service.go`
- All critical paths, error handling, and edge cases are tested
- CI pipeline enforces coverage threshold

## Troubleshooting
- If coverage is below target, review untested branches and add tests
- Use mocks to isolate logic and avoid external dependencies
- Document and justify any unreachable or excluded code
