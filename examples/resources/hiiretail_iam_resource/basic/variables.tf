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
