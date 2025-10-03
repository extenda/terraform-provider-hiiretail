#!/bin/bash

# Debug Verification Helper Script
# Helps verify that debug output is appearing as expected during development

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <debug_marker> [terraform_command...]"
    echo ""
    echo "Examples:"
    echo "  $0 'DEBUG Read METHOD' terraform plan"
    echo "  $0 'BUILD_VERIFICATION' terraform apply"
    echo "  $0 'CRITICAL' terraform plan -var-file=test.tfvars"
    echo ""
    echo "This script runs terraform and checks for the specified debug marker."
    echo "If the marker doesn't appear, it warns about potential build issues."
    exit 1
fi

DEBUG_MARKER="$1"
shift
TERRAFORM_CMD="$@"

if [ -z "$TERRAFORM_CMD" ]; then
    TERRAFORM_CMD="terraform plan"
fi

echo "🔍 Running: $TERRAFORM_CMD"
echo "🎯 Looking for debug marker: '$DEBUG_MARKER'"
echo ""

# Run terraform and capture output
TEMP_OUTPUT=$(mktemp)
$TERRAFORM_CMD > "$TEMP_OUTPUT" 2>&1

# Display the output
cat "$TEMP_OUTPUT"

echo ""
echo "🧪 Debug Verification Results:"

# Check if our debug marker appears
if grep -q "$DEBUG_MARKER" "$TEMP_OUTPUT"; then
    echo "✅ Debug marker '$DEBUG_MARKER' found in output!"
    echo "   Your code changes are active."
else
    echo "❌ Debug marker '$DEBUG_MARKER' NOT found in output!"
    echo ""
    echo "⚠️  This could indicate:"
    echo "   1. Your build didn't actually update the binary"
    echo "   2. The debug code isn't being executed"
    echo "   3. The debug output is being filtered out"
    echo ""
    echo "🔧 Suggested actions:"
    echo "   1. Check binary timestamp: ls -la terraform-provider-hiiretail"  
    echo "   2. Run: ./scripts/build-and-verify.sh"
    echo "   3. Verify your debug code is in the right execution path"
fi

# Check if build verification marker exists
if [ -f ".last-build-marker" ]; then
    BUILD_ID=$(grep "BUILD_ID=" ".last-build-marker" | cut -d'=' -f2)
    if grep -q "BUILD_VERIFICATION.*$BUILD_ID" "$TEMP_OUTPUT"; then
        echo "✅ Build verification marker found! Binary is active."
    else
        echo "⚠️  Build verification marker not found."
        echo "   Consider adding: fmt.Printf(\"[BUILD_VERIFICATION] Active build: $BUILD_ID\\n\")"
        echo "   to a frequently executed code path."
    fi
fi

# Clean up
rm "$TEMP_OUTPUT"

echo ""
echo "💡 Tip: Always add unique debug markers when making changes to verify they're active!"