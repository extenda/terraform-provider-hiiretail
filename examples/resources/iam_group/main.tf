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

# Basic group
resource "hiiretail_iam_group" "finance_team" {
  name        = "FinanceTeam"
  description = "Finance team members"
}

# Group with custom ID
resource "hiiretail_iam_group" "admin_group" {
  id          = "admin-group-001"
  name        = "AdminGroup"
  description = "System administrators group"
}

# Group with minimal configuration
resource "hiiretail_iam_group" "developers" {
  name = "Developers"
  # description will be computed if not provided
}

# Variables
variable "client_id" {
  description = "HiiRetail OAuth2 client ID"
  type        = string
  sensitive   = true
}

variable "client_secret" {
  description = "HiiRetail OAuth2 client secret"
  type        = string
  sensitive   = true
}

variable "tenant_id" {
  description = "Tenant ID for scoping HiiRetail resources"
  type        = string
  default     = "your-tenant-id"
}

# Outputs
output "finance_team_id" {
  description = "ID of the finance team group"
  value       = hiiretail_iam_group.finance_team.id
}

output "admin_group_id" {
  description = "ID of the admin group"
  value       = hiiretail_iam_group.admin_group.id
}

output "developers_group_id" {
  description = "ID of the developers group"
  value       = hiiretail_iam_group.developers.id
}