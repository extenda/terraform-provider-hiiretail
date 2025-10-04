#!/bin/bash

# Test V2 API role assignment manually with correct group ID

# Configuration  
CLIENT_ID="b3duZXI6IHNoYXluZQpzZnc6IGhpaXRmQDAuMUBDSVI3blF3dFMwckE2dDBTNmVqZAp0aWQ6IENJUjduUXd0UzByQTZ0MFM2ZWpkCg"
CLIENT_SECRET="726143f664f0a38efa96abe33bc0a7487d745ee725171101231c454ea9faa1ba"
TENANT_ID="CIR7nQwtS0rA6t0S6ejd"
GROUP_ID="9efOXfSsxwzK7AMddPtZ"  # Correct group ID from terraform state
ROLE_ID="TerraformTest" # Just the role name for API

echo "=== Testing V2 API Role Assignment ==="
echo "Group ID: $GROUP_ID"
echo "Role ID: $ROLE_ID"
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
    echo "Response: $TOKEN_RESPONSE"
    exit 1
fi

echo

# Step 2: Test V2 API POST with exact payload from debug output
echo "Step 2: Testing V2 API POST /api/v2/tenants/${TENANT_ID}/groups/${GROUP_ID}/roles"
echo "Payload: {\"roleId\": \"${ROLE_ID}\", \"isCustom\": true}"

V2_POST_RESPONSE=$(curl -s -X POST "https://iam.retailsvc.com/api/v2/tenants/${TENANT_ID}/groups/${GROUP_ID}/roles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"roleId\": \"${ROLE_ID}\", \"isCustom\": true}")

echo "POST Response:"
echo "$V2_POST_RESPONSE" | jq .

echo

# Step 3: Check if role was assigned by querying the group roles
echo "Step 3: Verifying role assignment with GET /api/v2/tenants/${TENANT_ID}/groups/${GROUP_ID}/roles"

V2_GET_RESPONSE=$(curl -s -X GET "https://iam.retailsvc.com/api/v2/tenants/${TENANT_ID}/groups/${GROUP_ID}/roles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "GET Response:"
echo "$V2_GET_RESPONSE" | jq .

echo
echo "Raw GET Response (no jq):"
echo "$V2_GET_RESPONSE"

echo

# Step 4: Summary  
echo "Checking for roleId '${ROLE_ID}' in response..."

# Check if response is empty or null
if [ -z "$V2_GET_RESPONSE" ] || [ "$V2_GET_RESPONSE" = "null" ] || [ "$V2_GET_RESPONSE" = "[]" ]; then
    echo "❌ Role assignment failed - API returned empty response"
    echo "Response was: '$V2_GET_RESPONSE'"
else
    # Check if role exists in response
    ROLE_CHECK=$(echo "$V2_GET_RESPONSE" | jq -e ".[] | select(.roleId == \"${ROLE_ID}\")" 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$ROLE_CHECK" ]; then
        echo "✅ Role assignment successful - role found in group roles"
        echo "Found role: $ROLE_CHECK"
    else
        echo "❌ Role assignment failed - role not found in group roles"
        echo "Available roles in response:"
        echo "$V2_GET_RESPONSE" | jq -r '.[] | .roleId' 2>/dev/null || echo "No roles found"
    fi
fi