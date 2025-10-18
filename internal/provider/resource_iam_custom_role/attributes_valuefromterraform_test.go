package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAttributesType_ValueFromTerraform_VariousPaths(t *testing.T) {
	ctx := context.Background()
	tObj := AttributesType{ObjectType: types.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}}

	// nil type -> null
	v, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(nil, nil))
	if err != nil {
		t.Fatalf("unexpected error for nil type: %v", err)
	}
	if v.IsNull() == false {
		t.Fatalf("expected null value for nil type")
	}

	// type mismatch (Number vs Object)
	_, err = tObj.ValueFromTerraform(ctx, tftypes.NewValue(tftypes.Number, 123))
	if err == nil {
		t.Fatalf("expected error for type mismatch")
	}

	// unknown and null paths
	objectType := types.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}.TerraformType(ctx)
	v2, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(objectType, tftypes.UnknownValue))
	if err != nil {
		t.Fatalf("unexpected error for unknown: %v", err)
	}
	if !v2.IsUnknown() {
		t.Fatalf("expected unknown for unknown input")
	}

	v3, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(objectType, nil))
	if err != nil {
		t.Fatalf("unexpected error for null: %v", err)
	}
	if !v3.IsNull() {
		t.Fatalf("expected null for null input")
	}

	// empty object roundtrip
	v4, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(objectType, map[string]tftypes.Value{}))
	if err != nil {
		t.Fatalf("unexpected error for empty object: %v", err)
	}
	if v4.IsUnknown() || v4.IsNull() {
		t.Fatalf("expected known value for empty object")
	}
}
