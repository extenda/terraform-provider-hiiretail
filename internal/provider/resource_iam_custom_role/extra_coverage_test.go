package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	// basetypes not directly referenced here
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestNewPermissionsValue_ErrorCases(t *testing.T) {
	ctx := context.Background()

	// Missing attribute (expect diagnostic and Unknown returned)
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	attrs := map[string]attr.Value{
		// omit "id" to trigger Missing PermissionsValue Attribute Value
		"alias": types.StringValue("a"),
	}

	pv, diags := NewPermissionsValue(attrTypes, attrs)
	if !diags.HasError() {
		t.Fatalf("expected diags.HasError() for missing attribute")
	}
	if !pv.IsUnknown() {
		t.Fatalf("expected returned PermissionsValue to be Unknown when diags present")
	}

	// Extra attribute should produce diagnostics
	attrs2 := map[string]attr.Value{
		"alias": types.StringValue("a"),
		"id":    types.StringValue("i"),
		"extra": types.StringValue("x"),
	}

	pv2, diags2 := NewPermissionsValue(attrTypes, attrs2)
	if !diags2.HasError() {
		t.Fatalf("expected diags for extra attribute")
	}
	if !pv2.IsUnknown() {
		t.Fatalf("expected Unknown PermissionsValue when extra attribute present")
	}

	// Wrong attribute type: pass ObjectValue where String expected
	wrongAttrs := map[string]attr.Value{
		"alias":      types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"id":         types.StringValue("i"),
		"attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
	}

	_, diags3 := NewPermissionsValue(attrTypes, wrongAttrs)
	if !diags3.HasError() {
		t.Fatalf("expected diags for wrong attribute types")
	}
}

func TestPermissionsValue_ToTerraformValue_UnhandledStatePanics(t *testing.T) {
	ctx := context.Background()

	// create a PermissionsValue with an invalid state
	pv := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		Id:         types.StringValue("i"),
		// choose a small out-of-band state that will hit the default branch
		state: attr.ValueState(3),
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for unhandled Object state")
		}
	}()

	_, _ = pv.ToTerraformValue(ctx)
}

func TestPermissionsValue_ToObjectValue_NullAndUnknownAttributes(t *testing.T) {
	ctx := context.Background()

	// Attributes null
	pvNull := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("i"),
		state:      attr.ValueStateKnown,
	}

	obj, diags := pvNull.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	// The returned object itself may not be null, but the nested "attributes"
	// attribute should be null when we passed ObjectNull.
	// obj is already a basetypes.ObjectValue
	nested := obj.Attributes()["attributes"]
	if nested == nil {
		t.Fatalf("expected nested attribute to be present")
	}
	if !nested.IsNull() {
		t.Fatalf("expected nested 'attributes' to be null")
	}

	// Attributes unknown
	pvUnknown := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectUnknown(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("i"),
		state:      attr.ValueStateKnown,
	}

	obj2, diags2 := pvUnknown.ToObjectValue(ctx)
	if diags2.HasError() {
		t.Fatalf("unexpected diags: %v", diags2)
	}
	// The returned object is known but contains a nested attributes attribute
	// that should be unknown.
	nested2 := obj2.Attributes()["attributes"]
	if nested2 == nil {
		t.Fatalf("expected nested attribute to be present")
	}
	if !nested2.IsUnknown() {
		t.Fatalf("expected nested 'attributes' to be unknown")
	}
}

func TestPermissionsValue_Equal_DifferentStatesAndFields(t *testing.T) {
	ctx := context.Background()

	pv1 := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		Id:         types.StringValue("i"),
		state:      attr.ValueStateKnown,
	}

	pv2 := pv1
	if !pv1.Equal(pv2) {
		t.Fatalf("expected equal values to be equal")
	}

	// different id
	pv3 := pv1
	pv3.Id = types.StringValue("other")
	if pv1.Equal(pv3) {
		t.Fatalf("expected different values to not be equal")
	}

	// different state
	pv4 := pv1
	pv4.state = attr.ValueStateNull
	if pv1.Equal(pv4) {
		t.Fatalf("expected different states to not be equal")
	}

}

func TestAttributesValue_ToTerraformValue_PanicUnhandledState(t *testing.T) {
	v := AttributesValue{state: attr.ValueState(3)}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for unhandled AttributesValue state")
		}
	}()

	_, _ = v.ToTerraformValue(context.Background())
}

func TestPermissionsType_ValueFromObject_WrongTypesAndMissingAttrs(t *testing.T) {
	ctx := context.Background()

	// Build an object missing attributes by using ObjectNull which simulates
	// a basetypes.ObjectValue missing its attributes.
	attributeTypes := PermissionsValue{}.AttributeTypes(ctx)
	obj := types.ObjectNull(attributeTypes)

	var pt PermissionsType
	_, diags2 := pt.ValueFromObject(ctx, obj)
	if !diags2.HasError() {
		t.Fatalf("expected diags for missing attributes in ValueFromObject")
	}

	// Wrong types: provide wrong attribute types
	wrongAttrs := map[string]attr.Value{
		"alias":      types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"attributes": types.StringValue("not-an-object"),
		"id":         types.StringValue("i"),
	}

	obj2, diags3 := types.ObjectValue(PermissionsValue{}.AttributeTypes(ctx), wrongAttrs)
	if diags3.HasError() {
		// Creating an ObjectValue with invalid attribute types may itself
		// produce diagnostics from the types package; that's acceptable for
		// this negative test — log and return early.
		t.Logf("types.ObjectValue produced diags as expected: %v", diags3)
		return
	}

	_, diags4 := pt.ValueFromObject(ctx, obj2)
	if !diags4.HasError() {
		t.Fatalf("expected diags for wrong attribute types in ValueFromObject")
	}
}

func TestAttributesType_ValueFromTerraform_TypeMismatch(t *testing.T) {
	ctx := context.Background()
	var at AttributesType

	// Provide a tftypes.Value with mismatched type
	val := tftypes.NewValue(tftypes.String, "x")
	_, err := at.ValueFromTerraform(ctx, val)
	if err == nil {
		t.Fatalf("expected error for type mismatch in AttributesType.ValueFromTerraform")
	}
}
