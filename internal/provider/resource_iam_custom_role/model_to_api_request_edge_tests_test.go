package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestModelToAPIRequest_ElementsAsTypeError(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Construct a list where elements are plain strings instead of PermissionsValue
	// Use a string list type to force ElementsAs to fail when converting to PermissionsValue
	listType := types.ListType{ElemType: types.StringType}
	lst, _ := types.ListValueFrom(ctx, listType, []string{"not-a-permission"})

	data := IamCustomRoleModel{Id: types.StringValue("r1"), Name: types.StringNull(), Permissions: lst}

	_, err := r.modelToAPIRequest(ctx, data)
	require.Error(t, err)
}

// Note: Update missing-ID case is covered in resource_method_edge_cases_test.go
