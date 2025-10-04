# HIIRETAIL API AUTHENTICATION ENDPOINTS

## CRITICAL: OAuth2 Token Endpoint
**Token URL: https://auth.retailsvc.com/oauth2/token**

## API Base URLs
- **IAM API Base**: https://iam-api.retailsvc.com
- **Auth Service**: https://auth.retailsvc.com

## REMEMBER: 
- Token endpoint is on auth.retailsvc.com (NOT iam-api.retailsvc.com)
- API calls go to iam-api.retailsvc.com
- Don't confuse these two domains!

## Authentication Flow
1. GET token from: https://auth.retailsvc.com/oauth2/token
2. USE token for API calls to: https://iam-api.retailsvc.com/iam/...

## Example curl for token:
```bash
curl -X POST "https://auth.retailsvc.com/oauth2/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials&client_id=$CLIENT_ID&client_secret=$CLIENT_SECRET"
```