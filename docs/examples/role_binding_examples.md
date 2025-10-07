# HiiRetail IAM Role Binding Examples

This document provides comprehensive examples of using the `hiiretail_iam_role_binding` resource in various scenarios.

## Table of Contents

1. [Basic Examples](#basic-examples)
2. [Advanced Use Cases](#advanced-use-cases)
3. [Integration with Other Resources](#integration-with-other-resources)
4. [Best Practices](#best-practices)
5. [Common Patterns](#common-patterns)

## Basic Examples

### Single User Binding

```hcl
# Bind a single user to a custom role
resource "hiiretail_iam_role_binding" "developer_access" {
  role_id = "developer-role-123"
  bindings = [
    {
      type = "user"
      id   = "jane.smith@company.com"
    }
  ]
  description = "Developer access for Jane Smith"
}
```

### Group-Based Access

```hcl
# Bind an entire group to a role
resource "hiiretail_iam_role_binding" "team_access" {
  role_id = "team-viewer-role"
  bindings = [
    {
      type = "group"
      id   = "development-team"
    }
  ]
  description = "Development team viewer access"
}
```

### Service Account Binding

```hcl
# Bind a service account for automation
resource "hiiretail_iam_role_binding" "automation_access" {
  role_id = "automation-role"
  bindings = [
    {
      type = "service_account"
      id   = "ci-cd-service-account"
    }
  ]
  description = "CI/CD automation access"
}
```

## Advanced Use Cases

### Mixed Binding Types

```hcl
# Combine different binding types in one role binding
resource "hiiretail_iam_role_binding" "admin_team" {
  role_id = "admin-role"
  bindings = [
    # Primary administrator
    {
      type = "user"
      id   = "admin@company.com"
    },
    # Admin group for escalations
    {
      type = "group"  
      id   = "admin-group"
    },
    # Monitoring service account
    {
      type = "service_account"
      id   = "monitoring-sa"
    }
  ]
  description = "Combined admin access for users, groups, and services"
}
```

### Maximum Bindings (10)

```hcl
# Example with maximum allowed bindings
resource "hiiretail_iam_role_binding" "large_team" {
  role_id = "team-role"
  bindings = [
    { type = "user", id = "user1@company.com" },
    { type = "user", id = "user2@company.com" },
    { type = "user", id = "user3@company.com" },
    { type = "user", id = "user4@company.com" },
    { type = "user", id = "user5@company.com" },
    { type = "group", id = "contractors-group" },
    { type = "group", id = "interns-group" },
    { type = "service_account", id = "app-sa-1" },
    { type = "service_account", id = "app-sa-2" },
    { type = "service_account", id = "backup-sa" }
  ]
  description = "Large team with maximum bindings"
}
```

### System Role Binding

```hcl
# Bind to a system-defined role
resource "hiiretail_iam_role_binding" "system_admin" {
  role_id = "system-administrator"
  bindings = [
    {
      type = "user"
      id   = "sysadmin@company.com"
    }
  ]
  description = "System administrator privileges"
}
```

## Integration with Other Resources

### Using with Custom Roles

```hcl
# Create a custom role first
resource "hiiretail_iam_custom_role" "data_analyst" {
  name = "data-analyst"
  description = "Data analysis permissions"
  permissions = [
    "data:read",
    "reports:create",
    "dashboards:view"
  ]
}

# Then bind users to the custom role
resource "hiiretail_iam_role_binding" "analysts" {
  role_id = hiiretail_iam_custom_role.data_analyst.id
  bindings = [
    {
      type = "user"
      id   = "analyst1@company.com"
    },
    {
      type = "user"
      id   = "analyst2@company.com"
    },
    {
      type = "group"
      id   = "data-team"
    }
  ]
  description = "Data analysts access binding"
}
```

### Using with Groups

```hcl
# Create a group first
resource "hiiretail_iam_group" "developers" {
  name = "developers"
  description = "Development team group"
}

# Bind the group to multiple roles
resource "hiiretail_iam_role_binding" "dev_read_access" {
  role_id = "read-only-role"
  bindings = [
    {
      type = "group"
      id   = hiiretail_iam_group.developers.id
    }
  ]
  description = "Developer read access"
}

resource "hiiretail_iam_role_binding" "dev_write_access" {
  role_id = "developer-role"
  bindings = [
    {
      type = "group"
      id   = hiiretail_iam_group.developers.id
    }
  ]
  description = "Developer write access"
}
```

## Best Practices

### 1. Use Descriptive Names and Descriptions

```hcl
resource "hiiretail_iam_role_binding" "prod_db_readonly_analysts" {
  role_id = "database-readonly"
  bindings = [
    {
      type = "group"
      id   = "data-analysts-group"
    }
  ]
  description = "Production database read-only access for data analysts"
}
```

### 2. Group Similar Bindings

```hcl
# Good: Group related users in one binding
resource "hiiretail_iam_role_binding" "qa_team_access" {
  role_id = "qa-tester-role"
  bindings = [
    { type = "user", id = "qa-lead@company.com" },
    { type = "user", id = "qa-tester1@company.com" },
    { type = "user", id = "qa-tester2@company.com" },
    { type = "group", id = "qa-contractors" }
  ]
  description = "QA team testing access"
}
```

### 3. Use Variables for Flexibility

```hcl
variable "team_members" {
  description = "List of team member email addresses"
  type        = list(string)
  default     = [
    "member1@company.com",
    "member2@company.com",
    "member3@company.com"
  ]
}

resource "hiiretail_iam_role_binding" "dynamic_team" {
  role_id = "team-role"
  bindings = [
    for member in var.team_members : {
      type = "user"
      id   = member
    }
  ]
  description = "Dynamic team role binding"
}
```

### 4. Separate Environments

```hcl
# Development environment
resource "hiiretail_iam_role_binding" "dev_access" {
  role_id = "developer-role"
  bindings = [
    {
      type = "group"
      id   = "dev-team"
    }
  ]
  description = "Development environment access"
}

# Production environment (more restrictive)
resource "hiiretail_iam_role_binding" "prod_access" {
  role_id = "production-operator"
  bindings = [
    {
      type = "user"
      id   = "prod-admin@company.com"
    }
  ]
  description = "Production environment access (restricted)"
}
```

## Common Patterns

### Hierarchical Access Pattern

```hcl
# Admin level - full access
resource "hiiretail_iam_role_binding" "admin_full_access" {
  role_id = "administrator"
  bindings = [
    { type = "user", id = "admin@company.com" }
  ]
  description = "Full administrative access"
}

# Manager level - management access  
resource "hiiretail_iam_role_binding" "manager_access" {
  role_id = "manager-role"
  bindings = [
    { type = "group", id = "managers-group" }
  ]
  description = "Management level access"
}

# Employee level - basic access
resource "hiiretail_iam_role_binding" "employee_access" {
  role_id = "employee-role"
  bindings = [
    { type = "group", id = "all-employees" }
  ]
  description = "Basic employee access"
}
```

### Temporary Access Pattern

```hcl
# Temporary contractor access
resource "hiiretail_iam_role_binding" "contractor_temp_access" {
  role_id = "contractor-role"
  bindings = [
    {
      type = "user"
      id   = "contractor@external.com"
    }
  ]
  description = "Temporary contractor access - Review quarterly"
}
```

### Service Account Pattern

```hcl
# Different service accounts for different purposes
resource "hiiretail_iam_role_binding" "monitoring_access" {
  role_id = "monitoring-role"
  bindings = [
    { type = "service_account", id = "prometheus-sa" },
    { type = "service_account", id = "grafana-sa" },
    { type = "service_account", id = "alertmanager-sa" }
  ]
  description = "Monitoring system service accounts"
}

resource "hiiretail_iam_role_binding" "backup_access" {
  role_id = "backup-role"
  bindings = [
    { type = "service_account", id = "backup-sa" }
  ]
  description = "Backup system access"
}
```

### Cross-Functional Team Pattern

```hcl
# Cross-functional team with mixed access needs
resource "hiiretail_iam_role_binding" "product_team" {
  role_id = "product-team-role"
  bindings = [
    # Product managers
    { type = "user", id = "pm1@company.com" },
    { type = "user", id = "pm2@company.com" },
    # Development group
    { type = "group", id = "dev-team" },
    # Design group
    { type = "group", id = "design-team" },
    # Analytics service account
    { type = "service_account", id = "analytics-sa" }
  ]
  description = "Cross-functional product team access"
}
```

## Import Examples

```bash
# Import existing role binding
terraform import hiiretail_iam_role_binding.example rb-12345678-1234-1234-1234-123456789012

# Import multiple role bindings
terraform import hiiretail_iam_role_binding.team_access rb-87654321-4321-4321-4321-210987654321
terraform import hiiretail_iam_role_binding.admin_access rb-11111111-2222-3333-4444-555555555555
```

## Troubleshooting

### Common Issues and Solutions

1. **Too Many Bindings Error**
   ```
   Error: Maximum of 10 bindings allowed per role binding
   ```
   Solution: Split into multiple role bindings or use groups instead of individual users.

2. **Invalid Binding Type**
   ```
   Error: Invalid binding type "invalid_type"
   ```
   Solution: Use only "user", "group", or "service_account".

3. **Empty Binding ID**
   ```
   Error: Binding ID cannot be empty
   ```
   Solution: Ensure all binding IDs are specified and non-empty.

4. **Role Not Found**
   ```
   Error: Role "non-existent-role" not found
   ```
   Solution: Verify the role exists and you have permission to bind to it.