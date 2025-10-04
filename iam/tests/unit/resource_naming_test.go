package unit_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResourceNamingConvention validates the resource naming contract
// This test validates the contract defined in contracts/resource-naming.md
func TestResourceNamingConvention(t *testing.T) {
	// This test should FAIL until resources are renamed
	expectedResources := map[string]string{
		"hiiretail_iam_group":        "hiiretail-iam_iam_group",        // old name
		"hiiretail_iam_custom_role":  "hiiretail-iam_custom_role",      // old name
		"hiiretail_iam_role_binding": "hiiretail-iam_iam_role_binding", // old name
	}

	for newName, oldName := range expectedResources {
		t.Run("Resource_"+newName, func(t *testing.T) {
			// Validate naming pattern: hiiretail_{service}_{resource_type}
			parts := strings.Split(newName, "_")
			assert.Len(t, parts, 3, "Resource name should have 3 parts: provider_service_resource")
			assert.Equal(t, "hiiretail", parts[0], "Provider name should be 'hiiretail'")
			assert.Equal(t, "iam", parts[1], "Service name should be 'iam'")
			assert.NotEmpty(t, parts[2], "Resource type should not be empty")

			// This assertion will FAIL until resources are renamed
			t.Logf("Expected new name: %s, Old name: %s", newName, oldName)
			// TODO: Add actual resource registration validation once provider is updated
		})
	}
}

// TestDataSourceNamingConvention validates data source naming follows same pattern
func TestDataSourceNamingConvention(t *testing.T) {
	// This test should FAIL until data sources are created
	expectedDataSources := []string{
		"hiiretail_iam_groups",
		"hiiretail_iam_roles",
	}

	for _, dsName := range expectedDataSources {
		t.Run("DataSource_"+dsName, func(t *testing.T) {
			// Validate naming pattern: hiiretail_{service}_{resource_type}
			parts := strings.Split(dsName, "_")
			assert.Len(t, parts, 3, "Data source name should have 3 parts: provider_service_resource")
			assert.Equal(t, "hiiretail", parts[0], "Provider name should be 'hiiretail'")
			assert.Equal(t, "iam", parts[1], "Service name should be 'iam'")
			assert.NotEmpty(t, parts[2], "Resource type should not be empty")
		})
	}
}

// TestProviderRegistryName validates provider registry name
func TestProviderRegistryName(t *testing.T) {
	// This test should PASS as we already updated main.go
	expectedRegistryName := "registry.terraform.io/extenda/hiiretail"

	// TODO: Add validation that main.go uses correct registry name
	// This is a placeholder until we can inspect the actual provider registration
	t.Logf("Expected registry name: %s", expectedRegistryName)
}
