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

# Business unit resource
resource "hiiretail_iam_resource" "store_001" {
  id   = "bu:store001"
  name = "Store 001"
  props = jsonencode({
    type     = "business-unit"
    location = "downtown"
  })
}

# Outputs
output "store_resource_id" {
  description = "ID of the store resource"
  value       = hiiretail_iam_resource.store_001.id
}