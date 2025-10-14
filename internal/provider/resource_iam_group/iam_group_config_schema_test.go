package resource_iam_group

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

func TestIamGroup_MetadataAndSchema(t *testing.T) {
	r := NewIamGroupResource()
	var mr resource.MetadataResponse
	r.Metadata(nil, resource.MetadataRequest{ProviderTypeName: "hiiretail"}, &mr)
	require.Contains(t, mr.TypeName, "hiiretail_iam_group")

	var sr resource.SchemaResponse
	r.Schema(nil, resource.SchemaRequest{}, &sr)
	// Schema should be non-empty
	require.NotNil(t, sr.Schema)
}

func TestIamGroup_Configure_NilAndInvalidAndSuccess(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)

	ctx := context.Background()
	// Nil provider data -> no diagnostics, should return safely
	var cresp resource.ConfigureResponse
	r.Configure(ctx, resource.ConfigureRequest{ProviderData: nil}, &cresp)
	require.False(t, cresp.Diagnostics.HasError())

	// Invalid provider data shape -> diagnostic error
	invalid := struct{ Foo string }{Foo: "x"}
	var cresp2 resource.ConfigureResponse
	r.Configure(ctx, resource.ConfigureRequest{ProviderData: &invalid}, &cresp2)
	require.True(t, cresp2.Diagnostics.HasError())

	// Valid provider data shape
	type Prov struct {
		BaseURL    string
		TenantID   string
		HTTPClient *http.Client
	}
	client := &http.Client{}
	prov := &Prov{BaseURL: "https://api.test", TenantID: "t1", HTTPClient: client}
	var cresp3 resource.ConfigureResponse
	r.Configure(ctx, resource.ConfigureRequest{ProviderData: prov}, &cresp3)
	require.False(t, cresp3.Diagnostics.HasError())
	// internal fields set
	if r.baseURL == "" || r.tenantID == "" || r.client == nil {
		t.Fatalf("Configure did not set internal fields")
	}
}
