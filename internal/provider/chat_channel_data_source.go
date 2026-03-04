package provider
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
	_ datasource.DataSource              = &chatChannelDataSource{}
	_ datasource.DataSourceWithConfigure = &chatChannelDataSource{}
	_ datasource.DataSource              = &chatChannelsDataSource{}
	_ datasource.DataSourceWithConfigure = &chatChannelsDataSource{}
)

func NewChatChannelDataSource() datasource.DataSource {
	return &chatChannelDataSource{}
}

type chatChannelDataSource struct {
	client *client.Client
}

func (d *chatChannelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chat_channel"
}

func (d *chatChannelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix chat channel by ID.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"channel_type":  schema.StringAttribute{Computed: true},
			"channel_id":    schema.StringAttribute{Computed: true},
			"channel_name":  schema.StringAttribute{Computed: true},
			"persona_id":    schema.StringAttribute{Computed: true},
			"connector_id":  schema.StringAttribute{Computed: true},
			"config":        schema.StringAttribute{Computed: true},
			"is_active":     schema.BoolAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *chatChannelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *chatChannelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp chatChannelAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/argos/chat-channels/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading chat channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, chatChannelAPIToState(apiResp))...)
}

// --- List chat channels data source ---

func NewChatChannelsDataSource() datasource.DataSource {
	return &chatChannelsDataSource{}
}

type chatChannelsDataSource struct {
	client *client.Client
}

type chatChannelsDataSourceModel struct {
	ChatChannels []chatChannelResourceModel `tfsdk:"chat_channels"`
}

func (d *chatChannelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chat_channels"
}

func (d *chatChannelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all chat channels in the organization.",
		Attributes: map[string]schema.Attribute{
			"chat_channels": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"channel_type":  schema.StringAttribute{Computed: true},
						"channel_id":    schema.StringAttribute{Computed: true},
						"channel_name":  schema.StringAttribute{Computed: true},
						"persona_id":    schema.StringAttribute{Computed: true},
						"connector_id":  schema.StringAttribute{Computed: true},
						"config":        schema.StringAttribute{Computed: true},
						"is_active":     schema.BoolAttribute{Computed: true},
						"date_created":  schema.StringAttribute{Computed: true},
						"date_modified": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *chatChannelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *chatChannelsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []chatChannelAPIModel
	err := d.client.List(ctx, "/argos/chat-channels/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing chat channels", err.Error())
		return
	}

	state := chatChannelsDataSourceModel{}
	for _, item := range apiResp {
		state.ChatChannels = append(state.ChatChannels, chatChannelAPIToState(item))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
