package auth

import (
	"testing"
)

func TestValidateConfig_Nil(t *testing.T) {
	if err := ValidateConfig(nil); err == nil {
		t.Fatalf("expected error validating nil config")
	}
}

func TestIsTestEnvironment_Various(t *testing.T) {
	if !IsTestEnvironment("test-tenant-123") {
		t.Fatalf("expected test tenant to be detected as test environment")
	}
	if IsTestEnvironment("prod-tenant") {
		t.Fatalf("expected prod tenant to not be test environment")
	}
}
