# Contract: IAM Service Coverage

## Purpose
Define the contract for achieving 90%+ code coverage in `internal/provider/iam/service.go`.

## Endpoints/Methods to Cover
- ListGroups
- GetGroup
- CreateGroup
- UpdateGroup
- DeleteGroup
- ListRoles
- GetRole
- CreateCustomRole
- GetCustomRole
- UpdateCustomRole
- DeleteCustomRole
- ListRoleBindings
- GetRoleBinding
- CreateRoleBinding
- UpdateRoleBinding
- DeleteRoleBinding
- SetResource
- GetResource
- DeleteResource
- GetResources
- AddRoleToGroup

## Test Requirements
- Each method must have unit tests covering normal and error scenarios
- All branches, conditions, and error handling must be exercised
- Mocks must be used for API client and network calls
- Table-driven tests for edge cases and parameter validation
- Coverage must be measured and reported for each method

## Acceptance Criteria
- Coverage report for `service.go` shows 90% or higher
- All critical paths and error handling are tested
- Exclusions are documented and justified

## Example Test Contract
```go
func TestService_ListGroups(t *testing.T) {
    // Test normal case
    // Test error case (API failure)
    // Test edge case (empty result)
}
```

## Next Steps
- Implement contract tests for each method
- Validate coverage and update contract as needed
