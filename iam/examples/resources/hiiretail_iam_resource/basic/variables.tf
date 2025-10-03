variable "tenant_id" {
  description = "HiiRetail tenant identifier"
  type        = string
  sensitive   = false
}

variable "client_id" {
  description = "OAuth2 client ID for authentication"
  type        = string
  sensitive   = false
}

variable "client_secret" {
  description = "OAuth2 client secret for authentication"
  type        = string
  sensitive   = true
}

variable "token_url" {
  description = "OAuth2 token endpoint URL"
  type        = string
  default     = "https://auth.hiiretail.com/oauth/token"
}