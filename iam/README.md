# HiiRetail IAM Terraform Provider

This Terraform provider enables management of HiiRetail IAM resources using Terraform.

## Features

- **OIDC Client Credentials Authentication**: Secure authentication using OIDC client credentials flow
- **Configurable Base URL**: Optional base_url parameter to connect to different API endpoints
- **Comprehensive Error Handling**: Clear error messages for configuration and authentication issues
- **Thorough Testing**: Full test coverage including unit tests and integration tests

## Requirements

- Terraform >= 1.0
- Go >= 1.21 (for development)

## Configuration

### Provider Configuration

```hcl
terraform {
  required_providers {
    hiiretail_iam = {
      source = "extenda/hiiretail_iam"
    }
  }
}

provider "hiiretail_iam" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-oidc-client-id"
  client_secret = "your-oidc-client-secret"
  base_url      = "https://custom-api.example.com" # Optional, defaults to https://iam-api.retailsvc-test.com
}
```

### Configuration Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tenant_id` | string | Yes | Tenant ID to use for all IAM API requests |
| `client_id` | string | Yes | OIDC client ID for IAM API authentication |
| `client_secret` | string | Yes | OIDC client secret for IAM API authentication (sensitive) |
| `base_url` | string | No | Base URL of the IAM API (defaults to https://iam-api.retailsvc-test.com) |

### Authentication

The provider uses OIDC client credentials flow for authentication. The authentication process:

1. Uses the provided `client_id` and `client_secret` to obtain an access token
2. Automatically refreshes tokens when they expire
3. Includes the access token in all API requests

### Environment Variables

You can also configure the provider using environment variables:

```bash
export TF_VAR_tenant_id="your-tenant-id"
export TF_VAR_client_id="your-oidc-client-id"
export TF_VAR_client_secret="your-oidc-client-secret"
export TF_VAR_base_url="https://custom-api.example.com"
```

## Development

### Building the Provider

```bash
go build .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./internal/provider

# Run specific test
go test -v ./internal/provider -run TestHiiRetailIamProvider
```

### Test Coverage

The provider includes comprehensive test coverage:

- **Unit Tests**: Test provider configuration, schema validation, and OIDC setup
- **Integration Tests**: Test actual OIDC authentication flow with mock servers
- **Configuration Validation Tests**: Test various configuration scenarios and error conditions

### Code Generation

This provider uses generated code. The main provider schema and resources are generated from configuration files. Manual modifications should be made carefully to avoid conflicts during regeneration.

## Error Handling

The provider provides clear error messages for common issues:

- **Missing Required Parameters**: Clear messages indicating which parameters are required
- **Invalid URLs**: Validation of base_url format with helpful error messages
- **Authentication Failures**: OIDC authentication errors with context
- **Network Issues**: Connection and timeout errors with diagnostic information

## Security Considerations

- The `client_secret` parameter is marked as sensitive and will not appear in logs
- OIDC tokens are managed internally and automatically refreshed
- All API communications use HTTPS
- Credentials are not stored persistently

## Contributing

1. Make changes to the code
2. Run tests: `go test ./...`
3. Build the provider: `go build .`
4. Test your changes with a real Terraform configuration

## License

This provider is developed for HiiRetail and follows internal licensing guidelines.