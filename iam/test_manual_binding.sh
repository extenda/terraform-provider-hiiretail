#!/bin/bash

# Test manual role binding creation to match exactly what NodeJS does

# Configuration
TENANT_ID="CIR7nQwtS0rA6t0S6ejd"
GROUP_ID="9efOXfSsxwzK7AMddPtZ"
ROLE_ID="TerraformTest"  # Just the role name, not "custom.TerraformTest"

echo "=== Manual Role Binding Test ==="
echo "Group ID: $GROUP_ID"
echo "Role ID: $ROLE_ID"
echo "Tenant ID: $TENANT_ID"

# Get JWT token using existing script
echo -e "\nStep 1: Getting JWT token..."
# Use the existing get_token function from extended_test_api.sh
get_token() {
    local response=$(curl -s -X POST https://auth.retailsvc.com/oauth2/token \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "grant_type=client_credentials&client_id=iam-cli&client_secret=NCtCHs02nGNNk5fhYyGQ5BzN3vgGgkgc&scope=iam:read%20iam:write")

    local token=$(echo "$response" | jq -r '.access_token')
    if [ "$token" != "null" ] && [ -n "$token" ]; then
        echo "$token"
        return 0
    else
        echo "Failed to get token: $response" >&2
        return 1
    fi
}

TOKEN=$(get_token)
if [ $? -ne 0 ]; then
    echo "❌ Failed to get token"
    exit 1
fi
echo "✅ Token acquired successfully"

# Test exactly what NodeJS does
echo -e "\nStep 2: Testing V2 API with exact NodeJS payload..."
V2_URL="https://iam-api.retailsvc.com/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles"

# Try different binding patterns
for BINDING_PATTERN in "bu:*" "bu:global" "*" "global"; do
    echo -e "\n--- Testing bindings: [\"$BINDING_PATTERN\"] ---"
    
    V2_RESPONSE=$(curl -s -X POST "$V2_URL" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"isCustom\": true,
            \"roleId\": \"$ROLE_ID\",
            \"bindings\": [\"$BINDING_PATTERN\"]
        }")
    
    echo "V2 Response: $V2_RESPONSE"
    
    # Check if it worked by querying the group roles
    echo "Checking if role binding was created..."
    CHECK_RESPONSE=$(curl -s -X GET "https://iam-api.retailsvc.com/api/v2/tenants/$TENANT_ID/groups/$GROUP_ID/roles" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "Group roles after creation attempt: $CHECK_RESPONSE"
    echo "---"
done

echo -e "\n=== Test Complete ==="