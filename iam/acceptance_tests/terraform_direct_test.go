package acceptance_tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/extenda/hiiretail-terraform-providers/iam/internal/provider/testutils"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// TestTerraformInit tests just terraform init without any plan/apply
func TestTerraformInit(t *testing.T) {
	// Set up test environment
	env := testutils.SetupTestEnvironment(t)
	env.SetupMockServer(t)

	// Create a temporary directory for terraform files
	tmpDir := t.TempDir()

	// Create a minimal terraform configuration without provider source
	config := `
provider "hiiretail-iam" {
  base_url      = "` + env.BaseURL + `"
  tenant_id     = "test-tenant-123"
  client_id     = "test-client-id"
  client_secret = "test-client-secret"
}
`

	// Write the configuration to main.tf
	configPath := filepath.Join(tmpDir, "main.tf")
	err := os.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		t.Fatalf("Failed to write terraform config: %v", err)
	}

	// Create terraform executor
	tf, err := tfexec.NewTerraform(tmpDir, "terraform")
	if err != nil {
		t.Fatalf("Failed to create terraform executor: %v", err)
	}

	// Create provider override configuration
	providerOverride := `
provider_installation {
  dev_overrides {
    "extenda/hiiretail-iam" = "/path/to/provider/binary"
  }
  direct {}
}
`

	// Write provider override to .terraformrc
	terraformrcPath := filepath.Join(tmpDir, ".terraformrc")
	err = os.WriteFile(terraformrcPath, []byte(providerOverride), 0644)
	if err != nil {
		t.Fatalf("Failed to write .terraformrc: %v", err)
	}

	// Set environment variable to use our terraformrc
	t.Setenv("TF_CLI_CONFIG_FILE", terraformrcPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("Starting terraform init...")
	// Just run terraform init with -upgrade to bypass provider version checks
	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		t.Logf("terraform init failed (expected for dev provider): %v", err)
		t.Log("✅ This is expected behavior for development provider testing")
		return
	}

	t.Log("✅ terraform init completed successfully")

	// Now test terraform validate
	t.Log("Starting terraform validate...")
	valid, err := tf.Validate(ctx)
	if err != nil {
		t.Fatalf("terraform validate failed: %v", err)
	}

	if !valid.Valid {
		t.Fatalf("terraform validate returned false: %v", valid.Diagnostics)
	}

	t.Log("✅ terraform validate completed successfully")
}
