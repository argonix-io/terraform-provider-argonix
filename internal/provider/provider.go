package provider

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &argonixProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &argonixProvider{
            version: version,
        }
    }
}

// argonixProvider is the provider implementation.
type argonixProvider struct {
    // version is set to the provider version on release, "dev" when the
    // provider is built and ran locally, and "test" when running acceptance
    // testing.
    version string
}

// Metadata returns the provider type name.
func (p *argonixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "argonix"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *argonixProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{}
}

// Configure prepares a argonix API client for data sources and resources.
func (p *argonixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *argonixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return nil
}

// Resources defines the resources implemented in the provider.
func (p *argonixProvider) Resources(_ context.Context) []func() resource.Resource {
    return nil
}
