package resource_iam_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

// Passing a non-struct providerData should return nil from extractAPIClientFields
func Test_extractAPIClientFields_NonStructAndMissingField(t *testing.T) {
	var v interface{} = 1234
	got := extractAPIClientFields(v)
	require.Nil(t, got)

	// struct missing expected fields
	type partial struct {
		Something string
	}
	p := partial{Something: "x"}
	got2 := extractAPIClientFields(p)
	require.Nil(t, got2)
}

// Configure should add a diagnostic error when ProviderData isn't the expected APIClient
func TestIamGroup_Configure_InvalidProviderData_AddsError(t *testing.T) {
	ctx := context.Background()
	r := NewIamGroupResource().(*IamGroupResource)

	req := resource.ConfigureRequest{ProviderData: 123}
	var resp resource.ConfigureResponse

	r.Configure(ctx, req, &resp)
	require.True(t, resp.Diagnostics.HasError(), "expected Configure to add an error diagnostic when ProviderData invalid")
}
