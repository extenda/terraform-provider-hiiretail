# End-User Testing Guide

This guide provides comprehensive testing scenarios for the HiiRetail Terraform provider after publication to the Terraform Registry.

## Pre-Testing Requirements

Before conducting end-user tests, ensure:
- [ ] Provider published to Terraform Registry
- [ ] All platform binaries available
- [ ] GPG signatures valid
- [ ] Documentation accessible

## Test Scenarios

### 1. Fresh Installation Test

**Objective**: Verify provider installs correctly from scratch

**Setup**:
```bash
mkdir terraform-provider-test && cd terraform-provider-test
```

**Test Configuration**:
```hcl
# main.tf
terraform {
  required_version = ">= 1.0"
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 0.1.1-beta"
    }
  }
}

provider "hiiretail" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
}

# Test a basic resource
data "hiiretail_iam_groups" "test" {}

output "groups_count" {
  value       = length(data.hiiretail_iam_groups.test.groups)
  description = "Number of IAM groups"
}
```

**Test Variables**:
```hcl
# variables.tf
variable "tenant_id" {
  description = "HiiRetail tenant ID"
  type        = string
}

variable "client_id" {
  description = "OAuth2 client ID"
  type        = string
}

variable "client_secret" {
  description = "OAuth2 client secret"
  type        = string
  sensitive   = true
}
```

**Test Commands**:
```bash
terraform init
terraform validate
terraform plan
terraform apply
```

**Success Criteria**:
- [ ] `terraform init` completes without errors
- [ ] Provider downloads automatically
- [ ] Version constraint respected
- [ ] Authentication works correctly
- [ ] Basic data source functions

### 2. Version Constraint Testing

**Objective**: Verify version constraints work properly

**Test Configurations**:

```hcl
# Exact version
terraform {
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "0.1.1-beta"
    }
  }
}
```

```hcl
# Pessimistic constraint
terraform {
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 0.1.0"
    }
  }
}
```

```hcl
# Range constraint
terraform {
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = ">= 0.1.0, < 1.0.0"
    }
  }
}
```

**Success Criteria**:
- [ ] Each constraint resolves to correct version
- [ ] Invalid constraints fail appropriately
- [ ] Version resolution is consistent

### 3. Multi-Platform Testing

**Objective**: Verify provider works on different platforms

**Platforms to Test**:
- [ ] Linux x86_64
- [ ] Linux ARM64
- [ ] macOS x86_64 (Intel)
- [ ] macOS ARM64 (Apple Silicon)
- [ ] Windows x86_64
- [ ] Windows ARM64

**Test Method**:
```bash
# On each platform
terraform init
terraform --version
terraform providers
```

**Success Criteria**:
- [ ] Provider installs on all platforms
- [ ] Correct binary architecture used
- [ ] No platform-specific errors

### 4. Authentication Testing

**Objective**: Verify OAuth2 authentication works properly

**Test Scenarios**:

```hcl
# Valid credentials
provider "hiiretail" {
  tenant_id     = "valid-tenant"
  client_id     = "valid-client"
  client_secret = "valid-secret"
}
```

```hcl
# Invalid credentials - should fail gracefully
provider "hiiretail" {
  tenant_id     = "invalid-tenant"
  client_id     = "invalid-client"
  client_secret = "invalid-secret"
}
```

```hcl
# Environment variables
provider "hiiretail" {
  # Uses HIIRETAIL_TENANT_ID, HIIRETAIL_CLIENT_ID, HIIRETAIL_CLIENT_SECRET
}
```

**Success Criteria**:
- [ ] Valid credentials authenticate successfully
- [ ] Invalid credentials fail with clear error messages
- [ ] Environment variables work correctly
- [ ] Sensitive values not logged

### 5. Resource CRUD Testing

**Objective**: Test full resource lifecycle management

**Test Configuration**:
```hcl
# Create an IAM group
resource "hiiretail_iam_group" "test" {
  name        = "terraform-test-group-${random_id.test.hex}"
  description = "Test group created by Terraform"
}

resource "random_id" "test" {
  byte_length = 4
}

# Create a custom role
resource "hiiretail_iam_custom_role" "test" {
  name        = "terraform-test-role-${random_id.test.hex}"
  description = "Test custom role created by Terraform"
  permissions = [
    "groups.read",
    "groups.write"
  ]
}

# Create a role binding
resource "hiiretail_iam_role_binding" "test" {
  tenant_id = var.tenant_id
  role_id   = hiiretail_iam_custom_role.test.id
  group_id  = hiiretail_iam_group.test.id
}
```

**Test Commands**:
```bash
# Create
terraform apply

# Update (modify descriptions)
terraform apply

# Import existing resource
terraform import hiiretail_iam_group.existing existing-group-id

# Destroy
terraform destroy
```

**Success Criteria**:
- [ ] Resources create successfully
- [ ] Updates apply correctly
- [ ] State management works
- [ ] Import functionality works
- [ ] Cleanup completes properly

### 6. Error Handling Testing

**Objective**: Verify graceful error handling

**Test Scenarios**:

```hcl
# Duplicate resource names (should fail)
resource "hiiretail_iam_group" "duplicate1" {
  name = "same-name"
}

resource "hiiretail_iam_group" "duplicate2" {
  name = "same-name"
}
```

```hcl
# Invalid permissions (should fail)
resource "hiiretail_iam_custom_role" "invalid" {
  name = "invalid-role"
  permissions = [
    "nonexistent.permission"
  ]
}
```

**Success Criteria**:
- [ ] API errors reported clearly
- [ ] Validation errors caught early
- [ ] No stack traces exposed to users
- [ ] Helpful error messages provided

### 7. Documentation Testing

**Objective**: Verify documentation accuracy and completeness

**Areas to Check**:
- [ ] Provider configuration examples work
- [ ] Resource documentation matches actual schemas
- [ ] Data source examples are correct
- [ ] Import instructions are accurate
- [ ] Version compatibility noted correctly

**Test Method**:
- Copy examples from registry documentation
- Run examples without modification
- Verify all attributes work as documented

### 8. Performance Testing

**Objective**: Verify provider performs adequately

**Test Scenarios**:
```hcl
# Multiple resources
resource "hiiretail_iam_group" "test" {
  count = 10
  name  = "test-group-${count.index}"
}

# Bulk operations
data "hiiretail_iam_groups" "all" {}

locals {
  group_count = length(data.hiiretail_iam_groups.all.groups)
}
```

**Success Criteria**:
- [ ] Reasonable response times (< 30s for most operations)
- [ ] No memory leaks during long operations
- [ ] Concurrent operations handled properly
- [ ] Rate limiting respected

## Test Environment Setup

### Required Credentials
```bash
export HIIRETAIL_TENANT_ID="your-tenant-id"
export HIIRETAIL_CLIENT_ID="your-client-id"
export HIIRETAIL_CLIENT_SECRET="your-client-secret"
```

### Test Data Cleanup
```bash
# Script to clean up test resources
#!/bin/bash
terraform destroy -auto-approve
```

## Automated Testing

### CI/CD Integration
```yaml
# .github/workflows/integration-test.yml
name: Integration Tests
on:
  release:
    types: [published]
    
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: hashicorp/setup-terraform@v3
    - name: Test Provider Installation
      run: |
        cd test/
        terraform init
        terraform validate
        terraform plan
```

### Test Matrix
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    terraform: ['1.0', '1.5', 'latest']
```

## Success Metrics

### Quantitative Metrics
- [ ] Installation success rate: >99%
- [ ] Test pass rate: >95%
- [ ] Documentation accuracy: >98%
- [ ] Performance: All operations <30s

### Qualitative Metrics
- [ ] User experience smooth
- [ ] Error messages helpful
- [ ] Documentation clear
- [ ] Examples work correctly

## Issue Reporting

When issues are found:

1. **Document the Issue**:
   - Terraform version
   - Provider version
   - Operating system
   - Full error message
   - Minimal reproduction case

2. **Create GitHub Issue**:
   - Use issue template
   - Include environment details
   - Provide reproduction steps
   - Tag appropriately

3. **Workaround Documentation**:
   - Document any workarounds
   - Update known issues list
   - Communicate to users

## Post-Testing Actions

### Successful Testing
- [ ] Mark provider as production-ready
- [ ] Update documentation with any findings
- [ ] Announce availability to users
- [ ] Monitor for user feedback

### Issues Found
- [ ] Document all issues
- [ ] Prioritize fixes
- [ ] Create follow-up releases
- [ ] Update testing procedures

## Long-term Monitoring

### Registry Analytics
- Monitor download statistics
- Track version adoption
- Analyze usage patterns

### User Feedback
- Monitor GitHub issues
- Collect user testimonials
- Survey user satisfaction

### Continuous Improvement
- Regular testing with new Terraform versions
- Update tests for new features
- Refine based on user feedback