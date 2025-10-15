package resource_iam_resource

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/client"
)

func Test_validateResourceID_cases(t *testing.T) {
	cases := []struct {
		id      string
		wantErr bool
	}{
		{"", true},
		{".", true},
		{"..", true},
		{"has/slash", true},
		{"has__double", true},
		{string(make([]byte, 1600)), true},
		{"valid-id-123", false},
	}
	for _, c := range cases {
		err := validateResourceID(c.id)
		if (err != nil) != c.wantErr {
			t.Fatalf("id=%q: unexpected error presence: %v", c.id, err)
		}
	}
}

func Test_handleAPIError_various(t *testing.T) {
	tests := []struct {
		err       error
		wantTitle string
	}{
		{errors.New("400 Bad Request: invalid"), "Invalid Request"},
		{errors.New("401 Unauthorized"), "Authentication Failed"},
		{errors.New("403 Forbidden"), "Permission Denied"},
		{errors.New("404 Not Found"), "Resource Not Found"},
		{errors.New("409 Conflict"), "Resource Conflict"},
		{errors.New("429 Too Many Requests"), "Rate Limit Exceeded"},
		{errors.New("500 Internal Server Error"), "Server Error"},
		{errors.New("502 Bad Gateway"), "Service Unavailable"},
		{errors.New("503 Service Unavailable"), "Service Maintenance"},
		{errors.New("context deadline exceeded"), "Request Timeout"},
	}
	for _, tt := range tests {
		title, detail := handleAPIError(tt.err, "op", "res1")
		if title == "" || detail == "" {
			t.Fatalf("expected non-empty title/detail for %v", tt.err)
		}
		if tt.wantTitle != "" && tt.wantTitle != title {
			t.Fatalf("unexpected title for %v: got %q want %q", tt.err, title, tt.wantTitle)
		}
	}
}

func Test_mapAPIErrorToDiagnostic_clientErrors(t *testing.T) {
	cases := []struct {
		status  int
		wantSub string
	}{
		{400, "Invalid Request"},
		{401, "Authentication Failed"},
		{403, "Access Denied"},
		{404, "Resource Not Found"},
		{409, "Resource Conflict"},
		{429, "Rate Limited"},
		{500, "Server Error"},
		{502, "Server Error"},
	}
	for _, c := range cases {
		cerr := &client.Error{StatusCode: c.status, Message: fmt.Sprintf("msg-%d", c.status)}
		title, detail := mapAPIErrorToDiagnostic(cerr, "create", "rid")
		if title == "" || detail == "" {
			t.Fatalf("expected non-empty for status %d", c.status)
		}
		if !strings.Contains(title, c.wantSub) && !strings.Contains(detail, c.wantSub) {
			t.Fatalf("expected %q in title or detail for status %d: got title=%q detail=%q", c.wantSub, c.status, title, detail)
		}
	}
}

func Test_ValidateJSONString(t *testing.T) {
	if err := ValidateJSONString(""); err != nil {
		t.Fatalf("empty string should be valid: %v", err)
	}
	if err := ValidateJSONString("{\"a\":1}"); err != nil {
		t.Fatalf("valid json reported error: %v", err)
	}
	if err := ValidateJSONString("not-json"); err == nil {
		t.Fatalf("expected error for invalid json")
	}
}
