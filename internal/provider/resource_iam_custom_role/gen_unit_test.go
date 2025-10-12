package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	basetypes "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestGenerated_PermissionsValue_Conversions(t *testing.T) {
	ctx := context.Background()

	// Test NewPermissionsValueNull/Unknown
	n := NewPermissionsValueNull()
	require.True(t, n.IsNull())
	u := NewPermissionsValueUnknown()
	require.True(t, u.IsUnknown())

	// Build a known PermissionsValue using simple types
	pv := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("perm.id"),
		state:      attr.ValueStateKnown,
	}

	// ToTerraformValue (Known)
	_, err := pv.ToTerraformValue(ctx)
	require.NoError(t, err)

	// ToObjectValue
	obj, diags := pv.ToObjectValue(ctx)
	require.False(t, diags.HasError())
	require.NotNil(t, obj)

	// Equal with itself
	require.True(t, pv.Equal(pv))

	// Type and AttributeTypes
	tp := pv.Type(ctx)
	require.Equal(t, "PermissionsType", tp.String())
	ats := pv.AttributeTypes(ctx)
	require.Contains(t, ats, "id")

	// Test Null and Unknown ToTerraformValue
	pv2 := PermissionsValue{state: attr.ValueStateNull}
	_, err = pv2.ToTerraformValue(ctx)
	require.NoError(t, err)
	pv3 := PermissionsValue{state: attr.ValueStateUnknown}
	_, err = pv3.ToTerraformValue(ctx)
	require.NoError(t, err)
}

func TestGenerated_AttributesValue_Conversions(t *testing.T) {
	ctx := context.Background()

	// AttributesValue Known/Null/Unknown
	aKnown := AttributesValue{state: attr.ValueStateKnown}
	_, err := aKnown.ToTerraformValue(ctx)
	require.NoError(t, err)

	aNull := NewAttributesValueNull()
	_, err = aNull.ToTerraformValue(ctx)
	require.NoError(t, err)

	aUnknown := NewAttributesValueUnknown()
	_, err = aUnknown.ToTerraformValue(ctx)
	require.NoError(t, err)

	// ToObjectValue for known should produce an object value
	obj, diags := aKnown.ToObjectValue(ctx)
	require.False(t, diags.HasError())
	require.NotNil(t, obj)
}

func TestPermissions_ValueFromObject_And_ValueFromTerraform(t *testing.T) {
	ctx := context.Background()

	// Build basetypes object for PermissionsValue
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	attrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("perm.id"),
	}

	// Create basetypes.ObjectValue via types.ObjectValueMust
	objVal, diags := types.ObjectValue(attrTypes, attrs)
	require.False(t, diags.HasError())

	// ValueFromObject
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	pvVal, diags := permType.ValueFromObject(ctx, objVal)
	require.False(t, diags.HasError())

	// pvVal should be a PermissionsValue
	pv, ok := pvVal.(PermissionsValue)
	require.True(t, ok)

	// NewPermissionsValue using attribute map
	n, diags := NewPermissionsValue(attrTypes, attrs)
	require.False(t, diags.HasError())
	require.Equal(t, n.Id.ValueString(), pv.Id.ValueString())

	// NewPermissionsValueMust should not panic for valid input
	m := NewPermissionsValueMust(attrTypes, attrs)
	require.Equal(t, m.Id.ValueString(), "perm.id")

	// ToTerraformValue on PermissionsValue (known)
	tfVal, err := pv.ToTerraformValue(ctx)
	require.NoError(t, err)

	// Call ValueFromTerraform using the same value with a properly configured PermissionsType
	v, err := permType.ValueFromTerraform(ctx, tfVal)
	require.NoError(t, err)
	// Should return an attr.Value (PermissionsValue)
	_, ok = v.(PermissionsValue)
	require.True(t, ok)
}

func TestNewPermissionsValue_Errors(t *testing.T) {
	ctx := context.Background()
	// attributeTypes expects alias, attributes, id
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Provide attributes map missing 'id' to trigger diagnostics
	attrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		// "id" is intentionally missing
	}

	_, diags := NewPermissionsValue(attrTypes, attrs)
	require.True(t, diags.HasError(), "Expected diagnostics when missing required attribute")
}

func TestNewAttributesValue_Errors(t *testing.T) {
	// attributeTypes requires a key 'k'
	attributeTypes := map[string]attr.Type{"k": basetypes.StringType{}}
	// Provide attributes map missing 'k'
	attributes := map[string]attr.Value{}

	_, diags := NewAttributesValue(attributeTypes, attributes)
	require.True(t, diags.HasError(), "Expected diagnostics when missing required attribute for AttributesValue")
}
