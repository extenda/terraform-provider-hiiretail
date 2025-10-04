const fs = require('fs');
const csv = require('csv-parser');
const axios = require('axios');

const tenantId = 'KTDfxPSHBPzAjFAOjaTg';
const roleId = 'financial.users';
const tokenFile = 'token.txt';

fs.createReadStream('input.csv')
  .pipe(csv())
  .on('data', async (data) => {
    const groupId = data['Group_ID'];
    const token = fs.readFileSync(tokenFile, 'utf-8');
    const url = `https://iam-api.retailsvc.com/api/v1/tenants/${tenantId}/groups/${groupId}/roles/${roleId}?isCustom=true`;

    try {
      const response = await axios.delete(url, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      console.log(`Deleted role binding for group ${groupId}`);
    } catch (error) {
      console.error(`Failed to delete role binding for group ${groupId}: ${error.message}`);
    }
  });
