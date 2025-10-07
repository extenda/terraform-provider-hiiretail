#!/bin/bash

# Test script to verify API permissions using same credentials as Terraform
# This will prove whether the 403 error is truly a permissions issue

set -e

# Credentials from terraform.tfvars
CLIENT_ID="your-oauth2-client-id"
CLIENT_SECRET="your-oauth2-client-secret"
TENANT_ID="your-tenant-id"
# CORRECT ENDPOINTS: Token from auth.retailsvc.com, API calls to iam-api.retailsvc.com
TOKEN_URL="https://auth.retailsvc.com/oauth2/token"
API_BASE_URL="https://iam-api.retailsvc.com"

# Values from the terraform apply attempt
GROUP_ID="g9ODaCKjNmRdljHNMWCe"
ROLE_ID="custom.ReconciliationApprover"
RESOURCE_ID="bu:tf01"

echo "=== Testing API call with same credentials as Terraform ==="
echo "Group ID: $GROUP_ID"
echo "Role ID: $ROLE_ID"
echo "Resource ID: $RESOURCE_ID"
echo ""

# Step 1: Get OAuth2 token
echo "Step 1: Getting OAuth2 token..."
echo "Request URL: $TOKEN_URL"
echo "Client ID: ${CLIENT_ID:0:20}..."
echo ""

TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_URL" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET")

echo "Response: $TOKEN_RESPONSE"
echo ""

# Extract access token from JSON response
ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" = "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo "❌ Failed to get access token"
    echo "Response: $TOKEN_RESPONSE"
    exit 1
fi

echo "✅ Got access token: ${ACCESS_TOKEN:0:20}..."
echo ""

# Step 2: Test the exact API call that Terraform is making
echo "Step 2: Making the same API call as Terraform..."
echo "CORRECT URL (V2 API): $API_BASE_URL/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles"

# This is the EXACT API call the terraform provider is making (V2 API with correct payload)
RESPONSE=$(curl -s -w "HTTP_STATUS:%{http_code}" -X POST \
    "$API_BASE_URL/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"roleId\": \"$ROLE_ID\",
        \"isCustom\": true,
        \"bindings\": [\"$RESOURCE_ID\"]
    }")

HTTP_STATUS=$(echo "$RESPONSE" | grep -o "HTTP_STATUS:[0-9]*" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed 's/HTTP_STATUS:[0-9]*$//')

echo "HTTP Status: $HTTP_STATUS"
echo "Response Body: $BODY"
echo ""

if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "201" ]; then
    echo "✅ SUCCESS: API call worked with same credentials!"
    echo "This proves it's NOT a permissions issue - there's a bug in the terraform provider"
elif [ "$HTTP_STATUS" = "403" ]; then
    echo "❌ 403 Forbidden: Same error as Terraform"
    echo "This would confirm it's a permissions issue"
else
    echo "❓ Unexpected status code: $HTTP_STATUS"
    echo "Need to investigate further"
fi

echo ""
echo "=== Testing completed ==="