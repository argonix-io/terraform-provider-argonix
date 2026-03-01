package provider

import (
	"context"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &syntheticTestDataSource{}
	_ datasource.DataSourceWithConfigure = &syntheticTestDataSource{}
	_ datasource.DataSource              = &syntheticTestsDataSource{}
	_ datasource.DataSourceWithConfigure = &syntheticTestsDataSource{}
)

// --- Single synthetic test data source ---

func NewSyntheticTestDataSource() datasource.DataSource {
	return &syntheticTestDataSource{}
}

type syntheticTestDataSource struct {
	client *client.Client
}

func (d *syntheticTestDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synthetic_test"
}

func (d *syntheticTestDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix synthetic test by ID.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true},
			"name":           schema.StringAttribute{Computed: true},
			"description":    schema.StringAttribute{Computed: true},
			"is_active":      schema.BoolAttribute{Computed: true},
			"test_type":      schema.StringAttribute{Computed: true},
			"steps":          schema.StringAttribute{Computed: true},
			"check_interval": schema.Int64Attribute{Computed: true},
			"timeout":        schema.Int64Attribute{Computed: true},
			"tags":           schema.StringAttribute{Computed: true},
			"locations":      schema.StringAttribute{Computed: true},
			"current_status": schema.StringAttribute{Computed: true},
			"date_created":   schema.StringAttribute{Computed: true},
			"date_modified":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *syntheticTestDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *syntheticTestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp syntheticTestAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/synthetic-tests/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading synthetic test", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, syntheticTestAPIToState(apiResp))...)
}

// --- List synthetic tests data source ---

func NewSyntheticTestsDataSource() datasource.DataSource {
	return &syntheticTestsDataSource{}
}

type syntheticTestsDataSource struct {
	client *client.Client
}

type syntheticTestsDataSourceModel struct {
	SyntheticTests []syntheticTestResourceModel `tfsdk:"synthetic_tests"`
}

func (d *syntheticTestsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synthetic_tests"
}

func (d *syntheticTestsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all synthetic tests in the organization.",
		Attributes: map[string]schema.Attribute{
			"synthetic_tests": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":             schema.StringAttribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"description":    schema.StringAttribute{Computed: true},
						"is_active":      schema.BoolAttribute{Computed: true},
						"test_type":      schema.StringAttribute{Computed: true},
						"steps":          schema.StringAttribute{Computed: true},
						"check_interval": schema.Int64Attribute{Computed: true},
						"timeout":        schema.Int64Attribute{Computed: true},
						"tags":           schema.StringAttribute{Computed: true},
						"locations":      schema.StringAttribute{Computed: true},
						"current_status": schema.StringAttribute{Computed: true},
						"date_created":   schema.StringAttribute{Computed: true},
						"date_modified":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *syntheticTestsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *syntheticTestsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []syntheticTestAPIModel
	err := d.client.List(ctx, "/synthetic-tests/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing synthetic tests", err.Error())
		return
	}

	state := syntheticTestsDataSourceModel{}
	for _, t := range apiResp {
		state.SyntheticTests = append(state.SyntheticTests, syntheticTestAPIToState(t))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
