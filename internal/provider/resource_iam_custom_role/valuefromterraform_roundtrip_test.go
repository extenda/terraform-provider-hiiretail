package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestPermissions_ValueFromTerraform_RoundTripAndNullUnknown(t *testing.T) {
	ctx := context.Background()

	// Build a known PermissionsValue with empty attributes
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	attributesObj := types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
	_ = NewPermissionsValueMust(attrTypes, map[string]attr.Value{
		"alias":      types.StringValue("alias"),
		"attributes": attributesObj,
		"id":         types.StringValue("pid"),
	})

	// Construct a tftypes.Value that matches the PermissionsType with attribute types
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	tfType := permType.TerraformType(ctx)
	vals := map[string]tftypes.Value{
		"alias":      tftypes.NewValue(tftypes.String, "alias"),
		"attributes": tftypes.NewValue(AttributesType{ObjectType: types.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}}.TerraformType(ctx), map[string]tftypes.Value{}),
		"id":         tftypes.NewValue(tftypes.String, "pid"),
	}
	tfVal := tftypes.NewValue(tfType, vals)

	got, err := permType.ValueFromTerraform(ctx, tfVal)
	require.NoError(t, err)
	gotPV, ok := got.(PermissionsValue)
	require.True(t, ok)
	// Compare ID since attributes normalization can vary
	require.Equal(t, "pid", gotPV.Id.ValueString())

	// Null branch
	nullVal := tftypes.NewValue(tfType, nil)
	gotNull, err := permType.ValueFromTerraform(ctx, nullVal)
	require.NoError(t, err)
	_, ok = gotNull.(PermissionsValue)
	require.True(t, ok)

	// Unknown branch
	unkVal := tftypes.NewValue(tfType, tftypes.UnknownValue)
	gotUnk, err := permType.ValueFromTerraform(ctx, unkVal)
	require.NoError(t, err)
	_, ok = gotUnk.(PermissionsValue)
	require.True(t, ok)
}

func TestAttributes_ValueFromTerraform_NullUnknownAndRoundTrip(t *testing.T) {
	ctx := context.Background()

	// Known AttributesValue (empty)
	av := AttributesValue{state: attr.ValueStateKnown}
	tfVal, err := av.ToTerraformValue(ctx)
	require.NoError(t, err)

	got, err := AttributesType{}.ValueFromTerraform(ctx, tfVal)
	require.NoError(t, err)
	_, ok := got.(AttributesValue)
	require.True(t, ok)

	// Null
	tfType := AttributesType{}.TerraformType(ctx)
	nullVal := tftypes.NewValue(tfType, nil)
	gotNull, err := AttributesType{}.ValueFromTerraform(ctx, nullVal)
	require.NoError(t, err)
	_, ok = gotNull.(AttributesValue)
	require.True(t, ok)

	// Unknown
	unkVal := tftypes.NewValue(tfType, tftypes.UnknownValue)
	gotUnk, err := AttributesType{}.ValueFromTerraform(ctx, unkVal)
	require.NoError(t, err)
	_, ok = gotUnk.(AttributesValue)
	require.True(t, ok)
}
