#!/bin/bash

# Build and Verification Script for HiiRetail Terraform Provider
# This script codifies lessons learned about ensuring builds are working properly

set -e

PROVIDER_BINARY="terraform-provider-hiiretail"
BUILD_MARKER_FILE=".last-build-marker"

echo "🔨 Building HiiRetail Terraform Provider..."

# Clean any existing binary to force rebuild
if [ -f "$PROVIDER_BINARY" ]; then
    echo "📦 Removing existing binary to ensure clean build..."
    rm "$PROVIDER_BINARY"
fi

# Record build start time
BUILD_START=$(date +%s)
echo "⏰ Build started at: $(date)"

# Generate unique build ID
BUILD_ID=$(uuidgen | cut -d'-' -f1)

# Inject build ID into source code before building
echo "🏷️  Injecting build ID: $BUILD_ID"
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS version
    sed -i '.bak' "s/PLACEHOLDER_BUILD_ID/$BUILD_ID/g" internal/provider/provider.go
else
    # Linux version
    sed -i.bak "s/PLACEHOLDER_BUILD_ID/$BUILD_ID/g" internal/provider/provider.go
fi

# Build the provider
echo "🚀 Running go build..."
if ! go build -v -o "$PROVIDER_BINARY" .; then
    echo "❌ Build failed, restoring source code..."
    mv internal/provider/provider.go.bak internal/provider/provider.go
    exit 1
fi

# Restore original source code (remove build ID)
echo "🔄 Restoring source code..."
mv internal/provider/provider.go.bak internal/provider/provider.go

# Verify binary was created
if [ ! -f "$PROVIDER_BINARY" ]; then
    echo "❌ ERROR: Binary was not created! Build failed."
    exit 1
fi

# Check binary timestamp
BINARY_TIMESTAMP=$(stat -f "%m" "$PROVIDER_BINARY" 2>/dev/null || stat -c "%Y" "$PROVIDER_BINARY" 2>/dev/null)
BUILD_END=$(date +%s)

echo "✅ Binary created successfully"
echo "📅 Binary timestamp: $(date -r $BINARY_TIMESTAMP)"
echo "⏱️  Build duration: $((BUILD_END - BUILD_START)) seconds"

# Verify binary is newer than our build start
if [ "$BINARY_TIMESTAMP" -lt "$BUILD_START" ]; then
    echo "⚠️  WARNING: Binary timestamp is older than build start time!"
    echo "   This suggests the build didn't actually update the binary."
    echo "   Binary: $(date -r $BINARY_TIMESTAMP)"
    echo "   Build:  $(date -r $BUILD_START)"
    exit 1
fi

# Create build marker with timestamp and the build ID for verification
echo "BUILD_ID=$BUILD_ID" > "$BUILD_MARKER_FILE"
echo "BUILD_TIME=$BUILD_END" >> "$BUILD_MARKER_FILE"
echo "BINARY_SIZE=$(wc -c < $PROVIDER_BINARY)" >> "$BUILD_MARKER_FILE"

echo "🏷️  Build marker created with ID: $BUILD_ID"
echo "💾 Binary size: $(wc -c < $PROVIDER_BINARY) bytes"

# Add a debug verification function to the codebase
echo ""
echo "🧪 To verify this build is active, add this debug line to your code:"
echo "   fmt.Printf(\"[BUILD_VERIFICATION] Active build: $BUILD_ID\\n\")"
echo ""
echo "✅ Build completed successfully!"
echo "   Run 'terraform plan' and look for the build verification marker in the output."