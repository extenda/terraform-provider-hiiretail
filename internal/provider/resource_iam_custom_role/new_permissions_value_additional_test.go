package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNewPermissionsValue_AdditiveCases(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// missing both attributes and id
	missingBoth := map[string]attr.Value{"alias": types.StringValue("only-alias")}
	vMissing, diagMissing := NewPermissionsValue(attrTypes, missingBoth)
	if !diagMissing.HasError() {
		t.Fatalf("expected diagnostics when required fields missing")
	}
	if !vMissing.IsUnknown() {
		t.Fatalf("expected unknown when diags present")
	}

	// wrong type for attributes (string instead of object)
	wrongAttrType := map[string]attr.Value{"alias": types.StringValue("a"), "id": types.StringValue("x"), "attributes": types.StringValue("bad")}
	vWrong, diagWrong := NewPermissionsValue(attrTypes, wrongAttrType)
	if !diagWrong.HasError() {
		t.Fatalf("expected diagnostics for wrong attributes type")
	}
	if !vWrong.IsUnknown() {
		t.Fatalf("expected unknown when diags present (wrong type)")
	}

	// extra attribute present
	extra := map[string]attr.Value{"alias": types.StringValue("a"), "id": types.StringValue("x"), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), "unexpected": types.StringValue("y")}
	vExtra, diagExtra := NewPermissionsValue(attrTypes, extra)
	if !diagExtra.HasError() {
		t.Fatalf("expected diagnostics for extra attribute")
	}
	if !vExtra.IsUnknown() {
		t.Fatalf("expected unknown when diags present (extra)")
	}

	// nominal success path
	okAttrs := map[string]attr.Value{"alias": types.StringValue("a"), "id": types.StringValue("pos.payment.create"), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})}
	vOk, diagOk := NewPermissionsValue(attrTypes, okAttrs)
	if diagOk.HasError() {
		t.Fatalf("unexpected diagnostics on success path: %v", diagOk)
	}
	if vOk.IsUnknown() || vOk.IsNull() {
		t.Fatalf("expected known PermissionsValue on success path")
	}
}
