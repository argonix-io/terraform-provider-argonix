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
	_ resource.Resource              = &workflowResource{}
	_ resource.ResourceWithConfigure = &workflowResource{}
)

func NewWorkflowResource() resource.Resource {
	return &workflowResource{}
}

type workflowResource struct {
	client *client.Client
}

type workflowResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Slug                   types.String `tfsdk:"slug"`
	Description            types.String `tfsdk:"description"`
	Category               types.String `tfsdk:"category"`
	Steps                  types.String `tfsdk:"steps"`
	RequiredConnectorTypes types.String `tfsdk:"required_connector_types"`
	RequiresConfirmation   types.Bool   `tfsdk:"requires_confirmation"`
	Schedule               types.String `tfsdk:"schedule"`
	IsActive               types.Bool   `tfsdk:"is_active"`
	DateCreated            types.String `tfsdk:"date_created"`
	DateModified           types.String `tfsdk:"date_modified"`
}

type workflowAPIModel struct {
	ID                     string      `json:"id"`
	Name                   string      `json:"name"`
	Slug                   string      `json:"slug"`
	Description            *string     `json:"description"`
	Category               string      `json:"category"`
	Steps                  interface{} `json:"steps"`
	RequiredConnectorTypes interface{} `json:"required_connector_types"`
	RequiresConfirmation   bool        `json:"requires_confirmation"`
	Schedule               *string     `json:"schedule"`
	IsActive               bool        `json:"is_active"`
	DateCreated            string      `json:"date_created"`
	DateModified           string      `json:"date_modified"`
}

func workflowAPIToState(api workflowAPIModel) workflowResourceModel {
	stepsJSON, _ := json.Marshal(api.Steps)
	rctJSON, _ := json.Marshal(api.RequiredConnectorTypes)
	state := workflowResourceModel{
		ID:                     types.StringValue(api.ID),
		Name:                   types.StringValue(api.Name),
		Slug:                   types.StringValue(api.Slug),
		Category:               types.StringValue(api.Category),
		Steps:                  types.StringValue(string(stepsJSON)),
		RequiredConnectorTypes: types.StringValue(string(rctJSON)),
		RequiresConfirmation:   types.BoolValue(api.RequiresConfirmation),
		IsActive:               types.BoolValue(api.IsActive),
		DateCreated:            types.StringValue(api.DateCreated),
		DateModified:           types.StringValue(api.DateModified),
	}
	if api.Description != nil {
		state.Description = types.StringValue(*api.Description)
	} else {
		state.Description = types.StringValue("")
	}
	if api.Schedule != nil {
		state.Schedule = types.StringValue(*api.Schedule)
	} else {
		state.Schedule = types.StringValue("")
	}
	return state
}

func (r *workflowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (r *workflowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix Argos AI workflow.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the workflow.",
			},
			"slug": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "URL-friendly identifier. Auto-generated from name if not set.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"category": schema.StringAttribute{
				Required:    true,
				Description: "Workflow category: identity, incident, onboarding, devops, security, general.",
			},
			"steps": schema.StringAttribute{
				Required:    true,
				Description: "JSON-encoded workflow step definitions.",
			},
			"required_connector_types": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of required connector types.",
			},
			"requires_confirmation": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the workflow requires human confirmation before execution.",
			},
			"schedule": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Cron-style schedule expression for periodic execution.",
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

func (r *workflowResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workflowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := workflowStateToPayload(plan)
	var apiResp workflowAPIModel
	err := r.client.Create(ctx, "/argos/workflows/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating workflow", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, workflowAPIToState(apiResp))...)
}

func (r *workflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workflowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp workflowAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/argos/workflows/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, workflowAPIToState(apiResp))...)
}

func (r *workflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workflowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state workflowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := workflowStateToPayload(plan)
	var apiResp workflowAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/argos/workflows/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workflow", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, workflowAPIToState(apiResp))...)
}

func (r *workflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workflowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/argos/workflows/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workflow", err.Error())
	}
}

func workflowStateToPayload(plan workflowResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":                  plan.Name.ValueString(),
		"category":              plan.Category.ValueString(),
		"is_active":             plan.IsActive.ValueBool(),
		"requires_confirmation": plan.RequiresConfirmation.ValueBool(),
	}

	if v := plan.Slug.ValueString(); v != "" {
		payload["slug"] = v
	}
	if v := plan.Description.ValueString(); v != "" {
		payload["description"] = v
	}
	if v := plan.Schedule.ValueString(); v != "" {
		payload["schedule"] = v
	}

	var steps interface{}
	if err := json.Unmarshal([]byte(plan.Steps.ValueString()), &steps); err == nil {
		payload["steps"] = steps
	}

	if !plan.RequiredConnectorTypes.IsNull() && !plan.RequiredConnectorTypes.IsUnknown() {
		var rct interface{}
		if err := json.Unmarshal([]byte(plan.RequiredConnectorTypes.ValueString()), &rct); err == nil {
			payload["required_connector_types"] = rct
		}
	}

	return payload
}
