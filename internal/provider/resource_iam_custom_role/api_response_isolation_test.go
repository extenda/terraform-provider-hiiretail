package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiResponseToModel_Isolated(t *testing.T) {
	// Create a test API response that matches what our mock server returns
	apiResp := &CustomRoleResponse{
		ID:       "test-role-001",
		Name:     "Test Custom Role",
		TenantID: "test-tenant-123",
		Permissions: []Permission{
			{
				ID:    "pos.payment.create",
				Alias: "Create Payment",
				// No attributes - consistent with our test case
			},
		},
	}

	// Create empty model to populate
	var data IamCustomRoleModel

	// Create resource instance
	resource := &IamCustomRoleResource{}

	// Call apiResponseToModel
	ctx := context.Background()
	err := resource.apiResponseToModel(ctx, apiResp, &data)

	// Verify no error
	assert.NoError(t, err, "apiResponseToModel should not return error")

	// Verify basic fields
	assert.Equal(t, "test-role-001", data.Id.ValueString(), "ID should be set correctly")
	assert.Equal(t, "Test Custom Role", data.Name.ValueString(), "Name should be set correctly")
	assert.Equal(t, "test-tenant-123", data.TenantId.ValueString(), "TenantID should be set correctly")

	// Verify permissions are not null
	assert.False(t, data.Permissions.IsNull(), "Permissions should not be null")
	assert.False(t, data.Permissions.IsUnknown(), "Permissions should not be unknown")

	// Debug: Check what's in the list before extraction
	t.Logf("Permissions list length: %d", len(data.Permissions.Elements()))
	t.Logf("Permissions list is null: %v", data.Permissions.IsNull())
	t.Logf("Permissions list is unknown: %v", data.Permissions.IsUnknown())

	// Convert permissions to slice to inspect
	var permissions []PermissionsValue
	diags := data.Permissions.ElementsAs(ctx, &permissions, false)
	if diags.HasError() {
		t.Logf("Diagnostics errors: %v", diags)
	}
	assert.False(t, diags.HasError(), "Should be able to extract permissions")

	// Verify permissions content
	assert.Len(t, permissions, 1, "Should have exactly 1 permission")

	perm := permissions[0]

	// Debug: Log what we actually got
	t.Logf("Permission ID: '%s' (length: %d)", perm.Id.ValueString(), len(perm.Id.ValueString()))
	t.Logf("Permission Alias: '%s' (length: %d)", perm.Alias.ValueString(), len(perm.Alias.ValueString()))
	t.Logf("Permission ID is null: %v", perm.Id.IsNull())
	t.Logf("Permission ID is unknown: %v", perm.Id.IsUnknown())

	assert.Equal(t, "pos.payment.create", perm.Id.ValueString(), "Permission ID should be correct")
	assert.Equal(t, "", perm.Alias.ValueString(), "Permission alias should be empty string (null)")
	assert.True(t, perm.Attributes.IsNull(), "Attributes should be null when not provided")

	t.Logf("✅ apiResponseToModel works correctly in isolation")
}
