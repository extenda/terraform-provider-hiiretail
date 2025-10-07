package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// T016: Data Conversion Unit Tests
// These tests verify the conversion between Terraform models and API request/response formats

func TestModelToAPIRequest_ValidInput(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create test permission
	permission := PermissionsValue{
		Id:         types.StringValue("pos.payment.create"),
		Alias:      types.StringValue("Create Payment"),
		Attributes: types.ObjectNull(map[string]attr.Type{}),
		state:      attr.ValueStateKnown,
	}

	// Create permissions list
	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	permissionsList, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{permission})
	require.False(t, diags.HasError(), "Should create permissions list without error")

	// Create test model
	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-001"),
		Name:        types.StringValue("Test Custom Role"),
		Permissions: permissionsList,
		TenantId:    types.StringValue("test-tenant-123"),
	}

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	// Assert no error
	assert.NoError(t, err, "modelToAPIRequest should not error")

	// Verify converted data
	assert.Equal(t, "test-role-001", apiReq.ID)
	assert.Equal(t, "Test Custom Role", apiReq.Name)
	assert.Len(t, apiReq.Permissions, 1)
	assert.Equal(t, "pos.payment.create", apiReq.Permissions[0].ID)
	assert.Equal(t, "Create Payment", apiReq.Permissions[0].Alias)
}

func TestModelToAPIRequest_EmptyName(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create empty permissions list
	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	emptyPermissions, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{})
	require.False(t, diags.HasError(), "Should create empty permissions list without error")

	// Create test model with null name
	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-002"),
		Name:        types.StringNull(),
		Permissions: emptyPermissions,
		TenantId:    types.StringValue("test-tenant-123"),
	}

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	// Assert no error
	assert.NoError(t, err, "modelToAPIRequest should not error")

	// Verify converted data
	assert.Equal(t, "test-role-002", apiReq.ID)
	assert.Equal(t, "", apiReq.Name) // Should be empty string
	assert.Len(t, apiReq.Permissions, 0)
}

func TestModelToAPIRequest_MultiplePermissions(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create multiple test permissions
	permissions := []PermissionsValue{
		{
			Id:         types.StringValue("pos.payment.create"),
			Alias:      types.StringValue("Create Payment"),
			Attributes: types.ObjectNull(map[string]attr.Type{}),
			state:      attr.ValueStateKnown,
		},
		{
			Id:         types.StringValue("pos.payment.read"),
			Alias:      types.StringValue("Read Payment"),
			Attributes: types.ObjectNull(map[string]attr.Type{}),
			state:      attr.ValueStateKnown,
		},
		{
			Id:         types.StringValue("pos.refund.create"),
			Alias:      types.StringNull(),
			Attributes: types.ObjectNull(map[string]attr.Type{}),
			state:      attr.ValueStateKnown,
		},
	}

	// Create permissions list
	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	permissionsList, diags := types.ListValueFrom(context.Background(), permissionType, permissions)
	require.False(t, diags.HasError(), "Should create permissions list without error")

	// Create test model
	model := IamCustomRoleModel{
		Id:          types.StringValue("test-role-003"),
		Name:        types.StringValue("Multi Permission Role"),
		Permissions: permissionsList,
		TenantId:    types.StringValue("test-tenant-123"),
	}

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), model)

	// Assert no error
	assert.NoError(t, err, "modelToAPIRequest should not error")

	// Verify converted data
	assert.Equal(t, "test-role-003", apiReq.ID)
	assert.Equal(t, "Multi Permission Role", apiReq.Name)
	assert.Len(t, apiReq.Permissions, 3)

	// Check first permission
	assert.Equal(t, "pos.payment.create", apiReq.Permissions[0].ID)
	assert.Equal(t, "Create Payment", apiReq.Permissions[0].Alias)

	// Check second permission
	assert.Equal(t, "pos.payment.read", apiReq.Permissions[1].ID)
	assert.Equal(t, "Read Payment", apiReq.Permissions[1].Alias)

	// Check third permission (null alias)
	assert.Equal(t, "pos.refund.create", apiReq.Permissions[2].ID)
	assert.Equal(t, "", apiReq.Permissions[2].Alias)
}

func TestAPIResponseToModel_ValidResponse(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create test API response
	apiResp := &CustomRoleResponse{
		ID:       "test-role-001",
		Name:     "Test Custom Role",
		TenantID: "test-tenant-123",
		Permissions: []Permission{
			{
				ID:    "pos.payment.create",
				Alias: "Create Payment",
			},
			{
				ID:    "pos.payment.read",
				Alias: "Read Payment",
			},
		},
		CreatedAt: "2025-09-28T15:30:00Z",
		UpdatedAt: "2025-09-28T15:30:00Z",
	}

	// Create empty model
	var model IamCustomRoleModel

	// Test conversion
	err := r.apiResponseToModel(context.Background(), apiResp, &model)

	// Assert no error
	assert.NoError(t, err, "apiResponseToModel should not error")

	// Verify converted data
	assert.Equal(t, "test-role-001", model.Id.ValueString())
	assert.Equal(t, "Test Custom Role", model.Name.ValueString())
	assert.Equal(t, "test-tenant-123", model.TenantId.ValueString())

	// Check permissions were converted
	assert.False(t, model.Permissions.IsNull())
	assert.False(t, model.Permissions.IsUnknown())
}

func TestAPIResponseToModel_EmptyPermissions(t *testing.T) {
	r := &IamCustomRoleResource{}

	// Create test API response with no permissions
	apiResp := &CustomRoleResponse{
		ID:          "test-role-002",
		Name:        "Empty Role",
		TenantID:    "test-tenant-123",
		Permissions: []Permission{},
	}

	// Create empty model
	var model IamCustomRoleModel

	// Test conversion
	err := r.apiResponseToModel(context.Background(), apiResp, &model)

	// Assert no error
	assert.NoError(t, err, "apiResponseToModel should not error")

	// Verify converted data
	assert.Equal(t, "test-role-002", model.Id.ValueString())
	assert.Equal(t, "Empty Role", model.Name.ValueString())
	assert.Equal(t, "test-tenant-123", model.TenantId.ValueString())

	// Check permissions list is empty but not null
	assert.False(t, model.Permissions.IsNull())
	assert.False(t, model.Permissions.IsUnknown())
}
