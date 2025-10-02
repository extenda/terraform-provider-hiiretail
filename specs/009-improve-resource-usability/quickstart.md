# Quickstart: Improved Resource Usability

**Feature**: Improve Resource Usability  
**Phase**: 1 | **Date**: October 2, 2025

## Overview

This quickstart validates the enhanced usability improvements for HiiRetail Terraform provider IAM resources. The improvements focus on better validation, clearer error messages, and comprehensive documentation based on the resource declarations in `simple_test.tf`.

## Prerequisites

- Go 1.19+ installed
- Terraform 1.0+ installed  
- HiiRetail IAM API access (test environment)
- Valid OAuth2 credentials for testing

## Quick Validation Steps

### 1. Enhanced Validation Testing

```bash
# Navigate to the provider directory
cd iam/

# Test enhanced group resource validation
terraform plan -var="client_id=test" -var="client_secret=test" \
  -target=hiiretail_iam_group.test_group

# Expected: Clear validation messages for any configuration issues
```

### 2. Error Message Clarity Test

Create a test configuration with intentional errors:

```hcl
# test_validation.tf
resource "hiiretail_iam_group" "invalid_group" {
  name        = "INVALID@NAME!"  # Invalid characters
  description = "x"              # Too short description
}

resource "hiiretail_iam_custom_role" "invalid_role" {
  name        = ""                    # Empty name
  title       = "x"                  # Too short title  
  description = "Test role"
  permissions = [
    "invalid:permission:format",     # Invalid permission
    "iam:group:read"                # Missing 's' in groups
  ]
  stage = "INVALID_STAGE"           # Invalid stage value
}

resource "hiiretail_iam_role_binding" "invalid_binding" {
  name = "test-binding"
  role = "invalid-role-format"      # Invalid role format
  members = [
    "invalid-member-format",        # Invalid member format
    "user:not-an-email"            # Invalid email format
  ]
  condition = "invalid expression"  # Invalid condition syntax
}
```

Run validation:
```bash
terraform validate -json
```

**Expected Enhanced Error Output**:
```json
{
  "valid": false,
  "error_count": 8,
  "errors": [
    {
      "severity": "error",
      "summary": "Invalid group name format",
      "detail": "Field 'name' validation failed: Contains invalid characters\nCurrent value: 'INVALID@NAME!'\nExpected format: Lowercase letters, numbers, and hyphens only\nExample: 'analytics-team-dev'\nGuidance: Remove special characters and use lowercase letters with hyphens as separators",
      "range": {
        "filename": "test_validation.tf",
        "start": {"line": 3, "column": 15}
      }
    }
  ]
}
```

### 3. Reference Validation Test  

```hcl
# test_references.tf
resource "hiiretail_iam_role_binding" "reference_test" {
  name = "test-reference-binding"
  role = "roles/custom.nonexistent-role"  # Non-existent role
  members = [
    "group:nonexistent-group",            # Non-existent group
    "user:invalid-email-format"           # Invalid email
  ]
}
```

Run plan:
```bash
terraform plan
```

**Expected Reference Validation Output**:
```
Error: Role reference validation failed

  on test_references.tf line 4, in resource "hiiretail_iam_role_binding" "reference_test":
   4:   role = "roles/custom.nonexistent-role"

Field 'role' validation failed: Referenced role does not exist
Current value: 'roles/custom.nonexistent-role'
Expected: Valid role reference
Examples: 'roles/viewer', 'roles/custom.analytics-reader'
Guidance: Check available roles with 'data.hiiretail_iam_roles.all_roles' or verify the role name spelling

Available custom roles: test-custom-role-unique-id
```

### 4. Permission Validation Test

```hcl
# test_permissions.tf  
resource "hiiretail_iam_custom_role" "permission_test" {
  name        = "permission-test-role"
  title       = "Permission Test Role"
  description = "Testing permission validation"
  permissions = [
    "iam:group:read",        # Missing 's' in groups
    "invalid:format",        # Invalid format (missing action)
    "unknown:service:read",  # Unknown service
    "iam:groups:invalid"     # Invalid action
  ]
}
```

**Expected Permission Validation Output**:
```
Error: Invalid permission format

  on test_permissions.tf line 7, in resource "hiiretail_iam_custom_role" "permission_test":
   7:     "iam:group:read",

Field 'permissions[0]' validation failed: Invalid resource name
Current value: 'iam:group:read'
Expected format: service:resource:action
Example: 'iam:groups:read'
Guidance: Did you mean 'iam:groups:read'? The resource should be plural.

Similar permissions: iam:groups:read, iam:groups:write, iam:groups:list
```

### 5. Working Configuration Test

Verify that valid configurations work without issues:

```bash
# Use the working simple_test.tf configuration
terraform plan -var-file="terraform.tfvars"
```

**Expected**: Clean plan output with no validation errors.

## Success Criteria Validation

### ✅ Error Message Quality
- [ ] Error messages include specific field paths
- [ ] Current values are displayed in error messages  
- [ ] Expected formats are clearly described
- [ ] Working examples are provided in error messages
- [ ] Actionable guidance is included for resolution

### ✅ Validation Coverage  
- [ ] All required fields have appropriate validation
- [ ] Format validation works for names, emails, permissions
- [ ] Cross-resource references are validated
- [ ] Permission strings follow expected patterns
- [ ] Conditional expressions are validated

### ✅ User Experience
- [ ] Users can understand errors without consulting documentation
- [ ] Suggestions are provided for common typos
- [ ] Related valid options are shown when validation fails
- [ ] Plan-time validation catches issues before apply

### ✅ Documentation Quality
- [ ] Resource schema descriptions are comprehensive
- [ ] Examples demonstrate real-world usage patterns
- [ ] Field descriptions explain business purpose
- [ ] Troubleshooting guidance is accessible

## Performance Validation

### Response Time Test
```bash
# Measure validation performance
time terraform validate

# Expected: < 2 seconds for validation of complex configurations
```

### Reference Resolution Test  
```bash
# Test reference validation performance
time terraform plan -target=hiiretail_iam_role_binding.test_role_binding

# Expected: < 5 seconds including API calls for reference validation
```

## Integration Test Scenarios

### Scenario 1: New User Experience
**Context**: New user configuring IAM resources for the first time
**Test**: Follow simple_test.tf example step by step
**Success**: User can configure all resources correctly with helpful guidance

### Scenario 2: Error Recovery  
**Context**: User makes configuration mistakes
**Test**: Introduce errors and validate recovery guidance
**Success**: User can fix errors based on error messages alone

### Scenario 3: Complex Configuration
**Context**: Enterprise setup with multiple role bindings
**Test**: Configure complex multi-resource setup
**Success**: Clear validation and progress feedback throughout

## Rollback Plan

If validation improvements cause issues:

1. **Immediate**: Disable enhanced validation via feature flag
2. **Short-term**: Roll back to basic validation with simple error messages  
3. **Long-term**: Fix issues and re-enable enhancements

## Next Steps

After quickstart validation passes:

1. **Phase 2**: Generate detailed implementation tasks
2. **Implementation**: Enhance existing resource validation code
3. **Testing**: Add comprehensive test coverage for new validation
4. **Documentation**: Update provider documentation with examples
5. **Release**: Deploy with proper migration guidance

---

**Quickstart Status**: READY FOR EXECUTION ✅  
**Success Criteria**: DEFINED ✅  
**Test Scenarios**: COMPREHENSIVE ✅  
**Ready for Task Generation**: YES ✅