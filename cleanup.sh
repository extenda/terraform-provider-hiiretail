#!/bin/bash

# Repository Cleanup Script for Terraform Provider HiiRetail
# This script removes build artifacts and temporary files

set -e

echo "🧹 Terraform Provider HiiRetail - Repository Cleanup"
echo "=================================================="

# Function to remove files if they exist
remove_if_exists() {
    if [ -e "$1" ]; then
        echo "Removing: $1"
        rm -rf "$1"
    fi
}

# Function to remove files matching pattern
remove_pattern() {
    find . -name "$1" -type f -exec rm -f {} \; 2>/dev/null || true
    echo "Removed files matching: $1"
}

echo ""
echo "🗑️  Removing build artifacts..."

# Remove Go binaries
remove_if_exists "./terraform-provider-hiiretail"
remove_pattern "terraform-provider-*"
remove_pattern "hiiretail"
remove_pattern "*.test"

# Remove from iam directory
remove_if_exists "./iam/terraform-provider-hiiretail"
remove_if_exists "./iam/terraform-provider-hiiretail-iam"
remove_if_exists "./iam/hiiretail"
remove_if_exists "./iam/bin/"
remove_if_exists "./iam/demo/demo"
remove_if_exists "./iam/demo/oauth2-demo"

echo ""
echo "🗑️  Removing temporary files..."

# Remove temporary directories
remove_if_exists "./temp/"
remove_if_exists "./tmp/"
remove_if_exists "./build/"
remove_if_exists "./dist/"

# Remove Terraform artifacts
remove_pattern "*.tfstate*"
remove_pattern "*.tfplan*"
remove_if_exists ".terraform/"
remove_pattern ".terraform.*"

# Remove debug and log files
remove_pattern "*.log"
remove_pattern "*.debug"
remove_pattern "full_debug.txt"
remove_pattern "plan_debug.txt"
remove_pattern "plan_output.txt"

# Remove backup files
remove_pattern "*.backup"
remove_pattern "*.bak"
remove_pattern "*.tmp"
remove_pattern "*~"

# Remove OS files
remove_pattern ".DS_Store"
remove_pattern ".DS_Store?"
remove_pattern "._*"
remove_pattern "Thumbs.db"

echo ""
echo "🔍 Checking for any remaining large files..."
echo "Files larger than 1MB:"
find . -type f -size +1M -not -path "./.git/*" -exec ls -lh {} \; 2>/dev/null | head -10 || echo "No large files found outside .git directory"

echo ""
echo "📊 Repository size:"
du -sh . 2>/dev/null || echo "Could not calculate size"

echo ""
echo "✅ Cleanup complete!"
echo ""
echo "💡 Note: This script only cleans the working directory."
echo "   To reduce git repository size, consider cleaning git history"
echo "   (coordinate with team first as it rewrites history):"
echo "   git filter-branch --tree-filter 'rm -rf iam/temp' HEAD"
echo ""