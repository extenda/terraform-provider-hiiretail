---
page_title: "Getting Started with HiiRetail Provider"
subcategory: "Getting Started"
---

# Getting Started with HiiRetail Provider

This guide helps you get started with the HiiRetail Terraform Provider for managing IAM resources.

## Prerequisites

- Terraform 0.14 or later
- HiiRetail OAuth2 credentials
- Go 1.21+ (for development)

## Quick Setup

1. **Configure Provider**:
```hcl
terraform {
  required_providers {
    hiiretail = {
      source = "registry.terraform.io/extenda/hiiretail"
    }
  }
}

provider "hiiretail" {
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}
```

2. **Create Basic Resources**:
```hcl
# Create a group
resource "hiiretail_iam_group" "team" {
  name        = "MyTeam"
  description = "My team group"
}

# Create a custom role
resource "hiiretail_iam_custom_role" "viewer" {
  id          = "CustomViewer"
  name        = "Custom Viewer"
  description = "Custom viewing permissions"

  permissions = [
    {
      id = "iam.groups.view"
    }
  ]
}

# Bind role to group
resource "hiiretail_iam_role_binding" "team_viewer" {
  group_id  = hiiretail_iam_group.team.id
  role_id   = hiiretail_iam_custom_role.viewer.id
  is_custom = true
  bindings  = ["*"]
}
```

3. **Apply Configuration**:
```bash
terraform init
terraform plan
terraform apply
```

## Next Steps

- Review [Authentication Guide](authentication) for credential setup
- Explore [IAM Patterns](iam-patterns) for common use cases
- Check individual resource documentation for advanced features