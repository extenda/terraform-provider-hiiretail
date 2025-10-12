package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestGeneratedHelpers_EqualNegativeAndTerraformValueBranches(t *testing.T) {
	ctx := context.Background()

	// Compare AttributesValue to a different type
	av := AttributesValue{state: attr.ValueStateKnown}
	var other attr.Value = PermissionsValue{state: attr.ValueStateKnown}
	require.False(t, av.Equal(other))

	// PermissionsValue Equal negative: different state
	pvKnown := PermissionsValue{state: attr.ValueStateKnown, Alias: types.StringValue("a"), Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)), Id: types.StringValue("x")}
	pvNull := PermissionsValue{state: attr.ValueStateNull}
	require.False(t, pvKnown.Equal(pvNull))
}

func TestModelToAPIRequest_MixedAttributeTypes(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Build an empty attributes object using the generated AttributesValue attribute types
	attrsObj := types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})

	pv := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: attrsObj,
		Id:         types.StringValue("perm-mix"),
		state:      attr.ValueStateKnown,
	}

	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, diags := types.ListValueFrom(ctx, permType, []PermissionsValue{pv})
	require.False(t, diags.HasError())

	data := IamCustomRoleModel{Id: types.StringValue("r-mix"), Name: types.StringNull(), Permissions: list}
	req, err := r.modelToAPIRequest(ctx, data)
	require.NoError(t, err)
	require.Equal(t, 1, len(req.Permissions))
	// AttributesValue has no declared attribute types in the generated code,
	// so we expect the request to include no attributes (the provider filters by declared attr types).
	require.Equal(t, 0, len(req.Permissions[0].Attributes))
}
