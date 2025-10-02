# Data Model: Multi-API Provider User Experience Enhancement

## Provider Configuration Entity

**Entity**: Provider Configuration  
**Purpose**: Centralized authentication and service endpoint configuration for all HiiRetail APIs

**Attributes**:
- `client_id` (string, required, sensitive): OAuth2 client ID for API authentication
- `client_secret` (string, required, sensitive): OAuth2 client secret  
- `auth_url` (string, optional): OAuth2 token endpoint URL (default: https://auth.retailsvc.com/oauth2/token)
- `audience` (string, optional): OAuth2 audience parameter (default: https://iam-api.retailsvc.com)
- `timeout_seconds` (number, optional): HTTP request timeout (default: 30)
- `max_retries` (number, optional): Maximum retry attempts (default: 3)
- `iam_endpoint` (string, optional): IAM API base URL override
- `ccc_endpoint` (string, optional): CCC API base URL override (future)

**Validation Rules**:
- `client_id` and `client_secret` must be provided via variables or environment
- URLs must be valid HTTPS endpoints
- `timeout_seconds` must be between 5 and 300
- `max_retries` must be between 0 and 10

**State Management**: Provider-level configuration, not stored in Terraform state

## Service Module Entity

**Entity**: Service Module  
**Purpose**: Logical grouping of resources and data sources for each HiiRetail API

**Attributes**:
- `service_name` (string): API service identifier (e.g., "iam", "ccc")
- `base_url` (string): Service-specific API base URL
- `resource_types` ([]string): List of resource types provided by service
- `data_source_types` ([]string): List of data source types provided by service
- `api_version` (string): Supported API version

**IAM Service Module**:
- `service_name`: "iam"
- `base_url`: "https://iam-api.retailsvc.com"  
- `resource_types`: ["group", "custom_role", "role_binding"]
- `data_source_types`: ["groups", "roles"]
- `api_version`: "v1"

**CCC Service Module** (future):
- `service_name`: "ccc"
- `base_url`: "https://ccc-api.retailsvc.com"
- `resource_types`: ["kind"]
- `data_source_types`: ["kinds"]
- `api_version`: "v1"

## Resource Naming Entity

**Entity**: Resource Naming Convention  
**Purpose**: Standardized naming pattern for all provider resources

**Pattern**: `hiiretail_{service}_{resource_type}`

**Current IAM Resources**:
- `hiiretail_iam_group` (was: hiiretail-iam_iam_group)
- `hiiretail_iam_custom_role` (was: hiiretail-iam_custom_role)  
- `hiiretail_iam_role_binding` (was: hiiretail-iam_iam_role_binding)

**Future CCC Resources**:
- `hiiretail_ccc_kind`

**Data Sources**: Follow same pattern with plural forms
- `hiiretail_iam_groups`
- `hiiretail_iam_roles`
- `hiiretail_ccc_kinds`

## Resource Documentation Entity

**Entity**: Resource Documentation  
**Purpose**: Comprehensive documentation for each resource including examples and relationships

**Attributes**:
- `resource_name` (string): Full resource name (e.g., hiiretail_iam_group)
- `service` (string): Parent service (e.g., "iam")
- `description` (string): Resource purpose and functionality
- `required_attributes` ([]string): List of required schema attributes
- `optional_attributes` ([]string): List of optional schema attributes
- `computed_attributes` ([]string): List of computed schema attributes
- `examples` ([]Example): Usage examples
- `relationships` ([]Relationship): Dependencies and references

**Example Structure**:
```json
{
  "resource_name": "hiiretail_iam_group",
  "service": "iam",
  "description": "Manages IAM groups for user and permission organization",
  "required_attributes": ["name"],
  "optional_attributes": ["description"],
  "computed_attributes": ["id", "tenant_id", "status"],
  "examples": [
    {
      "title": "Basic IAM Group",
      "code": "resource \"hiiretail_iam_group\" \"example\" {\n  name = \"developers\"\n  description = \"Development team access group\"\n}"
    }
  ],
  "relationships": [
    {
      "type": "referenced_by",
      "resource": "hiiretail_iam_role_binding",
      "description": "Groups can be bound to roles via role bindings"
    }
  ]
}
```

## Migration Entity

**Entity**: Migration Guide  
**Purpose**: Documentation and tooling for migrating from old provider structure

**Attributes**:
- `old_resource_name` (string): Previous resource name
- `new_resource_name` (string): New resource name  
- `migration_steps` ([]string): Required migration actions
- `breaking_changes` ([]string): Changes that require user action
- `automation_available` (boolean): Whether automated migration exists

**Migration Mappings**:
```json
[
  {
    "old_resource_name": "hiiretail-iam_iam_group",
    "new_resource_name": "hiiretail_iam_group", 
    "migration_steps": [
      "Update provider source from 'extenda/hiiretail-iam' to 'extenda/hiiretail'",
      "Replace resource type 'hiiretail-iam_iam_group' with 'hiiretail_iam_group'",
      "Update any references in data sources and outputs"
    ],
    "breaking_changes": ["Resource type name change"],
    "automation_available": true
  },
  {
    "old_resource_name": "hiiretail-iam_custom_role",
    "new_resource_name": "hiiretail_iam_custom_role",
    "migration_steps": [
      "Update provider source",
      "Replace resource type name"
    ],
    "breaking_changes": ["Resource type name change"],
    "automation_available": true
  },
  {
    "old_resource_name": "hiiretail-iam_iam_role_binding", 
    "new_resource_name": "hiiretail_iam_role_binding",
    "migration_steps": [
      "Update provider source",
      "Replace resource type name"
    ],
    "breaking_changes": ["Resource type name change"],
    "automation_available": true
  }
]
```

## Getting Started Entity

**Entity**: Getting Started Guide
**Purpose**: Step-by-step documentation for new users

**Attributes**:
- `steps` ([]Step): Ordered list of setup steps
- `prerequisites` ([]string): Required setup before starting
- `estimated_time` (string): Expected completion time
- `validation_steps` ([]string): Steps to verify successful setup

**Structure**:
```json
{
  "prerequisites": [
    "HiiRetail OAuth2 credentials",
    "Terraform 1.0+ installed",
    "Access to HiiRetail APIs"
  ],
  "estimated_time": "15 minutes",
  "steps": [
    {
      "title": "Configure Provider",
      "description": "Set up the HiiRetail provider with authentication",
      "code_example": "terraform {\n  required_providers {\n    hiiretail = {\n      source = \"extenda/hiiretail\"\n    }\n  }\n}\n\nprovider \"hiiretail\" {\n  client_id = var.client_id\n  client_secret = var.client_secret\n}"
    },
    {
      "title": "Create First Resource", 
      "description": "Create a basic IAM group to verify setup",
      "code_example": "resource \"hiiretail_iam_group\" \"my_first_group\" {\n  name = \"my-team\"\n  description = \"My first HiiRetail IAM group\"\n}"
    },
    {
      "title": "Apply Configuration",
      "description": "Run terraform plan and apply to create resources",
      "code_example": "terraform plan\nterraform apply"
    }
  ],
  "validation_steps": [
    "Verify group appears in HiiRetail IAM console",
    "Run 'terraform show' to confirm state",
    "Test group deletion with 'terraform destroy'"
  ]
}
```

## Entity Relationships

```
Provider Configuration
├── Service Modules (1:N)
│   ├── IAM Service Module
│   │   ├── Resources (1:N)
│   │   │   ├── hiiretail_iam_group
│   │   │   ├── hiiretail_iam_custom_role  
│   │   │   └── hiiretail_iam_role_binding
│   │   └── Data Sources (1:N)
│   │       └── hiiretail_iam_groups
│   └── CCC Service Module (future)
│       └── Resources (1:N)
│           └── hiiretail_ccc_kind
├── Resource Documentation (1:N)
├── Migration Guides (1:N)
└── Getting Started Guide (1:1)
```

## State Transitions

**Provider Configuration**: 
- Initialize → Configured → Active → Error (with retry logic)

**Resource Lifecycle**:
- Plan → Create → Read → Update → Delete
- Error states trigger appropriate retry and rollback logic

**Service Module Lifecycle**:
- Registered → Available → Active → Deprecated

This data model provides the foundation for implementing the multi-API provider architecture while maintaining security, usability, and extensibility requirements.