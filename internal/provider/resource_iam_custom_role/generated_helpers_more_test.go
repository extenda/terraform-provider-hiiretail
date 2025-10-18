package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPermissionsValue_ValueFromObject_Success_And_EqualNegative(t *testing.T) {
	ctx := context.Background()

	// Build an object with alias (string), attributes (empty object), id (string)
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	objVals := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("sys.r.a"),
	}

	obj := types.ObjectValueMust(attrTypes, objVals)

	pt := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	v, diags := pt.ValueFromObject(ctx, obj)
	require.False(t, diags.HasError())

	pv, ok := v.(PermissionsValue)
	require.True(t, ok)
	require.Equal(t, pv.Id.ValueString(), "sys.r.a")

	// Equal negative: change id
	pv2 := pv
	pv2.Id = types.StringValue("other")
	require.False(t, pv.Equal(pv2))

	// Different state comparison
	pvNull := PermissionsValue{state: attr.ValueStateNull}
	require.False(t, pv.Equal(pvNull))
}
