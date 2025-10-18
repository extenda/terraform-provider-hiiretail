package resource_iam_custom_role

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestModelToAPIRequest_PermissionsAndNameConversion(t *testing.T) {
	ctx := context.Background()

	// Build two permission entries
	permType := PermissionsType{
		ObjectType: types.ObjectType{
			AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
		},
	}

	// First permission: has alias and one string attribute
	p1 := PermissionsValue{
		Alias:      types.StringValue("alias1"),
		Attributes: types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		Id:         types.StringValue("pos.payment.create"),
		state:      attr.ValueStateKnown,
	}

	// Second permission: alias null, attributes with a string value
	p2 := PermissionsValue{
		Alias:      types.StringNull(),
		Attributes: types.ObjectValueMust(AttributesValue{}.AttributeTypes(ctx), map[string]attr.Value{}),
		Id:         types.StringValue("sys.user.manage"),
		state:      attr.ValueStateKnown,
	}

	lst, diags := types.ListValueFrom(ctx, permType, []PermissionsValue{p1, p2})
	if diags.HasError() {
		t.Fatalf("failed to create list: %v", diags)
	}

	model := IamCustomRoleModel{
		Id:          types.StringValue("role-1"),
		Name:        types.StringValue("My Role"),
		Permissions: lst,
		TenantId:    types.StringNull(),
	}

	r := &IamCustomRoleResource{}

	req, err := r.modelToAPIRequest(ctx, model)
	if err != nil {
		t.Fatalf("unexpected error from modelToAPIRequest: %v", err)
	}

	if req.ID != "role-1" {
		t.Fatalf("expected id role-1, got %s", req.ID)
	}
	if req.Name != "My Role" {
		t.Fatalf("expected name My Role, got %s", req.Name)
	}
	if len(req.Permissions) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(req.Permissions))
	}
	// second permission should not contain attributes (empty attrs not promoted to map)
	if req.Permissions[1].Attributes != nil {
		t.Fatalf("expected no attributes map for second permission, got one")
	}
}
