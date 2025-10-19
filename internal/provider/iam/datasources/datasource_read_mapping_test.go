package datasources

import (
	"context"
	"testing"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/iam"
)

func TestGroupsDataSource_Read_MapsGroupsToList(t *testing.T) {
	// Prepare mock raw client returning two groups as JSON array
	groups := []iam.Group{{ID: "g1", Name: "Group1", Members: []string{"m1"}, CreatedAt: "c1"}, {ID: "g2", Name: "Group2", Members: []string{}, CreatedAt: "c2"}}
	_ = context.Background()

	// Call the mapping helper directly to avoid state conversion complexities
	listResp := &iam.ListGroupsResponse{Groups: groups}
	listValue, diags := mapGroupsToListValue(listResp)
	if diags.HasError() {
		t.Fatalf("mapping diagnostics: %v", diags)
	}

	var groupsOut []GroupDataModel
	diags2 := listValue.ElementsAs(context.Background(), &groupsOut, false)
	if diags2.HasError() {
		t.Fatalf("elements as diagnostics: %v", diags2)
	}
	if len(groupsOut) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groupsOut))
	}
	if groupsOut[0].ID.ValueString() != "g1" || groupsOut[0].Name.ValueString() != "Group1" {
		t.Fatalf("unexpected first group: %+v", groupsOut[0])
	}
}

func TestRolesDataSource_Read_MapsRolesToList(t *testing.T) {
	_ = context.Background()

	// Call mapping helper directly
	rolesSlice := []iam.Role{{ID: "r1", Name: "Role1", Title: "T1"}}
	listValue, diags := mapRolesToListValue(rolesSlice)
	if diags.HasError() {
		t.Fatalf("mapping diagnostics: %v", diags)
	}
	var out []RoleDataModel
	diags2 := listValue.ElementsAs(context.Background(), &out, false)
	if diags2.HasError() {
		t.Fatalf("elements as diagnostics: %v", diags2)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 role, got %d", len(out))
	}
	if out[0].ID.ValueString() != "r1" || out[0].Name.ValueString() != "Role1" {
		t.Fatalf("unexpected role: %+v", out[0])
	}
}
