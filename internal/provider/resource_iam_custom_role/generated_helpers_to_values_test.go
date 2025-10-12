package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPermissionsValue_ToTerraformValue_NullUnknownKnown(t *testing.T) {
	ctx := context.Background()

	// Null
	pNull := NewPermissionsValueNull()
	vNull, err := pNull.ToTerraformValue(ctx)
	require.NoError(t, err)
	require.True(t, vNull.IsNull())

	// Unknown
	pUnknown := NewPermissionsValueUnknown()
	vUnknown, err := pUnknown.ToTerraformValue(ctx)
	require.NoError(t, err)
	// Unknown values are not known
	require.False(t, vUnknown.IsKnown())

	// Known with valid inner values
	pv := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("sys.r.a"),
		state:      attr.ValueStateKnown,
	}

	vKnown, err := pv.ToTerraformValue(ctx)
	require.NoError(t, err)
	require.True(t, vKnown.IsKnown())
}

func TestAttributesValue_ToObjectValue_KnownNullUnknown(t *testing.T) {
	ctx := context.Background()

	// Null
	aNull := NewAttributesValueNull()
	obj, diags := aNull.ToObjectValue(ctx)
	require.False(t, diags.HasError())
	require.True(t, obj.IsNull())

	// Unknown
	aUnknown := NewAttributesValueUnknown()
	obj2, diags2 := aUnknown.ToObjectValue(ctx)
	require.False(t, diags2.HasError())
	require.True(t, obj2.IsUnknown())

	// Known
	aKnown := AttributesValue{state: attr.ValueStateKnown}
	obj3, diags3 := aKnown.ToObjectValue(ctx)
	require.False(t, diags3.HasError())
	require.False(t, obj3.IsNull())
}
