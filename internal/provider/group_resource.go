package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &groupResource{}
	_ resource.ResourceWithConfigure = &groupResource{}
)

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	client *client.Client
}

type groupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Tags         types.String `tfsdk:"tags"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateModified types.String `tfsdk:"date_modified"`
}

type groupAPIModel struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  *string     `json:"description"`
	Tags         interface{} `json:"tags"`
	DateCreated  string      `json:"date_created"`
	DateModified string      `json:"date_modified"`
}

func groupAPIToState(api groupAPIModel) groupResourceModel {
	tagsJSON, _ := json.Marshal(api.Tags)
	state := groupResourceModel{
		ID:           types.StringValue(api.ID),
		Name:         types.StringValue(api.Name),
		Tags:         types.StringValue(string(tagsJSON)),
		DateCreated:  types.StringValue(api.DateCreated),
		DateModified: types.StringValue(api.DateModified),
	}
	if api.Description != nil {
		state.Description = types.StringValue(*api.Description)
	} else {
		state.Description = types.StringValue("")
	}
	return state
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the group.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"tags": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("{}"),
				Description: "JSON-encoded tags object.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := groupStateToPayload(plan)
	var apiResp groupAPIModel
	err := r.client.Create(ctx, "/groups/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating group", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, groupAPIToState(apiResp))...)
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp groupAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/groups/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading group", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, groupAPIToState(apiResp))...)
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := groupStateToPayload(plan)
	var apiResp groupAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/groups/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating group", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, groupAPIToState(apiResp))...)
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/groups/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting group", err.Error())
	}
}

func groupStateToPayload(plan groupResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags interface{}
		json.Unmarshal([]byte(plan.Tags.ValueString()), &tags)
		payload["tags"] = tags
	}

	return payload
}
