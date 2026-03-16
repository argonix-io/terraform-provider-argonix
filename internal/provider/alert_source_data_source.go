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
	_ datasource.DataSource              = &alertSourceDataSource{}
	_ datasource.DataSourceWithConfigure = &alertSourceDataSource{}
	_ datasource.DataSource              = &alertSourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &alertSourcesDataSource{}
)

// --- Single alert source data source ---

func NewAlertSourceDataSource() datasource.DataSource {
	return &alertSourceDataSource{}
}

type alertSourceDataSource struct {
	client *client.Client
}

func (d *alertSourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_source"
}

func (d *alertSourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix alert source by ID.",
		Attributes: map[string]schema.Attribute{
			"id":                   schema.StringAttribute{Required: true},
			"name":                 schema.StringAttribute{Computed: true},
			"source_type":          schema.StringAttribute{Computed: true},
			"is_active":            schema.BoolAttribute{Computed: true},
			"connector":            schema.StringAttribute{Computed: true},
			"filters":              schema.StringAttribute{Computed: true},
			"auto_investigate":     schema.BoolAttribute{Computed: true},
			"auto_remediate":       schema.BoolAttribute{Computed: true},
			"remediation_strategy": schema.StringAttribute{Computed: true},
			"channels":             schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"webhook_secret":       schema.StringAttribute{Computed: true, Sensitive: true},
			"webhook_url":          schema.StringAttribute{Computed: true},
			"last_received_at":     schema.StringAttribute{Computed: true},
			"total_received":       schema.Int64Attribute{Computed: true},
			"date_created":         schema.StringAttribute{Computed: true},
			"date_modified":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *alertSourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertSourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp alertSourceAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/alert-sources/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading alert source", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertSourceAPIToState(ctx, apiResp))...)
}

// --- List alert sources data source ---

func NewAlertSourcesDataSource() datasource.DataSource {
	return &alertSourcesDataSource{}
}

type alertSourcesDataSource struct {
	client *client.Client
}

type alertSourcesDataSourceModel struct {
	AlertSources []alertSourceResourceModel `tfsdk:"alert_sources"`
}

func (d *alertSourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_sources"
}

func (d *alertSourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all alert sources in the organization.",
		Attributes: map[string]schema.Attribute{
			"alert_sources": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"source_type":          schema.StringAttribute{Computed: true},
						"is_active":            schema.BoolAttribute{Computed: true},
						"connector":            schema.StringAttribute{Computed: true},
						"filters":              schema.StringAttribute{Computed: true},
						"auto_investigate":     schema.BoolAttribute{Computed: true},
						"auto_remediate":       schema.BoolAttribute{Computed: true},
						"remediation_strategy": schema.StringAttribute{Computed: true},
						"channels":             schema.ListAttribute{Computed: true, ElementType: types.StringType},
						"webhook_secret":       schema.StringAttribute{Computed: true, Sensitive: true},
						"webhook_url":          schema.StringAttribute{Computed: true},
						"last_received_at":     schema.StringAttribute{Computed: true},
						"total_received":       schema.Int64Attribute{Computed: true},
						"date_created":         schema.StringAttribute{Computed: true},
						"date_modified":        schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *alertSourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertSourcesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []alertSourceAPIModel
	err := d.client.List(ctx, "/alert-sources/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing alert sources", err.Error())
		return
	}

	state := alertSourcesDataSourceModel{}
	for _, src := range apiResp {
		state.AlertSources = append(state.AlertSources, alertSourceAPIToState(ctx, src))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
