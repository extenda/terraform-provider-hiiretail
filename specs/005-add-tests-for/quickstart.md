# Quickstart: IAM Role Binding Resource

**Feature**: IAM Role Binding Resource Implementation and Testing  
**Date**: September 30, 2025  
**Purpose**: Validate implementation through end-to-end user scenarios

## Overview

This quickstart demonstrates the complete lifecycle of the IAM Role Binding resource from a user perspective. Each scenario validates core functionality and user acceptance criteria.

## Prerequisites

- Terraform 1.5+ installed
- HiiRetail IAM provider configured with OAuth2 credentials
- Valid tenant access with role binding permissions
- At least one custom role and group resource available

## Provider Configuration

```hcl
terraform {
  required_providers {
    hiiretail_iam = {
      source  = "extenda/hiiretail-iam"
      version = "~> 1.0"
    }
  }
}

provider "hiiretail_iam" {
  client_id     = var.client_id     # OAuth2 client ID
  client_secret = var.client_secret # OAuth2 client secret
  auth_url      = var.auth_url      # OAuth2 token endpoint
  api_base_url  = var.api_base_url  # HiiRetail API base URL
}
```

## Scenario 1: Basic Role Binding Creation

**User Story**: As a system administrator, I need to bind a custom role to a user so they have appropriate permissions.

```hcl
# Create custom role (prerequisite)
resource "hiiretail_iam_custom_role" "example" {
  name        = "example-role"
  description = "Example custom role for quickstart"
  permissions = ["read:users", "write:groups"]
}

# Create role binding
resource "hiiretail_iam_role_binding" "user_binding" {
  role_id   = hiiretail_iam_custom_role.example.id
  is_custom = true
  
  bindings = [
    {
      type       = "user"
      subject_id = "user-12345"
    }
  ]
}
```

**Validation Steps**:
1. Run `terraform plan` - should show 2 resources to create
2. Run `terraform apply` - should successfully create resources
3. Verify role binding exists with correct role_id and binding
4. Check that computed fields (id, tenant_id) are populated

**Expected Outcome**: Role binding created with single user binding, proper tenant isolation, and generated ID.

## Scenario 2: Multiple Bindings Management

**User Story**: As a system administrator, I need to assign the same role to multiple users and groups for efficient permission management.

```hcl
resource "hiiretail_iam_role_binding" "multiple_bindings" {
  role_id   = hiiretail_iam_custom_role.example.id
  is_custom = true
  
  bindings = [
    {
      type       = "user"
      subject_id = "user-12345"
    },
    {
      type       = "user"
      subject_id = "user-67890"
    },
    {
      type       = "group"
      subject_id = "group-admin"
    }
  ]
}
```

**Validation Steps**:
1. Run `terraform plan` - should show update to role binding resource
2. Apply changes and verify all 3 bindings are present
3. Test boundary condition with 10 bindings (maximum allowed)
4. Attempt 11 bindings - should fail validation

**Expected Outcome**: Atomic update of bindings list with proper validation of maximum limit.

## Scenario 3: Resource Import Testing

**User Story**: As a system administrator, I need to import existing role bindings into Terraform state for management.

```bash
# Import existing role binding
terraform import hiiretail_iam_role_binding.imported_binding "tenant-abc/rb-existing-uuid"
```

```hcl
# Configuration for imported resource
resource "hiiretail_iam_role_binding" "imported_binding" {
  role_id   = "imported-role-id"
  is_custom = false
  
  bindings = [
    {
      type       = "group"
      subject_id = "existing-group-id"
    }
  ]
}
```

**Validation Steps**:
1. Import existing role binding using tenant/ID format
2. Run `terraform plan` - should show no changes (state matches config)
3. Modify configuration and apply changes
4. Verify imported resource is now managed by Terraform

**Expected Outcome**: Successful import with proper state synchronization and ongoing management.

## Scenario 4: Error Handling Validation

**User Story**: As a system administrator, I need clear error messages when role binding operations fail.

```hcl
# Test various error conditions
resource "hiiretail_iam_role_binding" "error_cases" {
  role_id   = "non-existent-role"  # Should fail validation
  is_custom = true
  
  bindings = []  # Empty list should fail validation
}
```

**Validation Steps**:
1. Test invalid role_id reference - should get API error
2. Test empty bindings list - should get validation error  
3. Test invalid subject_id - should get API error
4. Test OAuth2 authentication failure - should get auth error

**Expected Outcome**: Clear, actionable error messages for all failure scenarios with appropriate error codes.

## Scenario 5: Update and Deletion Testing

**User Story**: As a system administrator, I need to modify and remove role bindings as organizational needs change.

```hcl
# Initial configuration
resource "hiiretail_iam_role_binding" "lifecycle_test" {
  role_id   = hiiretail_iam_custom_role.example.id
  is_custom = true
  
  bindings = [
    {
      type       = "user"
      subject_id = "temp-user-123"
    }
  ]
}

# Updated configuration (add binding)
resource "hiiretail_iam_role_binding" "lifecycle_test" {
  role_id   = hiiretail_iam_custom_role.example.id
  is_custom = true
  
  bindings = [
    {
      type       = "user"
      subject_id = "temp-user-123"
    },
    {
      type       = "group"
      subject_id = "temp-group-456"
    }
  ]
}
```

**Validation Steps**:
1. Create initial role binding with one binding
2. Update to add additional binding - verify atomic update
3. Remove one binding - verify atomic update  
4. Delete entire resource - verify clean removal
5. Verify Terraform state is properly updated at each step

**Expected Outcome**: Smooth lifecycle management with proper state tracking and atomic updates.

## Performance Validation

**Load Testing Scenario**: Create multiple role bindings concurrently to test provider performance.

```hcl
# Create multiple role bindings for performance testing
resource "hiiretail_iam_role_binding" "load_test" {
  count = 10
  
  role_id   = hiiretail_iam_custom_role.example.id
  is_custom = true
  
  bindings = [
    {
      type       = "user"
      subject_id = "load-test-user-${count.index}"
    }
  ]
}
```

**Performance Criteria**:
- Single role binding operations complete within 2 seconds
- 10 concurrent operations complete within 10 seconds
- No resource state corruption under concurrent load
- Proper OAuth2 token reuse and refresh handling

## Integration Testing Checklist

- [ ] Basic CRUD operations work correctly
- [ ] Validation rules are properly enforced
- [ ] Error handling provides clear messages
- [ ] Import/export functionality works
- [ ] Tenant isolation is maintained
- [ ] OAuth2 authentication flows work
- [ ] Performance meets requirements
- [ ] Concurrent operations are safe
- [ ] State management is consistent
- [ ] Provider registration is successful

## Troubleshooting

### Common Issues
1. **Authentication errors**: Verify OAuth2 credentials and token endpoint
2. **Validation failures**: Check role_id exists and bindings format
3. **Permission errors**: Ensure OAuth2 scope includes role binding permissions
4. **Import failures**: Verify correct ID format (tenant_id/role_binding_id)

### Debug Steps
1. Enable Terraform debug logging: `TF_LOG=DEBUG terraform apply`
2. Check provider logs for OAuth2 token refresh issues
3. Validate API connectivity with curl/Postman
4. Verify tenant context in API requests

This quickstart serves as both user documentation and implementation validation, ensuring the IAM Role Binding resource meets all functional requirements and user experience expectations.