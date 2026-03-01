package provider

import (
	"context"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &monitorDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorDataSource{}
	_ datasource.DataSource              = &monitorsDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorsDataSource{}
)

// --- Single monitor data source ---

func NewMonitorDataSource() datasource.DataSource {
	return &monitorDataSource{}
}

type monitorDataSource struct {
	client *client.Client
}

func (d *monitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *monitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix monitor by ID.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true, Description: "UUID of the monitor."},
			"name":           schema.StringAttribute{Computed: true},
			"monitor_type":   schema.StringAttribute{Computed: true},
			"is_active":      schema.BoolAttribute{Computed: true},
			"url":            schema.StringAttribute{Computed: true},
			"hostname":       schema.StringAttribute{Computed: true},
			"port":           schema.Int64Attribute{Computed: true},
			"check_interval": schema.Int64Attribute{Computed: true},
			"timeout":        schema.Int64Attribute{Computed: true},
			"retries":        schema.Int64Attribute{Computed: true},
			"http_method":    schema.StringAttribute{Computed: true},
			"current_status": schema.StringAttribute{Computed: true},
			"group_id":       schema.StringAttribute{Computed: true},
			"tags":           schema.StringAttribute{Computed: true, Description: "JSON-encoded tags."},
			"regions":        schema.StringAttribute{Computed: true, Description: "JSON-encoded regions."},
			"date_created":   schema.StringAttribute{Computed: true},
			"date_modified":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *monitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *monitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state monitorResourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp monitorAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/monitors/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading monitor", err.Error())
		return
	}

	state = monitorAPIToState(apiResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// --- List monitors data source ---

func NewMonitorsDataSource() datasource.DataSource {
	return &monitorsDataSource{}
}

type monitorsDataSource struct {
	client *client.Client
}

type monitorsDataSourceModel struct {
	Monitors []monitorResourceModel `tfsdk:"monitors"`
}

func (d *monitorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitors"
}

func (d *monitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all monitors in the organization.",
		Attributes: map[string]schema.Attribute{
			"monitors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":             schema.StringAttribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"monitor_type":   schema.StringAttribute{Computed: true},
						"is_active":      schema.BoolAttribute{Computed: true},
						"url":            schema.StringAttribute{Computed: true},
						"hostname":       schema.StringAttribute{Computed: true},
						"port":           schema.Int64Attribute{Computed: true},
						"check_interval": schema.Int64Attribute{Computed: true},
						"timeout":        schema.Int64Attribute{Computed: true},
						"retries":        schema.Int64Attribute{Computed: true},
						"http_method":    schema.StringAttribute{Computed: true},
						"current_status": schema.StringAttribute{Computed: true},
						"group_id":       schema.StringAttribute{Computed: true},
						"tags":           schema.StringAttribute{Computed: true},
						"regions":        schema.StringAttribute{Computed: true},
						"date_created":   schema.StringAttribute{Computed: true},
						"date_modified":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *monitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *monitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []monitorAPIModel
	err := d.client.List(ctx, "/monitors/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing monitors", err.Error())
		return
	}

	state := monitorsDataSourceModel{}
	for _, m := range apiResp {
		state.Monitors = append(state.Monitors, monitorAPIToState(m))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
