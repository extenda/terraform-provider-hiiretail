package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"
)

// The resource.Configure method expects provider data with a specific shape.
// exercise the negative paths of extractAPIClientFields by passing nil and
// wrong-typed provider data via resource.ConfigureRequest.
func TestConfigure_EmptyProviderData(t *testing.T) {
	ctx := context.Background()

	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// provider data nil should be handled gracefully and result in a diagnostic
	req := resource.ConfigureRequest{ProviderData: nil}
	resp := &resource.ConfigureResponse{}

	r.Configure(ctx, req, resp)

	// Configure intentionally returns without error when provider data is nil
	// (prevents panic for unconfigured provider). Ensure there are no errors.
	require.False(t, resp.Diagnostics.HasError())
}

type bogusProvider struct {
	X string
}

func TestConfigure_WrongProviderDataType(t *testing.T) {
	ctx := context.Background()

	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// pass provider data of the wrong type
	req := resource.ConfigureRequest{ProviderData: &bogusProvider{X: "nope"}}
	resp := &resource.ConfigureResponse{}

	r.Configure(ctx, req, resp)

	require.True(t, resp.Diagnostics.HasError())
}
