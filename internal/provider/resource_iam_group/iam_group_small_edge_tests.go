package resource_iam_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Test_marshalRequest_ChannelError(t *testing.T) {
	r := &IamGroupResource{}
	_, err := r.marshalRequest(make(chan int))
	if err == nil {
		t.Fatalf("expected error marshalling channel")
	}
}

func Test_contains_empty_substr_skipped(t *testing.T) {
	// Validate contains behavior for empty and edge-case substrings
	if !contains("hello", "") {
		t.Fatalf("expected empty substring to be contained in non-empty string")
	}
	if !contains("", "") {
		t.Fatalf("expected empty substring to be contained in empty string")
	}
	if contains("", "a") {
		t.Fatalf("expected non-empty substring not to be found in empty string")
	}
}

func Test_Configure_NilProviderData_NoPanic(t *testing.T) {
	r := &IamGroupResource{}
	// Construct a request with nil ProviderData
	req := resource.ConfigureRequest{ProviderData: nil}
	resp := &resource.ConfigureResponse{}
	// Should not panic
	r.Configure(context.Background(), req, resp)
	// ensure resource still unconfigured
	if r.client != nil || r.baseURL != "" || r.tenantID != "" {
		t.Fatalf("expected resource to remain unconfigured when ProviderData is nil")
	}
}

// Minimal stubs to avoid importing terraform framework types in test
// stubs removed - use real framework types in tests
