package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPermissionsType_ValueFromObject_Success_New(t *testing.T) {
	ctx := context.Background()

	// Build inner attributes object for the 'attributes' field
	attributesAttrTypes := AttributesValue{}.AttributeTypes(ctx)
	attributesVal, diags := types.ObjectValue(attributesAttrTypes, map[string]attr.Value{})
	if diags.HasError() {
		t.Fatalf("failed to construct attributes object: %v", diags)
	}

	// Build the permissions object value using the attribute types reported by PermissionsValue
	permAttrTypes := PermissionsValue{}.AttributeTypes(ctx)

	objMap := map[string]attr.Value{
		"alias":      types.StringValue("alias"),
		"attributes": attributesVal,
		"id":         types.StringValue("pid"),
	}

	objVal, diags := types.ObjectValue(permAttrTypes, objMap)
	if diags.HasError() {
		t.Fatalf("failed to construct permissions object: %v", diags)
	}

	v, diags := PermissionsType{}.ValueFromObject(ctx, objVal)
	if diags.HasError() {
		t.Fatalf("ValueFromObject returned diagnostics: %v", diags)
	}

	pv, ok := v.(PermissionsValue)
	if !ok {
		t.Fatalf("expected PermissionsValue, got: %T", v)
	}

	if !pv.Alias.Equal(types.StringValue("alias")) {
		t.Fatalf("alias mismatch: %v", pv.Alias)
	}

	if !pv.Id.Equal(types.StringValue("pid")) {
		t.Fatalf("id mismatch: %v", pv.Id)
	}

	if pv.Attributes.IsNull() || pv.Attributes.IsUnknown() {
		t.Fatalf("attributes should be known object")
	}
}

func TestNewPermissionsValue_Success_New(t *testing.T) {
	ctx := context.Background()

	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	attributesAttrTypes := AttributesValue{}.AttributeTypes(ctx)
	attributesVal, diags := types.ObjectValue(attributesAttrTypes, map[string]attr.Value{})
	if diags.HasError() {
		t.Fatalf("failed to construct inner attributes object: %v", diags)
	}

	attrs := map[string]attr.Value{
		"alias":      types.StringValue("alias"),
		"attributes": attributesVal,
		"id":         types.StringValue("pid"),
	}

	pv, diags := NewPermissionsValue(attrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("NewPermissionsValue returned diagnostics: %v", diags)
	}

	if pv.IsNull() || pv.IsUnknown() {
		t.Fatalf("expected known PermissionsValue, got null/unknown")
	}

	if !pv.Alias.Equal(types.StringValue("alias")) {
		t.Fatalf("alias mismatch: %v", pv.Alias)
	}
}

func TestAttributesType_ValueFromTerraform_Success_New(t *testing.T) {
	ctx := context.Background()

	tfType := AttributesType{}.TerraformType(ctx)
	// Expect an object type for attributes (may have empty attr types)
	vals := map[string]tftypes.Value{}
	tv := tftypes.NewValue(tfType, vals)

	v, err := AttributesType{}.ValueFromTerraform(ctx, tv)
	if err != nil {
		t.Fatalf("ValueFromTerraform returned error: %v", err)
	}

	av, ok := v.(AttributesValue)
	if !ok {
		t.Fatalf("expected AttributesValue, got: %T", v)
	}

	if av.IsNull() || av.IsUnknown() {
		t.Fatalf("expected known AttributesValue")
	}
}
