package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
)

// Reusable test helper: convert ListGroupsResponse to types.List (same logic as production)
func mapGroupsToListValue_test(listResp *iam.ListGroupsResponse) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	groupElements := make([]attr.Value, len(listResp.Groups))
	for i, group := range listResp.Groups {
		groupObj := map[string]attr.Value{
			"id":           types.StringValue(group.ID),
			"name":         types.StringValue(group.Name),
			"description":  types.StringValue(group.Description),
			"member_count": types.Int64Value(int64(len(group.Members))),
			"created_at":   types.StringValue(group.CreatedAt),
		}

		objType := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":           types.StringType,
				"name":         types.StringType,
				"description":  types.StringType,
				"member_count": types.Int64Type,
				"created_at":   types.StringType,
			},
		}

		objValue, ds := types.ObjectValue(objType.AttrTypes, groupObj)
		diags.Append(ds...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}

		groupElements[i] = objValue
	}

	listType := types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":           types.StringType,
		"name":         types.StringType,
		"description":  types.StringType,
		"member_count": types.Int64Type,
		"created_at":   types.StringType,
	}}}

	listValue, ds := types.ListValue(listType.ElemType, groupElements)
	diags.Append(ds...)
	if diags.HasError() {
		return types.ListNull(listType.ElemType), diags
	}
	return listValue, diags
}

// Reusable test helper for roles
func mapRolesToListValue_test(roles []iam.Role) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	roleElements := make([]attr.Value, len(roles))
	for i, role := range roles {
		roleObj := map[string]attr.Value{
			"id":          types.StringValue(role.ID),
			"name":        types.StringValue(role.Name),
			"title":       types.StringValue(role.Title),
			"description": types.StringValue(role.Description),
			"stage":       types.StringValue(role.Stage),
			"type":        types.StringValue(role.Type),
		}

		objType := types.ObjectType{AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"name":        types.StringType,
			"title":       types.StringType,
			"description": types.StringType,
			"stage":       types.StringType,
			"type":        types.StringType,
		}}

		objValue, ds := types.ObjectValue(objType.AttrTypes, roleObj)
		diags.Append(ds...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		roleElements[i] = objValue
	}

	listType := types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"title":       types.StringType,
		"description": types.StringType,
		"stage":       types.StringType,
		"type":        types.StringType,
	}}}

	listValue, ds := types.ListValue(listType.ElemType, roleElements)
	diags.Append(ds...)
	if diags.HasError() {
		return types.ListNull(listType.ElemType), diags
	}
	return listValue, diags
}

// Edge-case test: empty groups list
func TestGroupsMapping_EmptyList(t *testing.T) {
	lr := &iam.ListGroupsResponse{Groups: []iam.Group{}}
	listValue, diags := mapGroupsToListValue_test(lr)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	var out []GroupDataModel
	diags2 := listValue.ElementsAs(context.Background(), &out, false)
	if diags2.HasError() {
		t.Fatalf("elements as diagnostics: %v", diags2)
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(out))
	}
}

// Edge-case test: role missing optional fields
func TestRolesMapping_MissingOptionalFields(t *testing.T) {
	roles := []iam.Role{{ID: "r-missing", Name: "RoleMissing"}}
	listValue, diags := mapRolesToListValue_test(roles)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	var out []RoleDataModel
	diags2 := listValue.ElementsAs(context.Background(), &out, false)
	if diags2.HasError() {
		t.Fatalf("elements as diagnostics: %v", diags2)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 role, got %d", len(out))
	}
	if out[0].ID.ValueString() != "r-missing" || out[0].Name.ValueString() != "RoleMissing" {
		t.Fatalf("unexpected role: %+v", out[0])
	}
}
