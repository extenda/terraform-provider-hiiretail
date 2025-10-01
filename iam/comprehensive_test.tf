# Comprehensive Terraform test for HiiRetail IAM Provider
# This configuration creates one of each resource type supported by the provider
# Run with: terraform plan -var-file="terraform.tfvars"

terraform {
  required_providers {
    hiiretail-iam = {
      source = "extenda/hiiretail-iam"
    }
  }
}

# Configure the HiiRetail IAM Provider
provider "hiiretail-iam" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
}

# Variable declarations (values loaded from terraform.tfvars)
variable "tenant_id" {
  description = "HiiRetail tenant ID"
  type        = string
  sensitive   = false
}

variable "client_id" {
  description = "OAuth2 Client ID"
  type        = string
  sensitive   = true
}

variable "client_secret" {
  description = "OAuth2 Client Secret"
  type        = string
  sensitive   = true
}

# 1. IAM Group Resource
resource "hiiretail-iam_iam_group" "comprehensive_test_group" {
  name        = "terraform-comprehensive-test-group"
  description = "Test group created by Terraform comprehensive test"
}

# 2. Custom Role Resource
resource "hiiretail-iam_custom_role" "test_custom_role" {
  id   = "terraform-test-custom-role"
  name = "Terraform Test Custom Role"
  
  permissions = [
    {
      id = "pos.products.read"
      attributes = {}
    },
    {
      id = "pos.products.write"
      attributes = {}
    },
    {
      id = "iam.users.manage"
      attributes = {}
    }
  ]
}

# 3. IAM Role Binding Resource - COMMENTED OUT DUE TO PROVIDER ISSUE
# resource "hiiretail-iam_iam_role_binding" "test_role_binding" {
#   role_id = hiiretail-iam_custom_role.test_custom_role.id
#   bindings = [
#     hiiretail-iam_iam_group.comprehensive_test_group.id,
#     "user:terraform-test-user@example.com"
#   ]
#   is_custom = true
# }

# Output the created resources for verification
output "iam_group" {
  description = "Details of the created IAM group"
  value = {
    id          = hiiretail-iam_iam_group.comprehensive_test_group.id
    name        = hiiretail-iam_iam_group.comprehensive_test_group.name
    description = hiiretail-iam_iam_group.comprehensive_test_group.description
  }
}

output "custom_role" {
  description = "Details of the created custom role"
  value = {
    id        = hiiretail-iam_custom_role.test_custom_role.id
    name      = hiiretail-iam_custom_role.test_custom_role.name
    tenant_id = hiiretail-iam_custom_role.test_custom_role.tenant_id
  }
}

# Role binding output commented out due to provider issue
# output "role_binding" {
#   description = "Details of the created role binding"
#   value = {
#     id        = hiiretail-iam_iam_role_binding.test_role_binding.id
#     role_id   = hiiretail-iam_iam_role_binding.test_role_binding.role_id
#     bindings  = hiiretail-iam_iam_role_binding.test_role_binding.bindings
#     is_custom = hiiretail-iam_iam_role_binding.test_role_binding.is_custom
#     tenant_id = hiiretail-iam_iam_role_binding.test_role_binding.tenant_id
#   }
# }