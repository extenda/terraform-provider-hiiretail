package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAttributesType_ValueFromTerraform_NilTypeAndMismatchAndRoundTrip(t *testing.T) {
	ctx := context.Background()
	tObj := AttributesType{}

	// in.Type() == nil -> should return Null
	v, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(nil, nil))
	if err != nil {
		t.Fatalf("unexpected error for nil type: %v", err)
	}

	if !v.(AttributesValue).IsNull() {
		t.Fatalf("expected null AttributesValue for nil type")
	}

	// Type mismatch: pass a number where object expected
	_, err = tObj.ValueFromTerraform(ctx, tftypes.NewValue(tftypes.Number, 123))
	if err == nil {
		t.Fatalf("expected error for type mismatch")
	}

	// Unknown value -> should return Unknown
	objectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
	v2, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(objectType, tftypes.UnknownValue))
	if err != nil {
		t.Fatalf("unexpected error for unknown value: %v", err)
	}
	if !v2.(AttributesValue).IsUnknown() {
		t.Fatalf("expected unknown AttributesValue for UnknownValue")
	}

	// Null value -> should return Null
	v3, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(objectType, nil))
	if err != nil {
		t.Fatalf("unexpected error for null object value: %v", err)
	}
	if !v3.(AttributesValue).IsNull() {
		t.Fatalf("expected null AttributesValue for null object value")
	}

	// Success path: empty object map
	v4, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(objectType, map[string]tftypes.Value{}))
	if err != nil {
		t.Fatalf("unexpected error for empty object value: %v", err)
	}
	if v4.(AttributesValue).IsNull() || v4.(AttributesValue).IsUnknown() {
		t.Fatalf("expected known AttributesValue on success path")
	}
}

func TestPermissionsType_ValueFromObject_Success(t *testing.T) {
	ctx := context.Background()

	// Build a correct basetypes.ObjectValue for PermissionsType
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Build the attributes for the object
	attributes := map[string]attr.Value{}

	// alias and id are string basetypes; attributes is an object (empty)
	attributes["alias"] = types.StringValue("a")
	attributes["id"] = types.StringValue("perm.x.y")
	// For attributes field, create an empty object value matching AttributesValue types
	attributes["attributes"] = types.ObjectValueMust(
		AttributesValue{}.AttributeTypes(ctx),
		map[string]attr.Value{},
	)

	obj := types.ObjectValueMust(attrTypes, attributes)

	val, diags := PermissionsType{}.ValueFromObject(ctx, obj)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics from ValueFromObject: %v", diags)
	}

	pv, ok := val.(PermissionsValue)
	if !ok {
		t.Fatalf("expected PermissionsValue, got %T", val)
	}

	if pv.IsNull() || pv.IsUnknown() {
		t.Fatalf("expected known PermissionsValue")
	}
}

func TestPermissionsValue_Equal_OtherTypeAndState(t *testing.T) {
	ctx := context.Background()

	// Known value
	p1 := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		Id:         types.StringValue("perm.x.y"),
		state:      attr.ValueStateKnown,
	}

	// Different type should return false
	if p1.Equal(types.StringValue("nope")) {
		t.Fatalf("expected Equal to return false when comparing to other type")
	}

	// Different state -> false
	p2 := p1
	p2.state = attr.ValueStateUnknown
	if p1.Equal(p2) {
		t.Fatalf("expected Equal to return false when states differ")
	}

	// Same known values -> true
	p3 := p1
	if !p1.Equal(p3) {
		t.Fatalf("expected Equal to return true for identical values")
	}
}
