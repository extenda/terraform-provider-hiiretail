terraform {
  required_version = ">= 1.0"
  
  required_providers {
    hiiretail = {
      source  = "extenda/hiiretail"
      version = "~> 1.0"
    }
  }
}

provider "hiiretail" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  token_url     = var.token_url
}

# Basic resource with minimal configuration
resource "hiiretail_iam_resource" "basic_store" {
  id   = "store:001"
  name = "Basic Store Example"
}

# Example of updating the resource (change name)
resource "hiiretail_iam_resource" "basic_department" {
  id   = "dept:electronics"
  name = "Basic Electronics Department"
}