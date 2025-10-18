package iam

// NewServiceForTest constructs a Service using provided RawClient and clientService.
// This helper is intentionally small and designed to help tests create a Service
// without relying on a fully-initialized *client.Client. Use only in tests.
func NewServiceForTest(raw RawClient, svc clientService, tenantID string) *Service {
	return &Service{
		rawClient: raw,
		client:    svc,
		tenantID:  tenantID,
	}
}
