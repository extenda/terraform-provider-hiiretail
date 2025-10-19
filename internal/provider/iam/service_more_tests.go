//go:build ignore

package iam

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

// Test CreateRoleBinding when ListGroups (GET /groups) returns an error from the client
func TestService_CreateRoleBinding_ListGroupsDoError(t *testing.T) {
	mockRaw := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "GET" && req.Path != "" && contains(req.Path, "/groups") {
			return nil, errors.New("network failure")
		}
		return &client.Response{StatusCode: 200, Body: []byte(`{}`)}, nil
	}}

	svc := &Service{rawClient: mockRaw, tenantID: "t"}
	rb := &RoleBinding{Members: []string{"group:my-group"}, Role: "roles/R1"}
	_, err := svc.CreateRoleBinding(context.Background(), rb)
	if err == nil || !contains(err.Error(), "failed to list groups") {
		t.Fatalf("expected list groups error, got: %v", err)
	}
}

// Test GetRoleBinding when V2 roles endpoint returns malformed JSON
func TestService_GetRoleBinding_V2MalformedJSON(t *testing.T) {
	// Prepare group response
	group := Group{ID: "g1", Name: "grp"}
	gbody, _ := json.Marshal(group)

	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// V2 GET roles returns invalid JSON
		if req.Method == "GET" && contains(req.Path, "/api/v2/") && contains(req.Path, "/roles") {
			return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
		}
		// GET group
		if req.Method == "GET" && contains(req.Path, "/groups/g1") {
			return &client.Response{StatusCode: 200, Body: gbody}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}

	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.GetRoleBinding(context.Background(), "g1-RoleX")
	if err == nil || !contains(err.Error(), "failed to parse role bindings response") {
		t.Fatalf("expected parse error, got: %v", err)
	}
}

// Test SetResource returns a decode error when API returns invalid JSON
func TestService_SetResource_DecodeError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" && contains(req.Path, "/resources/") {
			return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}

	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.SetResource(context.Background(), "res1", &SetResourceDto{Name: "R"})
	if err == nil || !contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

// tiny helper to avoid importing strings package everywhere; replicates strings.Contains
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && (indexOf(s, substr) >= 0)))
}

// naive indexOf implementation
func indexOf(s, sep string) int {
	n := len(s)
	m := len(sep)
	if m == 0 {
		return 0
	}
	for i := 0; i <= n-m; i++ {
		if s[i:i+m] == sep {
			return i
		}
	}
	return -1
}
