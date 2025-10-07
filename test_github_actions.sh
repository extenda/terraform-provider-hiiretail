#!/bin/bash
# GitHub Actions Workflow Validation Test
# This test validates .github/workflows/release.yml syntax and required components
# MUST FAIL until release.yml is properly created

set -e

echo "=== GitHub Actions Workflow Validation Test ==="

# Check if .github/workflows directory exists
if [ ! -d ".github/workflows" ]; then
    echo "❌ FAIL: .github/workflows directory not found"
    exit 1
fi

echo "✅ .github/workflows directory exists"

# Check if release.yml exists
if [ ! -f ".github/workflows/release.yml" ]; then
    echo "❌ FAIL: .github/workflows/release.yml not found"
    exit 1
fi

echo "✅ release.yml workflow file exists"

# Validate YAML syntax (basic check)
if ! python3 -c "import yaml; yaml.safe_load(open('.github/workflows/release.yml'))" 2>/dev/null; then
    echo "❌ FAIL: release.yml is not valid YAML"
    exit 1
fi

echo "✅ release.yml is valid YAML"

# Check required top-level fields
required_fields=(
    "name"
    "on"
    "permissions"
    "jobs"
)

for field in "${required_fields[@]}"; do
    if ! grep -q "^$field:" .github/workflows/release.yml; then
        echo "❌ FAIL: Required field '$field' not found in release.yml"
        exit 1
    fi
    echo "✅ Required field '$field' found"
done

# Validate workflow name
if ! grep -q "name: release" .github/workflows/release.yml; then
    echo "❌ FAIL: Workflow name must be 'release'"
    exit 1
fi

echo "✅ Workflow name is correct"

# Validate trigger on version tags
if ! grep -A 5 "on:" .github/workflows/release.yml | grep -q "tags:" &&
   ! grep -A 5 "tags:" .github/workflows/release.yml | grep -q "'v\*'"; then
    echo "❌ FAIL: Workflow must trigger on version tags (v*)"
    exit 1
fi

echo "✅ Workflow triggers on version tags"

# Validate required permissions
required_permissions=("contents: write" "id-token: write")
for permission in "${required_permissions[@]}"; do
    if ! grep -A 10 "permissions:" .github/workflows/release.yml | grep -q "$permission"; then
        echo "❌ FAIL: Permission '$permission' not found"
        exit 1
    fi
    echo "✅ Permission '$permission' configured"
done

# Validate goreleaser job
if ! grep -q "goreleaser:" .github/workflows/release.yml; then
    echo "❌ FAIL: goreleaser job not found"
    exit 1
fi

echo "✅ goreleaser job configured"

# Validate Ubuntu runner
if ! grep -q "runs-on: ubuntu-latest" .github/workflows/release.yml; then
    echo "❌ FAIL: Must use ubuntu-latest runner"
    exit 1
fi

echo "✅ Ubuntu runner configured"

# Validate required actions
required_actions=(
    "actions/checkout"
    "actions/setup-go"
    "crazy-max/ghaction-import-gpg"
    "goreleaser/goreleaser-action"
)

for action in "${required_actions[@]}"; do
    if ! grep -q "uses: $action" .github/workflows/release.yml; then
        echo "❌ FAIL: Action '$action' not found"
        exit 1
    fi
    echo "✅ Action '$action' configured"
done

# Validate Go version
if ! grep -A 5 "actions/setup-go" .github/workflows/release.yml | grep -q "go-version.*1\.2[1-9]"; then
    echo "❌ FAIL: Go version must be 1.21 or higher"
    exit 1
fi

echo "✅ Go version is correctly configured"

# Validate GPG secrets usage
gpg_secrets=("GPG_PRIVATE_KEY" "PASSPHRASE")
for secret in "${gpg_secrets[@]}"; do
    if ! grep -q "\${{ secrets\.$secret }}" .github/workflows/release.yml; then
        echo "❌ FAIL: Secret '$secret' not used"
        exit 1
    fi
    echo "✅ Secret '$secret' is used"
done

# Validate environment variables
env_vars=("GITHUB_TOKEN" "GPG_FINGERPRINT")
for env_var in "${env_vars[@]}"; do
    if ! grep -A 10 "env:" .github/workflows/release.yml | grep -q "$env_var:"; then
        echo "❌ FAIL: Environment variable '$env_var' not found"
        exit 1
    fi
    echo "✅ Environment variable '$env_var' configured"
done

echo "🎉 All GitHub Actions workflow validation tests passed!"