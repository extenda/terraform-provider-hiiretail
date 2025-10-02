# Service Module Contract

## Service Module Architecture Contract

**Module Interface**:
```go
type ServiceModule interface {
    // Service identification
    ServiceName() string
    APIVersion() string
    BaseURL() string
    
    // Resource registration
    Resources() []func() resource.Resource
    DataSources() []func() datasource.DataSource
    
    // Service-specific configuration
    ConfigureClient(ctx context.Context, config ProviderConfig) error
    
    // Health and validation
    ValidateService(ctx context.Context) error
}
```

## IAM Service Module Contract

**Service Configuration**:
- Service Name: `"iam"`
- API Version: `"v1"`
- Base URL: `"https://iam-api.retailsvc.com"` (configurable via provider)
- Authentication: OAuth2 via shared provider configuration

**Resource Registration**:
```go
func (s *IAMService) Resources() []func() resource.Resource {
    return []func() resource.Resource{
        NewGroupResource,        // hiiretail_iam_group
        NewCustomRoleResource,   // hiiretail_iam_custom_role  
        NewRoleBindingResource,  // hiiretail_iam_role_binding
    }
}

func (s *IAMService) DataSources() []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewGroupsDataSource,     // hiiretail_iam_groups
    }
}
```

**Resource Naming Pattern**:
- Resource Type Name: `hiiretail_iam_{resource}`
- Go Struct Name: `{Resource}Resource`
- Go Constructor: `New{Resource}Resource()`

## CCC Service Module Contract (Future)

**Service Configuration**:
- Service Name: `"ccc"`
- API Version: `"v1"`
- Base URL: `"https://ccc-api.retailsvc.com"` (configurable via provider)
- Authentication: OAuth2 via shared provider configuration

**Resource Registration**:
```go
func (s *CCCService) Resources() []func() resource.Resource {
    return []func() resource.Resource{
        NewKindResource,         // hiiretail_ccc_kind
    }
}
```

## Service Registration Contract

**Provider Registration**:
```go
func (p *HiiRetailProvider) Resources(ctx context.Context) []func() resource.Resource {
    var resources []func() resource.Resource
    
    // Register IAM service resources
    iamService := NewIAMService()
    resources = append(resources, iamService.Resources()...)
    
    // Register CCC service resources (when available)
    if p.cccEnabled {
        cccService := NewCCCService()
        resources = append(resources, cccService.Resources()...)
    }
    
    return resources
}
```

## Service Lifecycle Contract

**Initialization Phase**:
1. Service module instantiated with provider configuration
2. Service validates its configuration and connectivity
3. Service registers its resources and data sources
4. Service reports readiness to provider

**Runtime Phase**:
1. Service handles resource CRUD operations via its registered resources
2. Service manages API client connections and authentication
3. Service implements retry logic and error handling
4. Service reports health status and metrics

**Error Handling**:
- Service unavailable: Log warning, disable service resources
- Authentication failure: Propagate error to provider level
- Configuration error: Fail fast with clear error message
- Runtime error: Implement appropriate retry and fallback logic

## Extensibility Contract

**Adding New Services**:
1. Implement `ServiceModule` interface
2. Register service in provider configuration
3. Add service documentation and examples
4. Include service in test suite
5. Update migration guides if needed

**Service Dependencies**:
- Services should be independent where possible
- Cross-service references via Terraform resource references
- Shared infrastructure (auth, logging) via provider-level services