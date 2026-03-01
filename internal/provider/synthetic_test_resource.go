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
	_ resource.Resource              = &syntheticTestResource{}
	_ resource.ResourceWithConfigure = &syntheticTestResource{}
)

func NewSyntheticTestResource() resource.Resource {
	return &syntheticTestResource{}
}

type syntheticTestResource struct {
	client *client.Client
}

type syntheticTestResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	IsActive      types.Bool   `tfsdk:"is_active"`
	TestType      types.String `tfsdk:"test_type"`
	Steps         types.String `tfsdk:"steps"`
	CheckInterval types.Int64  `tfsdk:"check_interval"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	Tags          types.String `tfsdk:"tags"`
	Locations     types.String `tfsdk:"locations"`
	CurrentStatus types.String `tfsdk:"current_status"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

type syntheticTestAPIModel struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	IsActive      bool        `json:"is_active"`
	TestType      string      `json:"test_type"`
	Steps         interface{} `json:"steps"`
	CheckInterval int64       `json:"check_interval"`
	Timeout       int64       `json:"timeout"`
	Tags          interface{} `json:"tags"`
	Locations     interface{} `json:"locations"`
	CurrentStatus string      `json:"current_status"`
	DateCreated   string      `json:"date_created"`
	DateModified  string      `json:"date_modified"`
}

func syntheticTestAPIToState(api syntheticTestAPIModel) syntheticTestResourceModel {
	stepsJSON, _ := json.Marshal(api.Steps)
	tagsJSON, _ := json.Marshal(api.Tags)
	locationsJSON, _ := json.Marshal(api.Locations)
	return syntheticTestResourceModel{
		ID:            types.StringValue(api.ID),
		Name:          types.StringValue(api.Name),
		Description:   types.StringValue(api.Description),
		IsActive:      types.BoolValue(api.IsActive),
		TestType:      types.StringValue(api.TestType),
		Steps:         types.StringValue(string(stepsJSON)),
		CheckInterval: types.Int64Value(api.CheckInterval),
		Timeout:       types.Int64Value(api.Timeout),
		Tags:          types.StringValue(string(tagsJSON)),
		Locations:     types.StringValue(string(locationsJSON)),
		CurrentStatus: types.StringValue(api.CurrentStatus),
		DateCreated:   types.StringValue(api.DateCreated),
		DateModified:  types.StringValue(api.DateModified),
	}
}

func (r *syntheticTestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synthetic_test"
}

func (r *syntheticTestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix synthetic test.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the synthetic test.",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"test_type": schema.StringAttribute{
				Required:    true,
				Description: "Type: api or browser.",
			},
			"steps": schema.StringAttribute{
				Required:    true,
				Description: "JSON-encoded array of step objects.",
			},
			"check_interval": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(300),
			},
			"timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(30),
			},
			"tags": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("[]"),
			},
			"locations": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of region codes.",
			},
			"current_status": schema.StringAttribute{
				Computed: true,
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *syntheticTestResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *syntheticTestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan syntheticTestResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := syntheticTestStateToPayload(plan)
	var apiResp syntheticTestAPIModel
	err := r.client.Create(ctx, "/synthetic-tests/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating synthetic test", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, syntheticTestAPIToState(apiResp))...)
}

func (r *syntheticTestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state syntheticTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp syntheticTestAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/synthetic-tests/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading synthetic test", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, syntheticTestAPIToState(apiResp))...)
}

func (r *syntheticTestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan syntheticTestResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state syntheticTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := syntheticTestStateToPayload(plan)
	var apiResp syntheticTestAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/synthetic-tests/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating synthetic test", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, syntheticTestAPIToState(apiResp))...)
}

func (r *syntheticTestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state syntheticTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/synthetic-tests/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting synthetic test", err.Error())
	}
}

func syntheticTestStateToPayload(plan syntheticTestResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":           plan.Name.ValueString(),
		"description":    plan.Description.ValueString(),
		"is_active":      plan.IsActive.ValueBool(),
		"test_type":      plan.TestType.ValueString(),
		"check_interval": plan.CheckInterval.ValueInt64(),
		"timeout":        plan.Timeout.ValueInt64(),
	}

	var steps interface{}
	json.Unmarshal([]byte(plan.Steps.ValueString()), &steps)
	payload["steps"] = steps

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags interface{}
		json.Unmarshal([]byte(plan.Tags.ValueString()), &tags)
		payload["tags"] = tags
	}

	if !plan.Locations.IsNull() && !plan.Locations.IsUnknown() {
		var locations interface{}
		json.Unmarshal([]byte(plan.Locations.ValueString()), &locations)
		payload["locations"] = locations
	}

	return payload
}
