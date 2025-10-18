package resource_iam_custom_role

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestRead_NotFound_RemovesState(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// client returns 404 for GET
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	// set state with id 'missing'
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	var req resource.ReadRequest
	req.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags := req.State.Set(ctx, IamCustomRoleModel{Id: types.StringValue("missing"), Permissions: list})
	require.False(t, diags.HasError())

	var resp resource.ReadResponse
	resp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Read(ctx, req, &resp)
	// Read should not produce diagnostics, and should remove the resource from state
	require.False(t, resp.Diagnostics.HasError())
}

func TestUpdate_NotFound_ProducesDiagnostic(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// client returns 404 for PUT
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	data := IamCustomRoleModel{Id: types.StringValue("x"), Permissions: list}

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

func TestDelete_NotFound_TreatedAsSuccess(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// client returns 404 for DELETE
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	var dreq resource.DeleteRequest
	dreq.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags := dreq.State.Set(ctx, IamCustomRoleModel{Id: types.StringValue("missing"), Permissions: list})
	require.False(t, diags.HasError())

	var dresp resource.DeleteResponse

	r.Delete(ctx, dreq, &dresp)
	// Delete should not produce diagnostics when resource already missing
	require.False(t, dresp.Diagnostics.HasError())
}
