package resource_iam_custom_role

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestExtractAPIClientFields_PositiveAndNegative(t *testing.T) {
	type good struct {
		BaseURL    string
		TenantID   string
		HTTPClient *http.Client
	}

	g := good{BaseURL: "http://x", TenantID: "tid", HTTPClient: &http.Client{}}
	got := extractAPIClientFields(&g)
	require.NotNil(t, got)
	require.Equal(t, "http://x", got.BaseURL)
	require.Equal(t, "tid", got.TenantID)

	type bad struct {
		BaseURL    int
		TenantID   int
		HTTPClient int
	}

	b := bad{BaseURL: 1, TenantID: 2, HTTPClient: 3}
	got2 := extractAPIClientFields(&b)
	require.Nil(t, got2)
}

func TestModelToAPIRequest_BadPermissionsElements(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Create a list whose elements are strings (wrong type) to force ElementsAs to fail
	list := types.ListValueMust(basetypes.StringType{}, []attr.Value{types.StringValue("x")})

	data := IamCustomRoleModel{Id: types.StringValue("r1"), Name: types.StringNull(), Permissions: list}

	_, err := r.modelToAPIRequest(ctx, data)
	require.Error(t, err)
}

func TestUpdate_MissingID_Error(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	// Id is known but empty string -> should trigger Invalid State error in Update
	data := IamCustomRoleModel{Id: types.StringValue(""), Name: types.StringNull(), Permissions: list}

	var ureq resource.UpdateRequest
	ureq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	diags := ureq.Plan.Set(ctx, data)
	require.False(t, diags.HasError())
	ureq.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags = ureq.State.Set(ctx, data)
	require.False(t, diags.HasError())

	var uresp resource.UpdateResponse
	uresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Update(ctx, ureq, &uresp)
	require.True(t, uresp.Diagnostics.HasError())
}
