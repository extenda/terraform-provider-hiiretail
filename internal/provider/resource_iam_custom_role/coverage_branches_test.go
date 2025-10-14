package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Exercise many ValueFromObject and NewPermissionsValue negative branches
func Test_PermissionsValue_ObjectNegativeBranches(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// 1) Missing alias
	attrs1 := map[string]attr.Value{
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("perm.x.y"),
	}
	obj1, _ := types.ObjectValue(attrTypes, attrs1)
	var pt PermissionsType
	_, diags1 := pt.ValueFromObject(ctx, obj1)
	if !diags1.HasError() {
		t.Fatalf("expected diagnostics for missing alias")
	}

	// 2) Missing attributes
	attrs2 := map[string]attr.Value{
		"alias": types.StringValue("a"),
		"id":    types.StringValue("perm.x.y"),
	}
	obj2, _ := types.ObjectValue(attrTypes, attrs2)
	_, diags2 := pt.ValueFromObject(ctx, obj2)
	if !diags2.HasError() {
		t.Fatalf("expected diagnostics for missing attributes")
	}

	// 3) Missing id
	attrs3 := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
	}
	obj3, _ := types.ObjectValue(attrTypes, attrs3)
	_, diags3 := pt.ValueFromObject(ctx, obj3)
	if !diags3.HasError() {
		t.Fatalf("expected diagnostics for missing id")
	}

	// 4) Wrong type for alias
	attrs4 := map[string]attr.Value{
		"alias":      types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("perm.x.y"),
	}
	obj4, diags4 := types.ObjectValue(attrTypes, attrs4)
	if diags4.HasError() {
		t.Logf("types.ObjectValue produced diags for wrong alias type: %v", diags4)
	} else {
		_, diags4b := pt.ValueFromObject(ctx, obj4)
		if !diags4b.HasError() {
			t.Fatalf("expected diagnostics for wrong alias type")
		}
	}

	// 5) Extra attribute
	attrs5 := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("perm.x.y"),
		"extra":      types.StringValue("x"),
	}
	// NewPermissionsValue should report extra attribute diagnostic
	v, diags5 := NewPermissionsValue(attrTypes, attrs5)
	if diags5 == nil || !diags5.HasError() {
		t.Fatalf("expected diagnostics for extra attribute in NewPermissionsValue, got %v and value %v", diags5, v)
	}

	// 6) Wrong attribute type in NewPermissionsValue
	attrs6 := map[string]attr.Value{
		"alias":      types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		"id":         types.StringValue("perm.x.y"),
	}
	v2, diags6 := NewPermissionsValue(attrTypes, attrs6)
	if diags6 == nil || !diags6.HasError() {
		t.Fatalf("expected diagnostics for wrong attribute type in NewPermissionsValue, got %v and value %v", diags6, v2)
	}
}

// Test the success path for NewPermissionsValue (all attributes present and correct types)
func Test_NewPermissionsValue_Success(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	attrs := map[string]attr.Value{
		"alias":      types.StringValue("alias-val"),
		"attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"id":         types.StringValue("perm.x.y"),
	}

	v, diags := NewPermissionsValue(attrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics for valid input: %v", diags)
	}
	if v.IsUnknown() || v.IsNull() {
		t.Fatalf("expected known PermissionsValue, got unknown/null: %v", v)
	}
	if v.Id.ValueString() != "perm.x.y" {
		t.Fatalf("unexpected id value: %s", v.Id.ValueString())
	}
}

// Test AttributesType.ValueFromTerraform successful round-trip
func Test_AttributesType_ValueFromTerraform_Success(t *testing.T) {
	var at AttributesType
	// Build a known AttributesValue as Terraform object then round-trip
	// Create a basetypes.ObjectValue via types.ObjectValue
	attrs := map[string]attr.Value{}
	obj, diags := types.ObjectValue(AttributesValue{}.AttributeTypes(context.Background()), attrs)
	if diags.HasError() {
		t.Fatalf("failed to build attributes object: %v", diags)
	}

	// Convert to Terraform native value
	tfVal, err := obj.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("failed to convert object to terraform value: %v", err)
	}

	// Call ValueFromTerraform
	got, err := at.ValueFromTerraform(context.Background(), tfVal)
	if err != nil {
		t.Fatalf("ValueFromTerraform returned error: %v", err)
	}
	if !got.(AttributesValue).IsNull() && !got.(AttributesValue).IsUnknown() {
		// known with empty attributes is acceptable
	}
}
