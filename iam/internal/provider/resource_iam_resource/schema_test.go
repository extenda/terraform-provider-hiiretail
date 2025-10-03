package resource_iam_resource_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/extenda/hiiretail-terraform-providers/hiiretail/internal/provider/resource_iam_resource"
)

// TestResourceCreation verifies that the resource can be created
func TestResourceCreation(t *testing.T) {
	// Test that NewIAMResourceResource creates a valid resource instance
	resource := resource_iam_resource.NewIAMResourceResource()
	assert.NotNil(t, resource, "Resource instance should not be nil")

	// This test should PASS as it only tests basic resource creation
}

// TestJSONValidation verifies the props JSON validation function
func TestJSONValidation(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		isValid bool
	}{
		{"empty string", "", true},
		{"null", "null", true},
		{"valid object", `{"key": "value"}`, true},
		{"valid array", `["item1", "item2"]`, true},
		{"valid number", "123", true},
		{"valid boolean", "true", true},
		{"valid string", `"hello"`, true},
		{"complex object", `{"location": "downtown", "active": true, "count": 42}`, true},
		{"invalid json", `{"invalid": json}`, false},
		{"unclosed object", `{"key": "value"`, false},
		{"invalid syntax", `{key: value}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test should PASS as the validateJSONString function should work
			err := resource_iam_resource.ValidateJSONString(tt.json) // This function doesn't exist yet, should FAIL

			if tt.isValid {
				assert.NoError(t, err, "Should be valid JSON: %s", tt.json)
			} else {
				assert.Error(t, err, "Should be invalid JSON: %s", tt.json)
			}
		})
	}
}

// TestResourceModelValidation verifies the resource data model structure
func TestResourceModelValidation(t *testing.T) {
	// Test that the model struct has the expected fields
	model := resource_iam_resource.IAMResourceResourceModel{}

	// Verify field types using reflection-like approach
	// Set some values to ensure the struct works as expected
	model.ID = types.StringValue("test:resource")
	model.Name = types.StringValue("Test Resource")
	model.Props = types.StringValue(`{"key": "value"}`)
	model.TenantID = types.StringValue("test-tenant")

	// Basic validation that the model can hold the expected data
	assert.False(t, model.ID.IsNull(), "ID should not be null")
	assert.False(t, model.Name.IsNull(), "Name should not be null")
	assert.False(t, model.Props.IsNull(), "Props should not be null")
	assert.False(t, model.TenantID.IsNull(), "TenantID should not be null")

	assert.Equal(t, "test:resource", model.ID.ValueString())
	assert.Equal(t, "Test Resource", model.Name.ValueString())
	assert.Equal(t, `{"key": "value"}`, model.Props.ValueString())
	assert.Equal(t, "test-tenant", model.TenantID.ValueString())
}
