#!/bin/bash

# List all groups to see what actually exists

# Configuration  
CLIENT_ID="b3duZXI6IHNoYXluZQpzZnc6IGhpaXRmQDAuMUBDSVI3blF3dFMwckE2dDBTNmVqZAp0aWQ6IENJUjduUXd0UzByQTZ0MFM2ZWpkCg"
CLIENT_SECRET="726143f664f0a38efa96abe33bc0a7487d745ee725171101231c454ea9faa1ba"
TENANT_ID="CIR7nQwtS0rA6t0S6ejd"

echo "=== Listing All Groups ==="
echo "Tenant ID: $TENANT_ID"
echo

# Step 1: Get JWT token
echo "Step 1: Getting JWT token..."
TOKEN_RESPONSE=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&scope=iam:read%20iam:write")

if [ $? -eq 0 ]; then
    ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')
    if [ "$ACCESS_TOKEN" != "null" ] && [ -n "$ACCESS_TOKEN" ]; then
        echo "✅ Token acquired successfully"
    else
        echo "❌ Failed to get access token"
        echo "Response: $TOKEN_RESPONSE"
        exit 1
    fi
else
    echo "❌ Failed to call token endpoint"
    exit 1
fi

echo

# Step 2: List all groups
echo "Step 2: Listing all groups..."
GROUPS_RESPONSE=$(curl -s -X GET "https://iam.retailsvc.com/api/v1/tenants/${TENANT_ID}/groups" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "Groups response:"
echo "$GROUPS_RESPONSE" | jq .

echo

# Step 3: Look for testShayneGroup specifically
echo "Step 3: Filtering for testShayneGroup..."
SHAYNE_GROUP=$(echo "$GROUPS_RESPONSE" | jq '.[] | select(.name == "testShayneGroup")')
echo "testShayneGroup details:"
echo "$SHAYNE_GROUP" | jq .

if [ -n "$SHAYNE_GROUP" ] && [ "$SHAYNE_GROUP" != "null" ]; then
    ACTUAL_GROUP_ID=$(echo "$SHAYNE_GROUP" | jq -r '.id')
    echo
    echo "Found testShayneGroup with ID: $ACTUAL_GROUP_ID"
else
    echo
    echo "testShayneGroup not found in API response"
fi