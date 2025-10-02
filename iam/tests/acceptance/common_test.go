package acceptance

import (
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

// randomString generates a random string of specified length
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

const testAccProviderConfig = `
provider "hiiretail" {
  client_id     = "test-client-id" 
  client_secret = "test-client-secret"
  base_url      = "https://api-test.hiiretail.com"
  iam_endpoint  = "/iam/v1"
}
`

// init initializes the test provider factories
func init() {
	// TODO: This will fail until provider is properly implemented
	// For now, initialize empty provider factories
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"hiiretail": func() (tfprotov6.ProviderServer, error) {
			// TODO: Return actual provider factory once implemented
			return nil, nil
		},
	}
}
