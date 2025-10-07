#!/bin/bash

# Get token first using the working credentials
echo "Getting token..."
TOKEN=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=your-oauth2-client-id&client_secret=your-oauth2-client-secret&scope=iam:read%20iam:write" | jq -r '.access_token')

echo "Token: ${TOKEN}"

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to get token"
    exit 1
fi

# Test the exact payload our Go code would send
echo -e "\nTesting V2 API with Go payload..."
V2_URL="https://iam-api.retailsvc.com/api/v2/tenants/your-tenant-id/groups/9efOXfSsxwzK7AMddPtZ/roles"

curl -v -X POST "$V2_URL" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "isCustom": true,
        "roleId": "TerraformTest",
        "bindings": ["bu:*"]
    }'