# Data Model: IAM Service Testing

## Entities

### IAM Service
- Represents the main client for IAM API operations
- Methods: ListGroups, GetGroup, CreateGroup, UpdateGroup, DeleteGroup, ListRoles, GetRole, CreateCustomRole, GetCustomRole, UpdateCustomRole, DeleteCustomRole, ListRoleBindings, GetRoleBinding, CreateRoleBinding, UpdateRoleBinding, DeleteRoleBinding, SetResource, GetResource, DeleteResource, GetResources, AddRoleToGroup
- Attributes: client, rawClient, tenantID

### Group
- ID, Name, Description, Members, CreatedAt, UpdatedAt

### CustomRole
- ID, Name, Title, Description, Permissions, Stage, CreatedAt, UpdatedAt

### Permission
- ID, Attributes

### RoleBinding
- ID, Name, Role, Members, Condition, CreatedAt, UpdatedAt

### Role
- ID, Name, Title, Description, Stage, Type

### Resource
- ID, Name, Props

## Relationships
- Service aggregates all IAM operations and models
- Group, CustomRole, RoleBinding, Role, Resource are used as parameters and return types in Service methods

## Validation Rules
- All methods must validate input parameters (e.g., non-empty IDs, valid names)
- Error handling must be tested for all API calls
- Mock clients must simulate API responses and errors

## State Transitions
- CRUD operations for each entity must be tested for correct state changes
- Error scenarios must be tested for proper error propagation

## Next Steps
- Generate contracts and quickstart for test coverage improvement
