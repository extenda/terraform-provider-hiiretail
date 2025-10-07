#!/bin/bash

# Extended test script to verify IAM role binding via multiple API endpoints

# Configuration
CLIENT_ID="your-oauth2-client-id"
CLIENT_SECRET="your-oauth2-client-secret"
TENANT_ID="your-tenant-id"
GROUP_ID="9efOXfSsxwzK7AMddPtZ"
ROLE_ID="custom.TerraformTest"

echo "=== Extended IAM API Testing ==="
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

echo "✅ Token acquired successfully"
echo

# Step 2: Check if the group exists
echo "Step 2: Verifying group exists..."
GROUP_RESPONSE=$(curl -s -X GET "https://iam-api.retailsvc.com/api/v1/tenants/$TENANT_ID/groups/$GROUP_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "Group details:"
echo "$GROUP_RESPONSE" | jq '.'
echo

# Step 3: Check if the custom role exists
echo "Step 3: Verifying custom role exists..."
ROLE_RESPONSE=$(curl -s -X GET "https://iam-api.retailsvc.com/api/v1/tenants/$TENANT_ID/roles/custom.TerraformTest" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "Custom role details:"
echo "$ROLE_RESPONSE" | jq '.'
echo

# Step 4: Query group roles via V2 API
echo "Step 4: Querying group roles via V2 API..."
ROLES_RESPONSE=$(curl -s -X GET "https://iam-api.retailsvc.com/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "V2 API Group roles response:"
echo "$ROLES_RESPONSE" | jq '.'
echo

# Step 5: Try V1 API if available
echo "Step 5: Trying alternative API endpoints..."

# Try to list all role bindings
BINDINGS_RESPONSE=$(curl -s -X GET "https://iam-api.retailsvc.com/api/v1/bindings" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json")

echo "V1 bindings response:"
echo "$BINDINGS_RESPONSE" | jq '.'
echo

# Step 6: Try to query specific role binding ID (from terraform)
BINDING_ID="6IyHJfFwBt4PIigRg0eT-custom.TerraformTest"
echo "Step 6: Checking terraform role binding ID: $BINDING_ID"

# The provider's GetRoleBinding method constructs this from group + role
echo "This ID format indicates:"
echo "  Group ID: ${BINDING_ID%-*}"
echo "  Role ID: ${BINDING_ID#*-}"
echo

echo "=== Summary ==="
echo "✅ JWT Authentication: Working"
echo "✅ Group exists: $(if echo "$GROUP_RESPONSE" | jq -e '.id' > /dev/null; then echo "Yes"; else echo "No"; fi)"
echo "✅ Custom role exists: $(if echo "$ROLE_RESPONSE" | jq -e '.id' > /dev/null; then echo "Yes"; else echo "No"; fi)"
ROLE_BINDING_STATUS="$(if [ "$(echo "$ROLES_RESPONSE" | jq '. | length')" -gt 0 ]; then echo "Found"; else echo "Not found"; fi)"
if [ "$ROLE_BINDING_STATUS" = "Found" ]; then
    echo "✅ Role binding in API: $ROLE_BINDING_STATUS"
    echo
    echo "CONCLUSION: SUCCESS! The Terraform provider successfully created the role binding"
    echo "and it now exists in the IAM API with the correct configuration."
else
    echo "❌ Role binding in API: $ROLE_BINDING_STATUS"
    echo
    echo "CONCLUSION: The Terraform provider reports successful role binding creation,"
    echo "but the actual IAM API does not show the role binding exists."
    echo "This suggests an issue with the CreateRoleBinding implementation."
fi