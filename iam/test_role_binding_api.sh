#!/bin/bash

# Test script to manually create role binding using V2 API
# Based on debug output from CreateRoleBinding

# Get JWT token
echo "Getting JWT token..."
TOKEN_RESPONSE=$(curl -s -X POST "https://auth.retailsvc.com/oauth2/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&scope=iam:read%20iam:write" \
  --user "$HIIRETAIL_CLIENT_ID:$HIIRETAIL_CLIENT_SECRET")

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" = "null" ]; then
  echo "Failed to get token: $TOKEN_RESPONSE"
  exit 1
fi

echo "Got token successfully"

# From debug output:
# - Group ID: mUOyAYL4AEPwTsURTAUF 
# - Role ID: TerraformTest
# - isCustom: true
# - API endpoint: ../v2/tenants/your-tenant-id/groups/mUOyAYL4AEPwTsURTAUF/roles
# - Payload: map[isCustom:true roleId:TerraformTest]

TENANT_ID="your-tenant-id"
GROUP_ID="mUOyAYL4AEPwTsURTAUF"

echo ""
echo "Testing POST to V2 API for role binding creation..."
echo "Endpoint: https://iam.retailsvc.com/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles"

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -X POST "https://iam.retailsvc.com/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "roleId": "TerraformTest",
    "isCustom": true
  }')

echo "Response:"
echo "$RESPONSE"

echo ""
echo "Now checking if role binding was created..."
VERIFY_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  "https://iam.retailsvc.com/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles")

echo "Verification - Group roles response:"
echo "$VERIFY_RESPONSE"
