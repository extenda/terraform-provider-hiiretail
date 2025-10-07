package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam/datasources"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam/resources"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/resource_iam_role_binding"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/auth"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/validators"
)

// APIClient represents the configuration for making API calls
// This is used by resources that need direct HTTP access
type APIClient struct {
	BaseURL    string
	TenantID   string
	HTTPClient *http.Client
}

// Ensure HiiRetailProvider satisfies various provider interfaces.
var _ provider.Provider = &HiiRetailProvider{}

// HiiRetailProvider defines the provider implementation.
type HiiRetailProvider struct {
	version string
}

// HiiRetailProviderModel describes the provider data model.
type HiiRetailProviderModel struct {
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	TenantID       types.String `tfsdk:"tenant_id"`
	BaseURL        types.String `tfsdk:"base_url"`
	IAMEndpoint    types.String `tfsdk:"iam_endpoint"`
	CCCEndpoint    types.String `tfsdk:"ccc_endpoint"`
	TokenURL       types.String `tfsdk:"token_url"`
	Scopes         types.Set    `tfsdk:"scopes"`
	TimeoutSeconds types.Int64  `tfsdk:"timeout_seconds"`
	MaxRetries     types.Int64  `tfsdk:"max_retries"`
}

func (p *HiiRetailProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hiiretail"
	resp.Version = p.version
}

func (p *HiiRetailProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The HiiRetail provider enables management of HiiRetail resources through Terraform.",
		MarkdownDescription: "The HiiRetail provider enables management of HiiRetail resources through Terraform. " +
			"It supports multiple services including IAM (Identity and Access Management) and CCC (Customer Care Center).",

		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description:         "OAuth2 client ID for authentication. Can also be set via HIIRETAIL_CLIENT_ID environment variable.",
				MarkdownDescription: "OAuth2 client ID for authentication. Can also be set via `HIIRETAIL_CLIENT_ID` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"client_secret": schema.StringAttribute{
				Description:         "OAuth2 client secret for authentication. Can also be set via HIIRETAIL_CLIENT_SECRET environment variable.",
				MarkdownDescription: "OAuth2 client secret for authentication. Can also be set via `HIIRETAIL_CLIENT_SECRET` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"tenant_id": schema.StringAttribute{
				Description:         "Tenant ID for resources. Can also be set via HIIRETAIL_TENANT_ID environment variable.",
				MarkdownDescription: "Tenant ID for resources. Can also be set via `HIIRETAIL_TENANT_ID` environment variable.",
				Optional:            true,
			},
			"base_url": schema.StringAttribute{
				Description:         "Base URL for the HiiRetail APIs. Can also be set via HIIRETAIL_BASE_URL environment variable.",
				MarkdownDescription: "Base URL for the HiiRetail APIs. Can also be set via `HIIRETAIL_BASE_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					validators.StringIsURL(),
				},
			},
			"iam_endpoint": schema.StringAttribute{
				Description:         "IAM service endpoint path. Defaults to '/iam/v1'.",
				MarkdownDescription: "IAM service endpoint path. Defaults to `/iam/v1`.",
				Optional:            true,
			},
			"ccc_endpoint": schema.StringAttribute{
				Description:         "CCC service endpoint path. Defaults to '/ccc/v1'.",
				MarkdownDescription: "CCC service endpoint path. Defaults to `/ccc/v1`.",
				Optional:            true,
			},
			"token_url": schema.StringAttribute{
				Description:         "OAuth2 token URL. Can also be set via HIIRETAIL_TOKEN_URL environment variable.",
				MarkdownDescription: "OAuth2 token URL. Can also be set via `HIIRETAIL_TOKEN_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					validators.StringIsURL(),
				},
			},
			"scopes": schema.SetAttribute{
				ElementType:         types.StringType,
				Description:         "OAuth2 scopes to request. Defaults to ['iam:read', 'iam:write'].",
				MarkdownDescription: "OAuth2 scopes to request. Defaults to `['iam:read', 'iam:write']`.",
				Optional:            true,
			},
			"timeout_seconds": schema.Int64Attribute{
				Description:         "Request timeout in seconds. Defaults to 30.",
				MarkdownDescription: "Request timeout in seconds. Defaults to 30.",
				Optional:            true,
			},
			"max_retries": schema.Int64Attribute{
				Description:         "Maximum number of retries for failed requests. Defaults to 3.",
				MarkdownDescription: "Maximum number of retries for failed requests. Defaults to 3.",
				Optional:            true,
			},
		},
	}
}

func (p *HiiRetailProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HiiRetailProviderModel

	// Build verification marker - this proves the binary is active
	// This should be updated by the build script with a unique ID
	fmt.Printf("[BUILD_VERIFICATION] HiiRetail Provider binary is active - build: 371FFBEA\n")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build OAuth2 configuration from provider data and environment variables
	authConfig, diags := buildAuthConfig(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build client configuration with hardcoded URLs and defaults
	clientConfig := &client.Config{
		BaseURL:      "https://iam-api.retailsvc.com", // Hardcoded IAM API URL base
		IAMEndpoint:  "/api/v1",                       // Most resources use V1 API - role bindings will bypass this
		CCCEndpoint:  "/ccc/v1",                       // Default CCC endpoint
		Timeout:      30 * time.Second,                // Default timeout
		MaxRetries:   3,                               // Default retries
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
	}

	// Override with user-provided values if present
	if !data.CCCEndpoint.IsNull() && !data.CCCEndpoint.IsUnknown() {
		clientConfig.CCCEndpoint = data.CCCEndpoint.ValueString()
	}
	if !data.TimeoutSeconds.IsNull() && !data.TimeoutSeconds.IsUnknown() {
		clientConfig.Timeout = time.Duration(data.TimeoutSeconds.ValueInt64()) * time.Second
	}
	if !data.MaxRetries.IsNull() && !data.MaxRetries.IsUnknown() {
		clientConfig.MaxRetries = int(data.MaxRetries.ValueInt64())
	}

	// Convert AuthClientConfig to auth.Config with hardcoded URLs
	authConfigV2 := &auth.Config{
		ClientID:         authConfig.ClientID,
		ClientSecret:     authConfig.ClientSecret,
		TenantID:         authConfig.TenantID,
		AuthURL:          authConfig.TokenURL,             // Already set to hardcoded auth URL
		APIURL:           "https://iam-api.retailsvc.com", // Hardcoded IAM API URL base (auth client handles path separately)
		Scopes:           authConfig.Scopes,
		Timeout:          authConfig.Timeout,
		MaxRetries:       authConfig.MaxRetries,
		DisableDiscovery: authConfig.DisableDiscovery,
	}

	// Create unified API client
	apiClient, err := client.New(authConfigV2, clientConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Setup Failed",
			fmt.Sprintf("Failed to initialize API client: %s", err.Error()),
		)
		return
	}

	// Make the client available to resources and data sources
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *HiiRetailProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// IAM resources
		resources.NewGroupResource,
		resources.NewCustomRoleResource,
		resource_iam_role_binding.NewSimpleIamRoleBindingResource, // Use simple 1:1 role binding resource
		resources.NewResourceResource,
	}
}

func (p *HiiRetailProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// IAM data sources
		datasources.NewGroupsDataSource,
		datasources.NewRolesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HiiRetailProvider{
			version: version,
		}
	}
}

// buildAuthConfig creates an AuthClientConfig from provider configuration and environment variables
// Precedence order: 1. terraform.tfvars 2. TF_VAR_* env vars 3. HIIRETAIL_* env vars 4. error
func buildAuthConfig(ctx context.Context, data *HiiRetailProviderModel) (*auth.AuthClientConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := &auth.AuthClientConfig{}

	// Get tenant ID with precedence: terraform.tfvars → TF_VAR_* → HIIRETAIL_* → error
	if !data.TenantID.IsNull() && !data.TenantID.IsUnknown() {
		config.TenantID = data.TenantID.ValueString()
	} else if tfVarTenantID := os.Getenv("TF_VAR_tenant_id"); tfVarTenantID != "" {
		config.TenantID = tfVarTenantID
	} else if hiiRetailTenantID := os.Getenv("HIIRETAIL_TENANT_ID"); hiiRetailTenantID != "" {
		config.TenantID = hiiRetailTenantID
	} else {
		diags.AddError(
			"Missing Tenant ID",
			"Tenant ID must be configured via terraform.tfvars, TF_VAR_tenant_id environment variable, or HIIRETAIL_TENANT_ID environment variable",
		)
	}

	// Get client ID with precedence: terraform.tfvars → TF_VAR_* → HIIRETAIL_* → error
	if !data.ClientID.IsNull() && !data.ClientID.IsUnknown() {
		config.ClientID = data.ClientID.ValueString()
	} else if tfVarClientID := os.Getenv("TF_VAR_client_id"); tfVarClientID != "" {
		config.ClientID = tfVarClientID
	} else if hiiRetailClientID := os.Getenv("HIIRETAIL_CLIENT_ID"); hiiRetailClientID != "" {
		config.ClientID = hiiRetailClientID
	} else {
		diags.AddError(
			"Missing Client ID",
			"Client ID must be configured via terraform.tfvars, TF_VAR_client_id environment variable, or HIIRETAIL_CLIENT_ID environment variable",
		)
	}

	// Get client secret with precedence: terraform.tfvars → TF_VAR_* → HIIRETAIL_* → error
	if !data.ClientSecret.IsNull() && !data.ClientSecret.IsUnknown() {
		config.ClientSecret = data.ClientSecret.ValueString()
	} else if tfVarClientSecret := os.Getenv("TF_VAR_client_secret"); tfVarClientSecret != "" {
		config.ClientSecret = tfVarClientSecret
	} else if hiiRetailClientSecret := os.Getenv("HIIRETAIL_CLIENT_SECRET"); hiiRetailClientSecret != "" {
		config.ClientSecret = hiiRetailClientSecret
	} else {
		diags.AddError(
			"Missing Client Secret",
			"Client Secret must be configured via terraform.tfvars, TF_VAR_client_secret environment variable, or HIIRETAIL_CLIENT_SECRET environment variable",
		)
	}

	// Set hardcoded auth URL for HiiRetail
	config.TokenURL = "https://auth.retailsvc.com/oauth2/token"

	// Get scopes with precedence: terraform.tfvars → TF_VAR_* → HIIRETAIL_* → default
	if !data.Scopes.IsNull() && !data.Scopes.IsUnknown() {
		scopes := make([]string, 0, len(data.Scopes.Elements()))
		diags.Append(data.Scopes.ElementsAs(ctx, &scopes, false)...)
		config.Scopes = scopes
	} else if tfVarScopes := os.Getenv("TF_VAR_scopes"); tfVarScopes != "" {
		config.Scopes = strings.Split(tfVarScopes, ",")
	} else if hiiRetailScopes := os.Getenv("HIIRETAIL_SCOPES"); hiiRetailScopes != "" {
		config.Scopes = strings.Split(hiiRetailScopes, ",")
	} else {
		config.Scopes = []string{
			"IAM:create:roles", "IAM:read:roles", "IAM:update:roles", "IAM:delete:roles",
			"IAM:create:groups", "IAM:read:groups", "IAM:update:groups", "IAM:delete:groups",
			"IAM:create:role_bindings", "IAM:read:role_bindings", "IAM:update:role_bindings", "IAM:delete:role_bindings",
			"iam.group.list-roles", // Specific permission for V2 GET /api/v2/tenants/{tenantId}/groups/{id}/roles
		} // Default scopes with granular IAM permissions
	}

	// Set timeout with precedence: terraform.tfvars → TF_VAR_* → HIIRETAIL_* → default
	if !data.TimeoutSeconds.IsNull() && !data.TimeoutSeconds.IsUnknown() {
		config.Timeout = time.Duration(data.TimeoutSeconds.ValueInt64()) * time.Second
	} else if tfVarTimeout := os.Getenv("TF_VAR_timeout_seconds"); tfVarTimeout != "" {
		if timeoutVal, err := strconv.ParseInt(tfVarTimeout, 10, 64); err == nil {
			config.Timeout = time.Duration(timeoutVal) * time.Second
		} else {
			config.Timeout = 30 * time.Second // Default on parse error
		}
	} else if hiiRetailTimeout := os.Getenv("HIIRETAIL_TIMEOUT_SECONDS"); hiiRetailTimeout != "" {
		if timeoutVal, err := strconv.ParseInt(hiiRetailTimeout, 10, 64); err == nil {
			config.Timeout = time.Duration(timeoutVal) * time.Second
		} else {
			config.Timeout = 30 * time.Second // Default on parse error
		}
	} else {
		config.Timeout = 30 * time.Second // Default timeout
	}

	// Set max retries with precedence: terraform.tfvars → TF_VAR_* → HIIRETAIL_* → default
	if !data.MaxRetries.IsNull() && !data.MaxRetries.IsUnknown() {
		config.MaxRetries = int(data.MaxRetries.ValueInt64())
	} else if tfVarRetries := os.Getenv("TF_VAR_max_retries"); tfVarRetries != "" {
		if retriesVal, err := strconv.Atoi(tfVarRetries); err == nil {
			config.MaxRetries = retriesVal
		} else {
			config.MaxRetries = 3 // Default on parse error
		}
	} else if hiiRetailRetries := os.Getenv("HIIRETAIL_MAX_RETRIES"); hiiRetailRetries != "" {
		if retriesVal, err := strconv.Atoi(hiiRetailRetries); err == nil {
			config.MaxRetries = retriesVal
		} else {
			config.MaxRetries = 3 // Default on parse error
		}
	} else {
		config.MaxRetries = 3 // Default max retries
	}

	// Set default discovery and headers (not configurable in current model)
	config.DisableDiscovery = false
	config.CustomHeaders = make(map[string]string)

	return config, diags
}

// resolveBaseURL determines the appropriate base URL for API calls
func resolveBaseURL(config *auth.AuthClientConfig) string {
	if config.BaseURL != "" {
		// Convert OAuth2 discovery URL to API base URL
		if strings.Contains(config.BaseURL, "auth.retailsvc.com") {
			return "https://iam-api.retailsvc.com"
		}
		return config.BaseURL
	}

	// Default API base URL
	return "https://iam-api.retailsvc.com"
}
