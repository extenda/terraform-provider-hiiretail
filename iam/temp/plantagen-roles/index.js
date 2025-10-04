const axios = require('axios');
const csv = require('csv-parser');
const fs = require('fs');

const tenantId = 'KTDfxPSHBPzAjFAOjaTg';
const roleId = 'financial.users';
const apiUrl = `https://iam-api.retailsvc.com/api/v2/tenants/${tenantId}/groups/{id}/roles`;

// Read the JWT from the token.txt file
fs.readFile('token.txt', 'utf8', async (err, token) => {
  if (err) throw err;

  // Read the CSV file
  fs.createReadStream('input.csv')
    .pipe(csv())
    .on('data', async (data) => {
      const groupId = data.Group_ID;
      const bindings = [`bu:${data.Store_ID}`];

      const payload = {
        isCustom: true,
        roleId,
        bindings,
      };

      try {
        const response = await axios.post(apiUrl.replace('{id}', groupId), payload, {
          headers: {
            Authorization: `Bearer ${token.trim()}`,
          },
        });

        if (response.status === 200) {
          console.log(`Successfully posted data for Group ID: ${groupId} with ${roleId} and bu:${data.Store_ID}`);
        } else {
          console.error(`Failed to post data for Group ID: ${groupId}`);
        }
      } catch (error) {
        console.error(error.response.data);
      }
    });
});

