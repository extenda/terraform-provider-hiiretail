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

# Create supporting resources
resource "hiiretail_iam_group" "finance_team" {
  name        = "FinanceTeam"
  description = "Finance team members"
}

resource "hiiretail_iam_custom_role" "approver" {
  id          = "ReconciliationApprover"
  name        = "Reconciliation Approver"
  description = "Role for approving reconciliations"

  permissions = [
    {
      id = "rec.reconciliations.approve"
    }
  ]
}

resource "hiiretail_iam_resource" "store" {
  id   = "bu:store001"
  name = "Store 001"
  props = jsonencode({
    type = "business-unit"
  })
}

# Custom role binding to specific resource
resource "hiiretail_iam_role_binding" "finance_approver" {
  group_id  = hiiretail_iam_group.finance_team.id
  role_id   = hiiretail_iam_custom_role.approver.id
  is_custom = true
  bindings  = [hiiretail_iam_resource.store.id]
}

# Built-in role binding to all resources
resource "hiiretail_iam_role_binding" "manager_access" {
  group_id  = hiiretail_iam_group.finance_team.id
  role_id   = "rec.manager"
  is_custom = false
  bindings  = ["*"]
}

# Outputs
output "custom_role_binding_id" {
  description = "ID of the custom role binding"
  value       = hiiretail_iam_role_binding.finance_approver.id
}

output "builtin_role_binding_id" {
  description = "ID of the built-in role binding"
  value       = hiiretail_iam_role_binding.manager_access.id
}