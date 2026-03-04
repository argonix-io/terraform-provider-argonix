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
	_ resource.Resource              = &connectorResource{}
	_ resource.ResourceWithConfigure = &connectorResource{}
)

func NewConnectorResource() resource.Resource {
	return &connectorResource{}
}

type connectorResource struct {
	client *client.Client
}

type connectorResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	ConnectorType types.String `tfsdk:"connector_type"`
	IsActive      types.Bool   `tfsdk:"is_active"`
	Config        types.String `tfsdk:"config"`
	Capabilities  types.String `tfsdk:"capabilities"`
	Tags          types.String `tfsdk:"tags"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

type connectorAPIModel struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	ConnectorType string      `json:"connector_type"`
	IsActive      bool        `json:"is_active"`
	Config        interface{} `json:"config"`
	Capabilities  interface{} `json:"capabilities"`
	Tags          interface{} `json:"tags"`
	DateCreated   string      `json:"date_created"`
	DateModified  string      `json:"date_modified"`
}

func connectorAPIToState(api connectorAPIModel) connectorResourceModel {
	configJSON, _ := json.Marshal(api.Config)
	capsJSON, _ := json.Marshal(api.Capabilities)
	tagsJSON, _ := json.Marshal(api.Tags)
	return connectorResourceModel{
		ID:            types.StringValue(api.ID),
		Name:          types.StringValue(api.Name),
		ConnectorType: types.StringValue(api.ConnectorType),
		IsActive:      types.BoolValue(api.IsActive),
		Config:        types.StringValue(string(configJSON)),
		Capabilities:  types.StringValue(string(capsJSON)),
		Tags:          types.StringValue(string(tagsJSON)),
		DateCreated:   types.StringValue(api.DateCreated),
		DateModified:  types.StringValue(api.DateModified),
	}
}

func (r *connectorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connector"
}

func (r *connectorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix Argos AI connector.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the connector.",
			},
			"connector_type": schema.StringAttribute{
				Required:    true,
				Description: "Connector type: slack, jira, pagerduty, opsgenie, discord, telegram, email, webhook, gitlab, github, kubernetes, argocd, terraform, prometheus, loki, elasticsearch, tempo, grafana, datadog, servicenow, aws, gcp, azure, google_workspace, microsoft_365, keycloak, confluence, notion.",
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"config": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "JSON-encoded connector configuration (varies by type).",
			},
			"capabilities": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of capabilities.",
			},
			"tags": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded tags list.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *connectorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *connectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan connectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := connectorStateToPayload(plan)
	var apiResp connectorAPIModel
	err := r.client.Create(ctx, "/argos/connectors/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating connector", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, connectorAPIToState(apiResp))...)
}

func (r *connectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state connectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp connectorAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/argos/connectors/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading connector", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, connectorAPIToState(apiResp))...)
}

func (r *connectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan connectorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state connectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := connectorStateToPayload(plan)
	var apiResp connectorAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/argos/connectors/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating connector", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, connectorAPIToState(apiResp))...)
}

func (r *connectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state connectorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/argos/connectors/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting connector", err.Error())
	}
}

func connectorStateToPayload(plan connectorResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":           plan.Name.ValueString(),
		"connector_type": plan.ConnectorType.ValueString(),
		"is_active":      plan.IsActive.ValueBool(),
	}

	var config interface{}
	if err := json.Unmarshal([]byte(plan.Config.ValueString()), &config); err == nil {
		payload["config"] = config
	}

	if !plan.Capabilities.IsNull() && !plan.Capabilities.IsUnknown() {
		var caps interface{}
		if err := json.Unmarshal([]byte(plan.Capabilities.ValueString()), &caps); err == nil {
			payload["capabilities"] = caps
		}
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags interface{}
		if err := json.Unmarshal([]byte(plan.Tags.ValueString()), &tags); err == nil {
			payload["tags"] = tags
		}
	}

	return payload
}
