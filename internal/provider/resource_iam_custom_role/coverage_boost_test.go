package resource_iam_custom_role

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Test that a nil Terraform type yields a null PermissionsValue
func Test_PermissionsType_ValueFromTerraform_NilType_ReturnsNull(t *testing.T) {
    var pt PermissionsType
    v := tftypes.NewValue(nil, nil)
    got, err := pt.ValueFromTerraform(context.Background(), v)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !got.(PermissionsValue).IsNull() {
        t.Fatalf("expected PermissionsValue to be null")
    }
}

// Test that a mismatched terraform type returns an error
func Test_PermissionsType_ValueFromTerraform_TypeMismatch(t *testing.T) {
    var pt PermissionsType
    v := tftypes.NewValue(tftypes.Number, 42)
    _, err := pt.ValueFromTerraform(context.Background(), v)
    if err == nil {
        t.Fatalf("expected type mismatch error")
    }
}

// Test unknown and null paths for PermissionsType.ValueFromTerraform
func Test_PermissionsType_ValueFromTerraform_UnknownAndNullPaths(t *testing.T) {
    var pt PermissionsType
    // Use a minimal object type; ValueFromTerraform only checks for non-nil type
    objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}

    unk := tftypes.NewValue(objType, tftypes.UnknownValue)
    got, err := pt.ValueFromTerraform(context.Background(), unk)
    if err != nil {
        t.Fatalf("unexpected error for unknown: %v", err)
    }
    if !got.(PermissionsValue).IsUnknown() {
        t.Fatalf("expected unknown PermissionsValue")
    }

    nul := tftypes.NewValue(objType, nil)
    got2, err := pt.ValueFromTerraform(context.Background(), nul)
    if err != nil {
        t.Fatalf("unexpected error for null: %v", err)
    }
    if !got2.(PermissionsValue).IsNull() {
        t.Fatalf("expected null PermissionsValue")
    }
}

// Test NewPermissionsValue produces diagnostics for extra and wrong-typed attributes
func Test_NewPermissionsValue_ExtraAndWrongTypes(t *testing.T) {
    // Reuse the generated attribute types to keep this test robust to changes
    attributeTypes := PermissionsValue{}.AttributeTypes(context.Background())

    // extra attribute should give diags
    attrsExtra := map[string]attr.Value{
        "alias":      types.StringValue("ok"),
        "id":         types.StringValue("perm.x.y"),
        "attributes": types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}),
        "extra":      types.StringValue("extra"),
    }

    v, diags := NewPermissionsValue(attributeTypes, attrsExtra)
    if diags == nil || !diags.HasError() {
        t.Fatalf("expected diagnostics for extra attribute")
    }
    if !v.IsUnknown() {
        t.Fatalf("expected Unknown PermissionsValue when diags present")
    }

    // wrong-typed alias should also produce diagnostics
    attrsWrong := map[string]attr.Value{
        "alias":      types.ObjectValueMust(AttributesValue{}.AttributeTypes(context.Background()), map[string]attr.Value{}),
        "id":         types.StringValue("perm.x.y"),
        "attributes": types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}),
    }

    v2, diags2 := NewPermissionsValue(attributeTypes, attrsWrong)
    if diags2 == nil || !diags2.HasError() {
        t.Fatalf("expected diagnostics for wrong-typed attribute")
    }
    if !v2.IsUnknown() {
        t.Fatalf("expected Unknown PermissionsValue when diags present (wrong type)")
    }
}

// Test AttributesType.ValueFromTerraform nil and mismatch
func Test_AttributesType_ValueFromTerraform_NilAndMismatch(t *testing.T) {
    var at AttributesType
    nilv := tftypes.NewValue(nil, nil)
    got, err := at.ValueFromTerraform(context.Background(), nilv)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !got.(AttributesValue).IsNull() {
        t.Fatalf("expected null AttributesValue")
    }

    // mismatch: number instead of object
    num := tftypes.NewValue(tftypes.Number, 1)
    _, err2 := at.ValueFromTerraform(context.Background(), num)
    if err2 == nil {
        t.Fatalf("expected type mismatch error for AttributesType")
    }
}
