terraform {
  required_providers {
    hiiretail = {
      source = "extenda/hiiretail"
    }
  }
}

provider "hiiretail" {
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}

# Fetch an existing IAM resource by ID
data "hiiretail_iam_resource" "example_store" {
  id = "store:001"
}

# Output the resource details
output "resource_name" {
  description = "The name of the resource"
  value       = data.hiiretail_iam_resource.example_store.name
}

output "resource_type" {
  description = "The type of the resource"
  value       = data.hiiretail_iam_resource.example_store.type
}

output "resource_permissions" {
  description = "Available permissions for the resource"
  value       = data.hiiretail_iam_resource.example_store.permissions
}

output "resource_properties" {
  description = "Additional properties of the resource"
  value       = data.hiiretail_iam_resource.example_store.properties
}
