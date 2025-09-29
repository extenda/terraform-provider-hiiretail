# Quickstart Guide: IAM Group Resource Testing

**Created**: September 28, 2025  
**Feature**: 002-ensure-that-the

## Overview
This guide demonstrates how to test the IAM Group resource in the HiiRetail Terraform Provider, covering unit tests, integration tests, and acceptance tests.

## Prerequisites

### Development Environment
- Go 1.21 or later
- Terraform 1.5.0 or later
- Git for version control
- Make (optional, for build automation)

### Dependencies
```bash
# Install test dependencies
go mod tidy
go install github.com/stretchr/testify/assert@latest
```

### Environment Setup
```bash
# Clone the repository
git clone https://github.com/extenda/hiiretail-terraform-providers
cd hiiretail-terraform-providers/iam

# Verify Go environment
go version
terraform version
```

## Test Structure Overview

The Group resource tests are organized in three layers:

```
internal/provider/resource_iam_group/
├── iam_group_resource.go              # Resource implementation
├── iam_group_resource_test.go         # Unit tests
└── iam_group_integration_test.go      # Integration tests

acceptance_tests/
└── group_resource_test.go             # Acceptance tests
```

## Running Tests

### Unit Tests
Unit tests validate schema, validation logic, and business rules:

```bash
# Run unit tests for Group resource
go test ./internal/provider/resource_iam_group -v

# Run with coverage
go test ./internal/provider/resource_iam_group -v -cover

# Run specific test
go test ./internal/provider/resource_iam_group -v -run TestGroupResourceSchema
```

Expected output:
```
=== RUN   TestGroupResourceSchema
--- PASS: TestGroupResourceSchema (0.00s)
=== RUN   TestGroupValidation
--- PASS: TestGroupValidation (0.01s)
=== RUN   TestGroupModel
--- PASS: TestGroupModel (0.00s)
PASS
coverage: 95.2% of statements
```

### Integration Tests
Integration tests use mock HTTP servers to test API interactions:

```bash
# Run integration tests
go test ./internal/provider/resource_iam_group -v -run Integration

# Run with detailed HTTP logging
DEBUG=1 go test ./internal/provider/resource_iam_group -v -run Integration
```

Expected output:
```
=== RUN   TestGroupResourceIntegration
=== RUN   TestGroupResourceIntegration/Create_group_success
--- PASS: TestGroupResourceIntegration/Create_group_success (0.02s)
=== RUN   TestGroupResourceIntegration/Update_group_success
--- PASS: TestGroupResourceIntegration/Update_group_success (0.01s)
=== RUN   TestGroupResourceIntegration/Delete_group_success
--- PASS: TestGroupResourceIntegration/Delete_group_success (0.01s)
--- PASS: TestGroupResourceIntegration (0.05s)
```

### Acceptance Tests
Acceptance tests validate end-to-end Terraform workflows:

```bash
# Set up test environment variables
export TF_ACC=1
export HIIRETAIL_TENANT_ID=test-tenant
export HIIRETAIL_CLIENT_ID=test-client
export HIIRETAIL_CLIENT_SECRET=test-secret

# Run acceptance tests
go test ./acceptance_tests -v -timeout 10m
```

Expected output:
```
=== RUN   TestAccGroupResource_basic
--- PASS: TestAccGroupResource_basic (3.45s)
=== RUN   TestAccGroupResource_update
--- PASS: TestAccGroupResource_update (4.12s)
=== RUN   TestAccGroupResource_import
--- PASS: TestAccGroupResource_import (2.34s)
```

## Test Scenarios

### Basic Group Creation
Test the simplest group configuration:

```hcl
resource "hiiretail_iam_group" "test" {
  name = "test-group"
}
```

**Validation Points:**
- Group is created with specified name
- ID is auto-generated
- Status is set to "active"
- Description is optional

### Group with Description
Test group creation with optional description:

```hcl
resource "hiiretail_iam_group" "test" {
  name        = "developers"
  description = "Development team members"
}
```

**Validation Points:**
- Description is properly stored
- Description can be updated without resource recreation
- Empty description is handled correctly

### Group Updates
Test group modification scenarios:

```hcl
resource "hiiretail_iam_group" "test" {
  name        = "senior-developers"  # Updated name
  description = "Senior development team members with elevated access"  # Updated description
}
```

**Validation Points:**
- Name changes trigger proper API calls
- Description updates work correctly
- Terraform state is maintained accurately

### Multi-tenant Scenario
Test groups in different tenant contexts:

```hcl
resource "hiiretail_iam_group" "tenant_a" {
  name      = "developers"
  tenant_id = "tenant-a"
}

resource "hiiretail_iam_group" "tenant_b" {
  name      = "developers"
  tenant_id = "tenant-b"
}
```

**Validation Points:**
- Same group name allowed in different tenants
- Tenant isolation is properly maintained
- Cross-tenant conflicts are prevented

## Error Testing Scenarios

### Validation Errors
Test various validation failure cases:

```bash
# Test oversized name
terraform plan # Should fail with validation error

# Test empty required field
terraform plan # Should fail with required field error
```

### API Error Handling
Test API error response handling:

- **401 Unauthorized**: Invalid credentials
- **403 Forbidden**: Insufficient permissions  
- **404 Not Found**: Group doesn't exist for read/update/delete
- **409 Conflict**: Duplicate group name in same tenant
- **500 Server Error**: Upstream service failures

### Network Error Scenarios
Test network failure handling:

- Connection timeouts
- DNS resolution failures
- Temporary service unavailability
- Retry logic validation

## Performance Testing

### Benchmark Tests
Run performance benchmarks:

```bash
# Run resource operation benchmarks
go test ./internal/provider/resource_iam_group -bench=. -benchmem

# Expected output:
# BenchmarkGroupCreate-8    100    15.2ms/op    2.1MB/op
# BenchmarkGroupRead-8      500     3.4ms/op    0.8MB/op
# BenchmarkGroupUpdate-8    200     8.7ms/op    1.5MB/op
# BenchmarkGroupDelete-8    300     5.1ms/op    0.6MB/op
```

### Load Testing
Test multiple concurrent operations:

```bash
# Run concurrent group operations
go test ./internal/provider/resource_iam_group -v -run TestConcurrentOperations
```

## Debugging and Troubleshooting

### Test Debugging
Enable verbose test output:

```bash
# Enable debug logging
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log

# Run tests with HTTP tracing
export DEBUG_HTTP=1
go test ./internal/provider/resource_iam_group -v
```

### Common Issues

#### Test Failures
1. **Authentication errors**: Verify test credentials are set correctly
2. **Network timeouts**: Check if mock servers are running properly
3. **State inconsistencies**: Ensure proper test cleanup

#### Performance Issues
1. **Slow tests**: Check for unnecessary HTTP calls
2. **Memory leaks**: Verify proper resource cleanup in tests
3. **Race conditions**: Ensure proper test isolation

### Test Data Cleanup
Ensure tests don't leave orphaned data:

```bash
# Verify test isolation
go test ./internal/provider/resource_iam_group -v -count=5

# Check for resource leaks
go test ./internal/provider/resource_iam_group -race
```

## Continuous Integration

### GitHub Actions Integration
Example CI configuration for automated testing:

```yaml
- name: Run Unit Tests
  run: go test ./internal/provider/resource_iam_group -v -cover

- name: Run Integration Tests  
  run: go test ./internal/provider/resource_iam_group -v -run Integration

- name: Run Acceptance Tests
  env:
    TF_ACC: 1
    HIIRETAIL_TENANT_ID: ${{ secrets.TEST_TENANT_ID }}
    HIIRETAIL_CLIENT_ID: ${{ secrets.TEST_CLIENT_ID }}
    HIIRETAIL_CLIENT_SECRET: ${{ secrets.TEST_CLIENT_SECRET }}
  run: go test ./acceptance_tests -v -timeout 10m
```

### Test Coverage Requirements
- **Unit tests**: Minimum 90% code coverage
- **Integration tests**: All API endpoints covered
- **Acceptance tests**: All user scenarios validated
- **Error paths**: All error conditions tested

## Validation Checklist

Before completing Group resource testing:

- [ ] All unit tests pass with >90% coverage
- [ ] Integration tests cover all CRUD operations
- [ ] Acceptance tests validate end-to-end workflows
- [ ] Error scenarios are properly handled
- [ ] Performance benchmarks meet requirements
- [ ] Multi-tenant scenarios work correctly
- [ ] Concurrent operations are safe
- [ ] Test cleanup prevents resource leaks
- [ ] CI/CD pipeline integration works
- [ ] Documentation is complete and accurate

## Next Steps

After completing Group resource testing:

1. **Code Review**: Submit tests for peer review
2. **Performance Review**: Validate benchmark results
3. **Integration**: Merge tests with main codebase
4. **Documentation**: Update provider documentation
5. **Release**: Include in next provider version

---
*Quickstart guide completed: September 28, 2025*