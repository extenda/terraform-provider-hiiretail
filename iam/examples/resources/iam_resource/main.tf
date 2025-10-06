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

# Department resource
resource "hiiretail_iam_resource" "finance_dept" {
  id   = "dept:finance"
  name = "Finance Department"
  props = jsonencode({
    type = "department"
    cost_center = "FC001"
  })
}

# Application resource
resource "hiiretail_iam_resource" "pos_system" {
  id   = "app:pos-system"
  name = "Point of Sale System"
  props = jsonencode({
    type = "application"
    version = "2.1.0"
  })
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
  default     = "CIR7nQwtS0rA6t0S6ejd"
}

# Outputs
output "store_resource_id" {
  description = "ID of the store resource"
  value       = hiiretail_iam_resource.store_001.id
}

output "finance_resource_id" {
  description = "ID of the finance department resource"
  value       = hiiretail_iam_resource.finance_dept.id
}

output "pos_resource_id" {
  description = "ID of the POS system resource"
  value       = hiiretail_iam_resource.pos_system.id
}