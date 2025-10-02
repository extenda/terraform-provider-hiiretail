# Quickstart: OAuth2 Authentication Setup
**Feature**: OAuth2 Authentication with Environment-Specific Endpoints  
**Time to Complete**: ~10 minutes

## Prerequisites

- Go 1.21+ installed
- HiiRetail IAM OAuth2 credentials (client_id and client_secret)
- Valid HiiRetail tenant ID
- Access to HiiRetail IAM API (live or test environment)

## Step 1: Configure OAuth2 Credentials

### Option A: Terraform Configuration
```hcl
# terraform/main.tf
terraform {
  required_providers {
    hiiretail-iam = {
      source = "local/hiiretail-iam"
    }
  }
}

provider "hiiretail-iam" {
  client_id     = "your-oauth2-client-id"
  client_secret = "your-oauth2-client-secret"
  tenant_id     = "your-tenant-id"
}

# Example IAM resource
resource "hiiretail-iam_custom_role" "example" {
  name        = "example-role"
  description = "Example custom role"
  permissions = ["iam:roles:read", "iam:users:read"]
}
```

### Option B: Environment Variables
```bash
export HIIRETAIL_CLIENT_ID="your-oauth2-client-id"
export HIIRETAIL_CLIENT_SECRET="your-oauth2-client-secret"
export HIIRETAIL_TENANT_ID="your-tenant-id"
```

```hcl
# terraform/main.tf (with env vars)
provider "hiiretail-iam" {
  # Credentials loaded from environment variables
}
```

## Step 2: Test Authentication

### Manual Test (Demo Program)
```bash
# Run the OAuth2 demo
cd demo/
go run oauth2_demo.go

# Expected output:
# ✅ OAuth2 authentication successful
# ✅ Token acquired: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
# ✅ Environment detected: Live (iam-api.retailsvc.com)
# ✅ IAM API connection verified
```

### Terraform Test
```bash
# Initialize and plan
terraform init
terraform plan

# Expected output should show:
# - Provider authentication successful
# - Resources planned without authentication errors
# - Correct API endpoint detected based on tenant
```

## Step 3: Verify Environment Detection

### Live Tenant Example
```bash
export HIIRETAIL_TENANT_ID="production-company-123"
go run oauth2_demo.go

# Expected output:
# Environment detected: Live (iam-api.retailsvc.com)
```

### Test Tenant Example
```bash
export HIIRETAIL_TENANT_ID="test-company-123"
go run oauth2_demo.go

# Expected output:
# Environment detected: Test (iam-api.retailsvc-test.com)
```

## Step 4: Test Mock Server Mode (Development)

### Setup Mock Server
```bash
# Start mock OAuth2 and IAM API server
export HIIRETAIL_AUTH_URL="http://localhost:8080/oauth/token"
export HIIRETAIL_API_URL="http://localhost:8080/api"
export HIIRETAIL_MOCK_MODE="true"

# Run tests
go test ./tests/auth/integration_test.go -v

# Expected output:
# ✅ Mock server authentication test passed
# ✅ Mock API endpoint test passed
# ✅ Token refresh test passed
```

## Step 5: Validate Security Configuration

### Check TLS Enforcement
```bash
# This should fail (non-TLS URL)
export HIIRETAIL_AUTH_URL="http://auth.retailsvc.com/oauth/token"
go run oauth2_demo.go

# Expected error:
# ❌ Error: TLS required for OAuth2 authentication
```

### Check Credential Protection
```bash
# Enable debug logging
export TF_LOG=DEBUG
terraform plan

# Verify logs:
# ✅ Client credentials should be redacted as [REDACTED]
# ✅ Access tokens should be redacted as [REDACTED]
# ✅ No sensitive data visible in debug output
```

## Common Issues and Solutions

### Issue: "Invalid client credentials"
**Solution**: Verify client_id and client_secret are correct and have proper IAM scope permissions.

### Issue: "Tenant not found"
**Solution**: Check tenant_id format and ensure it's accessible with your OAuth2 credentials.

### Issue: "Wrong environment detected"
**Solutions**:
- Use `HIIRETAIL_FORCE_TEST_ENV=true` for test environment override
- Verify tenant ID contains test/dev/staging identifiers for auto-detection

### Issue: "TLS certificate errors"
**Solution**: Ensure system certificates are up to date and network allows HTTPS connections.

### Issue: "Token expired" errors
**Solution**: This should be handled automatically. If persistent, check system clock synchronization.

## Success Criteria

After completing this quickstart, you should have:

✅ **Authentication Working**: OAuth2 tokens acquired successfully  
✅ **Environment Detection**: Correct API endpoint selected based on tenant  
✅ **Terraform Integration**: Provider authenticates in Terraform operations  
✅ **Security Validated**: Credentials protected and TLS enforced  
✅ **Error Handling**: Clear error messages for common issues  

## Next Steps

- **Production Deployment**: Remove debug configurations and test overrides
- **Resource Management**: Start creating IAM resources using the authenticated provider
- **Monitoring**: Set up logging to track authentication health and performance
- **Team Onboarding**: Share environment variable configurations with team members

## Troubleshooting Support

For additional support:
1. Check provider logs with `TF_LOG=DEBUG`
2. Test OAuth2 flow with demo program
3. Verify network connectivity to auth.retailsvc.com and iam-api endpoints
4. Review tenant permissions in HiiRetail IAM console