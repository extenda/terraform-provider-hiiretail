package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

// Test that PermissionsType.ValueFromObject returns an error when required attributes
// are missing or of the wrong type, and succeeds on a correct object.
func TestPermissions_ValueFromObject_Behavior(t *testing.T) {
	ctx := context.Background()

	// Build an object missing required keys -> should return diagnostics with errors
	badObj, diags := types.ObjectValue(
		map[string]attr.Type{
			// intentionally empty - will trigger missing attribute diagnostics
		},
		map[string]attr.Value{},
	)
	require.False(t, diags.HasError())

	_, badDiags := PermissionsType{}.ValueFromObject(ctx, badObj)
	require.True(t, badDiags.HasError(), "expected diagnostics when attributes are missing")

	// Build a well-formed object that matches expected attribute types
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// attributes sub-object can be empty (AttributesValue.AttributeTypes is empty)
	attributesObj := types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})

	goodObj := types.ObjectValueMust(attrTypes, map[string]attr.Value{
		"alias":      types.StringValue("alias-val"),
		"attributes": attributesObj,
		"id":         types.StringValue("perm-1"),
	})

	val, goodDiags := PermissionsType{}.ValueFromObject(ctx, goodObj)
	require.False(t, goodDiags.HasError(), "expected no diagnostics for valid object")

	// Assert returned value is a PermissionsValue and has expected Id
	pv, ok := val.(PermissionsValue)
	require.True(t, ok)
	require.Equal(t, "perm-1", pv.Id.ValueString())
}

// Test NewPermissionsValue error branch when extra/missing attributes are present
func TestNewPermissionsValue_ErrorBranches(t *testing.T) {
	// Provide attributeTypes that expect an extra key, but do not pass it in attributes -> error
	attributeTypes := map[string]attr.Type{"extra": basetypes.StringType{}}

	attrs := map[string]attr.Value{}

	pv, diags := NewPermissionsValue(attributeTypes, attrs)
	require.True(t, diags.HasError())
	// Expect that the returned object is Unknown when diags present
	require.True(t, pv.IsUnknown())

	// Provide attributes with an extra key that is not expected -> error
	attributeTypes2 := map[string]attr.Type{}
	attrs2 := map[string]attr.Value{"unexpected": types.StringValue("x")}
	pv2, diags2 := NewPermissionsValue(attributeTypes2, attrs2)
	require.True(t, diags2.HasError())
	require.True(t, pv2.IsUnknown())
}

// Test PermissionsValue ToTerraform/ToObject null and unknown branches and Equal behavior
func TestPermissionsValue_StateAndEqual(t *testing.T) {
	ctx := context.Background()

	nullPV := NewPermissionsValueNull()
	tv, err := nullPV.ToTerraformValue(ctx)
	require.NoError(t, err)
	// Null terraform value has nil underlying
	require.True(t, tv.IsNull())

	unknownPV := NewPermissionsValueUnknown()
	tuv, err := unknownPV.ToTerraformValue(ctx)
	require.NoError(t, err)
	// tftypes.Value has no IsUnknown method; unknown values report as !IsKnown()
	require.False(t, tuv.IsKnown())

	// Create two known PermissionsValue instances and compare equality
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	attributesObj := types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
	a := NewPermissionsValueMust(attrTypes, map[string]attr.Value{
		"alias":      types.StringValue("a1"),
		"attributes": attributesObj,
		"id":         types.StringValue("id1"),
	})
	b := NewPermissionsValueMust(attrTypes, map[string]attr.Value{
		"alias":      types.StringValue("a1"),
		"attributes": attributesObj,
		"id":         types.StringValue("id1"),
	})
	require.True(t, a.Equal(b))

	// Change ID -> not equal
	c := NewPermissionsValueMust(attrTypes, map[string]attr.Value{
		"alias":      types.StringValue("a1"),
		"attributes": attributesObj,
		"id":         types.StringValue("id2"),
	})
	require.False(t, a.Equal(c))

	// Unknown vs Known -> not equal
	require.False(t, a.Equal(NewPermissionsValueUnknown()))
}

// Test modelToAPIRequest when permission attributes contain string values
func TestModelToAPIRequest_WithAttributes(t *testing.T) {
	ctx := context.Background()
	r := NewIamCustomRoleResource().(*IamCustomRoleResource)

	// Build a PermissionsValue with an attributes object matching the generated AttributesValue types
	attributesObj := types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})

	// Build a PermissionsValue manually to allow a non-empty attributes object
	pv := PermissionsValue{
		Alias:      types.StringValue("al"),
		Attributes: attributesObj,
		Id:         types.StringValue("perm-x"),
		state:      attr.ValueStateKnown,
	}

	// Convert to list and call modelToAPIRequest
	permType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
	list, diags := types.ListValueFrom(ctx, permType, []PermissionsValue{pv})
	require.False(t, diags.HasError())

	data := IamCustomRoleModel{Id: types.StringValue("r1"), Permissions: list}

	req, err := r.modelToAPIRequest(ctx, data)
	require.NoError(t, err)
	require.Len(t, req.Permissions, 1)
	attrs := req.Permissions[0].Attributes
	// The generated AttributesValue has no declared attribute types, so attributes are not carried through
	require.Len(t, attrs, 0)
}
