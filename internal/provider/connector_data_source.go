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
	_ datasource.DataSource              = &connectorDataSource{}
	_ datasource.DataSourceWithConfigure = &connectorDataSource{}
	_ datasource.DataSource              = &connectorsDataSource{}
	_ datasource.DataSourceWithConfigure = &connectorsDataSource{}
)

func NewConnectorDataSource() datasource.DataSource {
	return &connectorDataSource{}
}

type connectorDataSource struct {
	client *client.Client
}

func (d *connectorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector"
}

func (d *connectorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix connector by ID.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true},
			"name":           schema.StringAttribute{Computed: true},
			"connector_type": schema.StringAttribute{Computed: true},
			"is_active":      schema.BoolAttribute{Computed: true},
			"config":         schema.StringAttribute{Computed: true, Sensitive: true},
			"capabilities":   schema.StringAttribute{Computed: true},
			"tags":           schema.StringAttribute{Computed: true},
			"date_created":   schema.StringAttribute{Computed: true},
			"date_modified":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *connectorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *connectorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp connectorAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/argos/connectors/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading connector", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, connectorAPIToState(apiResp))...)
}

// --- List connectors data source ---

func NewConnectorsDataSource() datasource.DataSource {
	return &connectorsDataSource{}
}

type connectorsDataSource struct {
	client *client.Client
}

type connectorsDataSourceModel struct {
	Connectors []connectorResourceModel `tfsdk:"connectors"`
}

func (d *connectorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connectors"
}

func (d *connectorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all connectors in the organization.",
		Attributes: map[string]schema.Attribute{
			"connectors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":             schema.StringAttribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"connector_type": schema.StringAttribute{Computed: true},
						"is_active":      schema.BoolAttribute{Computed: true},
						"config":         schema.StringAttribute{Computed: true, Sensitive: true},
						"capabilities":   schema.StringAttribute{Computed: true},
						"tags":           schema.StringAttribute{Computed: true},
						"date_created":   schema.StringAttribute{Computed: true},
						"date_modified":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *connectorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *connectorsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []connectorAPIModel
	err := d.client.List(ctx, "/argos/connectors/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing connectors", err.Error())
		return
	}

	state := connectorsDataSourceModel{}
	for _, item := range apiResp {
		state.Connectors = append(state.Connectors, connectorAPIToState(item))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
