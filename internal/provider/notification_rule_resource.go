package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &notificationRuleResource{}
	_ resource.ResourceWithConfigure = &notificationRuleResource{}
)

func NewNotificationRuleResource() resource.Resource {
	return &notificationRuleResource{}
}

type notificationRuleResource struct {
	client *client.Client
}

type notificationRuleResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	IsActive            types.Bool   `tfsdk:"is_active"`
	TriggerCondition    types.String `tfsdk:"trigger_condition"`
	ConsecutiveFailures types.Int64  `tfsdk:"consecutive_failures"`
	CooldownMinutes     types.Int64  `tfsdk:"cooldown_minutes"`
	AllMonitors         types.Bool   `tfsdk:"all_monitors"`
	AllSyntheticTests   types.Bool   `tfsdk:"all_synthetic_tests"`
	AutoInvestigate     types.Bool   `tfsdk:"auto_investigate"`
	MonitorTags         types.String `tfsdk:"monitor_tags"`
	Monitors            types.String `tfsdk:"monitors"`
	SyntheticTests      types.String `tfsdk:"synthetic_tests"`
	Channels            types.String `tfsdk:"channels"`
	DateCreated         types.String `tfsdk:"date_created"`
	DateModified        types.String `tfsdk:"date_modified"`
}

type notificationRuleAPIModel struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	IsActive            bool        `json:"is_active"`
	TriggerCondition    string      `json:"trigger_condition"`
	ConsecutiveFailures int64       `json:"consecutive_failures"`
	CooldownMinutes     int64       `json:"cooldown_minutes"`
	AllMonitors         bool        `json:"all_monitors"`
	AllSyntheticTests   bool        `json:"all_synthetic_tests"`
	AutoInvestigate     bool        `json:"auto_investigate"`
	MonitorTags         interface{} `json:"monitor_tags"`
	Monitors            interface{} `json:"monitors"`
	SyntheticTests      interface{} `json:"synthetic_tests"`
	Channels            interface{} `json:"channels"`
	DateCreated         string      `json:"date_created"`
	DateModified        string      `json:"date_modified"`
}

func notificationRuleAPIToState(api notificationRuleAPIModel) notificationRuleResourceModel {
	monitorTagsJSON, _ := json.Marshal(api.MonitorTags)
	monitorsJSON, _ := json.Marshal(api.Monitors)
	syntheticTestsJSON, _ := json.Marshal(api.SyntheticTests)
	channelsJSON, _ := json.Marshal(api.Channels)

	return notificationRuleResourceModel{
		ID:                  types.StringValue(api.ID),
		Name:                types.StringValue(api.Name),
		IsActive:            types.BoolValue(api.IsActive),
		TriggerCondition:    types.StringValue(api.TriggerCondition),
		ConsecutiveFailures: types.Int64Value(api.ConsecutiveFailures),
		CooldownMinutes:     types.Int64Value(api.CooldownMinutes),
		AllMonitors:         types.BoolValue(api.AllMonitors),
		AllSyntheticTests:   types.BoolValue(api.AllSyntheticTests),
		AutoInvestigate:     types.BoolValue(api.AutoInvestigate),
		MonitorTags:         types.StringValue(string(monitorTagsJSON)),
		Monitors:            types.StringValue(string(monitorsJSON)),
		SyntheticTests:      types.StringValue(string(syntheticTestsJSON)),
		Channels:            types.StringValue(string(channelsJSON)),
		DateCreated:         types.StringValue(api.DateCreated),
		DateModified:        types.StringValue(api.DateModified),
	}
}

func (r *notificationRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rule"
}

func (r *notificationRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix notification rule.",
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
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"trigger_condition": schema.StringAttribute{
				Required:    true,
				Description: "Trigger condition: status_change, goes_down, goes_up, degraded, ssl_expiry, test_failing, test_passing, test_run_complete.",
			},
			"consecutive_failures": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
			},
			"cooldown_minutes": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(5),
			},
			"all_monitors": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Apply to all monitors in the organization.",
			},
			"all_synthetic_tests": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Apply to all synthetic tests in the organization.",
			},
			"monitor_tags": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of tags to match monitors.",
			},
			"monitors": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of monitor UUIDs.",
			},
			"synthetic_tests": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of synthetic test UUIDs.",
			},
			"channels": schema.StringAttribute{
				Required:    true,
				Description: "JSON-encoded list of notification channel UUIDs.",
			},
			"auto_investigate": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "When triggered, Argos AI automatically investigates the root cause and posts analysis to channels.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *notificationRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *notificationRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := notificationRuleStateToPayload(plan)
	var apiResp notificationRuleAPIModel
	err := r.client.Create(ctx, "/notification-rules/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating notification rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, notificationRuleAPIToState(apiResp))...)
}

func (r *notificationRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp notificationRuleAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/notification-rules/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading notification rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, notificationRuleAPIToState(apiResp))...)
}

func (r *notificationRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state notificationRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := notificationRuleStateToPayload(plan)
	var apiResp notificationRuleAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/notification-rules/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating notification rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, notificationRuleAPIToState(apiResp))...)
}

func (r *notificationRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/notification-rules/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting notification rule", err.Error())
	}
}

func notificationRuleStateToPayload(plan notificationRuleResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":                 plan.Name.ValueString(),
		"is_active":            plan.IsActive.ValueBool(),
		"trigger_condition":    plan.TriggerCondition.ValueString(),
		"consecutive_failures": plan.ConsecutiveFailures.ValueInt64(),
		"cooldown_minutes":     plan.CooldownMinutes.ValueInt64(),
		"all_monitors":         plan.AllMonitors.ValueBool(),
		"all_synthetic_tests":  plan.AllSyntheticTests.ValueBool(),
		"auto_investigate":     plan.AutoInvestigate.ValueBool(),
	}

	unmarshalJSONField(plan.MonitorTags.ValueString(), "monitor_tags", payload)
	unmarshalJSONField(plan.Monitors.ValueString(), "monitors", payload)
	unmarshalJSONField(plan.SyntheticTests.ValueString(), "synthetic_tests", payload)
	unmarshalJSONField(plan.Channels.ValueString(), "channels", payload)

	return payload
}

func unmarshalJSONField(jsonStr, key string, payload map[string]interface{}) {
	if jsonStr != "" {
		var v interface{}
		if err := json.Unmarshal([]byte(jsonStr), &v); err == nil {
			payload[key] = v
		}
	}
}
