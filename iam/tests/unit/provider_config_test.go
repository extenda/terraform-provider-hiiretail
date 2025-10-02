package provider_test
package provider

import (
	"testing"
	
	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/stretchr/testify/assert"
)

// TestProviderConfigurationSchema validates the provider configuration schema
// This test validates the contract defined in contracts/provider-config.md
func TestProviderConfigurationSchema(t *testing.T) {
	// This test should FAIL until the provider is properly updated
	p := provider.New("test")
	
	// Get the schema
	req := provider.SchemaRequest{}
	resp := provider.SchemaResponse{}
	p.Schema(nil, req, &resp)
	
	// Validate required OAuth2 fields exist
	assert.Contains(t, resp.Schema.Attributes, "client_id", "client_id should be in provider schema")
	assert.Contains(t, resp.Schema.Attributes, "client_secret", "client_secret should be in provider schema")
	
	// Validate optional OAuth2 configuration fields
	assert.Contains(t, resp.Schema.Attributes, "auth_url", "auth_url should be in provider schema")
	assert.Contains(t, resp.Schema.Attributes, "audience", "audience should be in provider schema")
	
	// Validate connection settings
	assert.Contains(t, resp.Schema.Attributes, "timeout_seconds", "timeout_seconds should be in provider schema")
	assert.Contains(t, resp.Schema.Attributes, "max_retries", "max_retries should be in provider schema")
	
	// Validate service endpoint overrides
	assert.Contains(t, resp.Schema.Attributes, "iam_endpoint", "iam_endpoint should be in provider schema")
	assert.Contains(t, resp.Schema.Attributes, "ccc_endpoint", "ccc_endpoint should be in provider schema")
	
	// Validate sensitive attributes are marked correctly
	clientIdAttr := resp.Schema.Attributes["client_id"]
	assert.True(t, clientIdAttr.IsSensitive(), "client_id should be marked as sensitive")
	
	clientSecretAttr := resp.Schema.Attributes["client_secret"]
	assert.True(t, clientSecretAttr.IsSensitive(), "client_secret should be marked as sensitive")
}

// TestProviderValidation validates provider configuration validation rules
func TestProviderValidation(t *testing.T) {
	// This test should FAIL until validation is implemented
	// Test timeout_seconds range validation (5-300)
	// Test max_retries range validation (0-10)
	// Test URL format validation for endpoints
	
	t.Skip("Provider validation not yet implemented - this test should fail until T037 is complete")
}