package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestExercise_GeneratedPaths(t *testing.T) {
	ctx := context.Background()

	// PermissionsValue known
	pv := PermissionsValue{state: attr.ValueStateKnown, Alias: types.StringValue("a"), Attributes: types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), Id: types.StringValue("id")}

	if pv.IsNull() || pv.IsUnknown() {
		t.Fatalf("pv should be known")
	}

	if _, err := pv.ToTerraformValue(ctx); err != nil {
		t.Fatalf("ToTerraformValue known returned error: %v", err)
	}

	if _, diags := pv.ToObjectValue(ctx); diags.HasError() {
		t.Fatalf("ToObjectValue known returned diags: %v", diags)
	}

	// PermissionsValue null
	pvNull := PermissionsValue{state: attr.ValueStateNull}
	if !pvNull.IsNull() {
		t.Fatalf("pvNull should be null")
	}
	if _, err := pvNull.ToTerraformValue(ctx); err != nil {
		t.Fatalf("ToTerraformValue null returned error: %v", err)
	}

	// PermissionsValue unknown
	pvUnk := PermissionsValue{state: attr.ValueStateUnknown}
	if !pvUnk.IsUnknown() {
		t.Fatalf("pvUnk should be unknown")
	}
	if _, err := pvUnk.ToTerraformValue(ctx); err != nil {
		t.Fatalf("ToTerraformValue unknown returned error: %v", err)
	}

	// AttributesValue known (empty)
	av := AttributesValue{state: attr.ValueStateKnown}
	if _, err := av.ToTerraformValue(ctx); err != nil {
		t.Fatalf("Attributes ToTerraformValue returned error: %v", err)
	}

	if _, diags := av.ToObjectValue(ctx); diags.HasError() {
		t.Fatalf("Attributes ToObjectValue returned diags: %v", diags)
	}

	// Roundtrip PermissionsType.ValueFromTerraform with known object
	tfType := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}.TerraformType(ctx)
	objTf := tfType.(tftypes.Object)

	vals := map[string]tftypes.Value{
		"alias":      tftypes.NewValue(objTf.AttributeTypes["alias"], "alias"),
		"attributes": tftypes.NewValue(objTf.AttributeTypes["attributes"], map[string]tftypes.Value{}),
		"id":         tftypes.NewValue(objTf.AttributeTypes["id"], "pid"),
	}

	tv := tftypes.NewValue(tfType, vals)
	_, err := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}.ValueFromTerraform(ctx, tv)
	if err != nil {
		t.Fatalf("PermissionsType.ValueFromTerraform failed: %v", err)
	}
}
