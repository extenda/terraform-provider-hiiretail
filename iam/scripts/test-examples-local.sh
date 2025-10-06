#!/bin/bash

# Test all Terraform examples using local development provider
# This script runs terraform plan on all example configurations

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🧪 Testing Terraform Examples (Local Development Mode)..."

# Build the provider first
echo "🔨 Building provider..."
cd "$PROJECT_ROOT"
go build -o terraform-provider-hiiretail

# Create test terraform.tfvars
TEST_TFVARS="$PROJECT_ROOT/test.tfvars"
cat > "$TEST_TFVARS" << EOF
client_id     = "b3duZXI6IHNoYXluZQpzZnc6IGhpaXRmQDAuMUBDSVI3blF3dFMwckE2dDBTNmVqZAp0aWQ6IENJUjduUXd0UzByQTZ0MFM2ZWpkCg"
client_secret = "726143f664f0a38efa96abe33bc0a7487d745ee725171101231c454ea9faa1ba"
tenant_id     = "CIR7nQwtS0rA6t0S6ejd"
EOF

# Test resource examples (these have all dependencies in one file)
RESOURCES=(
    "iam_custom_role"
    "iam_group"
    "iam_resource"
    "iam_role_binding"
)

for resource in "\${RESOURCES[@]}"; do
    echo "📋 Testing \${resource} example..."
    cd "$PROJECT_ROOT/examples/resources/\${resource}"
    
    # Copy .terraformrc for local development
    cp "$PROJECT_ROOT/.terraformrc" .
    
    export TF_CLI_CONFIG_FILE=".terraformrc"
    terraform plan -var-file="$TEST_TFVARS" > plan.out 2>&1 || {
        echo "❌ \${resource} example failed:"
        cat plan.out
        exit 1
    }
    echo "✅ \${resource} example plan successful"
    
    # Cleanup
    rm -f plan.out .terraformrc terraform.tfstate* .terraform.lock.hcl
    rm -rf .terraform/
    unset TF_CLI_CONFIG_FILE
done

# Cleanup test file
rm -f "$TEST_TFVARS"

echo ""
echo "🎉 All examples tested successfully with local provider!"
echo ""
echo "Example Test Summary:"
echo "✅ Custom role example"
echo "✅ Group example"
echo "✅ Resource example"
echo "✅ Role binding example"
echo ""
echo "All examples are syntactically valid and ready for use! 🚀"