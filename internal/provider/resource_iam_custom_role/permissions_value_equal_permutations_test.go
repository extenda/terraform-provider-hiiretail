package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPermissionsValue_Equal_Permutations(t *testing.T) {
	ctx := context.Background()

	// Known values equal
	a := PermissionsValue{state: attr.ValueStateKnown, Alias: types.StringValue("a"), Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)), Id: types.StringValue("x")}
	b := PermissionsValue{state: attr.ValueStateKnown, Alias: types.StringValue("a"), Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)), Id: types.StringValue("x")}
	if !a.Equal(b) {
		t.Fatalf("expected equal for identical known values")
	}

	// Known values different alias
	b.Alias = types.StringValue("different")
	if a.Equal(b) {
		t.Fatalf("expected not equal when alias differs")
	}

	// Attributes difference
	objAttrs, _ := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{"k": types.StringValue("v")})
	b = PermissionsValue{state: attr.ValueStateKnown, Alias: types.StringValue("a"), Attributes: objAttrs, Id: types.StringValue("x")}
	if a.Equal(b) {
		t.Fatalf("expected not equal when attributes differ")
	}

	// Id difference
	b = PermissionsValue{state: attr.ValueStateKnown, Alias: types.StringValue("a"), Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)), Id: types.StringValue("other")}
	if a.Equal(b) {
		t.Fatalf("expected not equal when id differs")
	}

	// State mismatch: one null, one known
	nullV := PermissionsValue{state: attr.ValueStateNull}
	if a.Equal(nullV) {
		t.Fatalf("expected not equal when one is null and other known")
	}

	// Both null -> equal
	nullV2 := PermissionsValue{state: attr.ValueStateNull}
	if !nullV.Equal(nullV2) {
		t.Fatalf("expected equal when both null")
	}

	// Both unknown -> equal
	unk1 := PermissionsValue{state: attr.ValueStateUnknown}
	unk2 := PermissionsValue{state: attr.ValueStateUnknown}
	if !unk1.Equal(unk2) {
		t.Fatalf("expected equal when both unknown")
	}

	// Different type -> not equal
	var other attr.Value = AttributesValue{state: attr.ValueStateKnown}
	if a.Equal(other) {
		t.Fatalf("expected not equal when comparing to different value type")
	}
}
