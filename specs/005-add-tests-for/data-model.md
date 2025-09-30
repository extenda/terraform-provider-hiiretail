# Data Model: IAM Role Binding Resource

**Date**: September 30, 2025  
**Feature**: IAM Role Binding Resource Implementation and Testing  
**Phase**: 1 - Data Model Design

## Core Entities

### IAM Role Binding
**Purpose**: Represents the assignment of roles to users or groups within a tenant boundary  
**Lifecycle**: Managed through Terraform CRUD operations with OAuth2 authentication

#### Fields
| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| `id` | String | Computed | Unique identifier for the role binding | Auto-generated UUID |
| `role_id` | String | Required | Reference to the role being bound | Must be valid role ID format |
| `bindings` | List[Binding] | Required | List of user/group assignments | Max 10 items, min 1 item |
| `is_custom` | Boolean | Optional | Whether binding applies to custom roles | Default: false |
| `tenant_id` | String | Computed | Tenant context for isolation | Auto-populated from auth context |

#### Binding Sub-Entity
| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| `type` | String | Required | Binding type (user/group) | Enum: "user", "group" |
| `subject_id` | String | Required | ID of user or group | Must be valid subject ID |
| `principal` | String | Optional | Principal identifier | Context-dependent format |

### State Transitions
```
[Non-existent] --CREATE--> [Active]
     ^                        |
     |                        |
   DELETE                   READ/UPDATE
     |                        |
     |                        v
[Deleted] <--DELETE-- [Active/Modified]
```

#### State Management Rules
1. **Create**: New role binding with initial bindings list
2. **Read**: Retrieve current state with tenant isolation
3. **Update**: Atomic replacement of entire bindings list (max 10)
4. **Delete**: Remove role binding and all associated bindings
5. **Import**: Support Terraform import with ID format: `{tenant_id}/{role_binding_id}`

## Validation Rules

### Business Rules
- Maximum 10 bindings per role binding resource
- Minimum 1 binding required (empty bindings list invalid)
- Role ID must reference existing role within tenant
- Subject IDs must be valid within tenant context
- Binding types must be supported values ("user", "group")

### Data Integrity
- Tenant isolation enforced on all operations
- Duplicate bindings within same resource prevented
- Atomic updates ensure consistency during modifications
- State drift detection through Read operations

### Error Conditions
| Condition | Error Type | Message |
|-----------|------------|---------|
| Bindings > 10 | Validation | "Maximum 10 bindings allowed per resource" |
| Empty bindings | Validation | "At least one binding required" |
| Invalid role_id | API Error | "Role not found or access denied" |
| Invalid subject_id | API Error | "Subject not found in tenant" |
| Duplicate binding | Validation | "Duplicate binding detected" |

## Integration Points

### With Custom Roles
- Role binding references custom roles via `role_id` field
- Custom role validation through `is_custom` flag
- Dependency: Custom role must exist before binding creation

### With Groups  
- Group bindings reference group resources via `subject_id`
- Group existence validated during binding operations
- Tenant boundary enforcement for group access

### With OAuth2 Authentication
- All operations require valid OAuth2 token
- Tenant context derived from authentication claims
- Token refresh handling for long-running operations

### With Terraform State
- State stored in Terraform state file with computed fields
- Import functionality for existing role bindings
- Drift detection through periodic Read operations

## API Data Mapping

### Terraform Schema → API Request
```go
// Terraform Resource
resource "hiiretail_iam_role_binding" "example" {
  role_id = "custom-role-123"
  is_custom = true
  bindings = [
    {
      type = "user"
      subject_id = "user-456"
    },
    {
      type = "group" 
      subject_id = "group-789"
    }
  ]
}

// API Request Body
{
  "role_id": "custom-role-123",
  "is_custom": true,
  "bindings": [
    {
      "type": "user",
      "subject_id": "user-456"
    },
    {
      "type": "group",
      "subject_id": "group-789" 
    }
  ]
}
```

### API Response → Terraform State
```go
// API Response
{
  "id": "rb-uuid-123",
  "role_id": "custom-role-123", 
  "is_custom": true,
  "bindings": [...],
  "tenant_id": "tenant-abc",
  "created_at": "2025-09-30T10:00:00Z",
  "updated_at": "2025-09-30T10:00:00Z"
}

// Terraform State
{
  "id": "rb-uuid-123",
  "role_id": "custom-role-123",
  "is_custom": true, 
  "bindings": [...],
  "tenant_id": "tenant-abc"
  // Note: timestamps not stored in Terraform state
}
```

## Testing Data Model

### Unit Test Scenarios
- Valid role binding creation with all field types
- Validation error testing (max bindings, required fields)
- State transformation testing (API ↔ Terraform)
- Edge cases (empty lists, boundary values)

### Integration Test Data
- Multiple tenant scenarios for isolation testing
- Role binding with both user and group bindings
- Update scenarios with binding additions/removals
- Error scenarios with invalid references

### Acceptance Test Models
- Complete lifecycle testing data sets
- Import/export test fixtures
- Performance testing with maximum bindings
- Concurrent operation test scenarios

This data model provides the foundation for implementing the IAM Role Binding resource with proper validation, state management, and integration with the existing HiiRetail IAM system.