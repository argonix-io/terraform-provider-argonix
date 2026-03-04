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
	_ resource.Resource              = &personaResource{}
	_ resource.ResourceWithConfigure = &personaResource{}
)

func NewPersonaResource() resource.Resource {
	return &personaResource{}
}

type personaResource struct {
	client *client.Client
}

type personaResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Template      types.String `tfsdk:"template"`
	SystemPrompt  types.String `tfsdk:"system_prompt"`
	IsActive      types.Bool   `tfsdk:"is_active"`
	ConnectorIDs  types.String `tfsdk:"connector_ids"`
	AllowedTools  types.String `tfsdk:"allowed_tools"`
	ApprovalRules types.String `tfsdk:"approval_rules"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

type personaAPIModel struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Description   *string     `json:"description"`
	Template      string      `json:"template"`
	SystemPrompt  *string     `json:"system_prompt"`
	IsActive      bool        `json:"is_active"`
	ConnectorIDs  interface{} `json:"connector_ids"`
	AllowedTools  interface{} `json:"allowed_tools"`
	ApprovalRules interface{} `json:"approval_rules"`
	DateCreated   string      `json:"date_created"`
	DateModified  string      `json:"date_modified"`
}

func personaAPIToState(api personaAPIModel) personaResourceModel {
	connJSON, _ := json.Marshal(api.ConnectorIDs)
	toolsJSON, _ := json.Marshal(api.AllowedTools)
	rulesJSON, _ := json.Marshal(api.ApprovalRules)
	state := personaResourceModel{
		ID:            types.StringValue(api.ID),
		Name:          types.StringValue(api.Name),
		Template:      types.StringValue(api.Template),
		IsActive:      types.BoolValue(api.IsActive),
		ConnectorIDs:  types.StringValue(string(connJSON)),
		AllowedTools:  types.StringValue(string(toolsJSON)),
		ApprovalRules: types.StringValue(string(rulesJSON)),
		DateCreated:   types.StringValue(api.DateCreated),
		DateModified:  types.StringValue(api.DateModified),
	}
	if api.Description != nil {
		state.Description = types.StringValue(*api.Description)
	} else {
		state.Description = types.StringValue("")
	}
	if api.SystemPrompt != nil {
		state.SystemPrompt = types.StringValue(*api.SystemPrompt)
	} else {
		state.SystemPrompt = types.StringValue("")
	}
	return state
}

func (r *personaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_persona"
}

func (r *personaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix Argos AI persona.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the persona.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"template": schema.StringAttribute{
				Required:    true,
				Description: "Persona template: devops, it_support, hr, security, custom.",
			},
			"system_prompt": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Custom system prompt for the persona.",
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"connector_ids": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of connector UUIDs to attach.",
			},
			"allowed_tools": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of allowed tool names.",
			},
			"approval_rules": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded approval rule definitions.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *personaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *personaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan personaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := personaStateToPayload(plan)
	var apiResp personaAPIModel
	err := r.client.Create(ctx, "/argos/personas/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating persona", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, personaAPIToState(apiResp))...)
}

func (r *personaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state personaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp personaAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/argos/personas/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading persona", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, personaAPIToState(apiResp))...)
}

func (r *personaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan personaResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state personaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := personaStateToPayload(plan)
	var apiResp personaAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/argos/personas/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating persona", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, personaAPIToState(apiResp))...)
}

func (r *personaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state personaResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/argos/personas/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting persona", err.Error())
	}
}

func personaStateToPayload(plan personaResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":      plan.Name.ValueString(),
		"template":  plan.Template.ValueString(),
		"is_active": plan.IsActive.ValueBool(),
	}

	if v := plan.Description.ValueString(); v != "" {
		payload["description"] = v
	}
	if v := plan.SystemPrompt.ValueString(); v != "" {
		payload["system_prompt"] = v
	}

	for _, field := range []struct{ key, val string }{
		{"connector_ids", plan.ConnectorIDs.ValueString()},
		{"allowed_tools", plan.AllowedTools.ValueString()},
		{"approval_rules", plan.ApprovalRules.ValueString()},
	} {
		if field.val != "" {
			var v interface{}
			if err := json.Unmarshal([]byte(field.val), &v); err == nil {
				payload[field.key] = v
			}
		}
	}

	return payload
}
