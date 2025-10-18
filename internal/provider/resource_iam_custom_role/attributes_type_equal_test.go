package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestAttributesType_Equal(t *testing.T) {
	ctx := context.Background()

	// Two AttributesType with same (empty) AttrTypes should be equal
	at1 := AttributesType{ObjectType: basetypes.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}}
	at2 := AttributesType{ObjectType: basetypes.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}}
	require.True(t, at1.Equal(at2))

	// Different attribute types should not be equal
	otherAttrTypes := map[string]attr.Type{"x": basetypes.StringType{}}
	at3 := AttributesType{ObjectType: basetypes.ObjectType{AttrTypes: otherAttrTypes}}
	require.False(t, at1.Equal(at3))
}
