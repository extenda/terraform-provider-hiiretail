terraform {
  required_version = ">= 1.0"
  
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 1.0"
    }
  }
}

# Configure the HiiRetail Provider
provider "hiiretail" {
  # Credentials are provided via environment variables:
  # - HIIRETAIL_CLIENT_ID
  # - HIIRETAIL_CLIENT_SECRET
  # - HIIRETAIL_TENANT_ID
}

# Example: Create an IAM group
resource "hiiretail_iam_group" "example" {
  name        = "terraform-registry-example"
  description = "Example group created using provider from Terraform Registry"
}

# Example: Create a custom role
resource "hiiretail_iam_custom_role" "example" {
  name        = "registry-example-role"
  description = "Example custom role created from Terraform Registry"
  
  permissions = [
    "iam:groups:read",
    "iam:groups:write"
  ]
}

# Example: Create a resource
resource "hiiretail_iam_resource" "example" {
  name        = "registry-example-resource"
  description = "Example resource created from Terraform Registry"
  type        = "application"
}

# Example: Create a role binding
resource "hiiretail_iam_role_binding" "example" {
  name        = "registry-example-binding"
  description = "Example role binding created from Terraform Registry"
  
  group_id    = hiiretail_iam_group.example.id
  role_id     = hiiretail_iam_custom_role.example.id
  resource_id = hiiretail_iam_resource.example.id
}

# Outputs
output "group_id" {
  description = "The ID of the created IAM group"
  value       = hiiretail_iam_group.example.id
}

output "role_id" {
  description = "The ID of the created custom role"
  value       = hiiretail_iam_custom_role.example.id
}

output "resource_id" {
  description = "The ID of the created resource"
  value       = hiiretail_iam_resource.example.id
}

output "role_binding_id" {
  description = "The ID of the created role binding"
  value       = hiiretail_iam_role_binding.example.id
}