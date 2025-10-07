#!/bin/bash

# Debug the ListGroups API call to see all groups named testShayneGroup

# Configuration  
CLIENT_ID="your-oauth2-client-id"
CLIENT_SECRET="your-oauth2-client-secret"
TENANT_ID="your-tenant-id"

echo "=== Debug ListGroups API ==="
echo "Tenant ID: $TENANT_ID"
echo

# Step 1: Get JWT token
echo "Step 1: Getting JWT token..."
TOKEN_RESPONSE=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&scope=iam:read%20iam:write")

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
if [ "$ACCESS_TOKEN" != "null" ] && [ -n "$ACCESS_TOKEN" ]; then
    echo "✅ Token acquired successfully"
else
    echo "❌ Failed to get access token"
    exit 1
fi

echo

# Step 2: Call ListGroups API (same endpoint the provider uses)
echo "Step 2: Calling ListGroups API /api/v1/tenants/${TENANT_ID}/groups"

GROUPS_RESPONSE=$(curl -s -X GET "https://iam.retailsvc.com/api/v1/tenants/${TENANT_ID}/groups" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "All groups:"
echo "$GROUPS_RESPONSE" | jq .

echo

# Step 3: Filter for testShayneGroup (multiple matches)
echo "Step 3: All groups named 'testShayneGroup':"
SHAYNE_GROUPS=$(echo "$GROUPS_RESPONSE" | jq '[.[] | select(.name == "testShayneGroup")]')
echo "$SHAYNE_GROUPS" | jq .

GROUP_COUNT=$(echo "$SHAYNE_GROUPS" | jq 'length')
echo
echo "Found $GROUP_COUNT groups named 'testShayneGroup'"

if [ "$GROUP_COUNT" -gt 1 ]; then
    echo "⚠️  Multiple groups with same name found - this could explain the wrong group ID issue!"
    echo
    echo "Group details:"
    echo "$SHAYNE_GROUPS" | jq '.[] | {id, name, status, created_at}'
fi