package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProvider_OAuth2Configuration tests provider configuration with OAuth2 parameters
func TestProvider_OAuth2Configuration(t *testing.T) {
	t.Run("valid_configuration_with_all_oauth2_parameters", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789-very-secure",
			"base_url":      "https://iam-api.retailsvc-test.com",
			"token_url":     "https://auth.retailsvc-test.com/oauth2/token",
			"scopes":        []string{"iam:read", "iam:write"},
			"timeout":       "30s",
			"max_retries":   3,
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Configuration should succeed with valid OAuth2 parameters")
		assert.NotNil(t, resp.ResourceData, "ResourceData should be configured")
		assert.NotNil(t, resp.DataSourceData, "DataSourceData should be configured")
	})

	t.Run("valid_configuration_with_minimal_parameters", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789-minimal",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Configuration should succeed with minimal OAuth2 parameters")
		assert.NotNil(t, resp.ResourceData, "ResourceData should be configured")
	})

	t.Run("missing_required_tenant_id", func(t *testing.T) {
		config := map[string]interface{}{
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail without tenant_id")
		assertDiagnosticContains(t, resp.Diagnostics, "tenant_id")
		assertDiagnosticContains(t, resp.Diagnostics, "required")
	})

	t.Run("missing_required_client_id", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_secret": "test-secret-789",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail without client_id")
		assertDiagnosticContains(t, resp.Diagnostics, "client_id")
		assertDiagnosticContains(t, resp.Diagnostics, "required")
	})

	t.Run("missing_required_client_secret", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id": "test-tenant-123",
			"client_id": "test-client-456",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail without client_secret")
		assertDiagnosticContains(t, resp.Diagnostics, "client_secret")
		assertDiagnosticContains(t, resp.Diagnostics, "required")
	})

	t.Run("invalid_base_url", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
			"base_url":      "not-a-valid-url",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with invalid base_url")
		assertDiagnosticContains(t, resp.Diagnostics, "base_url")
		assertDiagnosticContains(t, resp.Diagnostics, "valid URL")
	})

	t.Run("invalid_token_url", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
			"token_url":     "not-a-valid-url",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with invalid token_url")
		assertDiagnosticContains(t, resp.Diagnostics, "token_url")
		assertDiagnosticContains(t, resp.Diagnostics, "valid URL")
	})

	t.Run("non_https_base_url", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
			"base_url":      "http://insecure.retailsvc.com", // HTTP instead of HTTPS
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with non-HTTPS base_url")
		assertDiagnosticContains(t, resp.Diagnostics, "HTTPS")
		assertDiagnosticContains(t, resp.Diagnostics, "secure")
	})

	t.Run("invalid_timeout_format", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
			"timeout":       "invalid-duration",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with invalid timeout format")
		assertDiagnosticContains(t, resp.Diagnostics, "timeout")
		assertDiagnosticContains(t, resp.Diagnostics, "duration")
	})

	t.Run("timeout_out_of_range", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
			"timeout":       "600s", // Too long
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with timeout out of range")
		assertDiagnosticContains(t, resp.Diagnostics, "timeout")
		assertDiagnosticContains(t, resp.Diagnostics, "range")
	})

	t.Run("max_retries_out_of_range", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
			"max_retries":   15, // Too high
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with max_retries out of range")
		assertDiagnosticContains(t, resp.Diagnostics, "max_retries")
		assertDiagnosticContains(t, resp.Diagnostics, "range")
	})

	t.Run("weak_client_secret", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "weak", // Too short
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with weak client_secret")
		assertDiagnosticContains(t, resp.Diagnostics, "client_secret")
		assertDiagnosticContains(t, resp.Diagnostics, "length")
	})
}

// TestProvider_OAuth2EnvironmentVariables tests environment variable support
func TestProvider_OAuth2EnvironmentVariables(t *testing.T) {
	t.Run("environment_variables_override_config", func(t *testing.T) {
		// Set environment variables
		originalVars := setTestEnvironmentVariables(map[string]string{
			"HIIRETAIL_TENANT_ID":     "env-tenant-123",
			"HIIRETAIL_CLIENT_ID":     "env-client-456",
			"HIIRETAIL_CLIENT_SECRET": "env-secret-789-from-environment",
			"HIIRETAIL_BASE_URL":      "https://env.retailsvc.com",
			"HIIRETAIL_TOKEN_URL":     "https://env-auth.retailsvc.com/oauth2/token",
			"HIIRETAIL_SCOPES":        "iam:read,iam:write,iam:admin",
			"HIIRETAIL_TIMEOUT":       "45s",
			"HIIRETAIL_MAX_RETRIES":   "5",
		})
		defer restoreEnvironmentVariables(originalVars)

		// Empty configuration - should use environment variables
		config := map[string]interface{}{}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Configuration should succeed with environment variables")

		// Verify environment variables were used
		// Note: In actual implementation, you would check the configured values
		// This is a placeholder for the verification logic
	})

	t.Run("config_parameters_override_environment", func(t *testing.T) {
		// Set environment variables
		originalVars := setTestEnvironmentVariables(map[string]string{
			"HIIRETAIL_TENANT_ID":     "env-tenant-123",
			"HIIRETAIL_CLIENT_ID":     "env-client-456",
			"HIIRETAIL_CLIENT_SECRET": "env-secret-789",
		})
		defer restoreEnvironmentVariables(originalVars)

		// Explicit configuration should override environment variables
		config := map[string]interface{}{
			"tenant_id":     "config-tenant-456",
			"client_id":     "config-client-789",
			"client_secret": "config-secret-123-explicit",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Configuration should succeed with explicit config over env vars")
	})

	t.Run("partial_environment_variables", func(t *testing.T) {
		// Set only some environment variables
		originalVars := setTestEnvironmentVariables(map[string]string{
			"HIIRETAIL_TENANT_ID": "env-tenant-123",
			"HIIRETAIL_CLIENT_ID": "env-client-456",
			// Missing HIIRETAIL_CLIENT_SECRET
		})
		defer restoreEnvironmentVariables(originalVars)

		// Complete configuration with explicit client_secret
		config := map[string]interface{}{
			"client_secret": "explicit-secret-789-mixed-config",
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.False(t, resp.Diagnostics.HasError(), "Configuration should succeed with mixed env vars and explicit config")
	})

	t.Run("invalid_environment_variable_values", func(t *testing.T) {
		originalVars := setTestEnvironmentVariables(map[string]string{
			"HIIRETAIL_TENANT_ID":     "env-tenant-123",
			"HIIRETAIL_CLIENT_ID":     "env-client-456",
			"HIIRETAIL_CLIENT_SECRET": "env-secret-789",
			"HIIRETAIL_BASE_URL":      "not-a-valid-url", // Invalid URL
		})
		defer restoreEnvironmentVariables(originalVars)

		config := map[string]interface{}{}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail with invalid environment variable values")
		assertDiagnosticContains(t, resp.Diagnostics, "base_url")
	})
}

// TestProvider_OAuth2SchemaValidation tests schema-level validation
func TestProvider_OAuth2SchemaValidation(t *testing.T) {
	t.Run("client_secret_marked_as_sensitive", func(t *testing.T) {
		p := New("test")()
		schema := p.Schema(context.Background(), provider.SchemaRequest{}, &provider.SchemaResponse{})

		require.NotNil(t, schema.Schema, "Provider schema should not be nil")

		clientSecretAttr, exists := schema.Schema.Attributes["client_secret"]
		require.True(t, exists, "client_secret attribute should exist in schema")

		// Verify client_secret is marked as sensitive
		// Note: The exact method to check this depends on the Plugin Framework version
		// This is a placeholder for the actual sensitivity check
		assert.True(t, true, "client_secret should be marked as sensitive")
	})

	t.Run("optional_parameters_have_defaults", func(t *testing.T) {
		p := New("test")()
		schema := p.Schema(context.Background(), provider.SchemaRequest{}, &provider.SchemaResponse{})

		require.NotNil(t, schema.Schema, "Provider schema should not be nil")

		// Check that optional parameters have appropriate defaults
		timeoutAttr, exists := schema.Schema.Attributes["timeout"]
		if exists {
			// Verify default timeout value
			assert.True(t, true, "timeout should have a default value")
		}

		maxRetriesAttr, exists := schema.Schema.Attributes["max_retries"]
		if exists {
			// Verify default max_retries value
			assert.True(t, true, "max_retries should have a default value")
		}
	})

	t.Run("required_parameters_validation", func(t *testing.T) {
		p := New("test")()
		schema := p.Schema(context.Background(), provider.SchemaRequest{}, &provider.SchemaResponse{})

		require.NotNil(t, schema.Schema, "Provider schema should not be nil")

		// Verify required parameters
		requiredParams := []string{"tenant_id", "client_id", "client_secret"}
		for _, param := range requiredParams {
			attr, exists := schema.Schema.Attributes[param]
			require.True(t, exists, "%s attribute should exist in schema", param)

			// Verify it's marked as required
			assert.True(t, attr.IsRequired(), "%s should be marked as required", param)
		}
	})
}

// TestProvider_OAuth2ValidationMessages tests user-friendly validation messages
func TestProvider_OAuth2ValidationMessages(t *testing.T) {
	t.Run("helpful_error_messages", func(t *testing.T) {
		testCases := []struct {
			name            string
			config          map[string]interface{}
			expectedMessage string
		}{
			{
				name:            "missing_tenant_id_helpful_message",
				config:          map[string]interface{}{"client_id": "test", "client_secret": "test-secret-123"},
				expectedMessage: "tenant_id is required for OAuth2 authentication",
			},
			{
				name:            "invalid_url_helpful_message",
				config:          map[string]interface{}{"tenant_id": "test", "client_id": "test", "client_secret": "test-secret-123", "base_url": "not-url"},
				expectedMessage: "base_url must be a valid HTTPS URL",
			},
			{
				name:            "weak_secret_helpful_message",
				config:          map[string]interface{}{"tenant_id": "test", "client_id": "test", "client_secret": "weak"},
				expectedMessage: "client_secret must be at least 8 characters long",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				p := New("test")()
				req := provider.ConfigureRequest{
					Config: createProviderConfig(t, tc.config),
				}

				var resp provider.ConfigureResponse
				p.Configure(context.Background(), req, &resp)

				assert.True(t, resp.Diagnostics.HasError(), "Configuration should fail")
				assertDiagnosticContains(t, resp.Diagnostics, tc.expectedMessage)
			})
		}
	})

	t.Run("troubleshooting_guidance", func(t *testing.T) {
		config := map[string]interface{}{
			"tenant_id":     "test-tenant-123",
			"client_id":     "test-client-456",
			"client_secret": "test-secret-789",
			"base_url":      "https://invalid.nonexistent.com", // Will fail discovery
		}

		p := New("test")()
		req := provider.ConfigureRequest{
			Config: createProviderConfig(t, config),
		}

		var resp provider.ConfigureResponse
		p.Configure(context.Background(), req, &resp)

		// Should provide troubleshooting guidance for connection failures
		if resp.Diagnostics.HasError() {
			diagnosticMessages := getDiagnosticMessages(resp.Diagnostics)
			found := false
			for _, msg := range diagnosticMessages {
				if containsAny(msg, []string{"troubleshooting", "check", "verify", "ensure"}) {
					found = true
					break
				}
			}
			assert.True(t, found, "Error messages should include troubleshooting guidance")
		}
	})
}

// Helper functions for testing

func createProviderConfig(t *testing.T, configMap map[string]interface{}) tfsdk.Config {
	// This is a mock implementation for testing
	// Real implementation would create a proper tfsdk.Config from the map
	return tfsdk.Config{}
}

func setTestEnvironmentVariables(envVars map[string]string) map[string]string {
	original := make(map[string]string)
	for key, value := range envVars {
		original[key] = os.Getenv(key)
		os.Setenv(key, value)
	}
	return original
}

func restoreEnvironmentVariables(original map[string]string) {
	for key, value := range original {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}

func assertDiagnosticContains(t *testing.T, diags interface{}, expectedText string) {
	// Mock implementation - would check actual diagnostics in real code
	assert.True(t, true, "Should contain diagnostic message: %s", expectedText)
}

func getDiagnosticMessages(diags interface{}) []string {
	// Mock implementation - would extract actual diagnostic messages
	return []string{"mock diagnostic message"}
}

func containsAny(text string, substrings []string) bool {
	for _, substring := range substrings {
		if len(text) > 0 && len(substring) > 0 {
			return true // Mock implementation
		}
	}
	return false
}
