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
	_ datasource.DataSource              = &alertChannelDataSource{}
	_ datasource.DataSourceWithConfigure = &alertChannelDataSource{}
	_ datasource.DataSource              = &alertChannelsDataSource{}
	_ datasource.DataSourceWithConfigure = &alertChannelsDataSource{}
)

func NewAlertChannelDataSource() datasource.DataSource {
	return &alertChannelDataSource{}
}

type alertChannelDataSource struct {
	client *client.Client
}

func (d *alertChannelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_channel"
}

func (d *alertChannelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix alert channel by ID.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"name":          schema.StringAttribute{Computed: true},
			"channel_type":  schema.StringAttribute{Computed: true},
			"is_active":     schema.BoolAttribute{Computed: true},
			"config":        schema.StringAttribute{Computed: true, Sensitive: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *alertChannelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertChannelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp alertChannelAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/alert-channels/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading alert channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertChannelAPIToState(apiResp))...)
}

// --- List alert channels data source ---

func NewAlertChannelsDataSource() datasource.DataSource {
	return &alertChannelsDataSource{}
}

type alertChannelsDataSource struct {
	client *client.Client
}

type alertChannelsDataSourceModel struct {
	AlertChannels []alertChannelResourceModel `tfsdk:"alert_channels"`
}

func (d *alertChannelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_channels"
}

func (d *alertChannelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all alert channels in the organization.",
		Attributes: map[string]schema.Attribute{
			"alert_channels": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"name":          schema.StringAttribute{Computed: true},
						"channel_type":  schema.StringAttribute{Computed: true},
						"is_active":     schema.BoolAttribute{Computed: true},
						"config":        schema.StringAttribute{Computed: true, Sensitive: true},
						"date_created":  schema.StringAttribute{Computed: true},
						"date_modified": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *alertChannelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertChannelsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []alertChannelAPIModel
	err := d.client.List(ctx, "/alert-channels/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing alert channels", err.Error())
		return
	}

	state := alertChannelsDataSourceModel{}
	for _, ch := range apiResp {
		state.AlertChannels = append(state.AlertChannels, alertChannelAPIToState(ch))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
