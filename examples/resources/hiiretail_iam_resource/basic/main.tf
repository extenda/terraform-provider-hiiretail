terraform {
  required_providers {
    hiiretail = {
      source = "registry.terraform.io/extenda/hiiretail"
    }
  }
}

provider "hiiretail" {
  # Authentication will use precedence: terraform.tfvars → TF_VAR_* → HIIRETAIL_* → error
  # No explicit configuration needed when using environment variables
}

resource "hiiretail_iam_group" "some_group" {
  name        = "Store-tf01-Financial-Managers"
  description = "Members of this Group can manage Financial activities in Business Unit tf01"
  members     = []
}

resource "hiiretail_iam_custom_role" "some_custom_role" {
  id          = "ReconciliationApprover"
  name        = "ReconciliationApprover"
  description = "Role for approving reconciliations"

  permissions = [
    {
      id = "rec.reconciliations.approve"
    }
  ]
}

resource "hiiretail_iam_resource" "some_store" {
  id   = "bu:tf01"
  name = "Some Store"
  props = jsonencode({
    type     = "business-unit"
  })
}

resource "hiiretail_iam_role_binding" "custom_role_binding" {
  group_id  = hiiretail_iam_group.some_group.id
  role_id   = hiiretail_iam_custom_role.some_custom_role.id
  is_custom = true
  bindings  = [hiiretail_iam_resource.some_store.id]
}

resource "hiiretail_iam_role_binding" "builtin_role_binding" {
  group_id  = hiiretail_iam_group.some_group.id  
  role_id   = "rec.manager"
  is_custom = false
  bindings  = ["*"]
}

output "some_store_resource" {
  value = {
    id        = hiiretail_iam_resource.some_store.id
    name      = hiiretail_iam_resource.some_store.name
    tenant_id = hiiretail_iam_resource.some_store.tenant_id
  }
}

output "created_group_name" {
  value = hiiretail_iam_group.some_group.name
}
