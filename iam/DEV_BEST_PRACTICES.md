# Development Best Practices for HiiRetail Terraform Provider

This document codifies lessons learned during development to prevent common debugging pitfalls.

## 🔨 Building the Provider

### Always Use the Build Verification Script

```bash
./scripts/build-and-verify.sh
```

**Never use `go build .` directly** during development. The verification script ensures:
- Binary is actually rebuilt (not cached)
- Timestamps are correct
- Build markers are created for verification

### Manual Build Verification (if needed)

```bash
# Check if binary was actually updated
ls -la terraform-provider-hiiretail

# Look for timestamp after your last change
# If timestamp is old, your build didn't work!
```

## 🧪 Debug Verification

### Always Add Build Verification Markers

When making changes, add a unique debug marker to verify your code is active:

```go
// Add this to a frequently executed code path (like provider initialization)
fmt.Printf("[BUILD_VERIFICATION] Active build: $(BUILD_ID_FROM_MARKER)\n")
```

### Use the Debug Verification Script

```bash
./scripts/debug-verify.sh "YOUR_DEBUG_MARKER" terraform plan
```

This script:
- Runs terraform with your command
- Checks if your debug marker appears
- Warns if markers are missing
- Provides troubleshooting guidance

### Debug Marker Best Practices

1. **Use unique, descriptive markers**:
   ```go
   fmt.Printf("DEBUG_READ_METHOD_START: %s\n", resourceID)
   ```

2. **Include context in debug output**:
   ```go
   fmt.Printf("DEBUG_API_CALL: GET /groups/%s/roles -> %d bytes\n", groupID, len(response))
   ```

3. **Use different markers for different execution paths**:
   ```go
   // In Create method
   fmt.Printf("DEBUG_CREATE_BINDING: %s\n", bindingName)
   
   // In Read method  
   fmt.Printf("DEBUG_READ_BINDING: %s\n", bindingName)
   
   // In Update method
   fmt.Printf("DEBUG_UPDATE_BINDING: %s\n", bindingName)
   ```

## 🚨 Red Flags to Watch For

### 1. Debug Output Not Appearing

**Symptom**: You add debug statements but don't see them in terraform output.

**Likely Cause**: Your build isn't actually updating the binary.

**Solution**:
```bash
rm terraform-provider-hiiretail
./scripts/build-and-verify.sh
./scripts/debug-verify.sh "YOUR_MARKER" terraform plan
```

### 2. Inconsistent Debug Output

**Symptom**: Some debug statements appear, others don't.

**Likely Cause**: 
- Execution path not being hit
- Build partially successful
- Code in wrong location

**Solution**: Add debug markers at the entry point of each method to trace execution flow.

### 3. Changes Not Taking Effect

**Symptom**: Code changes don't seem to affect behavior.

**Likely Cause**: Binary not rebuilt or terraform using cached version.

**Solution**:
```bash
# Force clean rebuild
rm terraform-provider-hiiretail .last-build-marker
./scripts/build-and-verify.sh

# Verify with unique marker
./scripts/debug-verify.sh "BUILD_VERIFICATION" terraform plan
```

## 📋 Development Workflow

### 1. Before Making Changes
```bash
# Establish baseline
./scripts/build-and-verify.sh
```

### 2. Making Changes
```bash
# Add unique debug marker to your changes
fmt.Printf("DEBUG_FEATURE_XYZ_v1: Implementing new logic\n")
```

### 3. After Making Changes
```bash
# Rebuild and verify
./scripts/build-and-verify.sh

# Test that your changes are active
./scripts/debug-verify.sh "DEBUG_FEATURE_XYZ_v1" terraform plan
```

### 4. Before Committing
```bash
# Remove debug statements or convert to proper logging
# Run final test
terraform plan
```

## 🛠️ Troubleshooting Guide

### Problem: "My debug output never appears"

1. Check binary timestamp: `ls -la terraform-provider-hiiretail`
2. Force rebuild: `rm terraform-provider-hiiretail && ./scripts/build-and-verify.sh`
3. Verify execution path: Add debug at method entry points
4. Check terraform log level: `TF_LOG=DEBUG terraform plan`

### Problem: "Changes seem to have no effect"

1. Verify build with: `./scripts/debug-verify.sh "BUILD_VERIFICATION" terraform plan`
2. Check if terraform is using the right binary (dev_overrides in ~/.terraformrc)
3. Clear terraform cache: `rm -rf .terraform/`

### Problem: "Some code changes work, others don't"

1. Add debug markers to both working and non-working code paths
2. Check if both paths are actually being executed
3. Verify imports and build dependencies

## 🎯 Quick Reference

```bash
# Standard development cycle
./scripts/build-and-verify.sh                                    # Build
./scripts/debug-verify.sh "MY_DEBUG_MARKER" terraform plan      # Test
terraform plan                                                   # Final check
```

Remember: **If you can't verify your debug output is appearing, assume your changes aren't active!**