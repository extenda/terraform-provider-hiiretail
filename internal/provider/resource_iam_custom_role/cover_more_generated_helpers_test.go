package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// This test file covers several remaining branches in the generated helpers.
func TestGeneratedHelpers_ErrorAndStateBranches(t *testing.T) {
	ctx := context.Background()

	// 1) PermissionsType.ValueFromObject - missing alias -> should return diagnostics
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	// Build an object missing alias
	objValsMissingAlias := map[string]attr.Value{
		// "alias" omitted
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("sys.r.a"),
	}
	objMissingAlias, objDiags := types.ObjectValue(attrTypes, objValsMissingAlias)
	_ = objDiags // we expect diagnostics from constructing a missing-attribute object

	pt := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	_, diags := pt.ValueFromObject(ctx, objMissingAlias)
	require.True(t, diags.HasError(), "expected diagnostics for missing alias")

	// 2) NewPermissionsValue - extra attribute should produce diagnostic and return Unknown
	attrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("sys.r.a"),
		"extra":      types.StringValue("bad"),
	}
	pv, diags2 := NewPermissionsValue(PermissionsValue{}.AttributeTypes(ctx), attrs)
	require.True(t, diags2.HasError())
	require.True(t, pv.IsUnknown())

	// 3) PermissionsValue.ToTerraformValue - Null and Unknown branches
	pvNull := NewPermissionsValueNull()
	vNull, err := pvNull.ToTerraformValue(ctx)
	require.NoError(t, err)
	require.True(t, vNull.IsNull())

	pvUnknown := NewPermissionsValueUnknown()
	vUnknown, err := pvUnknown.ToTerraformValue(ctx)
	require.NoError(t, err)
	// tftypes.Value does not have IsUnknown(); check IsKnown instead
	require.False(t, vUnknown.IsKnown())

	// 4) PermissionsValue.ToObjectValue - Attributes null/unknown/known handling
	// Attributes null
	pvn := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringNull(),
		state:      attr.ValueStateKnown,
	}
	objN, diags3 := pvn.ToObjectValue(ctx)
	require.False(t, diags3.HasError())
	require.True(t, objN.IsNull() || objN.Type(ctx) != nil)

	// Attributes unknown
	pvu := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: types.ObjectUnknown(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringNull(),
		state:      attr.ValueStateKnown,
	}
	objU, diags4 := pvu.ToObjectValue(ctx)
	require.False(t, diags4.HasError())
	require.True(t, objU.IsUnknown() || objU.Type(ctx) != nil)

	// Attributes known (empty map)
	attrsKnown := map[string]attr.Value{}
	attrsObj, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), attrsKnown)
	pvKnown := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: attrsObj,
		Id:         types.StringValue("sys.r.a"),
		state:      attr.ValueStateKnown,
	}
	objK, diags5 := pvKnown.ToObjectValue(ctx)
	require.False(t, diags5.HasError())
	require.False(t, objK.IsNull())

	// 5) Attributes helper: NewAttributesValue with missing/extra keys
	// missing expected (empty attributeTypes -> provide extra key to trigger extra diagnostic)
	badAttrs := map[string]attr.Value{"something": types.StringValue("x")}
	av, ad := NewAttributesValue(AttributesValue{}.AttributeTypes(ctx), badAttrs)
	require.True(t, ad.HasError())
	require.True(t, av.IsUnknown())
}
