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
	_ resource.Resource              = &alertSourceResource{}
	_ resource.ResourceWithConfigure = &alertSourceResource{}
)

func NewAlertSourceResource() resource.Resource {
	return &alertSourceResource{}
}

type alertSourceResource struct {
	client *client.Client
}

type alertSourceResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	SourceType      types.String `tfsdk:"source_type"`
	IsActive        types.Bool   `tfsdk:"is_active"`
	Connector       types.String `tfsdk:"connector"`
	Filters         types.String `tfsdk:"filters"`
	AutoInvestigate types.Bool   `tfsdk:"auto_investigate"`
	AutoRemediate   types.Bool   `tfsdk:"auto_remediate"`
	Channels        types.List   `tfsdk:"channels"`
	WebhookSecret   types.String `tfsdk:"webhook_secret"`
	WebhookURL      types.String `tfsdk:"webhook_url"`
	LastReceivedAt  types.String `tfsdk:"last_received_at"`
	TotalReceived   types.Int64  `tfsdk:"total_received"`
	DateCreated     types.String `tfsdk:"date_created"`
	DateModified    types.String `tfsdk:"date_modified"`
}

type alertSourceAPIModel struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	SourceType      string      `json:"source_type"`
	IsActive        bool        `json:"is_active"`
	Connector       *string     `json:"connector"`
	Filters         interface{} `json:"filters"`
	AutoInvestigate bool        `json:"auto_investigate"`
	AutoRemediate   bool        `json:"auto_remediate"`
	Channels        []string    `json:"channels"`
	WebhookSecret   string      `json:"webhook_secret"`
	WebhookURL      string      `json:"webhook_url"`
	LastReceivedAt  *string     `json:"last_received_at"`
	TotalReceived   int64       `json:"total_received"`
	DateCreated     string      `json:"date_created"`
	DateModified    string      `json:"date_modified"`
}

func alertSourceAPIToState(ctx context.Context, api alertSourceAPIModel) alertSourceResourceModel {
	filtersJSON, _ := json.Marshal(api.Filters)

	channelElems := make([]types.String, 0, len(api.Channels))
	for _, ch := range api.Channels {
		channelElems = append(channelElems, types.StringValue(ch))
	}
	channelsList, _ := types.ListValueFrom(ctx, types.StringType, channelElems)

	connector := types.StringNull()
	if api.Connector != nil {
		connector = types.StringValue(*api.Connector)
	}

	lastReceived := types.StringNull()
	if api.LastReceivedAt != nil {
		lastReceived = types.StringValue(*api.LastReceivedAt)
	}

	return alertSourceResourceModel{
		ID:              types.StringValue(api.ID),
		Name:            types.StringValue(api.Name),
		SourceType:      types.StringValue(api.SourceType),
		IsActive:        types.BoolValue(api.IsActive),
		Connector:       connector,
		Filters:         types.StringValue(string(filtersJSON)),
		AutoInvestigate: types.BoolValue(api.AutoInvestigate),
		AutoRemediate:   types.BoolValue(api.AutoRemediate),
		Channels:        channelsList,
		WebhookSecret:   types.StringValue(api.WebhookSecret),
		WebhookURL:      types.StringValue(api.WebhookURL),
		LastReceivedAt:  lastReceived,
		TotalReceived:   types.Int64Value(api.TotalReceived),
		DateCreated:     types.StringValue(api.DateCreated),
		DateModified:    types.StringValue(api.DateModified),
	}
}

func (r *alertSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_source"
}

func (r *alertSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix alert source — an external alerting system (Datadog, Grafana, Alertmanager, etc.) that sends alerts via webhook for AI-powered investigation and remediation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":        schema.StringAttribute{Required: true},
			"source_type": schema.StringAttribute{Required: true, Description: "Source type: alertmanager, datadog, grafana, pagerduty, opsgenie, generic."},
			"is_active":   schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true)},
			"connector":   schema.StringAttribute{Optional: true, Description: "ID of a connector to the alert source (e.g. Datadog connector) for enriched investigation."},
			"filters": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("{}"),
				Description: `JSON-encoded filter rules: {"min_severity": "warning", "include_tags": ["prod"], "exclude_tags": ["test"]}.`,
			},
			"auto_investigate": schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(true), Description: "Automatically trigger Argos AI investigation on ingested alerts."},
			"auto_remediate":   schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(false), Description: "Automatically trigger Argos AI remediation after investigation."},
			"channels": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of alert channel IDs for investigation/remediation results.",
			},
			"webhook_secret":   schema.StringAttribute{Computed: true, Sensitive: true, Description: "Auto-generated secret token in the webhook URL."},
			"webhook_url":      schema.StringAttribute{Computed: true, Description: "Full webhook URL to configure in the external alert source."},
			"last_received_at": schema.StringAttribute{Computed: true},
			"total_received":   schema.Int64Attribute{Computed: true},
			"date_created":     schema.StringAttribute{Computed: true},
			"date_modified":    schema.StringAttribute{Computed: true},
		},
	}
}

func (r *alertSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *alertSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alertSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := alertSourceStateToPayload(ctx, plan)
	var apiResp alertSourceAPIModel
	err := r.client.Create(ctx, "/alert-sources/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating alert source", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertSourceAPIToState(ctx, apiResp))...)
}

func (r *alertSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp alertSourceAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/alert-sources/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading alert source", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertSourceAPIToState(ctx, apiResp))...)
}

func (r *alertSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alertSourceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state alertSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := alertSourceStateToPayload(ctx, plan)
	var apiResp alertSourceAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/alert-sources/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating alert source", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertSourceAPIToState(ctx, apiResp))...)
}

func (r *alertSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertSourceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/alert-sources/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting alert source", err.Error())
	}
}

func alertSourceStateToPayload(ctx context.Context, plan alertSourceResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":             plan.Name.ValueString(),
		"source_type":      plan.SourceType.ValueString(),
		"is_active":        plan.IsActive.ValueBool(),
		"auto_investigate": plan.AutoInvestigate.ValueBool(),
		"auto_remediate":   plan.AutoRemediate.ValueBool(),
	}

	if !plan.Connector.IsNull() && !plan.Connector.IsUnknown() {
		payload["connector"] = plan.Connector.ValueString()
	}

	if !plan.Filters.IsNull() && !plan.Filters.IsUnknown() {
		var filters interface{}
		if err := json.Unmarshal([]byte(plan.Filters.ValueString()), &filters); err == nil {
			payload["filters"] = filters
		}
	}

	if !plan.Channels.IsNull() && !plan.Channels.IsUnknown() {
		var channels []string
		plan.Channels.ElementsAs(ctx, &channels, false)
		payload["channels"] = channels
	}

	return payload
}
