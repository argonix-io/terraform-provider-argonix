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
	_ resource.Resource              = &securityPolicyResource{}
	_ resource.ResourceWithConfigure = &securityPolicyResource{}
)

func NewSecurityPolicyResource() resource.Resource {
	return &securityPolicyResource{}
}

type securityPolicyResource struct {
	client *client.Client
}

type securityPolicyResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Rules        types.String `tfsdk:"rules"`
	Environment  types.String `tfsdk:"environment"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateModified types.String `tfsdk:"date_modified"`
}

type securityPolicyAPIModel struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Rules        interface{} `json:"rules"`
	Environment  string      `json:"environment"`
	IsActive     bool        `json:"is_active"`
	DateCreated  string      `json:"date_created"`
	DateModified string      `json:"date_modified"`
}

func securityPolicyAPIToState(api securityPolicyAPIModel) securityPolicyResourceModel {
	rulesJSON, _ := json.Marshal(api.Rules)

	return securityPolicyResourceModel{
		ID:           types.StringValue(api.ID),
		Name:         types.StringValue(api.Name),
		Description:  types.StringValue(api.Description),
		Rules:        types.StringValue(string(rulesJSON)),
		Environment:  types.StringValue(api.Environment),
		IsActive:     types.BoolValue(api.IsActive),
		DateCreated:  types.StringValue(api.DateCreated),
		DateModified: types.StringValue(api.DateModified),
	}
}

func (r *securityPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_policy"
}

func (r *securityPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix security policy (deployment gate).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the security policy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Description of the policy purpose.",
			},
			"rules": schema.StringAttribute{
				Required:    true,
				Description: `JSON-encoded policy rules. Example: {"max_critical": 0, "max_high": 5, "scan_max_age_hours": 24, "required_scan_types": ["trivy", "gitleaks"]}`,
			},
			"environment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Target environment (e.g. production, staging). Empty applies to all environments.",
			},
			"is_active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the policy is active.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *securityPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *securityPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan securityPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := securityPolicyStateToPayload(plan)
	var apiResp securityPolicyAPIModel
	err := r.client.Create(ctx, "/security/policies/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating security policy", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityPolicyAPIToState(apiResp))...)
}

func (r *securityPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state securityPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp securityPolicyAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/security/policies/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading security policy", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityPolicyAPIToState(apiResp))...)
}

func (r *securityPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan securityPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state securityPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := securityPolicyStateToPayload(plan)
	var apiResp securityPolicyAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/security/policies/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating security policy", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityPolicyAPIToState(apiResp))...)
}

func (r *securityPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state securityPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/security/policies/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting security policy", err.Error())
	}
}

func securityPolicyStateToPayload(plan securityPolicyResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
		"environment": plan.Environment.ValueString(),
		"is_active":   plan.IsActive.ValueBool(),
	}

	unmarshalJSONField(plan.Rules.ValueString(), "rules", payload)

	return payload
}
