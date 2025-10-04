#!/bin/bash

# Get token first using the working credentials
echo "Getting token..."
TOKEN=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=b3duZXI6IHNoYXluZQpzZnc6IGhpaXRmQDAuMUBDSVI3blF3dFMwckE2dDBTNmVqZAp0aWQ6IENJUjduUXd0UzByQTZ0MFM2ZWpkCg&client_secret=726143f664f0a38efa96abe33bc0a7487d745ee725171101231c454ea9faa1ba&scope=iam:read%20iam:write" | jq -r '.access_token')

echo "Token: ${TOKEN}"

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to get token"
    exit 1
fi

# Test the exact payload our Go code would send
echo -e "\nTesting V2 API with Go payload..."
V2_URL="https://iam-api.retailsvc.com/api/v2/tenants/CIR7nQwtS0rA6t0S6ejd/groups/9efOXfSsxwzK7AMddPtZ/roles"

curl -v -X POST "$V2_URL" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "isCustom": true,
        "roleId": "TerraformTest",
        "bindings": ["bu:*"]
    }'