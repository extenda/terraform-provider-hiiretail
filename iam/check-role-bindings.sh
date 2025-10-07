#!/bin/bash

curl -X GET \
  'https://iam-api.retailsvc.com/api/v2/tenants/your-tenant-id/groups/3NGyxSchmH7RtF2pZPqC/roles' \
  -H 'Authorization: Bearer your-jwt-token-here' \
  -H 'X-Tenant-ID: your-tenant-id' \
  -H 'Content-Type: application/json' \
  -v