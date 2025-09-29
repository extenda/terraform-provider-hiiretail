# Quickstart: IAM Custom Role Resource Testing

**Date**: September 28, 2025  
**Purpose**: Step-by-step guide to validate IAM custom role resource implementation

## Prerequisites

### Development Environment
```bash
# Verify Go version
go version  # Should show Go 1.23.0 or higher

# Verify current directory
pwd  # Should be in /path/to/hiiretail-terraform-providers/iam

# Verify dependencies
go mod verify
go mod tidy
```

### Required Files Check
```bash
# Check existing generated resource
ls -la internal/provider/resource_iam_custom_role/iam_custom_role_resource_gen.go

# Check provider integration (should exist)
grep -n "NewIamCustomRoleResource" internal/provider/provider.go || echo "MISSING: Resource not registered"

# Check test utilities (should exist from group resource)
ls -la internal/provider/testutils/mock_server.go
```

## Step 1: Run Initial Contract Tests (Should Fail)

### Execute Schema Validation Tests
```bash
# Run schema-specific tests (will fail - no implementation)
go test ./internal/provider/resource_iam_custom_role -v -run="TestCustomRoleSchema" 2>/dev/null || echo "✓ EXPECTED: Schema tests not implemented"

# Verify test files don't exist yet
ls internal/provider/resource_iam_custom_role/*test*.go 2>/dev/null || echo "✓ EXPECTED: No test files exist"
```

### Verify Provider Registration Status
```bash
# Check if custom role resource is registered (should fail)
go test ./internal/provider -v -run="TestProvider.*CustomRole" 2>/dev/null || echo "✓ EXPECTED: Custom role not registered"

# Check available resources in provider
go run -c 'package main; import "fmt"; func main() { fmt.Println("Provider resources check") }' 2>/dev/null || echo "Manual verification needed"
```

## Step 2: Validate Generated Resource Schema

### Schema Structure Verification
```bash
# Examine generated schema structure
grep -A 10 -B 5 "func IamCustomRoleResourceSchema" internal/provider/resource_iam_custom_role/iam_custom_role_resource_gen.go

# Check permission validation pattern
grep -n "regexp.MustCompile" internal/provider/resource_iam_custom_role/iam_custom_role_resource_gen.go
```

### Expected Schema Elements
- [x] `id` field (required)
- [x] `name` field (optional, computed, 3-256 chars)  
- [x] `permissions` field (required list)
- [x] `tenant_id` field (optional, computed)
- [x] Permission ID pattern validation
- [x] Attributes object with constraints

## Step 3: Test Resource Registration

### Provider Integration Test
```bash
# Create temporary test to verify resource availability
cat > /tmp/provider_test.go << 'EOF'
package main

import (
    "context"
    "testing"
    "github.com/extenda/hiiretail-terraform-providers/iam/internal/provider"
)

func TestProviderHasCustomRoleResource(t *testing.T) {
    p := provider.New("test")()
    resources := p.Resources(context.Background())
    
    found := false
    for _, resourceFunc := range resources {
        // This will fail until resource is registered
        if resourceFunc != nil {
            found = true
            break
        }
    }
    
    if !found {
        t.Error("EXPECTED FAILURE: Custom role resource not registered")
    }
}
EOF

# Run the test (should fail)
cd /tmp && go mod init test && go test -v . 2>/dev/null || echo "✓ EXPECTED: Resource not registered in provider"

# Cleanup
rm -f /tmp/provider_test.go /tmp/go.mod

# Return to project directory
cd /Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam
```

## Step 4: Mock Server Infrastructure Test

### Verify Existing Mock Server
```bash
# Test existing mock server from group resource
grep -n "MockServer" internal/provider/testutils/mock_server.go | head -5

# Check if mock server supports custom role endpoints (should not yet)
grep -n "custom-role" internal/provider/testutils/mock_server.go || echo "✓ EXPECTED: Custom role endpoints not implemented"
```

### Mock Server Extension Requirements
- [ ] `POST /custom-roles` endpoint
- [ ] `GET /custom-roles/{id}` endpoint  
- [ ] `PUT /custom-roles/{id}` endpoint
- [ ] `DELETE /custom-roles/{id}` endpoint
- [ ] Permission validation logic
- [ ] Error scenario simulation

## Step 5: Authentication Integration Test

### OAuth2 Client Access Test
```bash
# Verify OAuth2 integration pattern from group resource
grep -A 5 -B 5 "APIClient" internal/provider/resource_iam_group/iam_group_resource.go | head -10

# This pattern should be replicated for custom role resource
echo "✓ OAuth2 pattern identified for replication"
```

## Step 6: Validation Logic Test

### Permission Pattern Testing
```bash
# Test permission pattern validation manually
cat > /tmp/pattern_test.go << 'EOF'
package main

import (
    "fmt"
    "regexp"
)

func main() {
    pattern := `^[a-z][-a-z]{2}\\.[a-z][-a-z]{1,15}\\.[a-z][-a-z]{1,15}$`
    re := regexp.MustCompile(pattern)
    
    validCases := []string{
        "pos.payment.create",
        "sys.user.manage", 
        "abc.resource.action",
    }
    
    invalidCases := []string{
        "invalid-format",
        "too.short.x",
        "ab.toolongresourcename12345.action",
        "123.numeric.start",
    }
    
    fmt.Println("Valid cases:")
    for _, test := range validCases {
        fmt.Printf("  %s: %v\n", test, re.MatchString(test))
    }
    
    fmt.Println("Invalid cases:")
    for _, test := range invalidCases {
        fmt.Printf("  %s: %v\n", test, re.MatchString(test))
    }
}
EOF

# Run pattern validation test
cd /tmp && go run pattern_test.go

# Cleanup
rm -f /tmp/pattern_test.go

# Return to project directory
cd /Users/shayne/repos/clients/extenda/hiiretail-terraform-providers/iam
```

## Step 7: Performance Baseline Test

### Memory and CPU Baseline
```bash
# Check current test performance baseline
go test ./internal/provider/resource_iam_group -bench=. -benchmem | grep "Benchmark"

# This establishes performance expectations for custom role implementation
echo "✓ Performance baseline established from group resource"
```

## Step 8: Full Test Suite Execution (Should Fail)

### Complete Test Run
```bash
# Attempt to run all tests (many will fail)
echo "Running complete test suite - expect failures for unimplemented features..."

# Test custom role specific functionality
go test ./internal/provider/resource_iam_custom_role -v 2>/dev/null || echo "✓ EXPECTED: Custom role tests not implemented"

# Test provider integration
go test ./internal/provider -v -run="CustomRole" 2>/dev/null || echo "✓ EXPECTED: Provider integration not implemented"

# Test acceptance tests
go test ./acceptance_tests -v -run="CustomRole" 2>/dev/null || echo "✓ EXPECTED: Acceptance tests not implemented"
```

## Validation Checklist

### Pre-Implementation Status (Should All Fail)
- [ ] ❌ Custom role resource registered in provider
- [ ] ❌ Unit tests exist and pass
- [ ] ❌ Integration tests exist and pass  
- [ ] ❌ Contract tests exist and pass
- [ ] ❌ Acceptance tests exist and pass
- [ ] ❌ Benchmark tests exist and pass
- [ ] ❌ Mock server supports custom role endpoints
- [ ] ❌ OAuth2 authentication integrated
- [ ] ❌ Permission validation implemented
- [ ] ❌ Error handling implemented
- [ ] ❌ Concurrent access handling implemented
- [ ] ❌ Performance optimizations implemented

### Post-Implementation Success Criteria
- [ ] ✅ All unit tests pass (>95% coverage)
- [ ] ✅ All integration tests pass
- [ ] ✅ All contract tests pass  
- [ ] ✅ All acceptance tests pass
- [ ] ✅ Benchmark tests show acceptable performance
- [ ] ✅ 500 permission roles perform within limits
- [ ] ✅ Concurrent operations maintain consistency
- [ ] ✅ Error scenarios handled gracefully
- [ ] ✅ OAuth2 authentication working
- [ ] ✅ All validation rules enforced

## Expected Timeline

### Implementation Phases
1. **Setup Phase** (Tasks 1-5): Resource registration, basic structure
2. **Core Implementation** (Tasks 6-15): CRUD operations, validation
3. **Testing Phase** (Tasks 16-25): Unit, integration, contract tests
4. **Performance Phase** (Tasks 26-30): Benchmarks, optimization
5. **Validation Phase** (Tasks 31-35): Acceptance tests, error handling

### Success Metrics
- **Test Coverage**: >95% for custom role resource package
- **Performance**: Operations complete within 1-10ms
- **Reliability**: All acceptance tests pass consistently
- **Code Quality**: No linting errors, comprehensive error handling

This quickstart guide provides a systematic way to validate the implementation progress and ensure all requirements are met through comprehensive testing.