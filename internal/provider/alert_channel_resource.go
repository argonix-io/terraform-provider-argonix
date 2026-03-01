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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &alertChannelResource{}
	_ resource.ResourceWithConfigure = &alertChannelResource{}
)

func NewAlertChannelResource() resource.Resource {
	return &alertChannelResource{}
}

type alertChannelResource struct {
	client *client.Client
}

type alertChannelResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ChannelType  types.String `tfsdk:"channel_type"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	Config       types.String `tfsdk:"config"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateModified types.String `tfsdk:"date_modified"`
}

type alertChannelAPIModel struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	ChannelType  string      `json:"channel_type"`
	IsActive     bool        `json:"is_active"`
	Config       interface{} `json:"config"`
	DateCreated  string      `json:"date_created"`
	DateModified string      `json:"date_modified"`
}

func alertChannelAPIToState(api alertChannelAPIModel) alertChannelResourceModel {
	configJSON, _ := json.Marshal(api.Config)
	return alertChannelResourceModel{
		ID:           types.StringValue(api.ID),
		Name:         types.StringValue(api.Name),
		ChannelType:  types.StringValue(api.ChannelType),
		IsActive:     types.BoolValue(api.IsActive),
		Config:       types.StringValue(string(configJSON)),
		DateCreated:  types.StringValue(api.DateCreated),
		DateModified: types.StringValue(api.DateModified),
	}
}

func (r *alertChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_channel"
}

func (r *alertChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix alert channel (notification destination).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"channel_type": schema.StringAttribute{
				Required:    true,
				Description: "Channel type: email, slack, webhook, pagerduty, opsgenie, telegram, discord, teams, jira.",
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"config": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "JSON-encoded channel configuration (varies by type).",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *alertChannelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *alertChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alertChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := alertChannelStateToPayload(plan)
	var apiResp alertChannelAPIModel
	err := r.client.Create(ctx, "/alert-channels/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating alert channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertChannelAPIToState(apiResp))...)
}

func (r *alertChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp alertChannelAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/alert-channels/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading alert channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertChannelAPIToState(apiResp))...)
}

func (r *alertChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alertChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state alertChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := alertChannelStateToPayload(plan)
	var apiResp alertChannelAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/alert-channels/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating alert channel", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertChannelAPIToState(apiResp))...)
}

func (r *alertChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/alert-channels/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting alert channel", err.Error())
	}
}

func alertChannelStateToPayload(plan alertChannelResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":         plan.Name.ValueString(),
		"channel_type": plan.ChannelType.ValueString(),
		"is_active":    plan.IsActive.ValueBool(),
	}

	var config interface{}
	if err := json.Unmarshal([]byte(plan.Config.ValueString()), &config); err == nil {
		payload["config"] = config
	}

	return payload
}
