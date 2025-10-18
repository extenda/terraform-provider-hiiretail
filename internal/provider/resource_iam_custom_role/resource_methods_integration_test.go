package resource_iam_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// Exercise framework-level Create/Read/Update/Delete methods to increase coverage
func TestResource_CreateReadUpdateDelete_Methods(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// API responses
	created := CustomRoleResponse{ID: "c1", Name: "C1", TenantID: "tid", Permissions: []Permission{{ID: "p1"}}}
	cb, _ := json.Marshal(created)
	rb, _ := json.Marshal(created)
	ub, _ := json.Marshal(created)

	// mock client to handle all CRUD methods
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		switch {
		case req.Method == "POST" && req.URL.Path == "/api/v1/tenants/tid/roles":
			return &http.Response{StatusCode: http.StatusCreated, Body: ioutil.NopCloser(bytes.NewBuffer(cb))}, nil
		case req.Method == "GET" && req.URL.Path == "/api/v1/tenants/tid/roles/c1":
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(rb))}, nil
		case req.Method == "PUT" && req.URL.Path == "/api/v1/tenants/tid/roles/c1":
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(ub))}, nil
		case req.Method == "DELETE" && req.URL.Path == "/api/v1/tenants/tid/roles/c1":
			return &http.Response{StatusCode: http.StatusNoContent, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
		default:
			return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer([]byte("err")))}, nil
		}
	}}}

	// Build an empty permissions list (required by schema)
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, diags := types.ListValueFrom(ctx, permType, []PermissionsValue{})
	require.False(t, diags.HasError(), "failed to build permissions list: %v", diags)

	// Prepare model for Create
	data := IamCustomRoleModel{Id: types.StringValue("c1"), Name: types.StringValue("C1"), Permissions: list}

	// Create
	var creq resource.CreateRequest
	creq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	diags = creq.Plan.Set(ctx, data)
	require.False(t, diags.HasError(), "plan set diags: %v", diags)

	var cresp resource.CreateResponse
	cresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Create(ctx, creq, &cresp)
	require.False(t, cresp.Diagnostics.HasError(), "Create produced diagnostics: %v", cresp.Diagnostics)

	// Verify state now contains ID
	var out IamCustomRoleModel
	diags = cresp.State.Get(ctx, &out)
	require.False(t, diags.HasError(), "state get diags: %v", diags)
	require.Equal(t, "c1", out.Id.ValueString())

	// Read
	var rreq resource.ReadRequest
	rreq.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags = rreq.State.Set(ctx, out)
	require.False(t, diags.HasError(), "state set diags: %v", diags)

	var rresp resource.ReadResponse
	rresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Read(ctx, rreq, &rresp)
	require.False(t, rresp.Diagnostics.HasError(), "Read produced diagnostics: %v", rresp.Diagnostics)

	var out2 IamCustomRoleModel
	diags = rresp.State.Get(ctx, &out2)
	require.False(t, diags.HasError(), "read state get diags: %v", diags)
	require.Equal(t, "c1", out2.Id.ValueString())

	// Update: prepare plan with same data
	var ureq resource.UpdateRequest
	ureq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	diags = ureq.Plan.Set(ctx, data)
	require.False(t, diags.HasError(), "update plan set diags: %v", diags)
	ureq.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags = ureq.State.Set(ctx, out2)
	require.False(t, diags.HasError(), "update state set diags: %v", diags)

	var uresp resource.UpdateResponse
	uresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Update(ctx, ureq, &uresp)
	require.False(t, uresp.Diagnostics.HasError(), "Update produced diagnostics: %v", uresp.Diagnostics)

	// Delete
	var dreq resource.DeleteRequest
	dreq.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags = dreq.State.Set(ctx, out2)
	require.False(t, diags.HasError(), "delete state set diags: %v", diags)

	var dresp resource.DeleteResponse

	r.Delete(ctx, dreq, &dresp)
	require.False(t, dresp.Diagnostics.HasError(), "Delete produced diagnostics: %v", dresp.Diagnostics)
}
