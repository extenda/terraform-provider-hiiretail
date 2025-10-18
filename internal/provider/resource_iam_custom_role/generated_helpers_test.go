package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPermissionsType_ValueFromObject_SuccessAndMissing(t *testing.T) {
	ctx := context.Background()

	// Build an object value with required attributes
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	// Build a proper object value and ensure ValueFromObject succeeds
	attrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("p1"),
	}
	obj, diags := types.ObjectValue(attrTypes, attrs)
	require.False(t, diags.HasError())

	var ptype PermissionsType
	val, diags := ptype.ValueFromObject(ctx, obj)
	require.False(t, diags.HasError())
	pv, ok := val.(PermissionsValue)
	require.True(t, ok)
	require.Equal(t, "p1", pv.Id.ValueString())

	// Now create an object missing the 'id' attribute to force diagnostic
	badAttrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
	}
	badObj, diagsObj := types.ObjectValue(attrTypes, badAttrs)
	// Either constructing the object or ValueFromObject should produce diagnostics
	_, diags2 := ptype.ValueFromObject(ctx, badObj)
	require.True(t, diagsObj.HasError() || diags2.HasError())
}

func TestNewPermissionsValue_DiagsAndMustPanic(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Missing attributes should produce diagnostics
	_, diags := NewPermissionsValue(attrTypes, map[string]attr.Value{})
	require.True(t, diags.HasError())

	// NewPermissionsValueMust should panic when diagnostics exist
	require.Panics(t, func() { NewPermissionsValueMust(attrTypes, map[string]attr.Value{}) })
}

func TestPermissionsValue_ToObjectValue_ToTerraformValue(t *testing.T) {
	ctx := context.Background()

	// Known permissions value
	pv := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("pid"),
		state:      attr.ValueStateKnown,
	}

	obj, diags := pv.ToObjectValue(ctx)
	require.False(t, diags.HasError())
	// Ensure id is present when converting back to object value
	// (we can't directly inspect basetypes.ObjectValue easily, but ensure no diags)
	_ = obj

	// ToTerraformValue should produce a tftypes.Value without error
	_, err := pv.ToTerraformValue(ctx)
	require.NoError(t, err)

	// Null and Unknown states
	pn := NewPermissionsValueNull()
	val, err := pn.ToTerraformValue(ctx)
	require.NoError(t, err)
	require.True(t, val.IsNull())

	pu := NewPermissionsValueUnknown()
	val2, err := pu.ToTerraformValue(ctx)
	require.NoError(t, err)
	require.True(t, !val2.IsKnown())
}

func TestAttributesValue_NewAndToTerraform(t *testing.T) {
	ctx := context.Background()
	// NewAttributesValue with no attribute types/values should succeed (known empty)
	av, diags := NewAttributesValue(map[string]attr.Type{}, map[string]attr.Value{})
	require.False(t, diags.HasError())
	// ToTerraformValue for known empty attributes should not error
	_, err := av.ToTerraformValue(ctx)
	require.NoError(t, err)

	// Null and Unknown
	an := NewAttributesValueNull()
	v, err := an.ToTerraformValue(ctx)
	require.NoError(t, err)
	require.True(t, v.IsNull())

	au := NewAttributesValueUnknown()
	v2, err := au.ToTerraformValue(ctx)
	require.NoError(t, err)
	require.True(t, !v2.IsKnown())
}
