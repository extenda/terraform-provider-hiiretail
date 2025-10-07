#!/bin/bash

# Script to list role bindings for the group to see what's there

set -e

# Credentials from terraform.tfvars
CLIENT_ID="your-oauth2-client-id"
CLIENT_SECRET="your-oauth2-client-secret"
TENANT_ID="your-tenant-id"

# CORRECT ENDPOINTS
TOKEN_URL="https://auth.retailsvc.com/oauth2/token"
API_BASE_URL="https://iam-api.retailsvc.com"

# Values from the terraform apply attempt
GROUP_ID="g9ODaCKjNmRdljHNMWCe"

echo "=== Listing role bindings for group $GROUP_ID ==="
echo ""

# Step 1: Get OAuth2 token
echo "Step 1: Getting OAuth2 token..."
TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_URL" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET")

# Extract access token from JSON response
ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

if [ "$ACCESS_TOKEN" = "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo "❌ Failed to get access token"
    echo "Response: $TOKEN_RESPONSE"
    exit 1
fi

echo "✅ Got access token: ${ACCESS_TOKEN:0:20}..."
echo ""

# Step 2: List role bindings using V2 API
echo "Step 2: Listing role bindings..."
echo "GET URL: $API_BASE_URL/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles"

RESPONSE=$(curl -s -w "HTTP_STATUS:%{http_code}" -X GET \
    "$API_BASE_URL/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Content-Type: application/json")

HTTP_STATUS=$(echo "$RESPONSE" | grep -o "HTTP_STATUS:[0-9]*" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed 's/HTTP_STATUS:[0-9]*$//')

echo "HTTP Status: $HTTP_STATUS"
echo "Response Body: $BODY"
echo ""

if [ "$HTTP_STATUS" = "200" ]; then
    echo "✅ SUCCESS: Retrieved role bindings"
    echo "Role bindings:"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
    echo "❌ Failed to get role bindings: HTTP $HTTP_STATUS"
    echo "Response: $BODY"
fi

echo ""
echo "=== Listing completed ==="