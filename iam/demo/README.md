# OAuth2 Authentication Demo

This demonstration program showcases the OAuth2 authentication capabilities of the HiiRetail IAM Terraform provider.

## Features Demonstrated

- ✅ **Configuration Loading**: Environment variable configuration
- ✅ **Configuration Validation**: Comprehensive validation rules  
- ✅ **OAuth2 Client Creation**: Client credentials flow setup
- ✅ **Endpoint Resolution**: Automatic environment-based endpoint discovery
- ✅ **Token Acquisition**: OAuth2 access token retrieval
- ✅ **Authenticated Requests**: HTTP client with automatic authentication
- ✅ **Retry Logic**: Automatic token refresh on 401 errors
- ✅ **Token Management**: Token validation and refresh capabilities
- ✅ **Security Features**: Credential redaction and secure handling

## Prerequisites

1. **HiiRetail IAM OAuth2 Credentials**: You need valid OAuth2 client credentials
2. **Network Access**: Access to HiiRetail IAM authentication and API endpoints
3. **Go 1.21+**: Required for building and running the demo

## Setup

### 1. Set Environment Variables

```bash
export HIIRETAIL_TENANT_ID="your-tenant-id"
export HIIRETAIL_CLIENT_ID="your-oauth2-client-id" 
export HIIRETAIL_CLIENT_SECRET="your-oauth2-client-secret"

# Optional: specify environment (defaults to production)
export HIIRETAIL_ENVIRONMENT="test"  # or "production", "dev"

# Optional: custom endpoints (overrides automatic resolution)
export HIIRETAIL_AUTH_URL="https://custom-auth.example.com/oauth/token"
export HIIRETAIL_API_URL="https://custom-api.example.com"

# Optional: configuration overrides
export HIIRETAIL_SCOPES="iam:read,iam:write,iam:admin"
export HIIRETAIL_TIMEOUT_SECONDS="30"
export HIIRETAIL_MAX_RETRIES="3"
```

### 2. Build the Demo

```bash
cd demo
go build -o oauth2-demo .
```

### 3. Run the Demo

```bash
./oauth2-demo
```

## Expected Output

The demo will run through 10 steps:

```
🚀 HiiRetail IAM OAuth2 Authentication Demo
==========================================

📋 Step 1: Loading OAuth2 Configuration
   ✅ Tenant ID: your-tenant-id
   ✅ Client ID: abcd...xyz
   ✅ Environment: test

🔍 Step 2: Validating Configuration
   ✅ Configuration is valid

🔐 Step 3: Creating OAuth2 Authentication Client
   ✅ OAuth2 client created successfully

🔗 Step 4: Demonstrating Endpoint Resolution
   ✅ Auth URL: https://auth.retailsvc-test.com/oauth/token
   ✅ API URL: https://iam-api.retailsvc-test.com
   ℹ️  Using test environment endpoints

🎫 Step 5: Acquiring OAuth2 Access Token
   ✅ Token acquired successfully
   ✅ Token type: Bearer
   ✅ Expires in: 59m59s
   ✅ Access token: eyJhbGci...k7XYZ9w

🌐 Step 6: Creating Authenticated HTTP Client
   ✅ Authenticated HTTP client created

📡 Step 7: Making Authenticated API Requests
   📊 Making GET request to list roles...
   ✅ Roles request successful (Status: 200)
   👥 Making GET request to list groups...
   ✅ Groups request successful (Status: 200)
   👤 Making GET request to get user info...
   ✅ User info request successful (Status: 200)

🔄 Step 8: Demonstrating Retry Logic with Token Refresh
   🔄 Testing retry logic with authenticated requests...
   ✅ Request with retry successful (Status: 200)
   ℹ️  The client will automatically refresh tokens on 401 errors

✅ Step 9: Demonstrating Token Validation
   🔄 Forcing token refresh for demonstration...
   ✅ Token refreshed successfully
   ✅ New token expires in: 59m59s

⚙️  Step 10: Demonstrating Configuration Variations
   🔧 Configuration Options:
   • Environment Variables: [list of variables]
   📝 Example Terraform Configuration: [example]
   🔐 Security Features: [features list]

🎉 Demo completed successfully!
   All OAuth2 authentication features are working correctly.
```

## Troubleshooting

### Authentication Errors

If you see authentication errors:

1. **Verify Credentials**: Ensure `HIIRETAIL_CLIENT_ID` and `HIIRETAIL_CLIENT_SECRET` are correct
2. **Check Tenant ID**: Verify `HIIRETAIL_TENANT_ID` matches your organization
3. **Network Access**: Ensure you can reach the authentication endpoints
4. **Environment**: Check if you need `HIIRETAIL_ENVIRONMENT=test` for test credentials

### Configuration Errors

If you see configuration validation errors:

1. **Required Fields**: Ensure all required environment variables are set
2. **Format Validation**: Check that tenant ID, client ID formats are correct
3. **URL Validation**: If using custom URLs, ensure they use HTTPS and are valid

### Network Errors

If you see network connectivity issues:

1. **Firewall**: Ensure outbound HTTPS access to `*.retailsvc.com` domains
2. **Proxy**: Configure HTTP proxy settings if required by your environment
3. **DNS**: Verify DNS resolution for HiiRetail endpoints

## Integration with Terraform

This demo shows the same authentication flow used by the Terraform provider. In Terraform:

```hcl
provider "hiiretail-iam" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-oauth2-client-id"
  client_secret = "your-oauth2-client-secret"
  
  # Optional configuration
  scopes          = ["iam:read", "iam:write"]
  timeout_seconds = 30
  max_retries     = 3
}

resource "hiiretail-iam_custom_role" "example" {
  id   = "custom-role-example"
  name = "Example Custom Role"
  
  permissions = [
    {
      id = "resource:action"
      attributes = {
        resource_type = "example"
      }
    }
  ]
}
```

## Security Notes

- 🔒 **Credentials**: Never commit credentials to version control
- 🔍 **Logging**: Credentials are automatically redacted in logs
- 🔄 **Token Rotation**: Tokens are automatically refreshed before expiration
- 🌐 **HTTPS**: All communication uses HTTPS encryption
- 💾 **Caching**: Tokens are cached securely with integrity validation

## Development

To modify the demo:

1. Edit `main.go` to add new demonstration scenarios
2. Build with `go build -o oauth2-demo .`
3. Test with your OAuth2 credentials
4. Check logs for detailed authentication flow information

The demo uses the same auth package as the Terraform provider, so any changes to the auth package will be reflected in the demo.