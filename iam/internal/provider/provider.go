package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/resource_iam_group"
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
	resp.TypeName = "hiiretail_iam"
	resp.Version = p.version
}

func (p *HiiRetailIamProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Hii Retail IAM API",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "Tenant ID to use for all IAM API requests",
				Required:    true,
			},
			"base_url": schema.StringAttribute{
				Description: "Base URL of the IAM API. Defaults to https://iam-api.retailsvc-test.com",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "OIDC client ID for IAM API authentication",
				Required:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "OIDC client secret for IAM API authentication",
				Required:    true,
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

	// Set default base URL if not provided
	baseUrl := "https://iam-api.retailsvc-test.com"
	if !data.BaseUrl.IsNull() && !data.BaseUrl.IsUnknown() {
		baseUrl = data.BaseUrl.ValueString()
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
	if data.TenantId.IsNull() || data.TenantId.IsUnknown() || data.TenantId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing tenant_id",
			"The tenant_id parameter is required",
		)
		return
	}

	if data.ClientId.IsNull() || data.ClientId.IsUnknown() || data.ClientId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing client_id",
			"The client_id parameter is required for OIDC authentication",
		)
		return
	}

	if data.ClientSecret.IsNull() || data.ClientSecret.IsUnknown() || data.ClientSecret.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing client_secret",
			"The client_secret parameter is required for OIDC authentication",
		)
		return
	}

	// Configure OIDC client credentials flow
	config := &clientcredentials.Config{
		ClientID:     data.ClientId.ValueString(),
		ClientSecret: data.ClientSecret.ValueString(),
		TokenURL:     fmt.Sprintf("%s/oauth/token", baseUrl),
	}

	// Test the OIDC authentication
	httpClient := config.Client(ctx)
	if httpClient == nil {
		resp.Diagnostics.AddError(
			"OIDC Configuration Error",
			"Failed to create OIDC client",
		)
		return
	}

	// Create API client configuration
	apiClient := &APIClient{
		BaseURL:    baseUrl,
		TenantID:   data.TenantId.ValueString(),
		HTTPClient: httpClient,
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *HiiRetailIamProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// T027: Register Group resource with provider
		resource_iam_group.NewIamGroupResource,
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
