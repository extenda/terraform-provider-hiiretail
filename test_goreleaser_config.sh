#!/bin/bash
# GoReleaser Configuration Validation Test
# This test validates .goreleaser.yml syntax and required fields
# MUST FAIL until .goreleaser.yml is properly created

set -e

echo "=== GoReleaser Configuration Validation Test ==="

# Check if .goreleaser.yml exists
if [ ! -f ".goreleaser.yml" ]; then
    echo "❌ FAIL: .goreleaser.yml not found"
    exit 1
fi

echo "✅ .goreleaser.yml file exists"

# Validate GoReleaser configuration syntax
if ! goreleaser check --config .goreleaser.yml 2>/dev/null; then
    echo "❌ FAIL: GoReleaser configuration syntax invalid"
    echo "Install GoReleaser first: brew install goreleaser"
    exit 1
fi

echo "✅ GoReleaser configuration syntax is valid"

# Check required fields in configuration
required_fields=(
    "project_name"
    "builds"
    "archives"
    "checksum"
    "signs"
    "release"
)

for field in "${required_fields[@]}"; do
    if ! grep -q "$field:" .goreleaser.yml; then
        echo "❌ FAIL: Required field '$field' not found in .goreleaser.yml"
        exit 1
    fi
    echo "✅ Required field '$field' found"
done

# Validate project name
if ! grep -q "project_name: terraform-provider-hiiretail" .goreleaser.yml; then
    echo "❌ FAIL: project_name must be 'terraform-provider-hiiretail'"
    exit 1
fi

echo "✅ Project name is correctly set"

# Validate build configuration
if ! grep -q 'binary: "{{ .ProjectName }}_v{{ .Version }}"' .goreleaser.yml; then
    echo "❌ FAIL: Binary naming template incorrect"
    exit 1
fi

echo "✅ Binary naming template is correct"

# Validate supported platforms
platforms=("linux" "darwin" "windows")
for platform in "${platforms[@]}"; do
    if ! grep -A 10 "goos:" .goreleaser.yml | grep -q "\- $platform"; then
        echo "❌ FAIL: Platform '$platform' not configured"
        exit 1
    fi
    echo "✅ Platform '$platform' configured"
done

# Validate architectures  
architectures=("amd64" "arm64")
for arch in "${architectures[@]}"; do
    if ! grep -A 10 "goarch:" .goreleaser.yml | grep -q "\- $arch"; then
        echo "❌ FAIL: Architecture '$arch' not configured"
        exit 1
    fi
    echo "✅ Architecture '$arch' configured"
done

# Validate GPG signing configuration
if ! grep -q "artifacts: checksum" .goreleaser.yml; then
    echo "❌ FAIL: GPG signing for checksums not configured"
    exit 1
fi

echo "✅ GPG signing configuration found"

# Validate GitHub release configuration
if ! grep -q "github:" .goreleaser.yml; then
    echo "❌ FAIL: GitHub release configuration not found"
    exit 1
fi

echo "✅ GitHub release configuration found"

echo "🎉 All GoReleaser configuration validation tests passed!"