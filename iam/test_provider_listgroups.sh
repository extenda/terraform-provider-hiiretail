#!/bin/bash

# Test the correct ListGroups endpoint that the provider uses

# Configuration  
CLIENT_ID="b3duZXI6IHNoYXluZQpzZnc6IGhpaXRmQDAuMUBDSVI3blF3dFMwckE2dDBTNmVqZAp0aWQ6IENJUjduUXd0UzByQTZ0MFM2ZWpkCg"
CLIENT_SECRET="726143f664f0a38efa96abe33bc0a7487d745ee725171101231c454ea9faa1ba"
TENANT_ID="CIR7nQwtS0rA6t0S6ejd"

echo "=== Test Provider ListGroups Endpoint ==="
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

# Step 2: Test provider's ListGroups endpoint (without /api/v1 prefix)
echo "Step 2: Testing provider ListGroups endpoint: /tenants/${TENANT_ID}/groups"

# Since the provider uses s.client.Get() which likely has a base URL, let me try both variations
echo "Trying with iam.retailsvc.com base:"
GROUPS_RESPONSE1=$(curl -s -X GET "https://iam.retailsvc.com/tenants/${TENANT_ID}/groups" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "Response 1:"
echo "$GROUPS_RESPONSE1" | jq .

echo
echo "Trying with iam.retailsvc.com/api/v1 base:"
GROUPS_RESPONSE2=$(curl -s -X GET "https://iam.retailsvc.com/api/v1/tenants/${TENANT_ID}/groups" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "Response 2:"
echo "$GROUPS_RESPONSE2" | jq .

echo

# Step 3: Check which one has groups
if echo "$GROUPS_RESPONSE1" | jq -e 'type == "array" and length > 0' > /dev/null 2>&1; then
    echo "✅ Found groups in response 1"
    GROUPS_RESPONSE="$GROUPS_RESPONSE1"
elif echo "$GROUPS_RESPONSE2" | jq -e 'type == "array" and length > 0' > /dev/null 2>&1; then
    echo "✅ Found groups in response 2" 
    GROUPS_RESPONSE="$GROUPS_RESPONSE2"
else
    echo "❌ No groups found in either response"
    exit 1
fi

# Step 4: Look for testShayneGroup
echo "Step 4: Looking for testShayneGroup..."
SHAYNE_GROUPS=$(echo "$GROUPS_RESPONSE" | jq '[.[] | select(.name == "testShayneGroup")]')
echo "$SHAYNE_GROUPS" | jq .

GROUP_COUNT=$(echo "$SHAYNE_GROUPS" | jq 'length')
echo
echo "Found $GROUP_COUNT groups named 'testShayneGroup'"

if [ "$GROUP_COUNT" -gt 0 ]; then
    echo
    echo "Group details:"
    echo "$SHAYNE_GROUPS" | jq '.[] | {id, name, status}'
    
    # Show which group would be selected (first one found)
    FIRST_GROUP_ID=$(echo "$SHAYNE_GROUPS" | jq -r '.[0].id')
    echo
    echo "First group found (what CreateRoleBinding would use): $FIRST_GROUP_ID"
    echo "Correct group from terraform state: mUOyAYL4AEPwTsURTAUF"
    
    if [ "$FIRST_GROUP_ID" = "mUOyAYL4AEPwTsURTAUF" ]; then
        echo "✅ CreateRoleBinding would use the correct group"
    else
        echo "❌ CreateRoleBinding would use the wrong group!"
    fi
fi