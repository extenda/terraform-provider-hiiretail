package resource_iam_custom_role

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExercise_AllHelpers(t *testing.T) {
    ctx := context.Background()

    // Types String/Equal
    pt := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
    at := AttributesType{ObjectType: types.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}}

    _ = pt.String()
    _ = at.String()

    if pt.Equal(pt) == false {
        t.Fatalf("PermissionsType should equal itself")
    }
    if at.Equal(at) == false {
        t.Fatalf("AttributesType should equal itself")
    }

    // Value Type and AttributeTypes
    pv := PermissionsValue{state: attr.ValueStateKnown, Alias: types.StringValue("a"), Attributes: types.ObjectNull(AttributesValue{}.AttributeTypes(ctx)), Id: types.StringValue("id")}
    _ = pv.Type(ctx)
    _ = pv.AttributeTypes(ctx)
    _ = pv.String()

    av := AttributesValue{state: attr.ValueStateKnown}
    _ = av.Type(ctx)
    _ = av.AttributeTypes(ctx)
    _ = av.String()

    // New Null/Unknown
    _ = NewPermissionsValueNull()
    _ = NewPermissionsValueUnknown()
    _ = NewAttributesValueNull()
    _ = NewAttributesValueUnknown()

    // NewPermissionsValueMust with valid attrs
    permAttrTypes := PermissionsValue{}.AttributeTypes(ctx)
    innerObj, diags := types.ObjectValue(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{})
    if diags.HasError() {
        t.Fatalf("failed to build inner object: %v", diags)
    }
    attrs := map[string]attr.Value{"alias": types.StringValue("a"), "attributes": innerObj, "id": types.StringValue("id")}
    _ = NewPermissionsValueMust(permAttrTypes, attrs)

    // NewAttributesValueMust (attributes have no internal attributes, this should succeed)
    _ = NewAttributesValueMust(map[string]attr.Type{}, map[string]attr.Value{})

    // ToTerraformValue and ToObjectValue for AttributesValue
    if _, err := av.ToTerraformValue(ctx); err != nil {
        t.Fatalf("AttributesValue.ToTerraformValue failed: %v", err)
    }
    if _, diags := av.ToObjectValue(ctx); diags.HasError() {
        t.Fatalf("AttributesValue.ToObjectValue failed: %v", diags)
    }

    // ToTerraformValue and ToObjectValue for PermissionsValue
    if _, err := pv.ToTerraformValue(ctx); err != nil {
        t.Fatalf("PermissionsValue.ToTerraformValue failed: %v", err)
    }
    if _, diags := pv.ToObjectValue(ctx); diags.HasError() {
        t.Fatalf("PermissionsValue.ToObjectValue failed: %v", diags)
    }

    // Equal against different type
    if pv.Equal(av) {
        t.Fatalf("PermissionsValue should not be equal to AttributesValue")
    }
}
