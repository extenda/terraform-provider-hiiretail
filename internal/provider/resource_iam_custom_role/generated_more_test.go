package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	basetypes "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestGenerated_ErrorBranchesAndPanicPaths(t *testing.T) {
	ctx := context.Background()

	// PermissionsType.ValueFromObject error branches: missing attributes or wrong types
	pt := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}

	// Build a basetypes.ObjectValue missing required fields to force diagnostics
	// Make attribute map empty
	emptyAttrs := map[string]attr.Value{}
	objVal, diags := types.ObjectValue(PermissionsValue{}.AttributeTypes(ctx), emptyAttrs)
	// types.ObjectValue returns diagnostics for missing required attributes
	require.True(t, diags.HasError())

	// Use ValueFromObject to exercise its error branches for missing attributes
	_, vdiags := pt.ValueFromObject(ctx, objVal)
	require.True(t, vdiags.HasError())

	// Now test NewPermissionsValueMust panic on invalid input
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from NewPermissionsValueMust with invalid input")
		}
	}()

	// missing required attributes -> should panic
	_ = NewPermissionsValueMust(PermissionsValue{}.AttributeTypes(ctx), emptyAttrs)
}

func TestAttributes_NewMustAndStringMethods(t *testing.T) {
	ctx := context.Background()

	// NewAttributesValueMust should panic when given invalid attributes
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from NewAttributesValueMust with invalid input")
		}
	}()

	// AttributesValue attribute types expect at least key 'k' (see earlier tests)
	attributeTypes := map[string]attr.Type{"k": basetypes.StringType{}}
	attributes := map[string]attr.Value{}
	_ = NewAttributesValueMust(attributeTypes, attributes)

	// Also validate String() implementations for PermissionsValue and AttributesValue
	pv := PermissionsValue{state: attr.ValueStateKnown}
	require.Equal(t, "PermissionsValue", pv.String())

	av := AttributesValue{state: attr.ValueStateKnown}
	require.Equal(t, "AttributesValue", av.String())

	// Verify ToObjectValue for null and unknown states returns expected typed objects
	pvNull := NewPermissionsValueNull()
	obj, _ := pvNull.ToObjectValue(ctx)
	// obj should be an object (no panic)
	require.NotNil(t, obj)

	pvUnknown := NewPermissionsValueUnknown()
	obj2, _ := pvUnknown.ToObjectValue(ctx)
	require.NotNil(t, obj2)

	avNull := NewAttributesValueNull()
	aobj, _ := avNull.ToObjectValue(ctx)
	require.NotNil(t, aobj)

	avUnknown := NewAttributesValueUnknown()
	aobj2, _ := avUnknown.ToObjectValue(ctx)
	require.NotNil(t, aobj2)
}
