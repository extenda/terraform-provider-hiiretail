# Data Model: Improve Resource Usability

**Phase**: 1 | **Date**: October 2, 2025  
**Feature**: Improve Resource Usability

## Core Entities

### ValidationMessage
**Purpose**: Enhanced error messages with actionable guidance
**Fields**:
- `field_path`: String - Terraform path to the invalid field
- `error_code`: String - Standardized error identifier
- `current_value`: String - User-provided value that failed validation
- `expected_format`: String - Description of valid format/constraints
- `example_value`: String - Working example for the field
- `guidance`: String - Specific steps to resolve the issue
- `severity`: String - ERROR, WARNING, INFO

**Validation Rules**:
- `field_path` must follow Terraform path syntax
- `error_code` must be from predefined error taxonomy
- `example_value` must pass validation when substituted
- `guidance` must be actionable and specific

**State Transitions**: 
Immutable - created per validation failure

### ResourceReference
**Purpose**: Standardized handling of cross-resource references
**Fields**:
- `reference_type`: String - "group", "role", "user", "service_account"
- `reference_value`: String - The actual reference (name, email, ID)
- `reference_format`: String - Expected format pattern
- `resolved_id`: String - API ID when reference is resolved
- `validation_status`: String - "unvalidated", "valid", "invalid", "pending"
- `suggestion`: String - Alternative suggestion if validation fails

**Validation Rules**:
- `reference_type` must be from supported types list
- `reference_value` format must match `reference_format` pattern
- `resolved_id` populated only when `validation_status` is "valid"
- `suggestion` provided only when `validation_status` is "invalid"

**State Transitions**:
unvalidated → pending → [valid|invalid]

### ConfigurationExample
**Purpose**: Working examples for resource configurations
**Fields**:
- `resource_type`: String - "hiiretail_iam_group", "hiiretail_iam_custom_role", etc.
- `scenario_name`: String - "basic", "enterprise", "multi_tenant", etc.
- `description`: String - What the example demonstrates
- `configuration`: String - Complete Terraform configuration
- `variables`: Map[String]String - Required variables and default values
- `expected_outputs`: Map[String]String - Expected output values
- `prerequisites`: []String - Required setup before running example

**Validation Rules**:
- `configuration` must be valid HCL syntax
- `configuration` must reference only declared variables
- All `variables` must be used in `configuration`
- `expected_outputs` must match configuration output blocks

**State Transitions**: 
Static - examples are versioned and tested

### PermissionPattern
**Purpose**: Validation patterns for IAM permissions
**Fields**:
- `service`: String - "iam", "ccc", etc.
- `resource_type`: String - "groups", "roles", "users", etc.
- `action`: String - "read", "write", "delete", "create", etc.
- `pattern`: String - Regex pattern for validation
- `description`: String - Human-readable description of permission
- `examples`: []String - Valid permission strings
- `common_typos`: Map[String]String - Typo → correct permission mapping

**Validation Rules**:
- `pattern` must be valid regex
- All `examples` must match `pattern`
- `common_typos` keys must NOT match `pattern`
- `common_typos` values must match `pattern`

**State Transitions**: 
Static - permissions are defined by API specification

## Entity Relationships

### ValidationMessage ↔ ResourceSchema
- ValidationMessage generated from ResourceSchema field validation
- One-to-many: ResourceSchema field → Multiple ValidationMessages (different error types)

### ResourceReference ↔ ValidationMessage  
- Invalid ResourceReference generates specific ValidationMessage
- One-to-one: ResourceReference validation failure → ValidationMessage with suggestions

### ConfigurationExample ↔ ResourceSchema
- ConfigurationExample demonstrates ResourceSchema usage
- Many-to-many: Examples can show multiple resources, resources appear in multiple examples

### PermissionPattern ↔ ValidationMessage
- Permission validation failures use PermissionPattern for suggestions
- One-to-many: PermissionPattern → Multiple potential ValidationMessages

## Validation Architecture

### Schema-Level Validation
```go
// Field-level validators
StringValidator.OneOf(validValues...)
StringValidator.LengthBetween(min, max)
StringValidator.RegexMatches(pattern, message)
StringValidator.NoneOf(reservedWords...)

// Custom validators for complex fields
PermissionStringValidator(permissionPatterns)
ResourceReferenceValidator(referenceTypes)
ConditionalExpressionValidator(allowedFields)
```

### Resource-Level Validation  
```go
// Cross-field validation in ValidateConfig
func (r *ResourceType) ValidateConfig(ctx, req, resp) {
    // Validate role binding members reference existing resources
    // Validate custom role permissions against known patterns
    // Validate group name uniqueness
    // Validate conditional expressions syntax
}
```

### Plan-Time Validation
```go
// Validation during planning phase
func (r *ResourceType) ModifyPlan(ctx, req, resp) {
    // Check resource references can be resolved
    // Validate API connectivity for existence checks
    // Provide warnings for potentially problematic configurations
}
```

## Error Classification

### Error Categories
- **Format Errors**: Invalid syntax, wrong data type, invalid characters
- **Constraint Errors**: Length limits, value ranges, enumeration violations  
- **Reference Errors**: Non-existent resources, circular dependencies, access issues
- **Business Logic Errors**: Policy violations, permission conflicts, security concerns

### Error Severity Levels
- **ERROR**: Configuration cannot be applied, must be fixed
- **WARNING**: Configuration may cause issues, should be reviewed
- **INFO**: Helpful guidance, no action required

### Error Context
- **Local Context**: Field value, expected format, validation rule
- **Resource Context**: Resource type, resource name, related fields
- **Global Context**: Provider configuration, API availability, permissions

---

**Data Model Status**: COMPLETE ✅  
**Entity Relationships**: DEFINED ✅  
**Validation Architecture**: DESIGNED ✅  
**Ready for Contract Generation**: YES ✅