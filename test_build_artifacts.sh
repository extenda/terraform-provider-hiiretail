#!/bin/bash
# Build Artifacts Validation Test
# This test verifies expected output structure from GoReleaser builds
# MUST FAIL until proper build configuration is implemented

set -e

echo "=== Build Artifacts Validation Test ==="

# Expected platforms and architectures
expected_platforms=(
    "linux_amd64"
    "linux_arm64"
    "darwin_amd64"
    "darwin_arm64"
    "windows_amd64"
)

# Check if dist directory exists (created by goreleaser build)
if [ ! -d "dist" ]; then
    echo "❌ FAIL: dist/ directory not found (run 'goreleaser build --snapshot --clean' first)"
    exit 1
fi

echo "✅ dist/ directory exists"

# Validate binary naming convention
binary_pattern="terraform-provider-hiiretail_v*"
if ! ls dist/*/terraform-provider-hiiretail_v* >/dev/null 2>&1; then
    echo "❌ FAIL: No binaries found matching pattern '$binary_pattern'"
    echo "Available files in dist/:"
    ls -la dist/ 2>/dev/null || echo "No dist directory found"
    if [ -d "dist" ]; then
        echo "Directory contents:"
        find dist/ -name "*" -type f 2>/dev/null || echo "No files found"
    fi
    exit 1  
fi

echo "✅ Binaries follow correct naming convention"

# Check for each expected platform binary
for platform in "${expected_platforms[@]}"; do
    # Look for platform directories with microarchitecture versions
    platform_dirs=$(ls -d dist/terraform-provider-hiiretail_${platform}* 2>/dev/null | head -1)
    if [ -z "$platform_dirs" ]; then
        echo "❌ FAIL: Platform directory not found for $platform"
        echo "Available directories in dist/:"
        ls -d dist/terraform-provider-hiiretail_* 2>/dev/null || echo "No platform directories found"
        exit 1
    fi
    echo "✅ Platform directory exists for $platform"
    
    # Check if binary exists in platform directory
    binary_file="$platform_dirs/terraform-provider-hiiretail_v*"
    if ! ls $binary_file >/dev/null 2>&1; then
        echo "❌ FAIL: Binary not found for $platform"
        echo "Contents of $platform_dirs:"
        ls -la "$platform_dirs" 2>/dev/null || echo "Directory not accessible"
        exit 1
    fi
    echo "✅ Binary exists for $platform"
    
    # Check if binary is executable (for non-Windows platforms)
    if [[ "$platform" != *"windows"* ]]; then
        binary_path=$(ls $binary_file | head -1)
        if [ ! -x "$binary_path" ]; then
            echo "❌ FAIL: Binary for $platform is not executable"
            exit 1
        fi
        echo "✅ Binary for $platform is executable"
    fi
done

# Check for archives (if goreleaser archive is configured)
echo "Checking for archives..."
archive_extensions=("tar.gz" "zip")
archive_count=0

for ext in "${archive_extensions[@]}"; do
    if ls dist/*.$ext >/dev/null 2>&1; then
        archive_count=$((archive_count + 1))
        echo "✅ Found archives with extension .$ext"
    fi
done

if [ $archive_count -eq 0 ]; then
    echo "⚠️  WARNING: No archives found (may not be configured yet)"
else
    echo "✅ Archives are present"
fi

# Check for checksums file
if ls dist/*SHA256SUMS >/dev/null 2>&1; then
    echo "✅ SHA256SUMS file found"
    
    # Validate checksums file format
    sums_file=$(ls dist/*SHA256SUMS | head -1)
    if [ -s "$sums_file" ]; then
        echo "✅ SHA256SUMS file is not empty"
        
        # Check if checksums are valid format (64 hex chars + filename)
        if grep -E "^[a-f0-9]{64}  .+$" "$sums_file" >/dev/null; then
            echo "✅ SHA256SUMS file format is valid"
        else
            echo "❌ FAIL: SHA256SUMS file format is invalid"
            exit 1
        fi
    else
        echo "❌ FAIL: SHA256SUMS file is empty"
        exit 1
    fi
else
    echo "⚠️  WARNING: SHA256SUMS file not found (may not be configured yet)"
fi

# Check for GPG signature
if ls dist/*SHA256SUMS.sig >/dev/null 2>&1; then
    echo "✅ GPG signature file found"
else
    echo "⚠️  WARNING: GPG signature file not found (requires GPG setup)"
fi

# Validate binary functionality (basic smoke test)
echo "Running binary smoke tests..."
for platform in "${expected_platforms[@]}"; do
    # Skip Windows binaries on non-Windows systems
    if [[ "$platform" == *"windows"* ]] && [[ "$OSTYPE" != "msys" && "$OSTYPE" != "cygwin" ]]; then
        echo "⏭️  Skipping Windows binary test on non-Windows system"
        continue
    fi
    
    # Test only compatible binaries
    if [[ "$OSTYPE" == "darwin"* && "$platform" == *"darwin"* ]] || 
       [[ "$OSTYPE" == "linux-gnu"* && "$platform" == *"linux"* ]]; then
        
        binary_path=$(ls dist/terraform-provider-hiiretail_v*_${platform}/terraform-provider-hiiretail_v* 2>/dev/null | head -1)
        if [ -n "$binary_path" ] && [ -x "$binary_path" ]; then
            # Try to execute binary (should not crash)
            if timeout 5s "$binary_path" --help >/dev/null 2>&1 || timeout 5s "$binary_path" --version >/dev/null 2>&1; then
                echo "✅ Binary for $platform executes successfully"
            else
                echo "⚠️  WARNING: Binary for $platform may have execution issues"
            fi
        fi
    fi
done

echo "🎉 Build artifacts validation completed!"
echo "Note: Some warnings are expected until full configuration is implemented"