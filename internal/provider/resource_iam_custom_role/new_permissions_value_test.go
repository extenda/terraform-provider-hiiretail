package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNewPermissionsValue_Exhaustive(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// 1) Missing attribute -> expect diags
	attrs1 := map[string]attr.Value{"alias": types.StringValue("a")} // missing id and attributes
	v1, diags1 := NewPermissionsValue(attrTypes, attrs1)
	if !diags1.HasError() {
		t.Fatalf("expected diagnostics error for missing attributes, got none")
	}
	if !v1.IsUnknown() {
		t.Fatalf("expected unknown value when diags present")
	}

	// 2) Extra attribute -> expect diags
	attrs2 := map[string]attr.Value{"alias": types.StringValue("a"), "id": types.StringValue("x"), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), "extra": types.StringValue("x")}
	v2, diags2 := NewPermissionsValue(attrTypes, attrs2)
	if !diags2.HasError() {
		t.Fatalf("expected diagnostics error for extra attribute, got none")
	}
	if !v2.IsUnknown() {
		t.Fatalf("expected unknown value when diags present (extra)")
	}

	// 3) Wrong attribute type -> alias provided as object
	attrs3 := map[string]attr.Value{"alias": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), "id": types.StringValue("x"), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})}
	v3, diags3 := NewPermissionsValue(attrTypes, attrs3)
	if !diags3.HasError() {
		t.Fatalf("expected diagnostics error for wrong attribute type, got none")
	}
	if !v3.IsUnknown() {
		t.Fatalf("expected unknown value when diags present (wrong type)")
	}

	// 4) Success
	attrs4 := map[string]attr.Value{"alias": types.StringValue("a"), "id": types.StringValue("pos.payment.create"), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})}
	v4, diags4 := NewPermissionsValue(attrTypes, attrs4)
	if diags4.HasError() {
		t.Fatalf("unexpected diagnostics on success path: %v", diags4)
	}
	if v4.IsUnknown() || v4.IsNull() {
		t.Fatalf("expected known value on success path")
	}
}
