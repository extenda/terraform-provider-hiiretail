package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPermissionsType_ValueFromObject_MissingAlias_Custom(t *testing.T) {
	ctx := context.Background()
	// Build object missing alias
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	attributes := map[string]attr.Value{
		// omit alias
		"attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"id":         types.StringValue("perm.x.y"),
	}

	obj, _ := types.ObjectValue(attrTypes, attributes)

	v, diags := PermissionsType{}.ValueFromObject(ctx, obj)
	if diags == nil || !diags.HasError() {
		t.Fatalf("expected diagnostics with error for missing alias")
	}
	if v != nil {
		t.Fatalf("expected nil value on error")
	}
}

func TestPermissionsType_ValueFromObject_WrongTypesProducesDiags_Custom(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Provide wrong types for alias and attributes
	attributes := map[string]attr.Value{
		"alias":      types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"attributes": types.StringValue("not-an-object"),
		"id":         types.StringValue("perm.x.y"),
	}

	obj, _ := types.ObjectValue(attrTypes, attributes)

	_, diags := PermissionsType{}.ValueFromObject(ctx, obj)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics errors for wrong types")
	}
}

func TestPermissionsValue_ToTerraformValue_DefaultPanic(t *testing.T) {
	// Force an unknown state value not handled by switch to trigger panic
	pv := PermissionsValue{state: attr.ValueState(3)}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from ToTerraformValue on unhandled state")
		}
	}()

	_, _ = pv.ToTerraformValue(context.Background())
}

func TestAttributesValue_ToTerraformValue_DefaultPanic(t *testing.T) {
	av := AttributesValue{state: attr.ValueState(3)}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from AttributesValue.ToTerraformValue on unhandled state")
		}
	}()

	_, _ = av.ToTerraformValue(context.Background())
}
