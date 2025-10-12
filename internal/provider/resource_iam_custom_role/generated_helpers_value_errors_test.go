package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPermissionsType_ValueFromObject_MissingAlias(t *testing.T) {
	ctx := context.Background()

	// Start with the expected attribute types
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Use an ObjectNull to simulate missing attributes in the object
	objMissing := types.ObjectNull(attrTypes)
	pt := PermissionsType{ObjectType: types.ObjectType{AttrTypes: attrTypes}}
	_, diags := pt.ValueFromObject(ctx, objMissing)
	require.True(t, diags.HasError(), "expected error diagnostics for missing alias")
}

func TestNewPermissionsValue_WrongAttributeTypeProducesDiagnostic(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Provide wrong type for alias (object instead of string)
	attrs := map[string]attr.Value{
		"alias":      types.ObjectNull(map[string]attr.Type{}),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("sys.r.a"),
	}

	_, diags := NewPermissionsValue(attrTypes, attrs)
	require.True(t, diags.HasError(), "expected diagnostic for wrong attribute type")
}

func TestNewPermissionsValue_MissingAndExtraAttributes(t *testing.T) {
	ctx := context.Background()

	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Missing: pass empty attributes map
	_, diags := NewPermissionsValue(attrTypes, map[string]attr.Value{})
	require.True(t, diags.HasError(), "expected diagnostics for missing attributes")

	// Extra: provide an unexpected attribute key
	attrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("sys.r.a"),
		"extra":      types.StringValue("x"),
	}

	_, diags2 := NewPermissionsValue(attrTypes, attrs)
	require.True(t, diags2.HasError(), "expected diagnostics for extra attributes")
}

func TestNewAttributesValue_ExtraAttributeProducesDiagnostic(t *testing.T) {
	// AttributesValue has empty AttributeTypes map; providing any attribute should trigger an error
	_, diags := NewAttributesValue(AttributesValue{}.AttributeTypes(context.Background()), map[string]attr.Value{"foo": types.StringValue("bar")})
	require.True(t, diags.HasError(), "expected diagnostics for extra attribute on AttributesValue")
}
