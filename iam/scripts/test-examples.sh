#!/bin/bash

# Test all Terraform examples execute successfully
# This script runs terraform plan on all example configurations

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🧪 Testing Terraform Examples..."

# Build the provider first
echo "🔨 Building provider..."
cd "$PROJECT_ROOT"
go build -o terraform-provider-hiiretail

# Create temporary terraform.tfvars for testing
TEST_TFVARS="$PROJECT_ROOT/test.tfvars"
cat > "$TEST_TFVARS" << EOF
client_id     = "test-client-id"
client_secret = "test-client-secret"
tenant_id     = "CIR7nQwtS0rA6t0S6ejd"
EOF

# Test provider example
echo "📋 Testing provider example..."
cd "$PROJECT_ROOT/examples/provider"
terraform init -upgrade
terraform plan -var-file="$TEST_TFVARS" -out=plan.out
echo "✅ Provider example plan successful"
rm -f plan.out terraform.tfstate* .terraform.lock.hcl
rm -rf .terraform/

# Test resource examples
RESOURCES=(
    "iam_custom_role"
    "iam_group"
    "iam_resource"
    "iam_role_binding"
)

for resource in "${RESOURCES[@]}"; do
    echo "📋 Testing ${resource} example..."
    cd "$PROJECT_ROOT/examples/resources/${resource}"
    
    terraform init -upgrade
    terraform plan -var-file="$TEST_TFVARS" -out=plan.out
    echo "✅ ${resource} example plan successful"
    
    # Cleanup
    rm -f plan.out terraform.tfstate* .terraform.lock.hcl
    rm -rf .terraform/
done

# Cleanup test file
rm -f "$TEST_TFVARS"

echo ""
echo "🎉 All examples tested successfully!"
echo ""
echo "Example Test Summary:"
echo "✅ Provider configuration example"
echo "✅ Custom role example"
echo "✅ Group example"
echo "✅ Resource example"
echo "✅ Role binding example"
echo ""
echo "All examples are functional and ready for use! 🚀"