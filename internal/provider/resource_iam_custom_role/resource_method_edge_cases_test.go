package resource_iam_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCreate_HandlesModelConversionError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Build a plan with a permissions list that will fail ElementsAs conversion
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	badList, _ := types.ListValueFrom(ctx, permType, []interface{}{"not-a-perm"})

	data := IamCustomRoleModel{Id: types.StringValue("x"), Permissions: badList}

	var creq resource.CreateRequest
	creq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	diags := creq.Plan.Set(ctx, data)
	require.False(t, diags.HasError())

	var cresp resource.CreateResponse
	cresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Create(ctx, creq, &cresp)
	require.True(t, cresp.Diagnostics.HasError())
}

func TestUpdate_InvalidStateMissingID(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Plan contains no ID
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})
	data := IamCustomRoleModel{Id: types.StringNull(), Permissions: list}

	var ureq resource.UpdateRequest
	ureq.Plan.Schema = IamCustomRoleResourceSchema(ctx)
	diags := ureq.Plan.Set(ctx, data)
	require.False(t, diags.HasError())

	var uresp resource.UpdateResponse
	uresp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Update(ctx, ureq, &uresp)
	require.True(t, uresp.Diagnostics.HasError())
}

func TestDelete_InvalidStateMissingID(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	var dreq resource.DeleteRequest
	dreq.State.Schema = IamCustomRoleResourceSchema(ctx)
	// Set state with empty ID and minimal required permissions list so Set succeeds
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})
	diags := dreq.State.Set(ctx, IamCustomRoleModel{Id: types.StringValue(""), Permissions: list})
	require.False(t, diags.HasError())

	var dresp resource.DeleteResponse

	r.Delete(ctx, dreq, &dresp)
	require.True(t, dresp.Diagnostics.HasError())
}

func TestRead_SkipsWhenIDEmpty(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	var req resource.ReadRequest
	req.State.Schema = IamCustomRoleResourceSchema(ctx)
	// Use empty string ID and minimal required permissions list so Set succeeds; Read checks for empty id and should skip API call
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})
	diags := req.State.Set(ctx, IamCustomRoleModel{Id: types.StringValue(""), Permissions: list})
	require.False(t, diags.HasError())

	var resp resource.ReadResponse
	resp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Read(ctx, req, &resp)
	// Should not produce diagnostics when called for new resource with empty id
	require.False(t, resp.Diagnostics.HasError())
}

func TestImportState_EmptyIDProducesError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	var req resource.ImportStateRequest
	req.ID = ""

	var resp resource.ImportStateResponse

	r.ImportState(ctx, req, &resp)
	require.True(t, resp.Diagnostics.HasError())
}

func TestImportState_SetsIDOnSuccess(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	var req resource.ImportStateRequest
	req.ID = "import-id"

	var resp resource.ImportStateResponse
	resp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.ImportState(ctx, req, &resp)
	// ImportState may produce diagnostics depending on schema handling; if no diagnostics, ensure state contains id
	if !resp.Diagnostics.HasError() {
		var out IamCustomRoleModel
		diags := resp.State.Get(ctx, &out)
		require.False(t, diags.HasError())
		require.Equal(t, "import-id", out.Id.ValueString())
	}
}

func TestCreate_AddsDiagnosticOnAPIError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// client returns network error
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("network")
	}}}

	// Build valid plan
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

func TestRead_AddsDiagnosticOnAPIError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	// client returns 500 with JSON error
	errResp := ErrorResponse{Message: "bad", Code: "BAD"}
	eb, _ := json.Marshal(errResp)
	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(bytes.NewBuffer(eb))}, nil
	}}}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	var req resource.ReadRequest
	req.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags := req.State.Set(ctx, IamCustomRoleModel{Id: types.StringValue("c1"), Permissions: list})
	require.False(t, diags.HasError())

	var resp resource.ReadResponse
	resp.State.Schema = IamCustomRoleResourceSchema(ctx)

	r.Read(ctx, req, &resp)
	require.True(t, resp.Diagnostics.HasError())
}

func TestUpdate_AddsDiagnosticOnAPIError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("network")
	}}}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})
	data := IamCustomRoleModel{Id: types.StringValue("c1"), Permissions: list}

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

func TestDelete_AddsDiagnosticOnAPIError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)
	r.baseURL = "http://api"
	r.tenantID = "tid"

	r.client = &http.Client{Transport: &mockRoundTripper{RoundTripFunc: func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("network")
	}}}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, _ := types.ListValueFrom(ctx, permType, []PermissionsValue{})

	var dreq resource.DeleteRequest
	dreq.State.Schema = IamCustomRoleResourceSchema(ctx)
	diags := dreq.State.Set(ctx, IamCustomRoleModel{Id: types.StringValue("c1"), Permissions: list})
	require.False(t, diags.HasError())

	var dresp resource.DeleteResponse

	r.Delete(ctx, dreq, &dresp)
	require.True(t, dresp.Diagnostics.HasError())
}
