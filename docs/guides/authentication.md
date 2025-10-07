------

page_title: "HiiRetail Provider Authentication"page_title: "HiiRetail Provider Authentication"

subcategory: "Authentication"subcategoryCreate a `terraform.tfvars` file for reusable credentials:

---

```hcl

# HiiRetail Terraform Provider Authentication Guideclient_id     = "your-oauth2-client-id"

client_secret = "your-oauth2-client-secret"

This guide provides detailed information about authenticating with the HiiRetail Terraform Provider.tenant_id     = "your-tenant-id"

```

## Authentication Method

**⚠️ IMPORTANT**: Add `terraform.tfvars` to your `.gitignore` file to prevent accidentally committing real credentials.entication"

The HiiRetail provider uses OAuth2 client credentials flow with Basic authentication.---



### Required Credentials# HiiRetail Terraform Provider Authentication Guide



The provider requires three authentication parameters:This guide provides detailed information about authenticating with the HiiRetail Terraform Provider.



- `client_id`: OAuth2 client identifier## Authentication Method

- `client_secret`: OAuth2 client secret  

- `tenant_id`: Tenant identifier for resourcesThe HiiRetail provider uses OAuth2 client credentials flow with Basic authentication.



### Authentication Schema### Required Credentials



**Important**: The Authorization header must use Basic schema and be a base64 encoded `client_id` and `client_secret`, concatenated with a colon.The provider requires three authentication parameters:



Format: `Authorization: Basic base64(client_id:client_secret)`- `client_id`: OAuth2 client identifier

- `client_secret`: OAuth2 client secret  

## Configuration Examples- `tenant_id`: Tenant identifier for resources



**⚠️ SECURITY WARNING**: Never use real credentials in documentation. Always use placeholder values and obtain actual credentials through secure channels.### Authentication Schema



```hcl**Important**: The Authorization header must use Basic schema and be a base64 encoded `client_id` and `client_secret`, concatenated with a colon.

provider "hiiretail" {

  client_id     = "your-oauth2-client-id"Format: `Authorization: Basic base64(client_id:client_secret)`

  client_secret = "your-oauth2-client-secret"

  tenant_id     = "your-tenant-id"## Configuration Examples

  

  # Optional configuration**⚠️ SECURITY WARNING**: Never use real credentials in documentation. Always use placeholder values and obtain actual credentials through secure channels.

  timeout_seconds = 30

  max_retries = 3```hcl

}provider "hiiretail" {

```  client_id     = "your-oauth2-client-id"

  client_secret = "your-oauth2-client-secret"

### Obtaining Credentials  tenant_id     = "your-tenant-id"

  

Contact your HiiRetail administrator to obtain:  # Optional configuration

  timeout_seconds = 30

- **Client ID**: OAuth2 client identifier for your application  max_retries = 3

- **Client Secret**: OAuth2 client secret (keep secure, never commit to version control)}

- **Tenant ID**: Your organization's tenant identifier```



## Environment Variables### Obtaining Credentials



You can also set credentials via environment variables:Contact your HiiRetail administrator to obtain:



```bash- **Client ID**: OAuth2 client identifier for your application

export HIIRETAIL_CLIENT_ID="your-oauth2-client-id"- **Client Secret**: OAuth2 client secret (keep secure, never commit to version control)

export HIIRETAIL_CLIENT_SECRET="your-oauth2-client-secret"- **Tenant ID**: Your organization's tenant identifier

export HIIRETAIL_TENANT_ID="your-tenant-id"

```## Environment Variables



## Terraform Variables FileYou can also set credentials via environment variables:



Create a `terraform.tfvars` file for reusable credentials:```bash

export HIIRETAIL_CLIENT_ID="your-oauth2-client-id"

```hclexport HIIRETAIL_CLIENT_SECRET="your-oauth2-client-secret"

client_id     = "your-oauth2-client-id"export HIIRETAIL_TENANT_ID="your-tenant-id"

client_secret = "your-oauth2-client-secret"```

tenant_id     = "your-tenant-id"

```## Terraform Variables File



**⚠️ IMPORTANT**: Add `terraform.tfvars` to your `.gitignore` file to prevent accidentally committing real credentials.Create a `terraform.tfvars` file for reusable credentials:



Then reference in your configuration:```hcl

```hcl

```hclclient_id     = "your-oauth2-client-id"

provider "hiiretail" {client_secret = "your-oauth2-client-secret"

  client_id     = var.client_idtenant_id     = "your-tenant-id"

  client_secret = var.client_secret```

  tenant_id     = var.tenant_id

}**⚠️ IMPORTANT**: Add `terraform.tfvars` to your `.gitignore` file to prevent accidentally committing real credentials.

``````



## OAuth2 Flow DetailsThen reference in your configuration:



The provider implements OAuth2 client credentials flow:```hcl

provider "hiiretail" {

1. **Token Request**: Provider sends Basic authenticated request to token endpoint  client_id     = var.client_id

2. **Token Response**: Receives access token for API calls  client_secret = var.client_secret

3. **API Requests**: Uses Bearer token for all subsequent API calls  tenant_id     = var.tenant_id

4. **Token Refresh**: Automatically handles token expiration and renewal}

```

## Troubleshooting

## OAuth2 Flow Details

### Common Authentication Issues

The provider implements OAuth2 client credentials flow:

1. **Invalid Credentials**: Verify client_id and client_secret are correct

2. **Tenant Access**: Ensure tenant_id corresponds to accessible tenant1. **Token Request**: Provider sends Basic authenticated request to token endpoint

3. **Network Issues**: Check connectivity to HiiRetail OAuth2 endpoints2. **Token Response**: Receives access token for API calls

4. **Token Expiration**: Provider handles this automatically, but check for clock skew3. **API Requests**: Uses Bearer token for all subsequent API calls

4. **Token Refresh**: Automatically handles token expiration and renewal

### Debug Authentication

## Troubleshooting

Enable debug logging to troubleshoot authentication issues:

### Common Authentication Issues

```bash

export TF_LOG=DEBUG1. **Invalid Credentials**: Verify client_id and client_secret are correct

terraform plan2. **Tenant Access**: Ensure tenant_id corresponds to accessible tenant

```3. **Network Issues**: Check connectivity to HiiRetail OAuth2 endpoints

4. **Token Expiration**: Provider handles this automatically, but check for clock skew

## Security Best Practices

### Debug Authentication

1. **Never commit credentials**: Use environment variables or encrypted storage

2. **Rotate secrets regularly**: Update client secrets periodicallyEnable debug logging to troubleshoot authentication issues:

3. **Limit scope**: Use tenant-specific credentials when possible

4. **Monitor access**: Review authentication logs regularly```bash

5. **Use .gitignore**: Always add credential files to .gitignoreexport TF_LOG=DEBUG

6. **Secure storage**: Store credentials in secure credential management systemsterraform plan

```

## API Endpoints

## Security Best Practices

The provider authenticates against these endpoints:

1. **Never commit credentials**: Use environment variables or encrypted storage

- **Token Endpoint**: `https://oauth2.hiiretail.com/token`2. **Rotate secrets regularly**: Update client secrets periodically

- **API Base**: `https://api.hiiretail.com/v1`3. **Limit scope**: Use tenant-specific credentials when possible

4. **Monitor access**: Review authentication logs regularly

## Support

## API Endpoints

For authentication issues, contact the HiiRetail support team or refer to the API documentation.

The provider authenticates against these endpoints:

## Security Notice

- **Token Endpoint**: `https://oauth2.hiiretail.com/token`

**CRITICAL**: If you suspect credentials have been exposed in version control:- **API Base**: `https://api.hiiretail.com/v1`



1. Immediately revoke and rotate the exposed credentials## Support

2. Purge the credentials from git history using tools like BFG Repo-Cleaner

3. Review all commits for potential credential exposureFor authentication issues, contact the HiiRetail support team or refer to the API documentation.
4. Update documentation to use placeholder values only
5. Implement git hooks to prevent future credential commits