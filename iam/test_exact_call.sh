#!/b# Get a fresh token using the exact same credentials as extended_test_api.sh
echo "Getting token..."
CLIENT_ID="your-oauth2-client-id"
CLIENT_SECRET="your-oauth2-client-secret"

TOKEN_RESPONSE=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET&scope=iam:read iam:write")

TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')sh

# Test the exact same API call our Go code is making

# Get a fresh token
echo "Getting token..."
TOKEN=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&client_id=your-oauth2-client-id&client_secret=your-oauth2-client-secret&scope=iam:read%20iam:write" | jq -r '.access_token')

echo "Token: ${TOKEN:0:50}..."

# Test the exact same call our Go code makes
echo -e "\nTesting exact Go payload..."
V2_URL="https://iam-api.retailsvc.com/api/v2/tenants/your-tenant-id/groups/9efOXfSsxwzK7AMddPtZ/roles"

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