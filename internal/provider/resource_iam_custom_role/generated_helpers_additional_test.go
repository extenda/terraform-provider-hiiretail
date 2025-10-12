package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestPermissionsValue_ToAndFromTerraform_And_ObjectConversions(t *testing.T) {
	ctx := context.Background()

	// Build a known PermissionsValue via NewPermissionsValueMust
	attributeTypes := PermissionsValue{}.AttributeTypes(ctx)

	attrs := map[string]attr.Value{}

	attrs["alias"] = types.StringValue("alias-val")
	// Attributes inner object uses AttributesValue attribute types (empty map)
	attrs["attributes"] = types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
	attrs["id"] = types.StringValue("sys.res.act")

	// Use NewPermissionsValueMust to create a known value
	perm := NewPermissionsValueMust(attributeTypes, attrs)

	// ToTerraformValue should succeed for known value
	tv, err := perm.ToTerraformValue(ctx)
	require.NoError(t, err)

	// ValueFromTerraform should round-trip
	pt := PermissionsType{
		ObjectType: basetypes.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)},
	}
	gotAttr, err := pt.ValueFromTerraform(ctx, tv)
	require.NoError(t, err)

	gotPerm, ok := gotAttr.(PermissionsValue)
	require.True(t, ok)
	require.True(t, perm.Equal(gotPerm))

	// ToObjectValue should return a basetypes.ObjectValue with no diagnostics
	objVal, diags := perm.ToObjectValue(ctx)
	require.False(t, diags.HasError())
	// Passing this object into ValueFromObject should reconstruct the value
	fromObj, diags := pt.ValueFromObject(ctx, objVal)
	require.False(t, diags.HasError())
	fromPerm, ok := fromObj.(PermissionsValue)
	require.True(t, ok)
	require.True(t, perm.Equal(fromPerm))
}

func TestPermissionsValue_NullAndUnknownTerraformConversions(t *testing.T) {
	ctx := context.Background()
	pt := PermissionsType{
		ObjectType: basetypes.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)},
	}

	// Unknown => ValueFromTerraform should return Unknown PermissionsValue
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	objectType := basetypes.ObjectType{AttrTypes: attrTypes}.TerraformType(ctx)
	unknownVal := tftypes.NewValue(objectType, tftypes.UnknownValue)
	v, err := pt.ValueFromTerraform(ctx, unknownVal)
	require.NoError(t, err)
	gotUnknown, ok := v.(PermissionsValue)
	require.True(t, ok)
	require.True(t, gotUnknown.IsUnknown())

	// Null => NewPermissionsValueNull via ValueFromTerraform
	nullVal := tftypes.NewValue(objectType, nil)
	v2, err := pt.ValueFromTerraform(ctx, nullVal)
	require.NoError(t, err)
	gotNull, ok := v2.(PermissionsValue)
	require.True(t, ok)
	require.True(t, gotNull.IsNull())

	// Wrong type should return error
	strVal := tftypes.NewValue(tftypes.String, "x")
	_, err = pt.ValueFromTerraform(ctx, strVal)
	require.Error(t, err)
}

func TestPermissionsValue_NewPermissionsValue_ErrorsAndMustPanics(t *testing.T) {
	ctx := context.Background()
	attributeTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Missing attributes entry should produce diagnostics and return Unknown
	missing := map[string]attr.Value{}
	res, diags := NewPermissionsValue(attributeTypes, missing)
	require.True(t, diags.HasError())
	require.True(t, res.IsUnknown())

	// NewPermissionsValueMust should panic when diagnostics present
	require.Panics(t, func() {
		_ = NewPermissionsValueMust(attributeTypes, missing)
	})

	// Extra attribute should cause diagnostics
	extra := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"id":         types.StringValue("sys.r.a"),
		"extra":      types.StringValue("x"),
	}
	res2, diags2 := NewPermissionsValue(attributeTypes, extra)
	require.True(t, diags2.HasError())
	require.True(t, res2.IsUnknown())
}

func TestAttributesType_ValueFromTerraform_And_ToTerraform(t *testing.T) {
	ctx := context.Background()
	at := AttributesType{
		ObjectType: basetypes.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)},
	}

	// Wrong type -> error
	_, err := at.ValueFromTerraform(ctx, tftypes.NewValue(tftypes.String, "x"))
	require.Error(t, err)

	// Known/empty attributes: use AttributesValue known
	attrTypes := AttributesValue{}.AttributeTypes(ctx)
	attrs := map[string]attr.Value{}
	a := NewAttributesValueMust(attrTypes, attrs)
	tv, err := a.ToTerraformValue(ctx)
	require.NoError(t, err)

	got, err := at.ValueFromTerraform(ctx, tv)
	require.NoError(t, err)
	_, ok := got.(AttributesValue)
	require.True(t, ok)

	// Null and Unknown handling
	objectType := basetypes.ObjectType{AttrTypes: attrTypes}.TerraformType(ctx)
	u := tftypes.NewValue(objectType, tftypes.UnknownValue)
	v, err := at.ValueFromTerraform(ctx, u)
	require.NoError(t, err)
	av, ok := v.(AttributesValue)
	require.True(t, ok)
	require.True(t, av.IsUnknown())

	n := tftypes.NewValue(objectType, nil)
	v2, err := at.ValueFromTerraform(ctx, n)
	require.NoError(t, err)
	av2, ok := v2.(AttributesValue)
	require.True(t, ok)
	require.True(t, av2.IsNull())
}
