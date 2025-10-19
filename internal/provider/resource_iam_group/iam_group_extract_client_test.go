package resource_iam_group

import (
	"net/http"
	"testing"
)

func Test_extractAPIClientFields_NonPointerAndWrongTypes(t *testing.T) {
	// non-pointer struct with correct types
	type Prov struct {
		BaseURL    string
		TenantID   string
		HTTPClient *http.Client
	}
	p := Prov{BaseURL: "https://api.example", TenantID: "t1", HTTPClient: &http.Client{}}
	if extractAPIClientFields(p) == nil {
		t.Fatalf("expected non-nil for non-pointer struct with correct fields")
	}

	// wrong types
	type WrongTypes struct {
		BaseURL    int
		TenantID   string
		HTTPClient *http.Client
	}
	if extractAPIClientFields(&WrongTypes{BaseURL: 1, TenantID: "t1", HTTPClient: &http.Client{}}) != nil {
		t.Fatalf("expected nil for struct with wrong field types")
	}
}
