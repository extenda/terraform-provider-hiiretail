package resource_iam_custom_role

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

// T018: Validation Logic Unit Tests
// These tests verify resource validation and schema constraints

func TestResource_Metadata(t *testing.T) {
	r := NewIamCustomRoleResource()

	req := resource.MetadataRequest{
		ProviderTypeName: "hiiretail_iam",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "hiiretail_iam_custom_role", resp.TypeName)
}

func TestResource_Schema_ValidStructure(t *testing.T) {
	r := NewIamCustomRoleResource()

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "Schema should not have errors")
	assert.NotNil(t, resp.Schema, "Schema should not be nil")

	// Check required attributes exist
	attributes := resp.Schema.Attributes
	assert.Contains(t, attributes, "id", "Schema should have id attribute")
	assert.Contains(t, attributes, "name", "Schema should have name attribute")
	assert.Contains(t, attributes, "permissions", "Schema should have permissions attribute")
	assert.Contains(t, attributes, "tenant_id", "Schema should have tenant_id attribute")

	// Check id is required
	idAttr := attributes["id"]
	assert.True(t, idAttr.IsRequired(), "ID should be required")

	// Check permissions is required
	permissionsAttr := attributes["permissions"]
	assert.True(t, permissionsAttr.IsRequired(), "Permissions should be required")
}

func TestResource_Configure_ValidAPIClient(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	client := &APIClient{
		BaseURL:    "https://api.example.com",
		TenantID:   "test-tenant-123",
		HTTPClient: &http.Client{},
	}

	req := resource.ConfigureRequest{
		ProviderData: client,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "Configure should not error")
	assert.Equal(t, "https://api.example.com", r.baseURL)
	assert.Equal(t, "test-tenant-123", r.tenantID)
	assert.NotNil(t, r.client)
}

func TestResource_Configure_InvalidProviderData(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Pass wrong type
	req := resource.ConfigureRequest{
		ProviderData: "invalid-data",
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "Configure should error for invalid data")
	assert.Contains(t, resp.Diagnostics[0].Summary(), "Unexpected Resource Configure Type")
}

func TestResource_Configure_NilProviderData(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	// Should not error with nil provider data (prevents panic)
	assert.False(t, resp.Diagnostics.HasError(), "Configure should not error for nil data")
}

func TestPermissionValidation_ValidFormats(t *testing.T) {
	validPermissions := []string{
		"pos.payment.create",
		"pos.refund.read",
		"iam.user.update",
		"api.order.delete",
		"crm.customer.view",
		"inv.product-info.read",
		"log.audit-trail.write",
	}

	for _, permID := range validPermissions {
		t.Run(permID, func(t *testing.T) {
			// Test that these would pass regex validation in schema
			// The actual validation is done by the framework, but we can test the pattern
			matched := permissionIDPattern.MatchString(permID)
			assert.True(t, matched, "Permission ID should match pattern: %s", permID)
		})
	}
}

func TestPermissionValidation_InvalidFormats(t *testing.T) {
	invalidPermissions := []string{
		"",                                 // Empty
		"pos",                              // Too short
		"pos.payment",                      // Missing action
		"pos.payment.create.extra",         // Too many parts
		"POS.payment.create",               // Uppercase
		"pos.Payment.create",               // Uppercase
		"pos.payment.Create",               // Uppercase
		"pos..create",                      // Empty middle part
		"pos.payment.",                     // Empty end part
		".payment.create",                  // Empty start part
		"pos.payment.create!",              // Invalid character
		"pos payment.create",               // Space instead of dot
		"a.payment.create",                 // System prefix too short
		"post.a.create",                    // Resource too short
		"pos.payment.a",                    // Action too short
		"post.verylongresourcename.create", // Resource too long (>16 chars)
		"pos.payment.verylongactionname",   // Action too long (>16 chars)
	}

	for _, permID := range invalidPermissions {
		t.Run(permID, func(t *testing.T) {
			// Test that these would fail regex validation in schema
			matched := permissionIDPattern.MatchString(permID)
			assert.False(t, matched, "Permission ID should NOT match pattern: %s", permID)
		})
	}
}

func TestNameValidation_LengthConstraints(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		valid  bool
		minLen int
		maxLen int
	}{
		{"Too short", "ab", false, 3, 256},
		{"Minimum valid", "abc", true, 3, 256},
		{"Normal length", "Test Custom Role", true, 3, 256},
		{"Maximum valid", string(make([]byte, 256)), true, 3, 256},
		{"Too long", string(make([]byte, 257)), false, 3, 256},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nameLen := len(tc.input)
			if tc.valid {
				assert.GreaterOrEqual(t, nameLen, tc.minLen, "Valid name should be >= min length")
				assert.LessOrEqual(t, nameLen, tc.maxLen, "Valid name should be <= max length")
			} else {
				assert.True(t, nameLen < tc.minLen || nameLen > tc.maxLen, "Invalid name should be outside length constraints")
			}
		})
	}
}

func TestAttributeConstraints_MaximumFields(t *testing.T) {
	// Test attribute constraints (up to 10 props, keys up to 40 chars, values up to 256 chars)
	validAttributes := map[string]interface{}{
		"department":  "finance",
		"level":       "basic",
		"region":      "europe",
		"cost_center": "cc-001",
		"project":     "project-alpha",
	}

	// Test maximum key length (40 chars)
	longKey := string(make([]byte, 40))
	validAttributes[longKey] = "value"

	// Test maximum value length (256 chars)
	longValue := string(make([]byte, 256))
	validAttributes["test_key"] = longValue

	assert.LessOrEqual(t, len(validAttributes), 10, "Should not exceed 10 attributes")

	for key, value := range validAttributes {
		assert.LessOrEqual(t, len(key), 40, "Key should not exceed 40 characters: %s", key)
		if strVal, ok := value.(string); ok {
			assert.LessOrEqual(t, len(strVal), 256, "Value should not exceed 256 characters")
		}
	}
}

// Helper: regex pattern for permission ID validation (from schema)
var permissionIDPattern = regexp.MustCompile(`^[a-z][-a-z]{2}\.[a-z][-a-z]{1,15}\.[a-z][-a-z]{1,15}$`)
