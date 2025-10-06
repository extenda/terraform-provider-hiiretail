#!/bin/bash

# Documentation validation script for HiiRetail Terraform Provider
# Validates Registry compliance and documentation quality

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🔍 Validating HiiRetail Provider Documentation..."

# Check if docs directory exists
if [ ! -d "$PROJECT_ROOT/docs" ]; then
    echo "❌ Error: docs/ directory not found"
    exit 1
fi

# Check mandatory files
echo "📋 Checking mandatory files..."

if [ ! -f "$PROJECT_ROOT/docs/index.md" ]; then
    echo "❌ Error: docs/index.md (provider overview) is required by Registry"
    exit 1
else
    echo "✅ docs/index.md exists"
fi

# Check resource documentation files
echo "📋 Checking resource documentation..."

RESOURCES=(
    "iam_custom_role"
    "iam_group" 
    "iam_resource"
    "iam_role_binding"
)

for resource in "${RESOURCES[@]}"; do
    if [ ! -f "$PROJECT_ROOT/docs/resources/${resource}.md" ]; then
        echo "❌ Error: docs/resources/${resource}.md not found"
        exit 1
    else
        echo "✅ docs/resources/${resource}.md exists"
    fi
done

# Check for incorrect naming (with provider prefix)
echo "📋 Checking filename conventions..."

INCORRECT_FILES=(
    "docs/resources/hiiretail_iam_custom_role.md"
    "docs/resources/hiiretail_iam_group.md"
    "docs/resources/hiiretail_iam_resource.md"
    "docs/resources/hiiretail_iam_role_binding.md"
)

for incorrect_file in "${INCORRECT_FILES[@]}"; do
    if [ -f "$PROJECT_ROOT/$incorrect_file" ]; then
        echo "❌ Error: $incorrect_file should not exist (Registry requires no provider prefix)"
        exit 1
    fi
done

echo "✅ All filenames follow Registry convention (no provider prefix)"

# Check document sizes (Registry 500KB limit)
echo "📋 Checking document sizes..."

MAX_SIZE_KB=500
MAX_SIZE_BYTES=$((MAX_SIZE_KB * 1024))

find "$PROJECT_ROOT/docs" -name "*.md" -type f | while read -r file; do
    file_size=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
    if [ "$file_size" -gt "$MAX_SIZE_BYTES" ]; then
        echo "❌ Error: $(basename "$file") is ${file_size} bytes (exceeds Registry 500KB limit)"
        exit 1
    fi
done

echo "✅ All documents under 500KB Registry limit"

# Check examples directory structure
echo "📋 Checking examples structure..."

if [ ! -d "$PROJECT_ROOT/examples" ]; then
    echo "❌ Error: examples/ directory not found"
    exit 1
fi

if [ ! -d "$PROJECT_ROOT/examples/provider" ]; then
    echo "❌ Error: examples/provider/ directory not found"
    exit 1
fi

for resource in "${RESOURCES[@]}"; do
    if [ ! -d "$PROJECT_ROOT/examples/resources/${resource}" ]; then
        echo "❌ Error: examples/resources/${resource}/ directory not found"
        exit 1
    fi
    
    if [ ! -f "$PROJECT_ROOT/examples/resources/${resource}/main.tf" ]; then
        echo "❌ Error: examples/resources/${resource}/main.tf not found"
        exit 1
    fi
done

echo "✅ Examples directory structure is correct"

# Check YAML frontmatter in guides
echo "📋 Checking YAML frontmatter..."

if [ -f "$PROJECT_ROOT/docs/guides/authentication.md" ]; then
    if ! grep -q "^---" "$PROJECT_ROOT/docs/guides/authentication.md"; then
        echo "❌ Error: docs/guides/authentication.md missing YAML frontmatter"
        exit 1
    else
        echo "✅ Authentication guide has YAML frontmatter"
    fi
fi

# Summary
echo ""
echo "🎉 Documentation validation completed successfully!"
echo ""
echo "Registry Compliance Summary:"
echo "✅ Provider overview (docs/index.md) exists"
echo "✅ All resource docs use correct naming (no provider prefix)"
echo "✅ All documents under 500KB limit"
echo "✅ Examples directory structure complete"
echo "✅ YAML frontmatter present in guides"
echo ""
echo "Ready for Terraform Registry publication! 🚀"