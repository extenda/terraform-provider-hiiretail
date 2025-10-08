terraform {
  required_providers {
    hiiretail = {
      source = "registry.terraform.io/extenda/hiiretail"
    }
  }
}

provider "hiiretail" {
  # Optional configuration
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
resource "hiiretail_iam_group" "cashiers" {
  name = "Cashiers"
  # description will be computed if not provided
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

output "cashiers_group_id" {
  description = "ID of the cashiers group"
  value       = hiiretail_iam_group.cashiers.id
}