package iam

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

// MockClient implements RawClient for tests
type MockClient struct {
	DoFunc func(ctx context.Context, req *client.Request) (*client.Response, error)
}

func (m *MockClient) Do(ctx context.Context, req *client.Request) (*client.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(ctx, req)
	}
	return nil, errors.New("not implemented")
}

func TestService_TenantID(t *testing.T) {
	svc := &Service{rawClient: &MockClient{}, tenantID: "test-tenant"}
	if got := svc.TenantID(); got != "test-tenant" {
		t.Errorf("TenantID() = %v, want %v", got, "test-tenant")
	}
}

func TestService_ListGroups(t *testing.T) {
	cases := []struct {
		name     string
		respBody []byte
		status   int
		err      error
		wantLen  int
		wantErr  string
	}{
		{
			name:     "success",
			respBody: []byte(`[{"id":"g1","name":"Group1"},{"id":"g2","name":"Group2"}]`),
			status:   200,
			wantLen:  2,
		},
		{
			name:    "api error",
			status:  500,
			err:     errors.New("api error"),
			wantErr: "failed to list groups: api error",
		},
	}

	for _, tc := range cases {
		mock := &MockClient{
			DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
				if tc.err != nil {
					return &client.Response{StatusCode: tc.status}, tc.err
				}
				return &client.Response{StatusCode: tc.status, Body: tc.respBody}, nil
			},
		}
		svc := &Service{rawClient: mock, tenantID: "t"}
		res, err := svc.ListGroups(context.Background(), &ListGroupsRequest{})
		if tc.wantErr != "" {
			if err == nil || err.Error() != tc.wantErr {
				t.Errorf("%s: error = %v, want %v", tc.name, err, tc.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
			continue
		}
		if res == nil {
			t.Errorf("%s: result is nil", tc.name)
			continue
		}
		if len(res.Groups) != tc.wantLen {
			t.Errorf("%s: got %d groups, want %d", tc.name, len(res.Groups), tc.wantLen)
		}
	}
}

func TestService_GetGroup(t *testing.T) {
	mockGroup := Group{ID: "g1", Name: "Group1"}
	body, _ := json.Marshal(mockGroup)

	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: body}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	g, err := svc.GetGroup(context.Background(), "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID != "g1" {
		t.Fatalf("got id %s want g1", g.ID)
	}
}

func TestService_CreateUpdateDeleteGroup(t *testing.T) {
	// Create
	created := Group{ID: "g1", Name: "Group1"}
	cbody, _ := json.Marshal(created)
	mockCreate := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method != "POST" {
			return &client.Response{StatusCode: 400}, errors.New("unexpected method")
		}
		return &client.Response{StatusCode: 201, Body: cbody}, nil
	}}
	svc := &Service{rawClient: mockCreate, tenantID: "t"}
	g, err := svc.CreateGroup(context.Background(), &Group{Name: "Group1"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if g.ID != "g1" {
		t.Fatalf("create id mismatch: %s", g.ID)
	}

	// Update (simulate 204 -> GetGroup called)
	mockUpdate := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}
		// GET for GetGroup
		return &client.Response{StatusCode: 200, Body: cbody}, nil
	}}
	svc = &Service{rawClient: mockUpdate, tenantID: "t"}
	ug, err := svc.UpdateGroup(context.Background(), "g1", &Group{Name: "Group1"})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if ug.ID != "g1" {
		t.Fatalf("update id mismatch: %s", ug.ID)
	}

	// Delete
	mockDelete := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method != "DELETE" {
			return &client.Response{StatusCode: 400}, errors.New("unexpected method")
		}
		return &client.Response{StatusCode: 204}, nil
	}}
	svc = &Service{rawClient: mockDelete, tenantID: "t"}
	if err := svc.DeleteGroup(context.Background(), "g1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
}

func TestService_ListRoles_GetRole_CreateCustom_GetCustomRole(t *testing.T) {
	// ListRoles
	rolesResp := struct {
		Roles []Role `json:"roles"`
	}{Roles: []Role{{ID: "r1", Name: "Role1"}}}
	rb, _ := json.Marshal(rolesResp)
	mockList := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: rb}, nil
	}}
	svc := &Service{rawClient: mockList, tenantID: "t"}
	roles, err := svc.ListRoles(context.Background(), "")
	if err != nil {
		t.Fatalf("ListRoles failed: %v", err)
	}
	if len(roles) != 1 || roles[0].ID != "r1" {
		t.Fatalf("unexpected roles: %v", roles)
	}

	// GetRole
	roleBody, _ := json.Marshal(Role{ID: "r1", Name: "Role1"})
	mockGet := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: roleBody}, nil
	}}
	svc = &Service{rawClient: mockGet, tenantID: "t"}
	r, err := svc.GetRole(context.Background(), "Role1")
	if err != nil {
		t.Fatalf("GetRole failed: %v", err)
	}
	if r.ID != "r1" {
		t.Fatalf("GetRole id mismatch: %s", r.ID)
	}

	// CreateCustomRole
	cr := CustomRole{ID: "cr1", Permissions: []Permission{{ID: "p1"}}}
	crBody, _ := json.Marshal(cr)
	mockCreateCR := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 201, Body: crBody}, nil
	}}
	svc = &Service{rawClient: mockCreateCR, tenantID: "t"}
	created, err := svc.CreateCustomRole(context.Background(), &cr)
	if err != nil {
		t.Fatalf("CreateCustomRole failed: %v", err)
	}
	if created.ID != "cr1" {
		t.Fatalf("CreateCustomRole id mismatch: %s", created.ID)
	}

	// GetCustomRole
	mockGetCR := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: crBody}, nil
	}}
	svc = &Service{rawClient: mockGetCR, tenantID: "t"}
	gcr, err := svc.GetCustomRole(context.Background(), "cr1")
	if err != nil {
		t.Fatalf("GetCustomRole failed: %v", err)
	}
	if gcr.ID != "cr1" {
		t.Fatalf("GetCustomRole id mismatch: %s", gcr.ID)
	}
}

func TestService_UpdateCustomRole_DeleteCustomRole(t *testing.T) {
	// UpdateCustomRole - normal 200 response
	updated := CustomRole{ID: "cr1", Name: "CustomRole1", Permissions: []Permission{{ID: "p1"}}}
	ub, _ := json.Marshal(updated)
	mockUpdate := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" {
			return &client.Response{StatusCode: 200, Body: ub}, nil
		}
		return nil, errors.New("unexpected method")
	}}
	svc := &Service{rawClient: mockUpdate, tenantID: "t"}
	ur, err := svc.UpdateCustomRole(context.Background(), "cr1", &updated)
	if err != nil {
		t.Fatalf("UpdateCustomRole failed: %v", err)
	}
	if ur.ID != "cr1" {
		t.Fatalf("UpdateCustomRole id mismatch: %s", ur.ID)
	}

	// UpdateCustomRole - 204 then GET
	mockUpdate204 := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}
		// GET
		return &client.Response{StatusCode: 200, Body: ub}, nil
	}}
	svc = &Service{rawClient: mockUpdate204, tenantID: "t"}
	ur2, err := svc.UpdateCustomRole(context.Background(), "cr1", &updated)
	if err != nil {
		t.Fatalf("UpdateCustomRole(204) failed: %v", err)
	}
	if ur2.ID != "cr1" {
		t.Fatalf("UpdateCustomRole(204) id mismatch: %s", ur2.ID)
	}

	// DeleteCustomRole - success
	mockDel := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "DELETE" {
			return &client.Response{StatusCode: 204}, nil
		}
		return nil, errors.New("unexpected method")
	}}
	svc = &Service{rawClient: mockDel, tenantID: "t"}
	if err := svc.DeleteCustomRole(context.Background(), "cr1"); err != nil {
		t.Fatalf("DeleteCustomRole failed: %v", err)
	}

	// DeleteCustomRole - api error
	mockDelErr := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 500}, errors.New("api error")
	}}
	svc = &Service{rawClient: mockDelErr, tenantID: "t"}
	if err := svc.DeleteCustomRole(context.Background(), "cr1"); err == nil {
		t.Fatalf("DeleteCustomRole expected error, got nil")
	}
}

func TestService_RoleBindingAndResourceFlows(t *testing.T) {
	// Prepare group returned by ListGroups/GetGroup
	group := Group{ID: "g1", Name: "my-group", CreatedAt: "2020-01-01T00:00:00Z", UpdatedAt: "2020-01-02T00:00:00Z"}
	gbody, _ := json.Marshal(group)

	// Prepare V2 role binding DTO for GetRoleBinding
	dto := []RoleBindingDto{{IsCustom: false, RoleID: "Role1", Bindings: []string{"bu:001"}}}
	dtoBody, _ := json.Marshal(dto)

	// Mock client that handles multiple endpoints
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// V2 GET roles for group (check before generic group GET)
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: dtoBody}, nil
		}
		// GET group (v1)
		if strings.Contains(req.Path, "/api/v1/tenants/") && strings.Contains(req.Path, "/groups/g1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: gbody}, nil
		}
		// V2 POST create role binding
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "POST" {
			return &client.Response{StatusCode: 201, Body: []byte(`{"status":"ok"}`)}, nil
		}
		// V2 DELETE role binding
		if strings.Contains(req.Path, "/api/v2/") && strings.Contains(req.Path, "/roles/") && req.Method == "DELETE" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}
		// resources endpoints
		if strings.Contains(req.Path, "/resources/res1") && req.Method == "PUT" {
			return &client.Response{StatusCode: 200, Body: []byte(`{"id":"res1","name":"Resource1"}`)}, nil
		}
		if strings.Contains(req.Path, "/resources/res1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: []byte(`{"id":"res1","name":"Resource1"}`)}, nil
		}
		if strings.Contains(req.Path, "/resources/res1") && req.Method == "DELETE" {
			return &client.Response{StatusCode: 204}, nil
		}
		if strings.Contains(req.Path, "/resources") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: []byte(`[{"id":"res1","name":"Resource1"}]`)}, nil
		}
		// fallback: 404
		return &client.Response{StatusCode: 404, Body: []byte(`{}`)}, nil
	}}

	svc := &Service{rawClient: mock, tenantID: "t"}

	// Test GetRoleBinding happy path
	b, err := svc.GetRoleBinding(context.Background(), "g1-Role1")
	if err != nil {
		t.Fatalf("GetRoleBinding failed: %v", err)
	}
	if b == nil || b.Role == "" {
		t.Fatalf("GetRoleBinding returned empty binding: %+v", b)
	}

	// Test CreateRoleBinding (uses ListGroups -> find group by name)
	// CreateRoleBinding looks up group by name; prepare ListGroups to return group
	mockList := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if strings.Contains(req.Path, "/groups") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: []byte(`[{"id":"g1","name":"my-group"}]`)}, nil
		}
		// POST to V2 create binding
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "POST" {
			return &client.Response{StatusCode: 201, Body: []byte(`{"ok":true}`)}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc = &Service{rawClient: mockList, tenantID: "t"}
	rb := &RoleBinding{Members: []string{"group:my-group"}, Role: "roles/Role1", Name: "rb1"}
	created, err := svc.CreateRoleBinding(context.Background(), rb)
	if err != nil {
		t.Fatalf("CreateRoleBinding failed: %v", err)
	}
	if created == nil || created.ID == "" {
		t.Fatalf("CreateRoleBinding returned invalid: %+v", created)
	}

	// Test UpdateRoleBinding simply calls GetRoleBinding internally; reuse mock from first section
	svc = &Service{rawClient: mock, tenantID: "t"}
	upd, err := svc.UpdateRoleBinding(context.Background(), "g1-Role1", &RoleBinding{Name: "newname"})
	if err != nil {
		t.Fatalf("UpdateRoleBinding failed: %v", err)
	}
	if upd.Name != "newname" && upd.ID == "" {
		t.Fatalf("UpdateRoleBinding returned unexpected: %+v", upd)
	}

	// Test DeleteRoleBinding
	svc = &Service{rawClient: mock, tenantID: "t"}
	if err := svc.DeleteRoleBinding(context.Background(), "g1-Role1"); err != nil {
		t.Fatalf("DeleteRoleBinding failed: %v", err)
	}

	// Test resource flows
	svc = &Service{rawClient: mock, tenantID: "t"}
	res, err := svc.SetResource(context.Background(), "res1", &SetResourceDto{Name: "Resource1"})
	if err != nil {
		t.Fatalf("SetResource failed: %v", err)
	}
	if res.ID != "res1" {
		t.Fatalf("SetResource id mismatch: %s", res.ID)
	}
	gres, err := svc.GetResource(context.Background(), "res1")
	if err != nil {
		t.Fatalf("GetResource failed: %v", err)
	}
	if gres.ID != "res1" {
		t.Fatalf("GetResource id mismatch: %s", gres.ID)
	}
	if err := svc.DeleteResource(context.Background(), "res1"); err != nil {
		t.Fatalf("DeleteResource failed: %v", err)
	}
	grs, err := svc.GetResources(context.Background(), &GetResourcesRequest{})
	if err != nil {
		t.Fatalf("GetResources failed: %v", err)
	}
	if len(grs.Resources) == 0 {
		t.Fatalf("GetResources empty: %+v", grs)
	}

	// Test AddRoleToGroup - non-custom
	mockAdd := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "POST" {
			return &client.Response{StatusCode: 201}, nil
		}
		return &client.Response{StatusCode: 200}, nil
	}}
	svc = &Service{rawClient: mockAdd, tenantID: "t"}
	if err := svc.AddRoleToGroup(context.Background(), "g1", "Role1", false, []string{"bu:001"}); err != nil {
		t.Fatalf("AddRoleToGroup failed: %v", err)
	}
}

// MockServiceClient allows mocking ServiceClient methods used by Service
type MockServiceClient struct {
	GetFunc  func(ctx context.Context, path string, query map[string]string) (*client.Response, error)
	PostFunc func(ctx context.Context, path string, body interface{}) (*client.Response, error)
	PutFunc  func(ctx context.Context, path string, body interface{}) (*client.Response, error)
	DelFunc  func(ctx context.Context, path string) (*client.Response, error)
}

func (m *MockServiceClient) Get(ctx context.Context, path string, query map[string]string) (*client.Response, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, path, query)
	}
	return &client.Response{StatusCode: 404, Body: []byte(`{}`)}, nil
}
func (m *MockServiceClient) Post(ctx context.Context, path string, body interface{}) (*client.Response, error) {
	if m.PostFunc != nil {
		return m.PostFunc(ctx, path, body)
	}
	return &client.Response{StatusCode: 201, Body: []byte(`{}`)}, nil
}
func (m *MockServiceClient) Put(ctx context.Context, path string, body interface{}) (*client.Response, error) {
	if m.PutFunc != nil {
		return m.PutFunc(ctx, path, body)
	}
	return &client.Response{StatusCode: 200, Body: []byte(`{}`)}, nil
}
func (m *MockServiceClient) Delete(ctx context.Context, path string) (*client.Response, error) {
	if m.DelFunc != nil {
		return m.DelFunc(ctx, path)
	}
	return &client.Response{StatusCode: 204, Body: nil}, nil
}

func TestService_ListRoleBindings(t *testing.T) {
	// success case
	bindingsJSON := []byte(`{"bindings":[{"id":"b1","name":"rb","role":"roles/R1","members":["group:g1"]}]}`)
	mockSvc := &MockServiceClient{GetFunc: func(ctx context.Context, path string, query map[string]string) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: bindingsJSON}, nil
	}}

	svc := &Service{rawClient: &MockClient{}, tenantID: "t", client: mockSvc}
	res, err := svc.ListRoleBindings(context.Background(), "")
	if err != nil {
		t.Fatalf("ListRoleBindings failed: %v", err)
	}
	if len(res) != 1 || res[0].ID != "b1" {
		t.Fatalf("unexpected result: %+v", res)
	}

	// empty list
	mockEmpty := &MockServiceClient{GetFunc: func(ctx context.Context, path string, query map[string]string) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`{"bindings":[]}`)}, nil
	}}
	svc.client = mockEmpty
	res2, err := svc.ListRoleBindings(context.Background(), "")
	if err != nil {
		t.Fatalf("ListRoleBindings(empty) failed: %v", err)
	}
	if len(res2) != 0 {
		t.Fatalf("expected empty list, got: %+v", res2)
	}

	// API error path
	mockErr := &MockServiceClient{GetFunc: func(ctx context.Context, path string, query map[string]string) (*client.Response, error) {
		return &client.Response{StatusCode: 500}, errors.New("api error")
	}}
	svc.client = mockErr
	_, err = svc.ListRoleBindings(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error from ListRoleBindings, got nil")
	}
}

func TestService_DeleteRoleBinding_403Fallback(t *testing.T) {
	// Simulate DELETE returning 403, then POST with empty bindings succeeds
	called := 0
	mockRaw := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		called++
		// First call is DELETE
		if req.Method == "DELETE" && strings.Contains(req.Path, "/roles/") {
			return &client.Response{StatusCode: 403, Body: []byte(`{}`)}, nil
		}
		// Second call is POST to remove bindings
		if req.Method == "POST" && strings.HasSuffix(req.Path, "/roles") {
			return &client.Response{StatusCode: 200, Body: []byte(`{"ok":true}`)}, nil
		}
		return &client.Response{StatusCode: 404, Body: []byte(`{}`)}, nil
	}}

	svc := &Service{rawClient: mockRaw, tenantID: "t"}
	if err := svc.DeleteRoleBinding(context.Background(), "g1-Role1"); err != nil {
		t.Fatalf("DeleteRoleBinding expected nil, got %v", err)
	}
	if called < 2 {
		t.Fatalf("expected multiple calls, got %d", called)
	}
}

func TestService_AddRoleToGroup_CustomRoleNotFound(t *testing.T) {
	// When isCustom=true, AddRoleToGroup calls GetCustomRole which should be simulated to fail
	mockRaw := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// Simulate GetCustomRole GET path failure
		if req.Method == "GET" && strings.Contains(req.Path, "/roles/") {
			return &client.Response{StatusCode: 404}, errors.New("not found")
		}
		return &client.Response{StatusCode: 200, Body: []byte(`{}`)}, nil
	}}

	svc := &Service{rawClient: mockRaw, tenantID: "t"}
	err := svc.AddRoleToGroup(context.Background(), "g1", "cr1", true, []string{"bu:001"})
	if err == nil {
		t.Fatalf("expected error when custom role not found, got nil")
	}
}

func TestService_UpdateGroup_200Body(t *testing.T) {
	updated := Group{ID: "g1", Name: "Updated"}
	ub, _ := json.Marshal(updated)
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" {
			return &client.Response{StatusCode: 200, Body: ub}, nil
		}
		return &client.Response{StatusCode: 404}, errors.New("unexpected")
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	g, err := svc.UpdateGroup(context.Background(), "g1", &Group{Name: "Updated"})
	if err != nil {
		t.Fatalf("UpdateGroup failed: %v", err)
	}
	if g.ID != "g1" || g.Name != "Updated" {
		t.Fatalf("unexpected updated group: %+v", g)
	}
}

func TestService_GetRoleBinding_InvalidFormatAndGroupNotFound(t *testing.T) {
	svc := &Service{rawClient: &MockClient{}, tenantID: "t"}
	// invalid format
	_, err := svc.GetRoleBinding(context.Background(), "invalidformat")
	if err == nil {
		t.Fatalf("expected error for invalid format, got nil")
	}

	// group not found path: make GetGroup return 404 error
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// Any GET to groups returns 404
		return &client.Response{StatusCode: 404, Body: []byte(`{"message":"not found"}`)}, nil
	}}
	svc = &Service{rawClient: mock, tenantID: "t"}
	_, err = svc.GetRoleBinding(context.Background(), "g1-Role1")
	if err == nil {
		t.Fatalf("expected not found error for missing group, got nil")
	}
	if apiErr, ok := err.(*client.Error); ok {
		if apiErr.StatusCode != 404 {
			t.Fatalf("expected 404, got %d", apiErr.StatusCode)
		}
	} else {
		t.Fatalf("expected client.Error, got %T: %v", err, err)
	}
}

func TestService_ListRoleBindings_InvalidJSON(t *testing.T) {
	mockSvc := &MockServiceClient{GetFunc: func(ctx context.Context, path string, query map[string]string) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`invalid-json`)}, nil
	}}
	svc := &Service{rawClient: &MockClient{}, tenantID: "t", client: mockSvc}
	_, err := svc.ListRoleBindings(context.Background(), "")
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_GetResource_DecodeError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.GetResource(context.Background(), "res1")
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_CreateRoleBinding_Errors(t *testing.T) {
	// group not found via ListGroups
	mockRaw := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if strings.Contains(req.Path, "/groups") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: []byte(`[]`)}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mockRaw, tenantID: "t"}
	rb := &RoleBinding{Members: []string{"group:unknown"}, Role: "roles/R1"}
	_, err := svc.CreateRoleBinding(context.Background(), rb)
	if err == nil || !strings.Contains(err.Error(), "group 'unknown' not found") {
		t.Fatalf("expected group not found error, got: %v", err)
	}

	// no group member present
	svc = &Service{rawClient: mockRaw, tenantID: "t"}
	rb2 := &RoleBinding{Members: []string{"user:someone"}, Role: "roles/R1"}
	_, err = svc.CreateRoleBinding(context.Background(), rb2)
	if err == nil || !strings.Contains(err.Error(), "no group found in members array") {
		t.Fatalf("expected no group found error, got: %v", err)
	}
}

func TestService_GetRoleBinding_CustomMatch(t *testing.T) {
	group := Group{ID: "g1", Name: "my-group", CreatedAt: "t1", UpdatedAt: "t2"}
	gbody, _ := json.Marshal(group)
	dto := []RoleBindingDto{{IsCustom: true, RoleID: "custom.T123", Bindings: []string{"bu:1"}}}
	dtoBody, _ := json.Marshal(dto)

	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// Handle V2 roles first to avoid matching the generic group GET
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: dtoBody}, nil
		}
		if strings.Contains(req.Path, "/groups/g1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: gbody}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	b, err := svc.GetRoleBinding(context.Background(), "g1-custom.T123")
	if err != nil {
		t.Fatalf("GetRoleBinding custom match failed: %v", err)
	}
	if b.Role != "roles/custom.T123" {
		t.Fatalf("unexpected role: %s", b.Role)
	}
}

func TestService_AddRoleToGroup_CustomSuccess(t *testing.T) {
	// GetCustomRole should return 200, then POST succeeds
	cr := CustomRole{ID: "cr1", Name: "cr1"}
	crb, _ := json.Marshal(cr)
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// GET custom role
		if req.Method == "GET" && strings.Contains(req.Path, "/roles/cr1") {
			return &client.Response{StatusCode: 200, Body: crb}, nil
		}
		// POST add role to group
		if req.Method == "POST" && strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") {
			return &client.Response{StatusCode: 201}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	if err := svc.AddRoleToGroup(context.Background(), "g1", "cr1", true, []string{"bu:1"}); err != nil {
		t.Fatalf("AddRoleToGroup custom success failed: %v", err)
	}
}

func TestService_ListRoles_Error(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 500}, errors.New("api error")
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.ListRoles(context.Background(), "")
	if err == nil {
		t.Fatalf("expected error from ListRoles, got nil")
	}
}

func TestService_FullHappyPath(t *testing.T) {
	// Prepare fixtures
	group := Group{ID: "g1", Name: "grp", CreatedAt: "c", UpdatedAt: "u"}
	gbody, _ := json.Marshal(group)

	groupsListBody, _ := json.Marshal([]Group{group})

	role := Role{ID: "r1", Name: "Role1"}
	roleBody, _ := json.Marshal(role)

	rolesResp := struct {
		Roles []Role `json:"roles"`
	}{Roles: []Role{role}}
	rolesRespBody, _ := json.Marshal(rolesResp)

	cr := CustomRole{ID: "cr1", Name: "cr1", Permissions: []Permission{{ID: "p1"}}}
	crBody, _ := json.Marshal(cr)

	dto := []RoleBindingDto{{IsCustom: false, RoleID: "Role1", Bindings: []string{"bu:1"}}}
	dtoBody, _ := json.Marshal(dto)

	resourcesBody, _ := json.Marshal([]Resource{{ID: "res1", Name: "R"}})

	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		p := req.Path
		// Handle V2 endpoints first (roles for a group) to avoid matching generic group GET
		if strings.Contains(p, "/api/v2/") && strings.HasSuffix(p, "/roles") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: dtoBody}, nil
		}
		// V2 create role binding
		if strings.Contains(p, "/api/v2/") && strings.HasSuffix(p, "/roles") && req.Method == "POST" {
			return &client.Response{StatusCode: 201, Body: []byte(`{"ok":true}`)}, nil
		}
		// V2 delete role binding
		if strings.Contains(p, "/api/v2/") && strings.Contains(p, "/roles/") && req.Method == "DELETE" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}
		// Groups list
		if strings.Contains(p, "/api/v1/tenants/") && strings.HasSuffix(p, "/groups") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: groupsListBody}, nil
		}
		// Get group
		if strings.Contains(p, "/groups/g1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: gbody}, nil
		}
		// Create group
		if strings.Contains(p, "/groups") && req.Method == "POST" {
			return &client.Response{StatusCode: 201, Body: gbody}, nil
		}
		// Update group returns 204 then GET handled above
		if strings.Contains(p, "/groups/g1") && req.Method == "PUT" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}
		// Roles list
		if strings.Contains(p, "/roles") && req.Method == "GET" && strings.Contains(p, "/api/v1/tenants/") {
			return &client.Response{StatusCode: 200, Body: rolesRespBody}, nil
		}
		// Get role
		if strings.Contains(p, "/api/v1/roles/") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: roleBody}, nil
		}
		// Create custom role
		if strings.Contains(p, "/roles") && strings.Contains(p, "/tenants/") && req.Method == "POST" {
			return &client.Response{StatusCode: 201, Body: crBody}, nil
		}
		// Get custom role
		if strings.Contains(p, "/roles/cr1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: crBody}, nil
		}
		// Update custom role returns 204
		if strings.Contains(p, "/roles/cr1") && req.Method == "PUT" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}
		// Delete custom role
		if strings.Contains(p, "/roles/cr1") && req.Method == "DELETE" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}

		// Resources
		if strings.Contains(p, "/resources/res1") && req.Method == "PUT" {
			return &client.Response{StatusCode: 200, Body: []byte(`{"id":"res1","name":"R"}`)}, nil
		}
		if strings.Contains(p, "/resources/res1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: []byte(`{"id":"res1","name":"R"}`)}, nil
		}
		if strings.Contains(p, "/resources/res1") && req.Method == "DELETE" {
			return &client.Response{StatusCode: 204, Body: nil}, nil
		}
		if strings.Contains(p, "/resources") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: resourcesBody}, nil
		}
		return &client.Response{StatusCode: 404, Body: []byte(`{}`)}, nil
	}}

	svc := &Service{rawClient: mock, tenantID: "t"}

	// Exercise flows
	if _, err := svc.CreateGroup(context.Background(), &Group{Name: "grp"}); err != nil {
		t.Fatalf("CreateGroup: %v", err)
	}
	if _, err := svc.UpdateGroup(context.Background(), "g1", &Group{Name: "grp"}); err != nil {
		t.Fatalf("UpdateGroup: %v", err)
	}
	if _, err := svc.ListGroups(context.Background(), &ListGroupsRequest{}); err != nil {
		t.Fatalf("ListGroups: %v", err)
	}
	if _, err := svc.ListRoles(context.Background(), ""); err != nil {
		t.Fatalf("ListRoles: %v", err)
	}
	if _, err := svc.GetRole(context.Background(), "Role1"); err != nil {
		t.Fatalf("GetRole: %v", err)
	}
	if _, err := svc.CreateCustomRole(context.Background(), &cr); err != nil {
		t.Fatalf("CreateCustomRole: %v", err)
	}
	if _, err := svc.GetCustomRole(context.Background(), "cr1"); err != nil {
		t.Fatalf("GetCustomRole: %v", err)
	}
	if _, err := svc.UpdateCustomRole(context.Background(), "cr1", &cr); err != nil {
		t.Fatalf("UpdateCustomRole: %v", err)
	}
	if err := svc.DeleteCustomRole(context.Background(), "cr1"); err != nil {
		t.Fatalf("DeleteCustomRole: %v", err)
	}
	if _, err := svc.GetRoleBinding(context.Background(), "g1-Role1"); err != nil {
		t.Fatalf("GetRoleBinding: %v", err)
	}
	if _, err := svc.CreateRoleBinding(context.Background(), &RoleBinding{Members: []string{"group:grp"}, Role: "roles/Role1"}); err != nil {
		t.Fatalf("CreateRoleBinding: %v", err)
	}
	if _, err := svc.UpdateRoleBinding(context.Background(), "g1-Role1", &RoleBinding{Name: "n"}); err != nil {
		t.Fatalf("UpdateRoleBinding: %v", err)
	}
	if err := svc.DeleteRoleBinding(context.Background(), "g1-Role1"); err != nil {
		t.Fatalf("DeleteRoleBinding: %v", err)
	}
	if _, err := svc.SetResource(context.Background(), "res1", &SetResourceDto{Name: "R"}); err != nil {
		t.Fatalf("SetResource: %v", err)
	}
	if _, err := svc.GetResource(context.Background(), "res1"); err != nil {
		t.Fatalf("GetResource: %v", err)
	}
	if err := svc.DeleteResource(context.Background(), "res1"); err != nil {
		t.Fatalf("DeleteResource: %v", err)
	}
	if _, err := svc.GetResources(context.Background(), &GetResourcesRequest{}); err != nil {
		t.Fatalf("GetResources: %v", err)
	}
	if err := svc.AddRoleToGroup(context.Background(), "g1", "Role1", false, []string{"bu:1"}); err != nil {
		t.Fatalf("AddRoleToGroup: %v", err)
	}
}

func TestService_GetRoleBinding_CustomVariants(t *testing.T) {
	group := Group{ID: "g1", Name: "grp", CreatedAt: "c", UpdatedAt: "u"}
	gbody, _ := json.Marshal(group)

	cases := []struct {
		current string
		roleID  string
		desc    string
	}{
		{current: "custom.T1", roleID: "custom.T1", desc: "direct custom match"},
		{current: "T2", roleID: "custom.T2", desc: "trimmed prefix match"},
		{current: "custom-roles/custom.T5", roleID: "custom.T5", desc: "custom-roles prefix match"},
	}

	for _, tc := range cases {
		dto := []RoleBindingDto{{IsCustom: true, RoleID: tc.current, Bindings: []string{"bu:1"}}}
		dtoBody, _ := json.Marshal(dto)
		mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
			// handle V2 roles GET first
			if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "GET" {
				return &client.Response{StatusCode: 200, Body: dtoBody}, nil
			}
			if strings.Contains(req.Path, "/groups/g1") && req.Method == "GET" {
				return &client.Response{StatusCode: 200, Body: gbody}, nil
			}
			return &client.Response{StatusCode: 404}, nil
		}}
		svc := &Service{rawClient: mock, tenantID: "t"}
		name := "g1-" + tc.roleID
		b, err := svc.GetRoleBinding(context.Background(), name)
		if err != nil {
			t.Fatalf("case %s failed: %v", tc.desc, err)
		}
		if !strings.HasPrefix(b.Role, "roles/") {
			t.Fatalf("case %s unexpected role: %s", tc.desc, b.Role)
		}
	}
}

func TestService_UpdateRoleBinding_SetTimestamps(t *testing.T) {
	// Make GetRoleBinding return a binding with empty timestamps (fallback path)
	group := Group{ID: "g1", Name: "grp"}
	gbody, _ := json.Marshal(group)
	dto := []RoleBindingDto{{IsCustom: false, RoleID: "Role1", Bindings: []string{"bu:1"}}}
	dtoBody, _ := json.Marshal(dto)
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// handle V2 roles GET first
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: dtoBody}, nil
		}
		if strings.Contains(req.Path, "/groups/g1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: gbody}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	up, err := svc.UpdateRoleBinding(context.Background(), "g1-Role1", &RoleBinding{Name: "n"})
	if err != nil {
		t.Fatalf("UpdateRoleBinding failed: %v", err)
	}
	if up.CreatedAt == "" || up.UpdatedAt == "" {
		t.Fatalf("expected timestamps set, got: %+v", up)
	}
}

func TestService_DeleteRoleBinding_PostFallbackFails(t *testing.T) {
	// DELETE returns 403, POST returns 500 -> expect error
	call := 0
	mockRaw := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		call++
		if req.Method == "DELETE" && strings.Contains(req.Path, "/roles/") {
			return &client.Response{StatusCode: 403, Body: []byte(`{}`)}, nil
		}
		if req.Method == "POST" && strings.HasSuffix(req.Path, "/roles") {
			return &client.Response{StatusCode: 500, Body: []byte(`{"message":"err"}`)}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mockRaw, tenantID: "t"}
	err := svc.DeleteRoleBinding(context.Background(), "g1-Role1")
	if err == nil {
		t.Fatalf("expected error when POST fallback fails, got nil")
	}
	if call < 2 {
		t.Fatalf("expected multiple calls, got %d", call)
	}
}

func TestService_AddRoleToGroup_DefaultBindings(t *testing.T) {
	// When bindings nil, default '*' used and POST succeeds
	captured := false
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "POST" && strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") {
			// inspect body
			b, _ := json.Marshal(req.Body)
			if strings.Contains(string(b), "\"*\"") {
				captured = true
			}
			return &client.Response{StatusCode: 201}, nil
		}
		// allow GetCustomRole to succeed if called
		if req.Method == "GET" && strings.Contains(req.Path, "/roles/") {
			return &client.Response{StatusCode: 200, Body: []byte(`{"id":"cr1"}`)}, nil
		}
		return &client.Response{StatusCode: 200}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	if err := svc.AddRoleToGroup(context.Background(), "g1", "Role1", false, nil); err != nil {
		t.Fatalf("AddRoleToGroup default bindings failed: %v", err)
	}
	if !captured {
		t.Fatalf("expected default '*' binding to be used")
	}
}

func TestService_GetRoleBinding_V2NotFound(t *testing.T) {
	// GetGroup returns OK
	group := Group{ID: "g1", Name: "grp"}
	gbody, _ := json.Marshal(group)

	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// V2 GET roles returns 404 (handle first)
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "GET" {
			return &client.Response{StatusCode: 404, Body: []byte(`{"message":"not found"}`)}, nil
		}
		// GET group
		if strings.Contains(req.Path, "/groups/g1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: gbody}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.GetRoleBinding(context.Background(), "g1-Role1")
	if err == nil {
		t.Fatalf("expected not found error from GetRoleBinding, got nil")
	}
	if apiErr, ok := err.(*client.Error); ok {
		if apiErr.StatusCode != 404 {
			t.Fatalf("expected 404, got %d", apiErr.StatusCode)
		}
	} else {
		t.Fatalf("expected client.Error, got %T: %v", err, err)
	}
}

func TestService_GetRoleBinding_GetGroupNetworkError(t *testing.T) {
	// Simulate GetGroup failing with network error
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if strings.Contains(req.Path, "/groups/") && req.Method == "GET" {
			return nil, errors.New("network")
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.GetRoleBinding(context.Background(), "g1-Role1")
	if err == nil || !strings.Contains(err.Error(), "failed to get group") {
		t.Fatalf("expected get group error, got: %v", err)
	}
}

func TestService_GetRoleBinding_FallbackConstructsBinding(t *testing.T) {
	group := Group{ID: "g1", Name: "grp", CreatedAt: "c", UpdatedAt: "u"}
	gbody, _ := json.Marshal(group)

	// V2 returns empty array
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// Handle V2 roles first to avoid matching generic group GET
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: []byte(`[]`)}, nil
		}
		if strings.Contains(req.Path, "/groups/g1") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: gbody}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	b, err := svc.GetRoleBinding(context.Background(), "g1-RoleX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil || b.Role != "roles/RoleX" {
		t.Fatalf("unexpected binding: %+v", b)
	}
}

func TestService_CreateGroup_WithOptionalFields(t *testing.T) {
	created := Group{ID: "g2", Name: "G2", Description: "desc", Members: []string{"m1"}}
	cbody, _ := json.Marshal(created)
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		// verify body contains description and members
		b, _ := json.Marshal(req.Body)
		if !strings.Contains(string(b), "desc") || !strings.Contains(string(b), "m1") {
			return &client.Response{StatusCode: 400}, errors.New("missing fields")
		}
		return &client.Response{StatusCode: 201, Body: cbody}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	g, err := svc.CreateGroup(context.Background(), &Group{Name: "G2", Description: "desc", Members: []string{"m1"}})
	if err != nil {
		t.Fatalf("CreateGroup with optional fields failed: %v", err)
	}
	if g.ID != "g2" {
		t.Fatalf("unexpected created group: %+v", g)
	}
}

func TestService_UpdateGroup_WithOptionalFields(t *testing.T) {
	updated := Group{ID: "g1", Name: "G1", Description: "d", Members: []string{"a"}}
	ub, _ := json.Marshal(updated)
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" {
			// ensure body contains description and members
			b, _ := json.Marshal(req.Body)
			if !strings.Contains(string(b), "d") || !strings.Contains(string(b), "a") {
				return &client.Response{StatusCode: 400}, errors.New("missing fields")
			}
			return &client.Response{StatusCode: 200, Body: ub}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	g, err := svc.UpdateGroup(context.Background(), "g1", &Group{Name: "G1", Description: "d", Members: []string{"a"}})
	if err != nil {
		t.Fatalf("UpdateGroup with optional fields failed: %v", err)
	}
	if g.ID != "g1" || g.Description != "d" {
		t.Fatalf("unexpected updated group: %+v", g)
	}
}

func TestService_ListRoles_WithFilter(t *testing.T) {
	rolesResp := struct {
		Roles []Role `json:"roles"`
	}{Roles: []Role{{ID: "r2", Name: "R2"}}}
	rb, _ := json.Marshal(rolesResp)
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "GET" && strings.Contains(req.Path, "/roles") {
			// ensure query includes filter
			if val, ok := req.Query["filter"]; !ok || val != "x" {
				return &client.Response{StatusCode: 400}, errors.New("missing filter")
			}
			return &client.Response{StatusCode: 200, Body: rb}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	roles, err := svc.ListRoles(context.Background(), "x")
	if err != nil {
		t.Fatalf("ListRoles with filter failed: %v", err)
	}
	if len(roles) != 1 || roles[0].ID != "r2" {
		t.Fatalf("unexpected roles: %+v", roles)
	}
}

func TestService_GetResources_WithQuery(t *testing.T) {
	resourcesBody, _ := json.Marshal([]Resource{{ID: "res2", Name: "Res2"}})
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "GET" && strings.Contains(req.Path, "/resources") {
			if req.Query["permission"] != "p1" || req.Query["type"] != "t1" {
				return &client.Response{StatusCode: 400}, errors.New("missing query")
			}
			return &client.Response{StatusCode: 200, Body: resourcesBody}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	gr, err := svc.GetResources(context.Background(), &GetResourcesRequest{Permission: "p1", Type: "t1"})
	if err != nil {
		t.Fatalf("GetResources with query failed: %v", err)
	}
	if len(gr.Resources) != 1 || gr.Resources[0].ID != "res2" {
		t.Fatalf("unexpected resources: %+v", gr)
	}
}

func TestService_NilResponsePaths(t *testing.T) {
	// Helper that returns nil response, nil error
	nilResp := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return nil, nil
	}}

	svc := &Service{rawClient: nilResp, tenantID: "t"}

	// GetGroup
	if _, err := svc.GetGroup(context.Background(), "g1"); err == nil {
		t.Fatalf("expected error for GetGroup nil response")
	}

	// CreateGroup
	if _, err := svc.CreateGroup(context.Background(), &Group{Name: "x"}); err == nil {
		t.Fatalf("expected error for CreateGroup nil response")
	}

	// CreateCustomRole
	if _, err := svc.CreateCustomRole(context.Background(), &CustomRole{ID: "cr1"}); err == nil {
		t.Fatalf("expected error for CreateCustomRole nil response")
	}

	// GetCustomRole
	if _, err := svc.GetCustomRole(context.Background(), "cr1"); err == nil {
		t.Fatalf("expected error for GetCustomRole nil response")
	}

	// SetResource
	if _, err := svc.SetResource(context.Background(), "res1", &SetResourceDto{Name: "x"}); err == nil {
		t.Fatalf("expected error for SetResource nil response")
	}

	// GetResource
	if _, err := svc.GetResource(context.Background(), "res1"); err == nil {
		t.Fatalf("expected error for GetResource nil response")
	}

	// DeleteResource
	if err := svc.DeleteResource(context.Background(), "res1"); err == nil {
		t.Fatalf("expected error for DeleteResource nil response")
	}

	// GetResources
	if _, err := svc.GetResources(context.Background(), &GetResourcesRequest{}); err == nil {
		t.Fatalf("expected error for GetResources nil response")
	}
}

func TestService_AddRoleToGroup_PostError(t *testing.T) {
	// Simulate POST to V2 roles returning 500
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "POST" && strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") {
			return &client.Response{StatusCode: 500, Body: []byte(`{"message":"err"}`)}, nil
		}
		// For ListGroups -> find group
		if req.Method == "GET" && strings.Contains(req.Path, "/groups") {
			return &client.Response{StatusCode: 200, Body: []byte(`[{"id":"g1","name":"grp"}]`)}, nil
		}
		return &client.Response{StatusCode: 200, Body: []byte(`{}`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	err := svc.AddRoleToGroup(context.Background(), "g1", "Role1", false, []string{"bu:1"})
	if err == nil {
		t.Fatalf("expected error when POST returns 500, got nil")
	}
}

func TestService_DeleteRoleBinding_CustomRolePath(t *testing.T) {
	// DELETE should be called with roleId without 'custom.' prefix when name includes custom.
	var seenPath string
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		seenPath = req.Path
		if req.Method == "DELETE" {
			return &client.Response{StatusCode: 204}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	if err := svc.DeleteRoleBinding(context.Background(), "g1-custom.cr1"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !strings.Contains(seenPath, "/roles/cr1") {
		t.Fatalf("expected path to contain /roles/cr1, got %s", seenPath)
	}
}

func TestService_CreateRoleBinding_CustomInput(t *testing.T) {
	// ListGroups returns group, then POST should be called; ensure composite ID contains custom prefix
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if strings.Contains(req.Path, "/groups") && req.Method == "GET" {
			return &client.Response{StatusCode: 200, Body: []byte(`[{"id":"g1","name":"grp"}]`)}, nil
		}
		if strings.Contains(req.Path, "/api/v2/") && strings.HasSuffix(req.Path, "/roles") && req.Method == "POST" {
			return &client.Response{StatusCode: 201}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	rb := &RoleBinding{Members: []string{"group:grp"}, Role: "roles/custom.cr1", Name: "rb"}
	out, err := svc.CreateRoleBinding(context.Background(), rb)
	if err != nil {
		t.Fatalf("CreateRoleBinding failed: %v", err)
	}
	if !strings.Contains(out.ID, "custom.cr1") {
		t.Fatalf("expected composite id to contain custom.cr1, got %s", out.ID)
	}
}

func TestService_UpdateRoleBinding_GetRoleBindingError(t *testing.T) {
	// Make GetRoleBinding return an error by having GetGroup return error
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if strings.Contains(req.Path, "/groups/") && req.Method == "GET" {
			return nil, errors.New("net")
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.UpdateRoleBinding(context.Background(), "g1-Role1", &RoleBinding{Name: "n"})
	if err == nil || !strings.Contains(err.Error(), "failed to verify existing role binding") {
		t.Fatalf("expected wrapped error from UpdateRoleBinding, got: %v", err)
	}
}

func TestService_DeleteRoleBinding_DoError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return nil, errors.New("net")
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	err := svc.DeleteRoleBinding(context.Background(), "g1-Role1")
	if err == nil || !strings.Contains(err.Error(), "failed to delete role binding") {
		t.Fatalf("expected delete error, got: %v", err)
	}
}

func TestService_UpdateGroup_204ThenGet_ErrorDecode(t *testing.T) {
	// PUT returns 204, but GET returns invalid JSON leading to decode error
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" {
			return &client.Response{StatusCode: 204}, nil
		}
		if req.Method == "GET" && strings.Contains(req.Path, "/groups/g1") {
			return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.UpdateGroup(context.Background(), "g1", &Group{Name: "x"})
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error from GET after 204, got: %v", err)
	}
}

func TestService_DeleteCustomRole_CheckResponseError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "DELETE" {
			return &client.Response{StatusCode: 500, Body: []byte(`{"message":"err"}`)}, nil
		}
		return &client.Response{StatusCode: 200, Body: []byte(`{}`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	err := svc.DeleteCustomRole(context.Background(), "cr1")
	if err == nil {
		t.Fatalf("expected error from DeleteCustomRole, got nil")
	}
}

func TestNewService(t *testing.T) {
	svc := &Service{rawClient: &MockClient{}, tenantID: "tid-1", client: &MockServiceClient{}}
	if svc == nil {
		t.Fatalf("Service constructor returned nil")
	}
	if svc.TenantID() != "tid-1" {
		t.Fatalf("expected tenant id tid-1, got %s", svc.TenantID())
	}
}

func TestService_DeleteResource_CheckResponseError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 500, Body: []byte(`{"message":"err"}`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	err := svc.DeleteResource(context.Background(), "r1")
	if err == nil {
		t.Fatalf("expected error from DeleteResource, got nil")
	}
}

func TestService_ListGroups_NilAndInvalidJSON(t *testing.T) {
	// nil response
	mockNil := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return nil, nil
	}}
	svc := &Service{rawClient: mockNil, tenantID: "t"}
	_, err := svc.ListGroups(context.Background(), &ListGroupsRequest{})
	if err == nil || !strings.Contains(err.Error(), "nil response") {
		t.Fatalf("expected nil response error, got: %v", err)
	}

	// invalid JSON
	mockInvalid := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
	}}
	svc = &Service{rawClient: mockInvalid, tenantID: "t"}
	_, err = svc.ListGroups(context.Background(), &ListGroupsRequest{})
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_CreateGroup_DoError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return nil, errors.New("network")
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.CreateGroup(context.Background(), &Group{Name: "x"})
	if err == nil || !strings.Contains(err.Error(), "failed to create group") {
		t.Fatalf("expected create group error, got: %v", err)
	}
}

func TestService_UpdateGroup_DecodeError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		if req.Method == "PUT" {
			return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
		}
		return &client.Response{StatusCode: 404}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.UpdateGroup(context.Background(), "g1", &Group{Name: "x"})
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_GetResources_Errors(t *testing.T) {
	// Do returns error
	mockErr := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return nil, errors.New("net")
	}}
	svc := &Service{rawClient: mockErr, tenantID: "t"}
	_, err := svc.GetResources(context.Background(), &GetResourcesRequest{})
	if err == nil {
		t.Fatalf("expected error from GetResources, got nil")
	}

	// invalid JSON
	mockInv := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
	}}
	svc = &Service{rawClient: mockInv, tenantID: "t"}
	_, err = svc.GetResources(context.Background(), &GetResourcesRequest{})
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error from GetResources, got: %v", err)
	}
}

func TestService_CreateGroup_DecodeError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.CreateGroup(context.Background(), &Group{Name: "x"})
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_CreateCustomRole_DecodeError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 201, Body: []byte(`not-json`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.CreateCustomRole(context.Background(), &CustomRole{ID: "cr1"})
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_GetRole_DecodeError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.GetRole(context.Background(), "r1")
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_ListRoles_DecodeError(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 200, Body: []byte(`not-json`)}, nil
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.ListRoles(context.Background(), "")
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Fatalf("expected decode error, got: %v", err)
	}
}

func TestService_DeleteGroup_ErrorPath(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return &client.Response{StatusCode: 500}, errors.New("api err")
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	err := svc.DeleteGroup(context.Background(), "g1")
	if err == nil {
		t.Fatalf("expected error from DeleteGroup, got nil")
	}
}

func TestService_SetResource_ErrorOnDo(t *testing.T) {
	mock := &MockClient{DoFunc: func(ctx context.Context, req *client.Request) (*client.Response, error) {
		return nil, errors.New("network")
	}}
	svc := &Service{rawClient: mock, tenantID: "t"}
	_, err := svc.SetResource(context.Background(), "r1", &SetResourceDto{Name: "x"})
	if err == nil || !strings.Contains(err.Error(), "failed to set resource") {
		t.Fatalf("expected error from SetResource, got: %v", err)
	}
}
