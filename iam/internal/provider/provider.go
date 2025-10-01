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
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/auth"
	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/provider_hiiretail_iam"
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
type HiiRetailIamProviderModel = provider_hiiretail_iam.HiiretailIamModel

func (p *HiiRetailIamProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hiiretail-iam"
	resp.Version = p.version
}

func (p *HiiRetailIamProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = provider_hiiretail_iam.HiiretailIamProviderSchema(ctx)
}

func (p *HiiRetailIamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data HiiRetailIamProviderModel

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

	// Validate configuration using the comprehensive validation system
	validationResult := auth.ValidateAuthConfig(authConfig, auth.DefaultValidationRules())
	if !validationResult.Valid {
		for _, err := range validationResult.Errors {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Configuration Error in %s", err.Field),
				fmt.Sprintf("%s. %s", err.Message, err.Suggestion),
			)
		}
		return
	}

	// Add warnings for configuration issues
	for _, warning := range validationResult.Warnings {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Configuration Warning for %s", warning.Field),
			fmt.Sprintf("%s. %s", warning.Message, warning.Suggestion),
		)
	}

	// Create OAuth2 authentication client
	authClient, err := auth.NewAuthClient(authConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"OAuth2 Authentication Setup Failed",
			fmt.Sprintf("Failed to initialize OAuth2 authentication client: %s", err.Error()),
		)
		return
	}

	// Validate authentication by attempting to acquire a token
	if _, err := authClient.GetToken(ctx); err != nil {
		authClient.Close() // Clean up resources
		resp.Diagnostics.AddError(
			"OAuth2 Authentication Failed",
			fmt.Sprintf("Failed to authenticate with OAuth2 provider: %s", err.Error()),
		)
		return
	}

	// Create authenticated HTTP client
	httpClient, err := authClient.HTTPClientWithRetry(ctx)
	if err != nil {
		authClient.Close() // Clean up resources
		resp.Diagnostics.AddError(
			"HTTP Client Setup Failed",
			fmt.Sprintf("Failed to create authenticated HTTP client: %s", err.Error()),
		)
		return
	}

	// Create API client configuration
	apiClient := &APIClient{
		BaseURL:    resolveBaseURL(authConfig),
		TenantID:   authConfig.TenantID,
		HTTPClient: httpClient,
		AuthClient: authClient,
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
	AuthClient *auth.AuthClient
}

// buildAuthConfig creates an AuthClientConfig from provider configuration and environment variables
func buildAuthConfig(ctx context.Context, data *HiiRetailIamProviderModel) (*auth.AuthClientConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	config := &auth.AuthClientConfig{}

	// Get tenant ID from config or environment
	if !data.TenantId.IsNull() && !data.TenantId.IsUnknown() {
		config.TenantID = data.TenantId.ValueString()
	} else {
		config.TenantID = os.Getenv("HIIRETAIL_TENANT_ID")
	}

	// Get client ID from config or environment
	if !data.ClientId.IsNull() && !data.ClientId.IsUnknown() {
		config.ClientID = data.ClientId.ValueString()
	} else {
		config.ClientID = os.Getenv("HIIRETAIL_CLIENT_ID")
	}

	// Get client secret from config or environment
	if !data.ClientSecret.IsNull() && !data.ClientSecret.IsUnknown() {
		config.ClientSecret = data.ClientSecret.ValueString()
	} else {
		config.ClientSecret = os.Getenv("HIIRETAIL_CLIENT_SECRET")
	}

	// Get base URL from config or environment
	if !data.BaseUrl.IsNull() && !data.BaseUrl.IsUnknown() {
		config.BaseURL = data.BaseUrl.ValueString()
	} else {
		baseURL := os.Getenv("HIIRETAIL_BASE_URL")
		if baseURL == "" {
			baseURL = "https://auth.retailsvc.com" // Default OAuth2 discovery endpoint
		}
		config.BaseURL = baseURL
	}

	// Get token URL from config or environment (optional)
	if !data.TokenUrl.IsNull() && !data.TokenUrl.IsUnknown() {
		config.TokenURL = data.TokenUrl.ValueString()
	} else {
		config.TokenURL = os.Getenv("HIIRETAIL_TOKEN_URL")
	}

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

	// Get timeout from config or environment
	if !data.TimeoutSeconds.IsNull() && !data.TimeoutSeconds.IsUnknown() {
		config.Timeout = time.Duration(data.TimeoutSeconds.ValueInt64()) * time.Second
	} else {
		timeoutEnv := os.Getenv("HIIRETAIL_TIMEOUT_SECONDS")
		if timeoutEnv != "" {
			if timeoutSeconds, err := strconv.Atoi(timeoutEnv); err == nil {
				config.Timeout = time.Duration(timeoutSeconds) * time.Second
			}
		}
		if config.Timeout == 0 {
			config.Timeout = 30 * time.Second // Default timeout
		}
	}

	// Get max retries from config or environment
	if !data.MaxRetries.IsNull() && !data.MaxRetries.IsUnknown() {
		config.MaxRetries = int(data.MaxRetries.ValueInt64())
	} else {
		retriesEnv := os.Getenv("HIIRETAIL_MAX_RETRIES")
		if retriesEnv != "" {
			if maxRetries, err := strconv.Atoi(retriesEnv); err == nil {
				config.MaxRetries = maxRetries
			}
		}
		if config.MaxRetries == 0 {
			config.MaxRetries = 3 // Default max retries
		}
	}

	// Get disable discovery from config or environment
	if !data.DisableDiscovery.IsNull() && !data.DisableDiscovery.IsUnknown() {
		config.DisableDiscovery = data.DisableDiscovery.ValueBool()
	} else {
		disableDiscoveryEnv := os.Getenv("HIIRETAIL_DISABLE_DISCOVERY")
		config.DisableDiscovery = strings.ToLower(disableDiscoveryEnv) == "true"
	}

	// Get custom headers from config
	if !data.CustomHeaders.IsNull() && !data.CustomHeaders.IsUnknown() {
		headers := make(map[string]string)
		diags.Append(data.CustomHeaders.ElementsAs(ctx, &headers, false)...)
		config.CustomHeaders = headers
	}

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
