package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// These tests exercise the branches in PermissionsType.ValueFromObject where
// attributes are present but of the wrong type.
func TestPermissions_ValueFromObject_WrongTypes(t *testing.T) {
	ctx := context.Background()

	// alias wrong type (should be basetypes.StringValue but we pass an ObjectValue)
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	wrongAliasObj, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
	obj, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"alias":      wrongAliasObj,
		"attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"id":         types.StringValue("id1"),
	})

	_, diags := PermissionsType{}.ValueFromObject(ctx, obj)
	require.True(t, diags.HasError())

	// attributes wrong type (should be ObjectValue but we pass a String)
	obj2, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.StringValue("not-object"),
		"id":         types.StringValue("id1"),
	})
	_, diags2 := PermissionsType{}.ValueFromObject(ctx, obj2)
	require.True(t, diags2.HasError())

	// id wrong type
	obj3, _ := types.ObjectValue(attrTypes, map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"id":         types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
	})
	_, diags3 := PermissionsType{}.ValueFromObject(ctx, obj3)
	require.True(t, diags3.HasError())
}
