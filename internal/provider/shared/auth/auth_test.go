package auth

import (
	"testing"
)

func TestValidateConfigNil(t *testing.T) {
	err := ValidateConfig(nil)
	if err == nil {
		t.Fatalf("expected error for nil config")
	}
}

func TestIsTestEnvironmentPatterns(t *testing.T) {
	if !IsTestEnvironment("test-tenant-123") {
		t.Fatalf("expected test environment for tenant containing 'test'")
	}
	if IsTestEnvironment("production-1") {
		t.Fatalf("did not expect production to be test environment")
	}
}
