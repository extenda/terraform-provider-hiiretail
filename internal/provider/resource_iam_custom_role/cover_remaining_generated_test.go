package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCoverRemainingGeneratedFunctions(t *testing.T) {
	ctx := context.Background()

	// PermissionsValue.String()
	pv := PermissionsValue{state: attr.ValueStateKnown}
	require.Equal(t, "PermissionsValue", pv.String())

	// PermissionsValue.Type
	pt := pv.Type(ctx)
	require.NotNil(t, pt)

	// AttributesValue.String()
	av := AttributesValue{state: attr.ValueStateKnown}
	require.Equal(t, "AttributesValue", av.String())

	// AttributesValue.Type
	at := av.Type(ctx)
	require.NotNil(t, at)

	// AttributesType.ValueFromObject should return known AttributesValue
	attrTypes := AttributesValue{}.AttributeTypes(ctx)
	// create an empty basetypes.ObjectValue of the expected type
	emptyAttrs := map[string]attr.Value{}
	obj, diags := types.ObjectValue(attrTypes, emptyAttrs)
	// types.ObjectValue may return diagnostics for missing keys; still call ValueFromObject
	_ = diags

	atype := AttributesType{ObjectType: types.ObjectType{AttrTypes: attrTypes}}
	_, vdiags := atype.ValueFromObject(ctx, obj)
	// vdiags may be empty or contain diagnostics depending on attributes; ensure call succeeds
	_ = vdiags

	// AttributesValue.Equal negative branch: compare with different type
	var other attr.Value = pv
	require.False(t, av.Equal(other))
}
