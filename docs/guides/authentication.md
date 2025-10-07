---
page_title: "HiiRetail Provider Authentication"
subcategory: "Authentication"
---

# HiiRetail Terraform Provider Authentication Guide

This guide provides detailed information about authenticating with the HiiRetail Terraform Provider.

## Authentication Method

The HiiRetail provider uses OAuth2 client credentials flow with Basic authentication.

### Required Credentials

The provider requires three authentication parameters:

- `client_id`: OAuth2 client identifier
- `client_secret`: OAuth2 client secret  
- `tenant_id`: Tenant identifier for resources

### Authentication Schema

**Important**: The Authorization header must use Basic schema and be a base64 encoded `client_id` and `client_secret`, concatenated with a colon.

Format: `Authorization: Basic base64(client_id:client_secret)`

## Configuration Examples

**⚠️ SECURITY WARNING**: Never use real credentials in documentation. Always use placeholder values and obtain actual credentials through secure channels.

```hcl
# Recommended: Empty provider block with environment variables
provider "hiiretail" {
  # Authentication via environment variables:
  # - HIIRETAIL_CLIENT_ID
  # - HIIRETAIL_CLIENT_SECRET
  # - HIIRETAIL_TENANT_ID
  
  # Optional configuration can still be specified
  timeout_seconds = 30
  max_retries = 3
}

# Alternative: Explicit configuration (not recommended for security)
# provider "hiiretail" {
#   client_id     = var.client_id
#   client_secret = var.client_secret
#   tenant_id     = var.tenant_id
#   timeout_seconds = 30
#   max_retries = 3
# }
```

### Obtaining Credentials

Contact your HiiRetail administrator to obtain:

- **Client ID**: OAuth2 client identifier for your application
- **Client Secret**: OAuth2 client secret (keep secure, never commit to version control)
- **Tenant ID**: Your organization's tenant identifier

## Environment Variables

You can set credentials via environment variables using multiple formats:

### HIIRETAIL Environment Variables
```bash
export HIIRETAIL_CLIENT_ID="your-oauth2-client-id"
export HIIRETAIL_CLIENT_SECRET="your-oauth2-client-secret"
export HIIRETAIL_TENANT_ID="your-tenant-id"
export HIIRETAIL_SCOPES="IAM:read:groups,IAM:read:roles"  # Optional
export HIIRETAIL_TIMEOUT_SECONDS="30"  # Optional
export HIIRETAIL_MAX_RETRIES="3"  # Optional
```

### Terraform Environment Variables
```bash
export TF_VAR_client_id="your-oauth2-client-id"
export TF_VAR_client_secret="your-oauth2-client-secret"
export TF_VAR_tenant_id="your-tenant-id"
export TF_VAR_scopes="IAM:read:groups,IAM:read:roles"  # Optional
export TF_VAR_timeout_seconds="30"  # Optional
export TF_VAR_max_retries="3"  # Optional
```

## Authentication Precedence

The provider follows a strict precedence order for authentication configuration:

1. **Terraform configuration** (terraform.tfvars or provider block)
2. **TF_VAR_* environment variables** (standard Terraform pattern)
3. **HIIRETAIL_* environment variables** (provider-specific pattern)
4. **Error** if no credentials are found

This precedence allows for flexible configuration management:
- Use `terraform.tfvars` for local development
- Use `TF_VAR_*` variables in CI/CD pipelines following Terraform conventions
- Use `HIIRETAIL_*` variables for provider-specific environments
- Provider will error with clear messages if no valid credentials are found

### Example Precedence Behavior
```bash
# If all three are set, terraform.tfvars takes precedence
export TF_VAR_tenant_id="tf-var-tenant"
export HIIRETAIL_TENANT_ID="hiiretail-tenant"
# terraform.tfvars contains: tenant_id = "tfvars-tenant"
# Result: "tfvars-tenant" is used

# If only environment variables are set, TF_VAR_* takes precedence
export TF_VAR_tenant_id="tf-var-tenant"
export HIIRETAIL_TENANT_ID="hiiretail-tenant"
# Result: "tf-var-tenant" is used

# If only HIIRETAIL_* is set, it will be used
export HIIRETAIL_TENANT_ID="hiiretail-tenant"
# Result: "hiiretail-tenant" is used
```

## Authentication Recommendations

**🌟 RECOMMENDED**: Use environment variables for better security:

```bash
export HIIRETAIL_CLIENT_ID="your-oauth2-client-id"
export HIIRETAIL_CLIENT_SECRET="your-oauth2-client-secret"
export HIIRETAIL_TENANT_ID="your-tenant-id"
```

**Alternative**: Create a `terraform.tfvars` file (less secure):

```hcl
client_id     = "your-oauth2-client-id"
client_secret = "your-oauth2-client-secret"
tenant_id     = "your-tenant-id"
```

**⚠️ IMPORTANT**: If using `terraform.tfvars`, add it to your `.gitignore` file to prevent accidentally committing real credentials.

Then reference in your configuration:

```hcl
# Recommended: Use environment variables instead
provider "hiiretail" {
  # Empty provider block - all authentication via environment variables
}

# If you must use terraform.tfvars, reference variables like this:
# provider "hiiretail" {
#   client_id     = var.client_id
#   client_secret = var.client_secret
#   tenant_id     = var.tenant_id
# }
```

## OAuth2 Flow Details

The provider implements OAuth2 client credentials flow:

1. **Token Request**: Provider sends Basic authenticated request to token endpoint
2. **Token Response**: Receives access token for API calls
3. **API Requests**: Uses Bearer token for all subsequent API calls
4. **Token Refresh**: Automatically handles token expiration and renewal

## Troubleshooting

### Common Authentication Issues

1. **Invalid Credentials**: Verify client_id and client_secret are correct
2. **Tenant Access**: Ensure tenant_id corresponds to accessible tenant
3. **Network Issues**: Check connectivity to HiiRetail OAuth2 endpoints
4. **Token Expiration**: Provider handles this automatically, but check for clock skew

### Debug Authentication

Enable debug logging to troubleshoot authentication issues:

```bash
export TF_LOG=DEBUG
terraform plan
```

## Security Best Practices

1. **Never commit credentials**: Use environment variables or encrypted storage
2. **Rotate secrets regularly**: Update client secrets periodically
3. **Limit scope**: Use tenant-specific credentials when possible
4. **Monitor access**: Review authentication logs regularly
5. **Use .gitignore**: Always add credential files to .gitignore
6. **Secure storage**: Store credentials in secure credential management systems

## API Endpoints

The provider authenticates against these endpoints:

- **Token Endpoint**: `https://auth.retailsvc.com/oauth2/token`
- **API Base**: `https://iam-api.retailsvc.com`

## Support

For authentication issues, contact the HiiRetail support team or refer to the API documentation.

## Security Notice

**CRITICAL**: If you suspect credentials have been exposed in version control:

1. Immediately revoke and rotate the exposed credentials
2. Purge the credentials from git history using tools like BFG Repo-Cleaner
3. Review all commits for potential credential exposure
4. Update documentation to use placeholder values only
5. Implement git hooks to prevent future credential commits