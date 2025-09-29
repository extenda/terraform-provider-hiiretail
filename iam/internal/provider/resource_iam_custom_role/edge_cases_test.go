package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T020: Edge Cases Unit Tests
// These tests verify handling of edge cases and boundary conditions

func TestEdgeCases_EmptyStringValues(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create model with empty strings
	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	emptyPermissions, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{})
	require.False(t, diags.HasError(), "Should create empty permissions list")

	model := IamCustomRoleModel{
		Id:          types.StringValue(""), // Empty ID
		Name:        types.StringValue(""), // Empty name
		Permissions: emptyPermissions,      // Empty permissions
		TenantId:    types.StringValue(""), // Empty tenant
	}

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	// Should not error, but should handle empty values
	assert.NoError(t, err, "Should handle empty strings without error")
	assert.Equal(t, "", apiReq.ID)
	assert.Equal(t, "", apiReq.Name)
	assert.Len(t, apiReq.Permissions, 0)
}

func TestEdgeCases_NullAndUnknownValues(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create model with null/unknown values
	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	emptyPermissions, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{})
	require.False(t, diags.HasError(), "Should create empty permissions list")

	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-001"),
		Name:        types.StringNull(), // Null name
		Permissions: emptyPermissions,
		TenantId:    types.StringUnknown(), // Unknown tenant
	}

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	assert.NoError(t, err, "Should handle null/unknown values")
	assert.Equal(t, "test-role-001", apiReq.ID)
	assert.Equal(t, "", apiReq.Name) // Null should become empty string
}

func TestEdgeCases_PermissionWithNullAlias(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create permission with null alias
	permission := PermissionsValue{
		Id:         types.StringValue("pos.payment.create"),
		Alias:      types.StringNull(), // Null alias
		Attributes: types.ObjectNull(map[string]attr.Type{}),
		state:      attr.ValueStateKnown,
	}

	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	permissionsList, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{permission})
	require.False(t, diags.HasError(), "Should create permissions list")

	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-001"),
		Name:        types.StringValue("Test Role"),
		Permissions: permissionsList,
		TenantId:    types.StringValue("test-tenant"),
	}

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	assert.NoError(t, err, "Should handle null alias")
	assert.Equal(t, "pos.payment.create", apiReq.Permissions[0].ID)
	assert.Equal(t, "", apiReq.Permissions[0].Alias) // Null should become empty string
}

func TestEdgeCases_PermissionWithUnknownAlias(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create permission with unknown alias
	permission := PermissionsValue{
		Id:         types.StringValue("pos.payment.create"),
		Alias:      types.StringUnknown(), // Unknown alias
		Attributes: types.ObjectNull(map[string]attr.Type{}),
		state:      attr.ValueStateKnown,
	}

	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	permissionsList, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{permission})
	require.False(t, diags.HasError(), "Should create permissions list")

	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-001"),
		Name:        types.StringValue("Test Role"),
		Permissions: permissionsList,
		TenantId:    types.StringValue("test-tenant"),
	}

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	assert.NoError(t, err, "Should handle unknown alias")
	assert.Equal(t, "pos.payment.create", apiReq.Permissions[0].ID)
	assert.Equal(t, "", apiReq.Permissions[0].Alias) // Unknown should become empty string
}

func TestEdgeCases_APIResponseWithMissingFields(t *testing.T) {
	r := &IamCustomRoleResource{}

	// API response with minimal fields
	apiResp := &CustomRoleResponse{
		ID: "test-role-001",
		// Name missing
		TenantID:    "test-tenant-123",
		Permissions: []Permission{}, // Empty permissions
		// CreatedAt, UpdatedAt missing
	}

	var model IamCustomRoleModel

	// Test conversion
	err := r.apiResponseToModel(context.Background(), apiResp, &model)

	assert.NoError(t, err, "Should handle missing fields")
	assert.Equal(t, "test-role-001", model.Id.ValueString())
	assert.Equal(t, "", model.Name.ValueString()) // Missing name should become empty
	assert.Equal(t, "test-tenant-123", model.TenantId.ValueString())
}

func TestEdgeCases_APIResponseWithEmptyStrings(t *testing.T) {
	r := &IamCustomRoleResource{}

	// API response with empty string fields
	apiResp := &CustomRoleResponse{
		ID:       "test-role-001",
		Name:     "", // Empty name
		TenantID: "test-tenant-123",
		Permissions: []Permission{
			{
				ID:    "pos.payment.create",
				Alias: "", // Empty alias
			},
		},
	}

	var model IamCustomRoleModel

	// Test conversion
	err := r.apiResponseToModel(context.Background(), apiResp, &model)

	assert.NoError(t, err, "Should handle empty strings")
	assert.Equal(t, "test-role-001", model.Id.ValueString())
	assert.Equal(t, "", model.Name.ValueString())
	assert.Equal(t, "test-tenant-123", model.TenantId.ValueString())
}

func TestEdgeCases_LargePermissionSets(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create maximum allowed permissions (simulate large set)
	const maxPermissions = 100
	permissions := make([]PermissionsValue, maxPermissions)

	for i := 0; i < maxPermissions; i++ {
		permissions[i] = PermissionsValue{
			Id:         types.StringValue("pos.payment.create"),
			Alias:      types.StringValue("Create Payment"),
			Attributes: types.ObjectNull(map[string]attr.Type{}),
			state:      attr.ValueStateKnown,
		}
	}

	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	permissionsList, diags := types.ListValueFrom(context.Background(), permissionType, permissions)
	require.False(t, diags.HasError(), "Should create large permissions list")

	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-large"),
		Name:        types.StringValue("Large Permission Set Role"),
		Permissions: permissionsList,
		TenantId:    types.StringValue("test-tenant"),
	}

	// Test conversion with large set
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	assert.NoError(t, err, "Should handle large permission sets")
	assert.Equal(t, "test-role-large", apiReq.ID)
	assert.Len(t, apiReq.Permissions, maxPermissions)
}

func TestEdgeCases_SpecialCharactersInStrings(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Test special characters in names and aliases
	permission := PermissionsValue{
		Id:         types.StringValue("pos.payment.create"),
		Alias:      types.StringValue("Create Payment with Special Chars: äöü €£¥"),
		Attributes: types.ObjectNull(map[string]attr.Type{}),
		state:      attr.ValueStateKnown,
	}

	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	permissionsList, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{permission})
	require.False(t, diags.HasError(), "Should create permissions list with special chars")

	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-special"),
		Name:        types.StringValue("Role with Special Chars: äöü €£¥ 中文 🚀"),
		Permissions: permissionsList,
		TenantId:    types.StringValue("test-tenant-äöü"),
	}

	// Test conversion with special characters
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	assert.NoError(t, err, "Should handle special characters")
	assert.Equal(t, "Role with Special Chars: äöü €£¥ 中文 🚀", apiReq.Name)
	assert.Equal(t, "Create Payment with Special Chars: äöü €£¥", apiReq.Permissions[0].Alias)
}

func TestEdgeCases_APIResponseConversionRobustness(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Test with various edge case API responses
	testCases := []struct {
		name     string
		response *CustomRoleResponse
		expectOK bool
	}{
		{
			name: "Normal response",
			response: &CustomRoleResponse{
				ID:       "test-001",
				Name:     "Test Role",
				TenantID: "tenant-123",
				Permissions: []Permission{
					{ID: "pos.payment.create", Alias: "Create"},
				},
			},
			expectOK: true,
		},
		{
			name: "Minimal response",
			response: &CustomRoleResponse{
				ID:          "test-002",
				TenantID:    "tenant-123",
				Permissions: []Permission{},
			},
			expectOK: true,
		},
		{
			name: "Response with nil permissions",
			response: &CustomRoleResponse{
				ID:          "test-003",
				Name:        "Test",
				TenantID:    "tenant-123",
				Permissions: nil, // nil slice
			},
			expectOK: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var model IamCustomRoleModel
			err := r.apiResponseToModel(context.Background(), tc.response, &model)

			if tc.expectOK {
				assert.NoError(t, err, "Should handle response without error")
				assert.Equal(t, tc.response.ID, model.Id.ValueString())
				assert.Equal(t, tc.response.TenantID, model.TenantId.ValueString())
			} else {
				assert.Error(t, err, "Should error for invalid response")
			}
		})
	}
}

func TestEdgeCases_ConcurrentDataAccess(t *testing.T) {
	// Test that data conversion is safe for concurrent access
	r := &IamCustomRoleResource{}

	permission := PermissionsValue{
		Id:         types.StringValue("pos.payment.create"),
		Alias:      types.StringValue("Create Payment"),
		Attributes: types.ObjectNull(map[string]attr.Type{}),
		state:      attr.ValueStateKnown,
	}

	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	permissionsList, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{permission})
	require.False(t, diags.HasError(), "Should create permissions list")

	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-concurrent"),
		Name:        types.StringValue("Concurrent Test Role"),
		Permissions: permissionsList,
		TenantId:    types.StringValue("test-tenant"),
	}

	// Run multiple conversions concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			apiReq, err := r.modelToAPIRequest(context.Background(), model)
			assert.NoError(t, err, "Concurrent conversion should not error")
			assert.Equal(t, "test-role-concurrent", apiReq.ID)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
