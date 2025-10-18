package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPermissionsType_ValueFromObject_MissingAttributesIndividually(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// Missing alias only
	attrs1 := map[string]attr.Value{"attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), "id": types.StringValue("x")}
	obj1, _ := types.ObjectValue(attrTypes, attrs1)
	v1, d1 := PermissionsType{}.ValueFromObject(ctx, obj1)
	if !d1.HasError() || v1 != nil {
		t.Fatalf("expected error for missing alias only")
	}

	// Missing attributes only
	attrs2 := map[string]attr.Value{"alias": types.StringValue("a"), "id": types.StringValue("x")}
	obj2, _ := types.ObjectValue(attrTypes, attrs2)
	v2, d2 := PermissionsType{}.ValueFromObject(ctx, obj2)
	if !d2.HasError() || v2 != nil {
		t.Fatalf("expected error for missing attributes only")
	}

	// Missing id only
	attrs3 := map[string]attr.Value{"alias": types.StringValue("a"), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})}
	obj3, _ := types.ObjectValue(attrTypes, attrs3)
	v3, d3 := PermissionsType{}.ValueFromObject(ctx, obj3)
	if !d3.HasError() || v3 != nil {
		t.Fatalf("expected error for missing id only")
	}
}

func TestPermissionsType_ValueFromObject_WrongTypeIndividually(t *testing.T) {
	ctx := context.Background()
	attrTypes := PermissionsValue{}.AttributeTypes(ctx)

	// alias wrong type
	attrs1 := map[string]attr.Value{"alias": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), "id": types.StringValue("x")}
	obj1, _ := types.ObjectValue(attrTypes, attrs1)
	v1, d1 := PermissionsType{}.ValueFromObject(ctx, obj1)
	if !d1.HasError() || v1 != nil {
		t.Fatalf("expected error for alias wrong type")
	}

	// attributes wrong type
	attrs2 := map[string]attr.Value{"alias": types.StringValue("a"), "attributes": types.StringValue("nope"), "id": types.StringValue("x")}
	obj2, _ := types.ObjectValue(attrTypes, attrs2)
	v2, d2 := PermissionsType{}.ValueFromObject(ctx, obj2)
	if !d2.HasError() || v2 != nil {
		t.Fatalf("expected error for attributes wrong type")
	}

	// id wrong type
	attrs3 := map[string]attr.Value{"alias": types.StringValue("a"), "attributes": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}), "id": types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})}
	obj3, _ := types.ObjectValue(attrTypes, attrs3)
	v3, d3 := PermissionsType{}.ValueFromObject(ctx, obj3)
	if !d3.HasError() || v3 != nil {
		t.Fatalf("expected error for id wrong type")
	}
}
