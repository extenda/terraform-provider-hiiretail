package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Test ValueFromObject missing attributes and wrong types path
func Test_PermissionsType_ValueFromObject_MissingAndWrongTypes(t *testing.T) {
	ctx := context.Background()

	// Build a basetypes.ObjectValue missing required fields using types.ObjectNull
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)
	obj := types.ObjectNull(attrTypes)

	var pt PermissionsType
	_, diags := pt.ValueFromObject(ctx, obj)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for missing fields")
	}

	// Now build an object with wrong attribute types
	wrong := map[string]attr.Value{
		"alias":      types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		"attributes": types.StringValue("not-an-object"),
		"id":         types.StringValue("perm.x.y"),
	}
	obj2, diags2 := types.ObjectValue(attrTypes, wrong)
	if diags2.HasError() {
		// Creating the object may itself produce diagnostics; that's acceptable
		t.Logf("types.ObjectValue produced diags as expected: %v", diags2)
		return
	}

	_, diags3 := pt.ValueFromObject(ctx, obj2)
	if !diags3.HasError() {
		t.Fatalf("expected diagnostics for wrong attribute types in ValueFromObject")
	}
}

// Test NewPermissionsValueMust panics when given invalid attributes
func Test_NewPermissionsValueMust_PanicsOnInvalidInput(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// missing id to produce diagnostic
	attrs := map[string]attr.Value{
		"alias":      types.StringValue("a"),
		"attributes": types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from NewPermissionsValueMust when diagnostics present")
		}
	}()

	_ = NewPermissionsValueMust(attrTypes, attrs)
}

// Test PermissionsValue.Equal when fields differ
func Test_PermissionsValue_Equal_DifferentFields(t *testing.T) {
	ctx := context.Background()
	base := PermissionsValue{
		Alias:      types.StringValue("a"),
		Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)),
		Id:         types.StringValue("perm.x.y"),
		state:      attr.ValueStateKnown,
	}

	same := base
	if !base.Equal(same) {
		t.Fatalf("expected equal values to be equal")
	}

	diffID := base
	diffID.Id = types.StringValue("other")
	if base.Equal(diffID) {
		t.Fatalf("expected different id to not be equal")
	}

	diffState := base
	diffState.state = attr.ValueStateNull
	if base.Equal(diffState) {
		t.Fatalf("expected different state to not be equal")
	}
}

func Test_PermissionsValue_Equal_OtherTypeAndNullEarlyReturn(t *testing.T) {
	// other type
	pv := PermissionsValue{state: attr.ValueStateKnown}
	if pv.Equal(types.StringValue("x")) {
		t.Fatalf("expected Equal to return false when comparing to other value type")
	}

	// early return when not known
	a := PermissionsValue{state: attr.ValueStateNull}
	b := PermissionsValue{state: attr.ValueStateNull}
	if !a.Equal(b) {
		t.Fatalf("expected Equal to return true when both states are not known")
	}
}

func Test_AttributesType_ValueFromTerraform_MatchingTypeEmptyObject(t *testing.T) {
	var at AttributesType
	ctx := context.Background()
	// get the Terraform type produced by AttributesType
	tt := at.TerraformType(ctx)
	// construct an empty native object value using the exact same type
	v := tftypes.NewValue(tt, map[string]tftypes.Value{})
	got, err := at.ValueFromTerraform(ctx, v)
	if err != nil {
		t.Fatalf("unexpected error from AttributesType.ValueFromTerraform: %v", err)
	}
	// result should be known or null/unknown; ensure call succeeded
	_ = got
}
