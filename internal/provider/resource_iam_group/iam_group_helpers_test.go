package resource_iam_group

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// errorTransport is a test transport that always returns an error
type errorTransport struct {
	err error
}

func (e *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, e.err
}

func Test_extractAPIClientFields_NilAndInvalid(t *testing.T) {
	if extractAPIClientFields(nil) != nil {
		t.Fatalf("expected nil for nil providerData")
	}

	// wrong shape
	type BadProvider struct{ Foo string }
	if extractAPIClientFields(&BadProvider{Foo: "x"}) != nil {
		t.Fatalf("expected nil for wrong providerData shape")
	}

	// non-pointer struct (should still work)
	type Prov struct {
		BaseURL    string
		TenantID   string
		HTTPClient *http.Client
	}
	p := Prov{BaseURL: "https://api.example", TenantID: "t1", HTTPClient: &http.Client{}}
	if extractAPIClientFields(p) == nil {
		t.Fatalf("expected non-nil for non-pointer struct with correct fields")
	}

	// missing fields
	type MissingFields struct {
		BaseURL string
		// missing TenantID and HTTPClient
	}
	if extractAPIClientFields(&MissingFields{BaseURL: "test"}) != nil {
		t.Fatalf("expected nil for struct missing required fields")
	}

	// wrong field types
	type WrongTypes struct {
		BaseURL    int // should be string
		TenantID   string
		HTTPClient *http.Client
	}
	if extractAPIClientFields(&WrongTypes{BaseURL: 123, TenantID: "t1", HTTPClient: &http.Client{}}) != nil {
		t.Fatalf("expected nil for struct with wrong field types")
	}
}

func Test_extractAPIClientFields_Success(t *testing.T) {
	type Prov struct {
		BaseURL    string
		TenantID   string
		HTTPClient *http.Client
	}

	client := &http.Client{}
	p := &Prov{BaseURL: "https://api.example", TenantID: "t1", HTTPClient: client}
	got := extractAPIClientFields(p)
	if got == nil {
		t.Fatalf("expected non-nil APIClient")
	}
	if got.BaseURL != p.BaseURL || got.TenantID != p.TenantID || got.HTTPClient != client {
		t.Fatalf("extracted fields do not match")
	}
}

func Test_validateGroupData_ErrorsAndSuccess(t *testing.T) {
	r := &IamGroupResource{}

	// missing name
	data := &IamGroupModel{Name: types.StringNull()}
	if err := r.validateGroupData(context.Background(), data); err == nil {
		t.Fatalf("expected error for empty name")
	}

	// too long name
	long := strings.Repeat("a", 300)
	data = &IamGroupModel{Name: types.StringValue(long)}
	if err := r.validateGroupData(context.Background(), data); err == nil {
		t.Fatalf("expected error for too long name")
	}

	// too long description
	desc := strings.Repeat("d", 300)
	data = &IamGroupModel{Name: types.StringValue("ok"), Description: types.StringValue(desc)}
	if err := r.validateGroupData(context.Background(), data); err == nil {
		t.Fatalf("expected error for too long description")
	}

	// success
	data = &IamGroupModel{Name: types.StringValue("valid"), Description: types.StringValue("ok")}
	if err := r.validateGroupData(context.Background(), data); err != nil {
		t.Fatalf("unexpected error for valid data: %v", err)
	}
}

func Test_mapHTTPError_Messages(t *testing.T) {
	r := &IamGroupResource{}
	cases := []struct {
		code int
		want string
	}{
		{http.StatusNotFound, "group not found"},
		{http.StatusUnauthorized, "authentication failed"},
		{http.StatusForbidden, "access denied"},
		{http.StatusConflict, "group already exists"},
		{http.StatusBadRequest, "invalid request"},
		{http.StatusInternalServerError, "server error"},
		{http.StatusServiceUnavailable, "service temporarily unavailable"},
	}
	for _, c := range cases {
		err := r.mapHTTPError(c.code, fmt.Errorf("x"))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Fatalf("expected error containing %q for code %d, got %v", c.want, c.code, err)
		}
	}
}

func Test_contains_and_findInString(t *testing.T) {
	if !contains("hello world", "hello") {
		t.Fatalf("expected contains true for prefix")
	}
	if !contains("hello world", "world") {
		t.Fatalf("expected contains true for suffix")
	}
	if !contains("hello world", "lo wo") {
		t.Fatalf("expected contains true for middle")
	}
	if contains("short", "longer") {
		t.Fatalf("expected contains false when substr longer than s")
	}
}

func Test_makeAPIRequest_MissingConfig(t *testing.T) {
	r := &IamGroupResource{}
	_, err := r.makeAPIRequest(context.Background(), "GET", "/", nil)
	if err == nil || !strings.Contains(err.Error(), "resource not properly configured") {
		t.Fatalf("expected missing-config error, got %v", err)
	}
}

func Test_makeAPIRequest_NetworkError(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	// Create a client that will fail with network error
	r.client = &http.Client{Transport: &errorTransport{err: fmt.Errorf("network error")}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	_, err := r.makeAPIRequest(context.Background(), "GET", "/test", nil)
	if err == nil {
		t.Fatalf("expected network error")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("expected network error in message, got: %s", err.Error())
	}
}

func Test_marshal_and_unmarshal_roundtrip(t *testing.T) {
	r := &IamGroupResource{}
	data := map[string]string{"a": "b"}
	body, err := r.marshalRequest(data)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	resp := &http.Response{Body: http.NoBody}
	// Replace Body with a reader over the marshalled bytes
	resp.Body = io.NopCloser(strings.NewReader(string(body)))
	var out map[string]string
	if err := r.unmarshalResponse(resp, &out); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if out["a"] != "b" {
		t.Fatalf("unexpected unmarshal result: %#v", out)
	}
}

func Test_isRetryableError_various_messages(t *testing.T) {
	r := &IamGroupResource{}
	cases := []struct {
		err    error
		expect bool
	}{
		{nil, false},
		{errors.New("timeout while dialing"), true},
		{errors.New("connection refused by peer"), true},
		{errors.New("service temporarily unavailable"), true},
		{errors.New("server error: status 500"), true},
		{errors.New("authentication failed"), false},
		{errors.New("random other error"), false},
	}
	for _, c := range cases {
		ok := r.isRetryableError(c.err)
		if ok != c.expect {
			t.Fatalf("unexpected isRetryableError(%v) = %v, want %v", c.err, ok, c.expect)
		}
	}
}
