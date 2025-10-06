terraform {
  required_providers {
    hiiretail = {
      source = "registry.terraform.io/extenda/hiiretail"
    }
  }
}

# Basic provider configuration
provider "hiiretail" {
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}

# Provider configuration with custom endpoints
provider "hiiretail" {
  alias = "custom"
  
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
  base_url      = "https://api.custom.hiiretail.com"
  iam_endpoint  = "/api/v2"
  token_url     = "https://oauth2.custom.hiiretail.com/token"
}

# Variables for provider configuration
variable "client_id" {
  description = "OAuth2 client ID for HiiRetail authentication"
  type        = string
  sensitive   = true
}

variable "client_secret" {
  description = "OAuth2 client secret for HiiRetail authentication"
  type        = string
  sensitive   = true
}

variable "tenant_id" {
  description = "Tenant ID for scoping HiiRetail resources"
  type        = string
  default     = "CIR7nQwtS0rA6t0S6ejd"
}