package client

import (
	"testing"

	authpkg "github.com/extenda/hiiretail-terraform-providers/internal/provider/shared/auth"
)

func TestBuildURL_Basic(t *testing.T) {
	ac := &authpkg.Config{TenantID: "tenant-1", TestToken: "tok"}
	c, err := New(ac, nil)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	// Build a path and ensure the base URL is used
	u := c.buildURL("/foo/bar")
	if u == nil {
		t.Fatalf("expected non-nil URL")
	}
	if u.Path == "" {
		t.Fatalf("expected non-empty path in built URL")
	}
}
