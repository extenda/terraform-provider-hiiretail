package resource_iam_group

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

// Ensure marshalRequest returns an error on unsupported types
func TestMarshalRequest_BadType(t *testing.T) {
	// create a channel which json.Marshal cannot encode
	bad := make(chan int)
	r := &IamGroupResource{}
	_, err := r.marshalRequest(bad)
	if err == nil {
		t.Fatalf("expected error marshalling unsupported type, got nil")
	}
}

// makeAPIRequest should return an error when resource is not configured
func TestMakeAPIRequest_Misconfigured(t *testing.T) {
	r := &IamGroupResource{}
	_, err := r.makeAPIRequest(context.Background(), http.MethodGet, "", nil)
	if err == nil {
		t.Fatalf("expected error when calling makeAPIRequest with empty baseURL and nil client")
	}

	// Now provide a base URL but keep the HTTP client nil inside
	r.baseURL = "https://example.com"
	_, err = r.makeAPIRequest(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatalf("expected error when calling makeAPIRequest with nil http client")
	}
}

// sanity check: unmarshalResponse returns error on invalid JSON
func TestUnmarshalResponse_InvalidJSON(t *testing.T) {
	var out map[string]interface{}
	invalid := []byte("not-json")
	resp := &http.Response{Body: io.NopCloser(bytes.NewReader(invalid))}
	r := &IamGroupResource{}
	err := r.unmarshalResponse(resp, &out)
	if err == nil {
		t.Fatalf("expected error unmarshalling invalid JSON, got nil")
	}
}
