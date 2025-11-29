package resource_iam_group

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

type errTransport struct{}

func (e *errTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("connection refused")
}

func Test_makeAPIRequest_RoundTripErrorMapsToGeneric(t *testing.T) {
	r := NewIamGroupResource().(*IamGroupResource)
	r.client = &http.Client{Transport: &errTransport{}}
	r.baseURL = "https://api.test"
	r.tenantID = "t1"

	resp, err := r.makeAPIRequest(context.Background(), http.MethodGet, "/x", nil)
	if err == nil {
		t.Fatalf("expected error for transport failure")
	}
	if resp != nil {
		t.Fatalf("expected nil resp on transport failure")
	}
}
