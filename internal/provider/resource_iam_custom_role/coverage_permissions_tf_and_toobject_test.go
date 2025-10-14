package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPermissionsType_ValueFromTerraform_NilMismatchUnknownNullSuccess(t *testing.T) {
	ctx := context.Background()
	tObj := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}

	// in.Type() == nil -> Null
	v, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(nil, nil))
	if err != nil {
		t.Fatalf("unexpected error for nil type: %v", err)
	}
	if !v.(PermissionsValue).IsNull() {
		t.Fatalf("expected null PermissionsValue for nil type")
	}

	// Type mismatch
	_, err = tObj.ValueFromTerraform(ctx, tftypes.NewValue(tftypes.Number, 5))
	if err == nil {
		t.Fatalf("expected error for type mismatch")
	}

	// Unknown value (use the Type reported by the PermissionsType)
	var tfType tftypes.Type = tObj.TerraformType(ctx)
	v2, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(tfType, tftypes.UnknownValue))
	if err != nil {
		t.Fatalf("unexpected error for unknown value: %v", err)
	}
	if !v2.(PermissionsValue).IsUnknown() {
		t.Fatalf("expected unknown PermissionsValue for UnknownValue")
	}

	// Null value
	v3, err := tObj.ValueFromTerraform(ctx, tftypes.NewValue(tfType, nil))
	if err != nil {
		t.Fatalf("unexpected error for null object value: %v", err)
	}
	if !v3.(PermissionsValue).IsNull() {
		t.Fatalf("expected null PermissionsValue for null object value")
	}

	// Success path: build proper tf object using reported Terraform type
	// reuse tfType declared above
	objTf := tfType.(tftypes.Object)

	vals := map[string]tftypes.Value{
		"alias":      tftypes.NewValue(objTf.AttributeTypes["alias"], "alias"),
		"attributes": tftypes.NewValue(objTf.AttributeTypes["attributes"], map[string]tftypes.Value{}),
		"id":         tftypes.NewValue(objTf.AttributeTypes["id"], "pid"),
	}

	tv := tftypes.NewValue(tfType, vals)
	v4, err := tObj.ValueFromTerraform(ctx, tv)
	if err != nil {
		t.Fatalf("unexpected error for success path: %v", err)
	}
	if v4.(PermissionsValue).IsNull() || v4.(PermissionsValue).IsUnknown() {
		t.Fatalf("expected known PermissionsValue on success path")
	}
}

func TestPermissionsValue_ToObjectValue_Branches(t *testing.T) {
	ctx := context.Background()

	// Null
	nullPV := PermissionsValue{state: attr.ValueStateNull}
	ov, diags := nullPV.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diags for null ToObjectValue: %v", diags)
	}
	if !ov.IsNull() {
		t.Fatalf("expected object null for null PermissionsValue")
	}

	// Unknown
	unkPV := PermissionsValue{state: attr.ValueStateUnknown}
	ov2, diags := unkPV.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diags for unknown ToObjectValue: %v", diags)
	}
	if !ov2.IsUnknown() {
		t.Fatalf("expected object unknown for unknown PermissionsValue")
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

	ov3, diags := knownPV.ToObjectValue(ctx)
	if diags.HasError() {
		t.Fatalf("unexpected diags for known ToObjectValue: %v", diags)
	}
	if ov3.IsNull() || ov3.IsUnknown() {
		t.Fatalf("expected known object value for known PermissionsValue")
	}
}

func TestNewPermissionsValueMust_PanicOnError(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// missing id to induce diagnostic
	attrs := map[string]attr.Value{
		"alias": types.StringValue("a"),
		"attributes": func() attr.Value {
			v, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
			return v
		}(),
	}

	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		_ = NewPermissionsValueMust(attrTypes, attrs)
	}()

	if !didPanic {
		t.Fatalf("expected NewPermissionsValueMust to panic on invalid attributes")
	}
}
