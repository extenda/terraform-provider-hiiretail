#!/bin/bash

# Test script to verify IAM role binding via API

# Configuration
CLIENT_ID="your-oauth2-client-id"
CLIENT_SECRET="your-oauth2-client-secret"
TENANT_ID="your-tenant-id"
GROUP_ID="6IyHJfFwBt4PIigRg0eT"
ROLE_ID="custom.TerraformTest"

echo "=== Testing IAM API with JWT Authentication ==="
echo "Group ID: $GROUP_ID"
echo "Role ID: $ROLE_ID"
echo "Tenant ID: $TENANT_ID"
echo

# Step 1: Get JWT token
echo "Step 1: Getting JWT token..."
TOKEN_RESPONSE=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET&scope=iam:read iam:write")

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" = "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo "Failed to get access token:"
    echo "$TOKEN_RESPONSE"
    exit 1
fi

echo "✅ Token acquired successfully (${#ACCESS_TOKEN} characters)"
echo

# Step 2: Query group roles via V2 API
echo "Step 2: Querying group roles via V2 API..."
ROLES_RESPONSE=$(curl -s -X GET "https://iam-api.retailsvc.com/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "API Response:"
echo "$ROLES_RESPONSE" | jq '.'
echo

# Step 3: Check if our role is in the response
echo "Step 3: Checking for our role binding..."
ROLE_COUNT=$(echo "$ROLES_RESPONSE" | jq '. | length')
echo "Found $ROLE_COUNT role bindings for this group"

# Look for our specific role
HAS_ROLE=$(echo "$ROLES_RESPONSE" | jq --arg role "$ROLE_ID" '.[] | select(.roleId == $role or .roleId == ("custom." + $role) or .roleId == ("TerraformTest"))')

if [ -n "$HAS_ROLE" ]; then
    echo "✅ FOUND: Role binding exists in API!"
    echo "Role details:"
    echo "$HAS_ROLE" | jq '.'
else
    echo "❌ NOT FOUND: Role binding not found in API response"
    echo "Available roles:"
    echo "$ROLES_RESPONSE" | jq '.[] | .roleId'
fi

echo
echo "=== Test Complete ==="