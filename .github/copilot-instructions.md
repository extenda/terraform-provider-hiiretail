# hiiretail-terraform-providers Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-09-28

## Active Technologies
- Go 1.21+ + HashiCorp Terraform Plugin Framework, golang.org/x/oauth2, terraform-plugin-framework-validators (001-ensure-that-the)
- HiiRetail IAM API (RESTful service) (002-ensure-that-the)
- Go 1.23.0 + Terraform Plugin Framework v1.4.2, terraform-plugin-testing v1.13.3, testify v1.8.4, golang.org/x/oauth2 v0.26.0 (004-ensure-that-the)
- N/A (Terraform provider for API resource management) (004-ensure-that-the)
- Go 1.21+ (terraform-plugin-framework v1.16.0) + HashiCorp terraform-plugin-framework, terraform-plugin-testing, golang.org/x/oauth2 (005-add-tests-for)
- RESTful API backend with OAuth2 client credentials authentication (005-add-tests-for)
- Go 1.21+ + OAuth2 Discovery, clientcredentials, HTTP client with connection pooling (006-correctly-handle-client)
- Hii Retail OAuth Client Management Service (OCMS) integration (006-correctly-handle-client)
- Go 1.21+ (matches existing Terraform provider) + Terraform Plugin Framework v1.16.0, golang.org/x/oauth2 v0.30.0 (007-oauth2-authentication-with)
- N/A (stateless authentication service) (007-oauth2-authentication-with)
- Go 1.19+ (Terraform Plugin Framework requirement) + HashiCorp Terraform Plugin Framework, golang.org/x/oauth2, existing HiiRetail APIs (008-the-iam-is)
- N/A (provider manages external API resources, not local storage) (008-the-iam-is)
- Go 1.19+ (Terraform Plugin Framework) + HashiCorp Terraform Plugin Framework, OAuth2 client libraries, HiiRetail IAM API clients (009-improve-resource-usability)
- N/A (provider manages remote API resources) (009-improve-resource-usability)

## Project Structure
```
src/
tests/
```

## Commands
# Go 1.21+ with OAuth2 Authentication Enhancement
go mod tidy                    # Update dependencies
go test ./internal/provider/auth -v     # Test OAuth2 components
make build                     # Build provider with OAuth2 support
make test-unit                 # Run unit tests including auth tests
terraform init && terraform plan       # Test OAuth2 provider configuration

# OAuth2 Testing
export HIIRETAIL_TENANT_ID="test-tenant"
export HIIRETAIL_CLIENT_ID="test-client"
export HIIRETAIL_CLIENT_SECRET="test-secret"
TF_LOG=DEBUG terraform plan            # Debug OAuth2 authentication flow

## Code Style
Go 1.21+: Follow standard conventions

### OAuth2 Implementation Guidelines
- Use golang.org/x/oauth2/clientcredentials for OAuth2 client credentials flow
- Implement OAuth2 discovery using /.well-known/openid-configuration endpoints
- Cache discovery responses for 1 hour to minimize network calls
- Use context.Background() for OAuth2 operations to prevent premature cancellation
- Mark client_secret fields as sensitive in Terraform schemas
- Never log or expose client credentials or access tokens
- Implement exponential backoff for retryable authentication errors
- Validate OAuth2 configuration before making API calls
- Use connection pooling for HTTP clients (MaxIdleConns: 100, MaxIdleConnsPerHost: 10)
- Handle token refresh automatically via oauth2.TokenSource
- Provide clear error messages distinguishing between credential, network, and server errors

## Recent Changes
- 009-improve-resource-usability: Added Go 1.19+ (Terraform Plugin Framework) + HashiCorp Terraform Plugin Framework, OAuth2 client libraries, HiiRetail IAM API clients
- 008-the-iam-is: Added Go 1.19+ (Terraform Plugin Framework requirement) + HashiCorp Terraform Plugin Framework, golang.org/x/oauth2, existing HiiRetail APIs
- 007-oauth2-authentication-with: Added Go 1.21+ (matches existing Terraform provider) + Terraform Plugin Framework v1.16.0, golang.org/x/oauth2 v0.30.0

### OAuth2 Authentication Enhancement (006-correctly-handle-client)
**Key Components**:

**Testing Strategy**:

**Configuration Support**:

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
