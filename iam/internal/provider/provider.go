package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/resource_iam_custom_role"
	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/resource_iam_group"
	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/resource_iam_role_binding"
)

// Ensure HiiRetailIamProvider satisfies various provider interfaces.
var _ provider.Provider = &HiiRetailIamProvider{}

// HiiRetailIamProvider defines the provider implementation.
type HiiRetailIamProvider struct {
	version string
}

// HiiRetailIamProviderModel describes the provider data model.
type HiiRetailIamProviderModel struct {
	TenantId     types.String `tfsdk:"tenant_id"`
	BaseUrl      types.String `tfsdk:"base_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *HiiRetailIamProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hiiretail-iam"
	resp.Version = p.version
}

func (p *HiiRetailIamProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Hii Retail IAM API",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "Tenant ID to use for all IAM API requests. Can be set via HIIRETAIL_TENANT_ID environment variable.",
				Optional:    true,
			},
			"base_url": schema.StringAttribute{
				Description: "Base URL of the IAM API. Defaults to https://iam-api.retailsvc-test.com. Can be set via HIIRETAIL_BASE_URL environment variable.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "OIDC client ID for IAM API authentication. Can be set via HIIRETAIL_CLIENT_ID environment variable.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "OIDC client secret for IAM API authentication. Can be set via HIIRETAIL_CLIENT_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *HiiRetailIamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HiiRetailIamProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get configuration values from attributes or environment variables
	var tenantId, baseUrl, clientId, clientSecret string

	// Tenant ID
	if !data.TenantId.IsNull() && !data.TenantId.IsUnknown() {
		tenantId = data.TenantId.ValueString()
	} else {
		tenantId = os.Getenv("HIIRETAIL_TENANT_ID")
	}

	// Base URL
	if !data.BaseUrl.IsNull() && !data.BaseUrl.IsUnknown() {
		baseUrl = data.BaseUrl.ValueString()
	} else {
		baseUrl = os.Getenv("HIIRETAIL_BASE_URL")
		if baseUrl == "" {
			baseUrl = "https://iam-api.retailsvc-test.com"
		}
	}

	// Client ID
	if !data.ClientId.IsNull() && !data.ClientId.IsUnknown() {
		clientId = data.ClientId.ValueString()
	} else {
		clientId = os.Getenv("HIIRETAIL_CLIENT_ID")
	}

	// Client Secret
	if !data.ClientSecret.IsNull() && !data.ClientSecret.IsUnknown() {
		clientSecret = data.ClientSecret.ValueString()
	} else {
		clientSecret = os.Getenv("HIIRETAIL_CLIENT_SECRET")
	}

	// Validate base URL format
	parsedURL, err := url.Parse(baseUrl)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		resp.Diagnostics.AddError(
			"Invalid base_url",
			fmt.Sprintf("The provided base_url is not a valid URL: %s", baseUrl),
		)
		return
	}

	// Validate required fields
	if tenantId == "" {
		resp.Diagnostics.AddError(
			"Missing tenant_id",
			"The tenant_id parameter is required. Set it in the provider configuration or via HIIRETAIL_TENANT_ID environment variable.",
		)
		return
	}

	if clientId == "" {
		resp.Diagnostics.AddError(
			"Missing client_id",
			"The client_id parameter is required for OIDC authentication. Set it in the provider configuration or via HIIRETAIL_CLIENT_ID environment variable.",
		)
		return
	}

	if clientSecret == "" {
		resp.Diagnostics.AddError(
			"Missing client_secret",
			"The client_secret parameter is required for OIDC authentication. Set it in the provider configuration or via HIIRETAIL_CLIENT_SECRET environment variable.",
		)
		return
	}

	// Configure OIDC client credentials flow
	config := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("%s/oauth2/token", baseUrl),
	}

	// Create HTTP client with OAuth2 configuration and shorter timeout for test environments
	// Note: Token acquisition is lazy - happens on first API call, not during provider configuration
	baseHTTPClient := &http.Client{
		Timeout: 10 * time.Second, // Shorter timeout for tests
	}

	// Create OAuth2 context that won't be canceled when Terraform operations timeout
	// Use background context to avoid context cancellation issues during token acquisition
	oauthCtx := context.WithValue(context.Background(), oauth2.HTTPClient, baseHTTPClient)
	httpClient := config.Client(oauthCtx)

	// Create API client configuration
	apiClient := &APIClient{
		BaseURL:    baseUrl,
		TenantID:   tenantId,
		HTTPClient: httpClient,
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *HiiRetailIamProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// T027: Register Group resource with provider
		resource_iam_group.NewIamGroupResource,
		// T001: Register Custom Role resource with provider
		resource_iam_custom_role.NewIamCustomRoleResource,
		// T014: Register Role Binding resource with provider
		resource_iam_role_binding.NewIamRoleBindingResource,
	}
}

func (p *HiiRetailIamProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HiiRetailIamProvider{
			version: version,
		}
	}
}

// APIClient represents the configuration for making API calls
type APIClient struct {
	BaseURL    string
	TenantID   string
	HTTPClient *http.Client
}
