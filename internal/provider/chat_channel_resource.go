package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &chatChannelResource{}
	_ resource.ResourceWithConfigure = &chatChannelResource{}
)

func NewChatChannelResource() resource.Resource {
	return &chatChannelResource{}
}

type chatChannelResource struct {
	client *client.Client
}

type chatChannelResourceModel struct {
	ID           types.String `tfsdk:"id"`
	ChannelType  types.String `tfsdk:"channel_type"`
	ChannelID    types.String `tfsdk:"channel_id"`
	ChannelName  types.String `tfsdk:"channel_name"`
	PersonaID    types.String `tfsdk:"persona_id"`
	ConnectorID  types.String `tfsdk:"connector_id"`
	Config       types.String `tfsdk:"config"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateModified types.String `tfsdk:"date_modified"`
}

type chatChannelAPIModel struct {
	ID           string      `json:"id"`
	ChannelType  string      `json:"channel_type"`
	ChannelID    *string     `json:"channel_id"`
	ChannelName  *string     `json:"channel_name"`
	Persona      *string     `json:"persona"`
	Connector    *string     `json:"connector"`
	Config       interface{} `json:"config"`
	IsActive     bool        `json:"is_active"`
	DateCreated  string      `json:"date_created"`
	DateModified string      `json:"date_modified"`
}

func chatChannelAPIToState(api chatChannelAPIModel) chatChannelResourceModel {
	configJSON, _ := json.Marshal(api.Config)
	state := chatChannelResourceModel{
		ID:           types.StringValue(api.ID),
		ChannelType:  types.StringValue(api.ChannelType),
		IsActive:     types.BoolValue(api.IsActive),
		Config:       types.StringValue(string(configJSON)),
		DateCreated:  types.StringValue(api.DateCreated),
		DateModified: types.StringValue(api.DateModified),
	}
	if api.ChannelID != nil {
		state.ChannelID = types.StringValue(*api.ChannelID)
	} else {
		state.ChannelID = types.StringValue("")
	}
	if api.ChannelName != nil {
		state.ChannelName = types.StringValue(*api.ChannelName)
	} else {
		state.ChannelName = types.StringValue("")
	}
	if api.Persona != nil {
		state.PersonaID = types.StringValue(*api.Persona)
	} else {
		state.PersonaID = types.StringValue("")
	}
	if api.Connector != nil {
		state.ConnectorID = types.StringValue(*api.Connector)
	} else {
		state.ConnectorID = types.StringValue("")
	}
	return state
}

func (r *chatChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_chat_channel"
}

func (r *chatChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix Argos AI chat channel binding (Slack/Teams/Jira → Persona).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_type": schema.StringAttribute{
				Required:    true,
				Description: "Chat channel type: slack, teams, discord, jira.",
			},
			"channel_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "External channel identifier (e.g. Slack channel ID).",
			},
			"channel_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Display name of the channel.",
			},
			"persona_id": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the persona to use in this channel.",
			},
			"connector_id": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the connector for this channel.",
			},
			"config": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("{}"),
				Description: "JSON-encoded channel configuration.",
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *chatChannelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
		return
	}
	r.client = c
}

func (r *chatChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan chatChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := chatChannelStateToPayload(plan)
	var apiResp chatChannelAPIModel
	err := r.client.Create(ctx, "/argos/chat-channels/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating chat channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, chatChannelAPIToState(apiResp))...)
}

func (r *chatChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state chatChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp chatChannelAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/argos/chat-channels/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading chat channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, chatChannelAPIToState(apiResp))...)
}

func (r *chatChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan chatChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state chatChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := chatChannelStateToPayload(plan)
	var apiResp chatChannelAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/argos/chat-channels/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating chat channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, chatChannelAPIToState(apiResp))...)
}

func (r *chatChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state chatChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/argos/chat-channels/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting chat channel", err.Error())
	}
}

func chatChannelStateToPayload(plan chatChannelResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"channel_type": plan.ChannelType.ValueString(),
		"persona":      plan.PersonaID.ValueString(),
		"connector":    plan.ConnectorID.ValueString(),
		"is_active":    plan.IsActive.ValueBool(),
	}
	if v := plan.ChannelID.ValueString(); v != "" {
		payload["channel_id"] = v
	}
	if v := plan.ChannelName.ValueString(); v != "" {
		payload["channel_name"] = v
	}

	var config interface{}
	if err := json.Unmarshal([]byte(plan.Config.ValueString()), &config); err == nil {
		payload["config"] = config
	}

	return payload
}
