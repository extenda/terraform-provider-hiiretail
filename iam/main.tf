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

variable "client_id" {
  description = "OAuth2 client ID"
  type        = string
}

variable "client_secret" {
  description = "OAuth2 client secret"
  type        = string
  sensitive   = true
}

variable "tenant_id" {
  description = "Tenant ID"
  type        = string
  default     = "CIR7nQwtS0rA6t0S6ejd"
}

resource "hiiretail_iam_group" "test_group" {
  name        = "testShayneGroup"
  description = "This is my second description"
  members     = []
}

resource "hiiretail_iam_custom_role" "test_custom_role" {
  id          = "custom.TerraformTest"
  name        = "TerraformTest"
  description = "Test custom role created via Terraform"
  
  # Define permissions for the custom role - using standard IAM permissions
  permissions = [
    {
      id = "rec.reconciliations.approve"
      attributes = {}
    }
  ]
}

resource "hiiretail_iam_resource" "test_bu" {
  id   = "bu:tf01"
  name = "Terraform Store"
  props = jsonencode({
    type     = "business-unit"
  })
}

resource "hiiretail_iam_role_binding" "test_role_binding" {
  # Name for the role binding
  name = "test-role-binding-shayne"
  
  # Assign the custom role to the group  
  role = "roles/${hiiretail_iam_custom_role.test_custom_role.id}"
  
  # Use members array with the hiiretail_iam_resource reference (array of strings as requested)
  members = [
    hiiretail_iam_resource.test_bu.id
  ]
}

output "test_bu_resource" {
  value = {
    id        = hiiretail_iam_resource.test_bu.id
    name      = hiiretail_iam_resource.test_bu.name
    tenant_id = hiiretail_iam_resource.test_bu.tenant_id
  }
}

# output "test_department_resource" {
#   value = {
#     id        = hiiretail_iam_resource.test_department.id
#     name      = hiiretail_iam_resource.test_department.name
#     tenant_id = hiiretail_iam_resource.test_department.tenant_id
#   }
# }

output "created_group_name" {
  value = hiiretail_iam_group.test_group.name
}
