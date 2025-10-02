# Test Terraform configuration for HiiRetail IAM Group resource
# This validates the complete OAuth2 authentication flow

terraform {
  required_providers {
    hiiretail-iam = {
      source = "extenda/hiiretail-iam"
    }
  }
}

# Configure the HiiRetail IAM provider with OAuth2 authentication
provider "hiiretail-iam" {
  # OAuth2 credentials (can be set via environment variables)
  # HIIRETAIL_TENANT_ID, HIIRETAIL_CLIENT_ID, HIIRETAIL_CLIENT_SECRET
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  
  # Optional configuration
  scopes = ["iam:read", "iam:write", "groups:read", "groups:write"]
  timeout_seconds = 30
  max_retries = 3
}

# Variables for OAuth2 credentials
variable "tenant_id" {
  description = "HiiRetail tenant ID"
  type        = string
  sensitive   = false
}

variable "client_id" {
  description = "OAuth2 client ID"
  type        = string
  sensitive   = false
}

variable "client_secret" {
  description = "OAuth2 client secret"
  type        = string
  sensitive   = true
}

# Create a test IAM group
resource "hiiretail-iam_iam_group" "test_group" {
  name        = "terraform-test-group"
  description = "Test group created by Terraform with OAuth2 authentication"
}

# Output the created group information
output "group_id" {
  description = "ID of the created group"
  value       = hiiretail-iam_iam_group.test_group.id
}

output "group_name" {
  description = "Name of the created group"
  value       = hiiretail-iam_iam_group.test_group.name
}

output "group_description" {
  description = "Description of the created group"
  value       = hiiretail-iam_iam_group.test_group.description
}