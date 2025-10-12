package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestGenerated_ExhaustiveHelpers(t *testing.T) {
	ctx := context.Background()

	// PermissionsType basic calls
	ptype := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	require.NotEmpty(t, ptype.String())
	v := ptype.ValueType(ctx)
	require.NotNil(t, v)

	// Build a known PermissionsValue and exercise Equal/ToObjectValue/ToTerraformValue
	pv := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("perm.id"),
		state:      attr.ValueStateKnown,
	}

	// Equal with same contents
	pv2 := pv
	require.True(t, pv.Equal(pv2))

	// Different id -> not equal
	pv3 := pv
	pv3.Id = types.StringValue("other")
	require.False(t, pv.Equal(pv3))

	// ToObjectValue should succeed
	obj, diags := pv.ToObjectValue(ctx)
	require.False(t, diags.HasError())
	require.NotNil(t, obj)

	// ToTerraformValue should succeed and then ValueFromTerraform roundtrip
	tfVal, err := pv.ToTerraformValue(ctx)
	require.NoError(t, err)
	out, err := ptype.ValueFromTerraform(ctx, tfVal)
	require.NoError(t, err)
	_, ok := out.(PermissionsValue)
	require.True(t, ok)

	// AttributesType basic calls
	atype := AttributesType{ObjectType: types.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}}
	require.NotEmpty(t, atype.String())
	av := atype.ValueType(ctx)
	require.NotNil(t, av)

	// AttributesValue equality for known/null/unknown
	aKnown := AttributesValue{state: attr.ValueStateKnown}
	aKnown2 := AttributesValue{state: attr.ValueStateKnown}
	require.True(t, aKnown.Equal(aKnown2))

	aNull := NewAttributesValueNull()
	aUnknown := NewAttributesValueUnknown()
	require.False(t, aKnown.Equal(aNull))
	require.False(t, aKnown.Equal(aUnknown))

	// ToObjectValue and ToTerraformValue for AttributesValue known/null/unknown
	objA, diags := aKnown.ToObjectValue(ctx)
	require.False(t, diags.HasError())
	require.NotNil(t, objA)

	tfA, err := aKnown.ToTerraformValue(ctx)
	require.NoError(t, err)

	// Use AttributesType.ValueFromTerraform to roundtrip the tf value
	// ensure no panic and no error when type matches
	outAttr, err := atype.ValueFromTerraform(ctx, tfA)
	require.NoError(t, err)
	_, ok = outAttr.(AttributesValue)
	require.True(t, ok)

	// Also exercise a tftypes.Value with nil Type to get New*Null branches
	nullVal := tftypes.Value{}
	_, err = ptype.ValueFromTerraform(ctx, nullVal)
	require.NoError(t, err)

	_, err = atype.ValueFromTerraform(ctx, nullVal)
	require.NoError(t, err)
}
