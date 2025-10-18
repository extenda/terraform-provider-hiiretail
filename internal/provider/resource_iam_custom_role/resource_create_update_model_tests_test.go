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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestResource_Create_HappyPath(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	created := CustomRoleResponse{ID: "c1", Name: "C1", TenantID: "tid", Permissions: []Permission{{ID: "p1"}}}
	cb, _ := json.Marshal(created)

	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		if req.Method == "POST" && req.URL.Path == "/api/v1/tenants/tid/roles" {
			return &http.Response{StatusCode: http.StatusCreated, Body: ioutil.NopCloser(bytes.NewBuffer(cb))}, nil
		}
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	// Build a plan with empty permissions list
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	data := IamCustomRoleModel{Id: types.StringValue("c1"), Name: types.StringValue("C1"), Permissions: list}

	var creq resource.CreateRequest
	creq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	diags := creq.Plan.Set(ctx, data)
	require.False(t, diags.HasError())

	var cresp resource.CreateResponse
	cresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Create(ctx, creq, &cresp)
	require.False(t, cresp.Diagnostics.HasError())
}

func TestResource_Create_APIError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	errResp := ErrorResponse{Message: "bad", Code: "BAD"}
	eb, _ := json.Marshal(errResp)
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(bytes.NewBuffer(eb))}, nil
	}}}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	data := IamCustomRoleModel{Id: types.StringValue("c1"), Name: types.StringValue("C1"), Permissions: list}

	var creq resource.CreateRequest
	creq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	diags := creq.Plan.Set(ctx, data)
	require.False(t, diags.HasError())

	var cresp resource.CreateResponse
	cresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Create(ctx, creq, &cresp)
	require.True(t, cresp.Diagnostics.HasError())
}

func TestResource_Update_InvalidStateAndHappyPath(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// Invalid state: missing ID
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	dataNoID := IamCustomRoleModel{Id: types.StringNull(), Name: types.StringValue("C1"), Permissions: list}
	var ureq resource.UpdateRequest
	ureq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	_ = ureq.Plan.Set(ctx, dataNoID)

	var uresp resource.UpdateResponse
	r.Update(ctx, ureq, &uresp)
	require.True(t, uresp.Diagnostics.HasError())

	// Happy path: mock PUT
	created := CustomRoleResponse{ID: "c2", Name: "C2", TenantID: "tid", Permissions: []Permission{{ID: "p1"}}}
	cb, _ := json.Marshal(created)
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		if req.Method == "PUT" {
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(cb))}, nil
		}
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}, nil
	}}}

	data := IamCustomRoleModel{Id: types.StringValue("c2"), Name: types.StringValue("C2"), Permissions: list}
	ureq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	_ = ureq.Plan.Set(ctx, data)
	ureq.State.Schema = IamCustomRoleResourceSchema(ctx)
	_ = ureq.State.Set(ctx, data)

	uresp = resource.UpdateResponse{}
	uresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Update(ctx, ureq, &uresp)
	require.False(t, uresp.Diagnostics.HasError())
}

func TestModelToAPIRequest_ElementsAsDiagnostic(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Build a list value that contains the wrong type (string instead of object)
	listVal, _ := types.ListValueFrom(ctx, basetypes.StringType{}, []types.String{types.StringValue("not-an-object")})

	data := IamCustomRoleModel{Id: types.StringValue("r1"), Name: types.StringNull(), Permissions: listVal}

	_, err := r.modelToAPIRequest(ctx, data)
	require.Error(t, err)
}
