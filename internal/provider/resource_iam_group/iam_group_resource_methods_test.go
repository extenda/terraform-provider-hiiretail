package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// Note: import above includes the package path for completeness, but tests run in the same package.

func TestIamGroup_Create_HappyAndValidation(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// Validation error: missing name
	var creq resource.CreateRequest
	creq.Plan.Schema = IamGroupResourceSchema(ctx)
	// empty model
	data := IamGroupModel{Id: types.StringNull(), Name: types.StringNull()}
	_ = creq.Plan.Set(ctx, data)

	var cresp resource.CreateResponse
	cresp.State.Schema = IamGroupResourceSchema(ctx)

	r.Create(ctx, creq, &cresp)
	require.True(t, cresp.Diagnostics.HasError(), "expected validation error when name missing")

	// Validation error: description too long
	var creq4 resource.CreateRequest
	creq4.Plan.Schema = IamGroupResourceSchema(ctx)
	longDesc := string(make([]byte, 256)) // 256 chars
	data4 := IamGroupModel{Id: types.StringNull(), Name: types.StringValue("test"), Description: types.StringValue(longDesc)}
	_ = creq4.Plan.Set(ctx, data4)
	var cresp4 resource.CreateResponse
	cresp4.State.Schema = IamGroupResourceSchema(ctx)
	r.Create(ctx, creq4, &cresp4)
	require.True(t, cresp4.Diagnostics.HasError(), "expected validation error when description too long")

	// Happy path
	var creq2 resource.CreateRequest
	creq2.Plan.Schema = IamGroupResourceSchema(ctx)
	data2 := IamGroupModel{Id: types.StringNull(), Name: types.StringValue("devs")}
	_ = creq2.Plan.Set(ctx, data2)
	var cresp2 resource.CreateResponse
	cresp2.State.Schema = IamGroupResourceSchema(ctx)

	r.Create(ctx, creq2, &cresp2)
	require.False(t, cresp2.Diagnostics.HasError(), "create should succeed for valid data")

	var out IamGroupModel
	_ = cresp2.State.Get(ctx, &out)
	require.NotEmpty(t, out.Id.ValueString(), "expected ID to be set")
	require.Equal(t, "active", out.Status.ValueString())
}

func TestIamGroup_Read_Update_Delete_Import(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// Prepare a state with a resource
	data := IamGroupModel{Id: types.StringValue("g1"), Name: types.StringValue("team"), Description: types.StringValue("desc"), Status: types.StringValue("active"), TenantId: types.StringValue("t1")}

	// Read should accept the existing state and set it back
	var rreq resource.ReadRequest
	rreq.State.Schema = IamGroupResourceSchema(ctx)
	_ = rreq.State.Set(ctx, data)
	var rresp resource.ReadResponse
	rresp.State.Schema = IamGroupResourceSchema(ctx)
	r.Read(ctx, rreq, &rresp)
	require.False(t, rresp.Diagnostics.HasError())

	// Update: invalid state (missing Name should trigger validation)
	var ureq resource.UpdateRequest
	ureq.Plan.Schema = IamGroupResourceSchema(ctx)
	_ = ureq.Plan.Set(ctx, IamGroupModel{Id: types.StringValue("g1"), Name: types.StringNull()})
	var uresp resource.UpdateResponse
	uresp.State.Schema = IamGroupResourceSchema(ctx)
	r.Update(ctx, ureq, &uresp)
	require.True(t, uresp.Diagnostics.HasError())

	// Update: happy path
	var ureq2 resource.UpdateRequest
	ureq2.Plan.Schema = IamGroupResourceSchema(ctx)
	ureq2.State.Schema = IamGroupResourceSchema(ctx)
	_ = ureq2.Plan.Set(ctx, data)
	_ = ureq2.State.Set(ctx, data)
	var uresp2 resource.UpdateResponse
	uresp2.State.Schema = IamGroupResourceSchema(ctx)
	r.Update(ctx, ureq2, &uresp2)
	require.False(t, uresp2.Diagnostics.HasError())

	// Delete should accept state and not error
	var dreq resource.DeleteRequest
	dreq.State.Schema = IamGroupResourceSchema(ctx)
	_ = dreq.State.Set(ctx, data)
	var dresp resource.DeleteResponse
	dresp.State.Schema = IamGroupResourceSchema(ctx)
	r.Delete(ctx, dreq, &dresp)
	require.False(t, dresp.Diagnostics.HasError())

	// ImportState should populate state with provided ID
	var ireq resource.ImportStateRequest
	ireq.ID = "imported-1"
	var iresp resource.ImportStateResponse
	iresp.State.Schema = IamGroupResourceSchema(ctx)
	r.ImportState(ctx, ireq, &iresp)
	require.False(t, iresp.Diagnostics.HasError())
	var imported IamGroupModel
	_ = iresp.State.Get(ctx, &imported)
	require.Equal(t, "imported-1", imported.Id.ValueString())
}

func Test_makeAPIRequest_UnexpectedStatusDefaultMaps(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// Return a 418 to hit the default case in mapHTTPError via makeAPIRequest
	seq := &seqRoundTripper{responses: []*http.Response{makeResp(418, `I'm a teapot`)}}
	r.client = &http.Client{Transport: seq}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(ctx, http.MethodGet, "/teapot", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected HTTP status 418")
}

func Test_Configure_InvalidProviderDataProducesDiagnostic(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	// Pass provider data that's not a struct/expected type to Configure
	var req resource.ConfigureRequest
	req.ProviderData = "not-a-struct"
	var resp resource.ConfigureResponse

	r.Configure(ctx, req, &resp)
	require.True(t, resp.Diagnostics.HasError())
}
