#!/bin/bash

# Test the exact same API call our Go code is making

# Get a fresh token using the exact same credentials as extended_test_api.sh
echo "Getting token..."
CLIENT_ID="b3duZXI2IHNoYXluZQpzZnc2IGhpaXRmQDAuMUBDSVI3blF3dFMwckE2dDBTNmVqZAp0aWQ2IENJUjduUXd0UzByQTZ0MFM2ZWpkCg"
CLIENT_SECRET="726143f664f0a38efa96abe33bc0a7487d745ee725171101231c454ea9faa1ba"

TOKEN_RESPONSE=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET&scope=iam:read iam:write")

TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

echo "Token: ${TOKEN:0:50}..."

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to get token: $TOKEN_RESPONSE"
    exit 1
fi

# Test the exact same call our Go code makes
echo -e "\nTesting exact Go payload..."
V2_URL="https://iam-api.retailsvc.com/api/v2/tenants/CIR7nQwtS0rA6t0S6ejd/groups/9efOXfSsxwzK7AMddPtZ/roles"

echo "POST request:"
curl -v -X POST "$V2_URL" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "isCustom": true,
        "roleId": "TerraformTest",
        "bindings": ["bu:001"]
    }'

echo -e "\n\nNow checking if role binding was created..."
curl -s -X GET "$V2_URL" \
    -H "Authorization: Bearer $TOKEN" | jq