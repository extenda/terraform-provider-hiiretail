package resource_iam_custom_role

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestConfigure_ExtractAPIClientFields_Success(t *testing.T) {
	ctx := context.Background()

	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Create a provider-like struct with the expected fields
	type providerLike struct {
		BaseURL    string
		TenantID   string
		HTTPClient *http.Client
	}

	p := &providerLike{BaseURL: "http://api", TenantID: "tid", HTTPClient: &http.Client{}}
	req := resource.ConfigureRequest{ProviderData: p}
	resp := &resource.ConfigureResponse{}

	// Call Configure via the resource method directly
	r.Configure(ctx, req, resp)

	// Ensure the resource captured the fields
	require.Equal(t, "http://api", r.baseURL)
	require.Equal(t, "tid", r.tenantID)
	require.NotNil(t, r.client)
}

// End of configure tests

// Test modelToAPIRequest converts PermissionsValue with alias and attributes into API request
func TestModelToAPIRequest_PermissionsWithAliasAndAttributes(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Build a PermissionsValue with alias and null attributes (AttributesValue has no attribute types)
	attributesObj := types.ObjectNull(AttributesValue{}.AttributeTypes(ctx))

	pv := PermissionsValue{
		Alias:      types.StringValue("alias"),
		Attributes: attributesObj,
		Id:         types.StringValue("perm.id"),
		state:      attr.ValueStateKnown,
	}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{pv})

	data := IamCustomRoleModel{Id: types.StringValue("r1"), Permissions: list}

	apiReq, err := r.modelToAPIRequest(ctx, data)
	require.NoError(t, err)
	require.Len(t, apiReq.Permissions, 1)
	pa := apiReq.Permissions[0]
	require.Equal(t, "perm.id", pa.ID)
	// attributes are null (AttributesValue has no attribute types)
	require.Nil(t, pa.Attributes)
}
