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
	_ resource.Resource              = &detectionRuleResource{}
	_ resource.ResourceWithConfigure = &detectionRuleResource{}
)

func NewDetectionRuleResource() resource.Resource {
	return &detectionRuleResource{}
}

type detectionRuleResource struct {
	client *client.Client
}

type detectionRuleResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	RuleType        types.String `tfsdk:"rule_type"`
	Severity        types.String `tfsdk:"severity"`
	IsActive        types.Bool   `tfsdk:"is_active"`
	Config          types.String `tfsdk:"config"`
	MitreTactic     types.String `tfsdk:"mitre_tactic"`
	MitreTechnique  types.String `tfsdk:"mitre_technique"`
	DetectionsCount types.Int64  `tfsdk:"detections_count"`
	DateCreated     types.String `tfsdk:"date_created"`
	DateModified    types.String `tfsdk:"date_modified"`
}

type detectionRuleAPIModel struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	RuleType        string      `json:"rule_type"`
	Severity        string      `json:"severity"`
	IsActive        bool        `json:"is_active"`
	Config          interface{} `json:"config"`
	MitreTactic     string      `json:"mitre_tactic"`
	MitreTechnique  string      `json:"mitre_technique"`
	DetectionsCount int64       `json:"detections_count"`
	DateCreated     string      `json:"date_created"`
	DateModified    string      `json:"date_modified"`
}

func detectionRuleAPIToState(api detectionRuleAPIModel) detectionRuleResourceModel {
	configJSON, _ := json.Marshal(api.Config)

	return detectionRuleResourceModel{
		ID:              types.StringValue(api.ID),
		Name:            types.StringValue(api.Name),
		Description:     types.StringValue(api.Description),
		RuleType:        types.StringValue(api.RuleType),
		Severity:        types.StringValue(api.Severity),
		IsActive:        types.BoolValue(api.IsActive),
		Config:          types.StringValue(string(configJSON)),
		MitreTactic:     types.StringValue(api.MitreTactic),
		MitreTechnique:  types.StringValue(api.MitreTechnique),
		DetectionsCount: types.Int64Value(api.DetectionsCount),
		DateCreated:     types.StringValue(api.DateCreated),
		DateModified:    types.StringValue(api.DateModified),
	}
}

func (r *detectionRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detection_rule"
}

func (r *detectionRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix threat detection rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the detection rule.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Description of the rule purpose.",
			},
			"rule_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of detection rule: threshold, pattern, or sequence.",
			},
			"severity": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("medium"),
				Description: "Severity level: critical, high, medium, or low.",
			},
			"is_active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the detection rule is active.",
			},
			"config": schema.StringAttribute{
				Required: true,
				Description: `JSON-encoded rule configuration. ` +
					`threshold: {"field": "severity", "value": "critical", "count": 5, "window_minutes": 10}. ` +
					`pattern: {"field": "title", "regex": "brute.force|credential.stuff"}. ` +
					`sequence: {"tactics": ["reconnaissance", "initial-access", "exfiltration"], "window_minutes": 60}.`,
			},
			"mitre_tactic": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "MITRE ATT&CK tactic ID (e.g. initial-access, lateral-movement).",
			},
			"mitre_technique": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "MITRE ATT&CK technique ID (e.g. T1078, T1059).",
			},
			"detections_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of threat detections triggered by this rule.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *detectionRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *detectionRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan detectionRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := detectionRuleStateToPayload(plan)
	var apiResp detectionRuleAPIModel
	err := r.client.Create(ctx, "/security/detection-rules/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating detection rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, detectionRuleAPIToState(apiResp))...)
}

func (r *detectionRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state detectionRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp detectionRuleAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/security/detection-rules/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading detection rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, detectionRuleAPIToState(apiResp))...)
}

func (r *detectionRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan detectionRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state detectionRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := detectionRuleStateToPayload(plan)
	var apiResp detectionRuleAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/security/detection-rules/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating detection rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, detectionRuleAPIToState(apiResp))...)
}

func (r *detectionRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state detectionRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/security/detection-rules/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting detection rule", err.Error())
	}
}

func detectionRuleStateToPayload(plan detectionRuleResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":            plan.Name.ValueString(),
		"description":     plan.Description.ValueString(),
		"rule_type":       plan.RuleType.ValueString(),
		"severity":        plan.Severity.ValueString(),
		"is_active":       plan.IsActive.ValueBool(),
		"mitre_tactic":    plan.MitreTactic.ValueString(),
		"mitre_technique": plan.MitreTechnique.ValueString(),
	}

	unmarshalJSONField(plan.Config.ValueString(), "config", payload)

	return payload
}
