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
	_ resource.Resource              = &environmentResource{}
	_ resource.ResourceWithConfigure = &environmentResource{}
)

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

type environmentResource struct {
	client *client.Client
}

type environmentResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Variables    types.String `tfsdk:"variables"`
	IsDefault    types.Bool   `tfsdk:"is_default"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateModified types.String `tfsdk:"date_modified"`
}

type environmentAPIModel struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Variables    interface{} `json:"variables"`
	IsDefault    bool        `json:"is_default"`
	DateCreated  string      `json:"date_created"`
	DateModified string      `json:"date_modified"`
}

func environmentAPIToState(api environmentAPIModel) environmentResourceModel {
	varsJSON, _ := json.Marshal(api.Variables)
	return environmentResourceModel{
		ID:           types.StringValue(api.ID),
		Name:         types.StringValue(api.Name),
		Variables:    types.StringValue(string(varsJSON)),
		IsDefault:    types.BoolValue(api.IsDefault),
		DateCreated:  types.StringValue(api.DateCreated),
		DateModified: types.StringValue(api.DateModified),
	}
}

func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix test environment (variables for synthetic tests).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the environment.",
			},
			"variables": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString("{}"),
				Description: "JSON-encoded key-value variables.",
			},
			"is_default": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether this is the default environment.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan environmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := environmentStateToPayload(plan)
	var apiResp environmentAPIModel
	err := r.client.Create(ctx, "/environments/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating environment", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, environmentAPIToState(apiResp))...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state environmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp environmentAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/environments/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, environmentAPIToState(apiResp))...)
}

func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan environmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state environmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := environmentStateToPayload(plan)
	var apiResp environmentAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/environments/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating environment", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, environmentAPIToState(apiResp))...)
}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state environmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/environments/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting environment", err.Error())
	}
}

func environmentStateToPayload(plan environmentResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":       plan.Name.ValueString(),
		"is_default": plan.IsDefault.ValueBool(),
	}

	if !plan.Variables.IsNull() && !plan.Variables.IsUnknown() {
		var vars interface{}
		if err := json.Unmarshal([]byte(plan.Variables.ValueString()), &vars); err == nil {
			payload["variables"] = vars
		}
	}

	return payload
}
