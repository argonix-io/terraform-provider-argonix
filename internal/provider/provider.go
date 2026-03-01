package provider

import (
	"context"
	"os"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	version string
}

// argonixProviderModel maps provider schema data to a Go type.
type argonixProviderModel struct {
	URL            types.String `tfsdk:"url"`
	APIKey         types.String `tfsdk:"api_key"`
	OrganizationID types.String `tfsdk:"organization_id"`
}

// Metadata returns the provider type name.
func (p *argonixProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "argonix"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *argonixProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Argonix monitoring platform.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "Base URL of the Argonix API. Defaults to https://api.argonix.io. Can also be set via the ARGONIX_URL environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "API key for authenticating with Argonix. Can also be set via the ARGONIX_API_KEY environment variable.",
				Required:    true,
				Sensitive:   true,
			},
			"organization_id": schema.StringAttribute{
				Description: "The UUID of the organization to manage resources in. Can also be set via the ARGONIX_ORGANIZATION_ID environment variable.",
				Required:    true,
			},
		},
	}
}

// Configure prepares an Argonix API client for data sources and resources.
func (p *argonixProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config argonixProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values from environment
	url := os.Getenv("ARGONIX_URL")
	apiKey := os.Getenv("ARGONIX_API_KEY")
	organizationID := os.Getenv("ARGONIX_ORGANIZATION_ID")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}
	if !config.OrganizationID.IsNull() {
		organizationID = config.OrganizationID.ValueString()
	}

	// Apply defaults
	if url == "" {
		url = "https://api.argonix.io"
	}

	// Validate required values
	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Argonix API Key",
			"The provider cannot create the Argonix API client as there is a missing or empty value for the Argonix API key. "+
				"Set the api_key value in the configuration or use the ARGONIX_API_KEY environment variable.",
		)
	}

	if organizationID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization_id"),
			"Missing Argonix Organization ID",
			"The provider cannot create the Argonix API client as there is a missing or empty value for the organization ID. "+
				"Set the organization_id value in the configuration or use the ARGONIX_ORGANIZATION_ID environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create client
	c := client.NewClient(url, apiKey, organizationID)
	resp.DataSourceData = c
	resp.ResourceData = c
}

// DataSources defines the data sources implemented in the provider.
func (p *argonixProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMonitorDataSource,
		NewMonitorsDataSource,
		NewSyntheticTestDataSource,
		NewSyntheticTestsDataSource,
		NewGroupDataSource,
		NewGroupsDataSource,
		NewAlertRuleDataSource,
		NewAlertRulesDataSource,
		NewAlertChannelDataSource,
		NewAlertChannelsDataSource,
		NewStatusPageDataSource,
		NewStatusPagesDataSource,
		NewEnvironmentDataSource,
		NewTestSuiteDataSource,
		NewManualTestCaseDataSource,
		NewTestPlanDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *argonixProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMonitorResource,
		NewSyntheticTestResource,
		NewGroupResource,
		NewAlertRuleResource,
		NewAlertChannelResource,
		NewStatusPageResource,
		NewEnvironmentResource,
		NewTestSuiteResource,
		NewManualTestCaseResource,
		NewTestPlanResource,
	}
}
