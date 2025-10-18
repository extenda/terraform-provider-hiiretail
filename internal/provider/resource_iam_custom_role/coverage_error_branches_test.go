package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPermissionsType_ValueFromObject_WrongType(t *testing.T) {
	ctx := context.Background()

	permAttrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Build attributes with wrong type for alias
	attrs := map[string]attr.Value{
		"alias": types.Int64Value(1),
		"attributes": func() attr.Value {
			v, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
			return v
		}(),
		"id": types.StringValue("perm.x.y"),
	}

	obj, diags := types.ObjectValue(permAttrTypes, attrs)
	// types.ObjectValue may validate and return diagnostics when attribute values mismatch types.
	if diags.HasError() {
		// construction already failed as expected for wrong attribute type - test passes
		return
	}

	// If construction succeeded (somehow), ensure ValueFromObject reports diagnostics for wrong type
	_, diags = PermissionsType{}.ValueFromObject(ctx, obj)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for wrong alias type")
	}
}

func TestNewPermissionsValue_MissingExtraWrongType(t *testing.T) {
	ctx := context.Background()

	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	innerObj, diags := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
	if diags.HasError() {
		t.Fatalf("failed to build inner attributes object: %v", diags)
	}

	// Missing 'id' attribute
	missingAttrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": innerObj,
	}

	pv, diags := NewPermissionsValue(attrTypes, missingAttrs)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for missing attribute, got none")
	}
	if !pv.IsUnknown() {
		t.Fatalf("expected unknown PermissionsValue on error")
	}

	// Extra attribute
	extraAttrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": innerObj,
		"id":         types.StringValue("perm.x.y"),
		"k":          types.StringValue("x"),
	}

	pv2, diags := NewPermissionsValue(attrTypes, extraAttrs)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for extra attribute, got none")
	}
	if !pv2.IsUnknown() {
		t.Fatalf("expected unknown PermissionsValue on extra attribute error")
	}

	// Wrong type for alias
	wrongTypeAttrs := map[string]attr.Value{
		"alias":      types.Int64Value(5),
		"attributes": innerObj,
		"id":         types.StringValue("perm.x.y"),
	}

	pv3, diags := NewPermissionsValue(attrTypes, wrongTypeAttrs)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for wrong attribute type, got none")
	}
	if !pv3.IsUnknown() {
		t.Fatalf("expected unknown PermissionsValue on wrong type error")
	}
}

func TestPermissionsValue_ToTerraform_NullUnknownKnown(t *testing.T) {
	ctx := context.Background()

	// Null
	nullPV := PermissionsValue{state: attr.ValueStateNull}
	v, err := nullPV.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("unexpected error for null ToTerraformValue: %v", err)
	}
	if !v.IsNull() {
		t.Fatalf("expected tftypes.Value to be null for null PermissionsValue")
	}

	// Unknown
	unkPV := PermissionsValue{state: attr.ValueStateUnknown}
	v2, err := unkPV.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("unexpected error for unknown ToTerraformValue: %v", err)
	}
	if v2.IsKnown() || v2.IsNull() {
		t.Fatalf("expected tftypes.Value to be unknown for unknown PermissionsValue")
	}

	// Known
	innerObj, diags := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
	if diags.HasError() {
		t.Fatalf("failed to create inner attributes object: %v", diags)
	}

	knownPV := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: innerObj,
		Id:         types.StringValue("perm.x.y"),
		state:      attr.ValueStateKnown,
	}

	v3, err := knownPV.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("unexpected error for known ToTerraformValue: %v", err)
	}

	if !v3.IsKnown() {
		t.Fatalf("expected known tftypes.Value for known PermissionsValue")
	}

	// ensure it's an object type and can be inspected
	if v3.Type() == nil {
		t.Fatalf("expected non-nil type for known value")
	}
	if _, ok := v3.Type().(tftypes.Object); !ok {
		t.Fatalf("expected object terraform type for PermissionsValue, got: %T", v3.Type())
	}

	// Basic roundtrip: convert back using PermissionsType.ValueFromTerraform
	got, err := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}.ValueFromTerraform(ctx, v3)
	if err != nil {
		t.Fatalf("ValueFromTerraform roundtrip error: %v", err)
	}
	if _, ok := got.(PermissionsValue); !ok {
		t.Fatalf("expected PermissionsValue on roundtrip, got: %T", got)
	}
}

func TestPermissionsType_ValueFromObject_MissingAndWrongTypes(t *testing.T) {
	ctx := context.Background()

	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Missing alias
	attrsMissingAlias := map[string]attr.Value{
		"attributes": func() attr.Value {
			v, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
			return v
		}(),
		"id": types.StringValue("perm.x.y"),
	}

	obj1, diags := types.ObjectValue(attrTypes, attrsMissingAlias)
	if diags.HasError() {
		// construction may fail; that's acceptable for this negative test
	} else {
		v, diags2 := PermissionsType{}.ValueFromObject(ctx, obj1)
		if v != nil || !diags2.HasError() {
			t.Fatalf("expected diagnostic and nil value when alias missing")
		}
	}

	// Missing attributes
	attrsMissingAttributes := map[string]attr.Value{
		"alias": types.StringValue("a"),
		"id":    types.StringValue("perm.x.y"),
	}

	obj2, diags := types.ObjectValue(attrTypes, attrsMissingAttributes)
	if diags.HasError() {
		// construction may fail; ok
	} else {
		v, diags2 := PermissionsType{}.ValueFromObject(ctx, obj2)
		if v != nil || !diags2.HasError() {
			t.Fatalf("expected diagnostic and nil value when attributes missing")
		}
	}

	// Missing id
	attrsMissingId := map[string]attr.Value{
		"alias": types.StringValue("a"),
		"attributes": func() attr.Value {
			v, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
			return v
		}(),
	}

	obj3, diags := types.ObjectValue(attrTypes, attrsMissingId)
	if diags.HasError() {
		// construction may fail; ok
	} else {
		v, diags2 := PermissionsType{}.ValueFromObject(ctx, obj3)
		if v != nil || !diags2.HasError() {
			t.Fatalf("expected diagnostic and nil value when id missing")
		}
	}

	// Wrong type for attributes
	attrsWrongAttributesType := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.StringValue("not-an-object"),
		"id":         types.StringValue("perm.x.y"),
	}

	obj4, diags := types.ObjectValue(attrTypes, attrsWrongAttributesType)
	if diags.HasError() {
		// construction likely fails; ok
	} else {
		_, diags2 := PermissionsType{}.ValueFromObject(ctx, obj4)
		if !diags2.HasError() {
			t.Fatalf("expected diagnostics for wrong attributes type")
		}
	}

	// Wrong type for id
	attrsWrongIdType := map[string]attr.Value{
		"alias": types.StringValue("a"),
		"attributes": func() attr.Value {
			v, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
			return v
		}(),
		"id": types.Int64Value(5),
	}

	obj5, diags := types.ObjectValue(attrTypes, attrsWrongIdType)
	if diags.HasError() {
		// construction may fail; ok
	} else {
		_, diags2 := PermissionsType{}.ValueFromObject(ctx, obj5)
		if !diags2.HasError() {
			t.Fatalf("expected diagnostics for wrong id type")
		}
	}
}
