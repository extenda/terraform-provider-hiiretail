# HiiRetail Terraform Provider

The HiiRetail Terraform provider enables management of HiiRetail platform resources through Infrastructure as Code. This unified provider supports multiple HiiRetail APIs, starting with Identity and Access Management (IAM) operations including user groups, custom roles, resources, and role bindings. Future versions will extend support to additional APIs like OCMS (OAuth Client Management Service).

## 📚 Documentation

**Complete documentation is available in the [`docs/`](docs/) directory:**

- **[Provider Overview](docs/index.md)** - Configuration and authentication
- **[Getting Started Guide](docs/guides/getting-started.md)** - Quick setup tutorial
- **[Authentication Guide](docs/guides/authentication.md)** - Detailed auth setup
- **Resource Documentation:**
  - [Custom Roles](docs/resources/iam_custom_role.md)
  - [Groups](docs/resources/iam_group.md)
  - [Resources](docs/resources/iam_resource.md) 
  - [Role Bindings](docs/resources/iam_role_binding.md)
- **[Examples](examples/)** - Working Terraform configurations

## 🚀 Features

- **Registry-Ready**: Terraform Registry compliant documentation and structure
- **OAuth2 Authentication**: Secure client credentials flow with automatic token management
- **Comprehensive IAM Management**: Groups, custom roles, resources, and role bindings
- **Auto-Generated Documentation**: Schema docs generated with terraform-plugin-docs
- **Working Examples**: Validated examples for all resources
- **Error Handling**: Clear error messages and validation

## 📋 Requirements

- Terraform >= 1.0
- Go >= 1.21 (for development)

## Configuration

### Provider Configuration

```hcl
terraform {
  required_providers {
    hiiretail = {
      source = "extenda/hiiretail"
    }
  }
}

provider "hiiretail" {
  tenant_id     = "your-tenant-id"
  client_id     = "your-oidc-client-id"
  client_secret = "your-oidc-client-secret"
  base_url      = "https://custom-api.example.com" # Optional, defaults to https://iam-api.retailsvc.com
}
```

### Configuration Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tenant_id` | string | Yes | Tenant ID to use for all IAM API requests |
| `client_id` | string | Yes | OIDC client ID for IAM API authentication |
| `client_secret` | string | Yes | OIDC client secret for IAM API authentication (sensitive) |
| `base_url` | string | No | Base URL of the IAM API (defaults to https://iam-api.retailsvc-test.com) |

> **📖 For detailed authentication information including test credentials, see the [Authentication Guide](docs/guides/authentication.md)**

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

## Resources

### hiiretail_iam_group

Manages IAM groups within a tenant.

#### Example Usage

```hcl
# Basic group creation
resource "hiiretail_iam_group" "developers" {
  name = "developers"
}

# Group with description
resource "hiiretail_iam_group" "admin_group" {
  name        = "administrators"
  description = "Administrative users with full system access"
}

# Group with explicit tenant
resource "hiiretail_iam_group" "tenant_specific" {
  name        = "tenant-users"
  description = "Users specific to this tenant"
  tenant_id   = "custom-tenant-id"
}
```

#### Arguments Reference

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | Yes | The name of the group (max 255 characters) |
| `description` | string | No | Description of the group (max 255 characters) |
| `tenant_id` | string | No | Tenant ID for the group (defaults to provider tenant_id) |

#### Attributes Reference

In addition to all arguments above, the following attributes are exported:

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | string | The unique identifier of the group |
| `status` | string | The current status of the group |

#### Import

Groups can be imported using their ID:

```bash
terraform import hiiretail_iam_group.example group-12345
```

#### Examples

**Complete group management:**
```hcl
# Create multiple related groups
resource "hiiretail_iam_group" "developers" {
  name        = "developers"
  description = "Software development team"
}

resource "hiiretail_iam_group" "qa_team" {
  name        = "qa-engineers"
  description = "Quality assurance engineers"
}

resource "hiiretail_iam_group" "devops" {
  name        = "devops-engineers"
  description = "DevOps and infrastructure team"
}

# Output group information
output "developer_group_id" {
  value = hiiretail_iam_group.developers.id
}

output "all_groups" {
  value = {
    developers = hiiretail_iam_group.developers
    qa         = hiiretail_iam_group.qa_team
    devops     = hiiretail_iam_group.devops
  }
}
```

**Multi-tenant setup:**
```hcl
# Groups for different tenants
resource "hiiretail_iam_group" "tenant_a_users" {
  name        = "users"
  description = "Standard users for Tenant A"
  tenant_id   = "tenant-a"
}

resource "hiiretail_iam_group" "tenant_b_users" {
  name        = "users"
  description = "Standard users for Tenant B"  
  tenant_id   = "tenant-b"
}
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