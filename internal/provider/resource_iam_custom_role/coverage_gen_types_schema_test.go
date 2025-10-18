package resource_iam_custom_role

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestGeneratedTypes_And_Schema_Exercise(t *testing.T) {
    ctx := context.Background()

    // Exercise schema creation
    s := IamCustomRoleResourceSchema(ctx)
    if s.Attributes == nil {
        t.Fatalf("schema attributes should not be nil")
    }

    // PermissionsType & AttributesType string and equality
    pt := PermissionsType{ObjectType: types.ObjectType{AttrTypes: PermissionsValue{}.AttributeTypes(ctx)}}
    at := AttributesType{ObjectType: types.ObjectType{AttrTypes: AttributesValue{}.AttributeTypes(ctx)}}

    _ = pt.String()
    _ = at.String()

    if !pt.Equal(pt) {
        t.Fatalf("PermissionsType should equal itself")
    }
    if !at.Equal(at) {
        t.Fatalf("AttributesType should equal itself")
    }

    // ValueFromTerraform: nil type should produce Null
    if v, err := pt.ValueFromTerraform(ctx, tftypes.NewValue(nil, nil)); err != nil || !v.(PermissionsValue).IsNull() {
        t.Fatalf("expected null PermissionsValue for nil tftype")
    }
    if v, err := at.ValueFromTerraform(ctx, tftypes.NewValue(nil, nil)); err != nil || !v.(AttributesValue).IsNull() {
        t.Fatalf("expected null AttributesValue for nil tftype")
    }
}
