package resource_iam_custom_role

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIamCustomRoleResource_Metadata(t *testing.T) {
	r := NewIamCustomRoleResource()

	req := resource.MetadataRequest{
		ProviderTypeName: "hiiretail_iam",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "hiiretail_iam_custom_role", resp.TypeName)
}

func TestIamCustomRoleResource_Schema(t *testing.T) {
	r := NewIamCustomRoleResource()

	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	// Verify schema has required fields
	assert.NotNil(t, resp.Schema)
	assert.Contains(t, resp.Schema.Attributes, "id")
	assert.Contains(t, resp.Schema.Attributes, "permissions")
	assert.Contains(t, resp.Schema.Attributes, "name")
	assert.Contains(t, resp.Schema.Attributes, "tenant_id")
}

func TestIamCustomRoleResource_Configure(t *testing.T) {
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Test with nil provider data
	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	// Should not error with nil provider data
	assert.False(t, resp.Diagnostics.HasError())
}

// CRUD operation tests - T006 Create operation implemented
func TestIamCustomRoleResource_Create(t *testing.T) {
	// For now, just test that Create method exists and can be called without panic
	// Full integration testing will be done in Phase 4

	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Configure with minimal data
	configureReq := resource.ConfigureRequest{
		ProviderData: &APIClient{
			BaseURL:    "http://localhost:8080",
			TenantID:   "test-tenant",
			HTTPClient: &http.Client{},
		},
	}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)

	assert.False(t, configureResp.Diagnostics.HasError(), "Configure should not error")

	// Test helper method directly - modelToAPIRequest
	testData := IamCustomRoleModel{
		Id:   types.StringValue("test-role"),
		Name: types.StringValue("Test Role"),
	}

	// Create empty permissions list for now
	permissionType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(context.Background()),
		},
	}

	emptyPermissions, diags := types.ListValueFrom(context.Background(), permissionType, []PermissionsValue{})
	require.False(t, diags.HasError(), "Should create empty permissions list without error")
	testData.Permissions = emptyPermissions

	// Test conversion
	apiReq, err := r.modelToAPIRequest(context.Background(), testData)
	assert.NoError(t, err, "modelToAPIRequest should not error")
	assert.Equal(t, "test-role", apiReq.ID)
	assert.Equal(t, "Test Role", apiReq.Name)
	assert.Equal(t, 0, len(apiReq.Permissions))
}

func TestIamCustomRoleResource_Read(t *testing.T) {
	// Test Read operation - T007 implemented
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Configure with minimal data
	configureReq := resource.ConfigureRequest{
		ProviderData: &APIClient{
			BaseURL:    "http://localhost:8080",
			TenantID:   "test-tenant",
			HTTPClient: &http.Client{},
		},
	}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)

	assert.False(t, configureResp.Diagnostics.HasError(), "Configure should not error")

	// Test helper method directly - readCustomRole
	// This will fail with connection error since no mock server, but tests the method exists
	_, err := r.readCustomRole(context.Background(), "test-role-001")
	assert.Error(t, err, "Should error when no server available")
	assert.Contains(t, err.Error(), "HTTP request failed", "Should be connection error")
}

func TestIamCustomRoleResource_Update(t *testing.T) {
	// Test Update operation - T008 implemented
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://localhost:8080"
	r.tenantID = "test-tenant"
	r.client = &http.Client{}

	// Test helper method directly - updateCustomRole
	apiReq := &CustomRoleRequest{
		ID:   "test-role-001",
		Name: "Updated Role",
	}

	_, err := r.updateCustomRole(context.Background(), "test-role-001", apiReq)
	assert.Error(t, err, "Should error when no server available")
	assert.Contains(t, err.Error(), "HTTP request failed", "Should be connection error")
}

func TestIamCustomRoleResource_Delete(t *testing.T) {
	// Test Delete operation - T009 implemented
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://localhost:8080"
	r.tenantID = "test-tenant"
	r.client = &http.Client{}

	// Test helper method directly - deleteCustomRole
	err := r.deleteCustomRole(context.Background(), "test-role-001")
	assert.Error(t, err, "Should error when no server available")
	assert.Contains(t, err.Error(), "HTTP request failed", "Should be connection error")
}

func TestIamCustomRoleResource_ImportState(t *testing.T) {
	// Test ImportState operation - T010 implemented
	// For now, just test that the method exists and handles empty ID
	r := NewIamCustomRoleResource().(resource.ResourceWithImportState)

	// Test with empty ID (should error)
	req := resource.ImportStateRequest{
		ID: "",
	}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), req, resp)

	// Should error for empty ID
	assert.True(t, resp.Diagnostics.HasError(), "ImportState should error for empty ID")
	assert.Contains(t, resp.Diagnostics[0].Summary(), "Invalid Import ID")
}
