package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/iam/datasources"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/iam/resources"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/auth"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/client"
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/shared/validators"
)

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

	// Build client configuration with hardcoded URLs
	clientConfig := &client.Config{
		BaseURL:      "https://iam-api.retailsvc.com", // Hardcoded IAM API URL
		IAMEndpoint:  data.IAMEndpoint.ValueString(),
		CCCEndpoint:  data.CCCEndpoint.ValueString(),
		Timeout:      time.Duration(data.TimeoutSeconds.ValueInt64()) * time.Second,
		MaxRetries:   int(data.MaxRetries.ValueInt64()),
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 30 * time.Second,
	}

	// Convert AuthClientConfig to auth.Config with hardcoded URLs
	authConfigV2 := &auth.Config{
		ClientID:         authConfig.ClientID,
		ClientSecret:     authConfig.ClientSecret,
		TenantID:         authConfig.TenantID,
		AuthURL:          authConfig.TokenURL,             // Already set to hardcoded auth URL
		APIURL:           "https://iam-api.retailsvc.com", // Hardcoded IAM API URL
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
		resources.NewRoleBindingResource,
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
func buildAuthConfig(ctx context.Context, data *HiiRetailProviderModel) (*auth.AuthClientConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := &auth.AuthClientConfig{}

	// Get tenant ID from config or environment
	if !data.TenantID.IsNull() && !data.TenantID.IsUnknown() {
		config.TenantID = data.TenantID.ValueString()
	} else {
		config.TenantID = os.Getenv("HIIRETAIL_TENANT_ID")
	}

	// Get client ID from config or environment
	if !data.ClientID.IsNull() && !data.ClientID.IsUnknown() {
		config.ClientID = data.ClientID.ValueString()
	} else {
		config.ClientID = os.Getenv("HIIRETAIL_CLIENT_ID")
	}

	// Get client secret from config or environment
	if !data.ClientSecret.IsNull() && !data.ClientSecret.IsUnknown() {
		config.ClientSecret = data.ClientSecret.ValueString()
	} else {
		config.ClientSecret = os.Getenv("HIIRETAIL_CLIENT_SECRET")
	}

	// Set hardcoded auth URL for HiiRetail
	config.TokenURL = "https://auth.retailsvc.com/oauth2/token"

	// Get scopes from config or default
	if !data.Scopes.IsNull() && !data.Scopes.IsUnknown() {
		scopes := make([]string, 0, len(data.Scopes.Elements()))
		diags.Append(data.Scopes.ElementsAs(ctx, &scopes, false)...)
		config.Scopes = scopes
	} else {
		scopesEnv := os.Getenv("HIIRETAIL_SCOPES")
		if scopesEnv != "" {
			config.Scopes = strings.Split(scopesEnv, ",")
		} else {
			config.Scopes = []string{"iam:read", "iam:write"} // Default scopes
		}
	}

	// Set default timeout and retries (using provider-level settings from client config)
	if !data.TimeoutSeconds.IsNull() && !data.TimeoutSeconds.IsUnknown() {
		config.Timeout = time.Duration(data.TimeoutSeconds.ValueInt64()) * time.Second
	} else {
		config.Timeout = 30 * time.Second // Default timeout
	}

	if !data.MaxRetries.IsNull() && !data.MaxRetries.IsUnknown() {
		config.MaxRetries = int(data.MaxRetries.ValueInt64())
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
