package resource_iam_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestModelToAPIRequest_AliasAndNameOmitted(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}

	// Permission with alias set
	pv := PermissionsValue{
		Alias:      types.StringValue("ali"),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("perm1"),
		state:      attr.ValueStateKnown,
	}

	list, diags := types.ListValueFrom(ctx, permType, []PermissionsValue{pv})
	require.False(t, diags.HasError())

	data := IamCustomRoleModel{Id: types.StringValue("r1"), Name: types.StringNull(), Permissions: list}

	req, err := r.modelToAPIRequest(ctx, data)
	require.NoError(t, err)
	require.Equal(t, "r1", req.ID)
	// Name omitted should be empty
	require.Equal(t, "", req.Name)
	require.Equal(t, "perm1", req.Permissions[0].ID)
}

func TestUpdate_HappyPath_NameOmitted(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// Update response
	created := CustomRoleResponse{ID: "c1", Name: "C1", TenantID: "tid", Permissions: []Permission{{ID: "p1"}}}
	cb, _ := json.Marshal(created)

	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		if req.Method == "PUT" {
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(cb))}, nil
		}
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	// Plan has Name omitted
	data := IamCustomRoleModel{Id: types.StringValue("c1"), Name: types.StringNull(), Permissions: list}

	// Prepare UpdateRequest and State
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
	require.False(t, uresp.Diagnostics.HasError())
}

func TestModelToAPIRequest_AttributesStringAndNonStringValues(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}

	// Attributes object is known but empty (AttributesValue has no declared attribute types)
	attrVals := map[string]attr.Value{}

	pv := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), attrVals),
		Id:         types.StringValue("perm-1"),
		state:      attr.ValueStateKnown,
	}

	list, diags := types.ListValueFrom(ctx, permType, []PermissionsValue{pv})
	require.False(t, diags.HasError())

	data := IamCustomRoleModel{Id: types.StringValue("r1"), Name: types.StringNull(), Permissions: list}

	req, err := r.modelToAPIRequest(ctx, data)
	require.NoError(t, err)
	require.Equal(t, 1, len(req.Permissions))
	// No attributes should be included because the attributes object is empty
	require.Equal(t, 0, len(req.Permissions[0].Attributes))
}

func TestModelToAPIRequest_AttributesEmptyAndNull(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}

	// Attributes null should result in no attributes in request
	pv := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("perm-2"),
		state:      attr.ValueStateKnown,
	}

	list, diags := types.ListValueFrom(ctx, permType, []PermissionsValue{pv})
	require.False(t, diags.HasError())

	data := IamCustomRoleModel{Id: types.StringValue("r2"), Name: types.StringNull(), Permissions: list}

	req, err := r.modelToAPIRequest(ctx, data)
	require.NoError(t, err)
	require.Equal(t, 0, len(req.Permissions[0].Attributes))
}
