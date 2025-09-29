# Data Model: IAM Custom Role Resource

**Date**: September 28, 2025  
**Feature**: Add Comprehensive Tests for IAM Custom Role Resource

## Entity Definitions

### Custom Role
**Purpose**: Represents a custom IAM role with specific permissions and attributes
**Scope**: Tenant-scoped resource managed by Terraform

**Fields**:
- `id` (string, required): Unique identifier for the custom role
- `name` (string, optional, computed): Human-readable name for the role (3-256 characters)
- `permissions` (list, required): List of permission objects (max 100 general, 500 POS)  
- `tenant_id` (string, optional, computed): Tenant scope inherited from provider

**Validation Rules**:
- ID must be unique within tenant scope
- Name length between 3-256 characters when provided
- Permissions list cannot be empty
- Permission count limits: 100 general, 500 for POS permissions

**State Transitions**:
```
null -> planned -> created -> [updated] -> destroyed
```

### Permission
**Purpose**: Individual permission entry within a custom role
**Scope**: Nested object within Custom Role

**Fields**:
- `id` (string, required): Permission identifier following pattern validation
- `alias` (string, computed): Auto-generated alias for the permission
- `attributes` (object, optional): Key-value attributes for permission metadata

**Validation Rules**:
- ID pattern: `^[a-z][-a-z]{2}\\.[a-z][-a-z]{1,15}\\.[a-z][-a-z]{1,15}$`
- Pattern format: `{systemPrefix}.{resource}.{action}`
- systemPrefix: exactly 3 characters, a-z and hyphens
- resource: 2-16 characters, a-z and hyphens
- action: 2-16 characters, a-z and hyphens

**State Behavior**:
- Alias computed server-side, read-only in Terraform
- Attributes optional, subject to separate validation rules

### Attributes
**Purpose**: Extensible metadata for permissions
**Scope**: Nested object within Permission

**Constraints**:
- Maximum 10 properties per attributes object
- Keys: maximum 40 characters each
- Values: maximum 256 characters each, string type only
- Optional - can be null or omitted

**Validation Rules**:
- Property count validation (≤ 10)
- Key length validation (≤ 40 chars)
- Value length validation (≤ 256 chars)
- Value type validation (strings only)

## Data Relationships

```
Provider (1) --> (n) Custom Roles
  ↓ tenant_id, auth context
  
Custom Role (1) --> (n) Permissions
  ↓ role ownership
  
Permission (1) --> (0..1) Attributes
  ↓ optional metadata
```

## Schema Validation Matrix

| Field | Required | Type | Constraints | Computed | Validation |
|-------|----------|------|-------------|----------|------------|
| CustomRole.id | ✓ | string | unique | ❌ | required validation |
| CustomRole.name | ❌ | string | 3-256 chars | ✓ | length validation |
| CustomRole.permissions | ✓ | list | 1-500 items | ❌ | count + item validation |
| CustomRole.tenant_id | ❌ | string | provider context | ✓ | inherited validation |
| Permission.id | ✓ | string | pattern match | ❌ | regex validation |
| Permission.alias | ❌ | string | server computed | ✓ | read-only |
| Permission.attributes | ❌ | object | size limits | ❌ | nested validation |
| Attributes.* | ❌ | string | length limits | ❌ | key/value validation |

## Test Data Scenarios

### Valid Role Examples
```go
// Minimal valid role
{
  id: "test-role-001",
  permissions: [{
    id: "pos.payment.create"
  }]
}

// Role with attributes
{
  id: "admin-role",
  name: "Administrator Role",
  permissions: [{
    id: "sys.user.manage",
    attributes: {
      "department": "IT",
      "level": "admin"
    }
  }]
}

// Maximum POS permissions (500)
{
  id: "pos-full-role",
  permissions: [
    {id: "pos.payment.create"},
    {id: "pos.payment.read"},
    // ... up to 500 pos.* permissions
  ]
}
```

### Invalid Role Examples
```go
// Missing required permissions
{
  id: "invalid-role-001"
  // permissions missing - should fail validation
}

// Invalid permission pattern
{
  id: "test-role-002",
  permissions: [{
    id: "invalid-permission-format"  // doesn't match pattern
  }]
}

// Exceeded permission limits
{
  id: "too-many-perms",
  permissions: [
    // 101 non-POS permissions - should fail
  ]
}

// Invalid attribute constraints
{
  id: "attr-violation",
  permissions: [{
    id: "sys.test.action",
    attributes: {
      "very-long-key-exceeding-forty-character-limit": "value",  // key too long
      "normal-key": "very long value exceeding 256 character limit..." // value too long
    }
  }]
}
```

## State Management

### Import Scenarios
- Import existing role by ID
- Validate imported state matches API response
- Handle missing optional fields during import
- Preserve computed fields (alias, tenant_id)

### Update Scenarios  
- Add permissions to existing role
- Remove permissions from role
- Modify permission attributes
- Update role name
- Handle partial update failures

### Error Recovery
- API unavailable during operations
- Permission conflicts during updates
- Tenant context changes
- Network timeouts and retries

## Performance Considerations

### Large Permission Sets
- Roles with 500 POS permissions
- Bulk permission validation
- API request batching if needed
- Memory usage with large attribute sets

### Concurrent Operations
- Multiple Terraform operations on same role
- Race condition handling
- State consistency validation
- Lock contention scenarios

This data model provides the foundation for comprehensive test scenarios covering all validation rules, constraints, and edge cases identified in the feature specification.