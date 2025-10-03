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
  id          = "custom.TerraformTestShayne"
  name        = "TerraformTestShayne"
  description = "Test custom role created via Terraform"
  
  # Define permissions for the custom role
  permissions = [
    {
      id = "iam.group.list"
      attributes = {}
    }
  ]
}

resource "hiiretail_iam_role_binding" "test_role_binding" {
  # Name for the role binding
  name = "test-role-binding-shayne"
  
  # Assign the custom role to the group
  role = "roles/custom.${hiiretail_iam_custom_role.test_custom_role.id}"
  
  # Bind the role to the group (using group name as member)
  members = [
    "group:${hiiretail_iam_group.test_group.name}"
  ]
}

output "created_group_name" {
  value = hiiretail_iam_group.test_group.name
}

output "role_binding_id" {
  value = hiiretail_iam_role_binding.test_role_binding.id
}