package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ===========================================================================
// Test Plan Resource
// ===========================================================================

var (
	_ resource.Resource              = &testPlanResource{}
	_ resource.ResourceWithConfigure = &testPlanResource{}
)

func NewTestPlanResource() resource.Resource {
	return &testPlanResource{}
}

type testPlanResource struct {
	client *client.Client
}

type testPlanResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Suites       types.String `tfsdk:"suites"`
	Tags         types.String `tfsdk:"tags"`
	EndDate      types.String `tfsdk:"end_date"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateModified types.String `tfsdk:"date_modified"`
}

type testPlanAPIModel struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Suites       interface{} `json:"suites"`
	Tags         interface{} `json:"tags"`
	EndDate      *string     `json:"end_date"`
	DateCreated  string      `json:"date_created"`
	DateModified string      `json:"date_modified"`
}

func testPlanAPIToState(api testPlanAPIModel) testPlanResourceModel {
	suitesJSON, _ := json.Marshal(api.Suites)
	tagsJSON, _ := json.Marshal(api.Tags)

	state := testPlanResourceModel{
		ID:           types.StringValue(api.ID),
		Name:         types.StringValue(api.Name),
		Description:  types.StringValue(api.Description),
		Suites:       types.StringValue(string(suitesJSON)),
		Tags:         types.StringValue(string(tagsJSON)),
		DateCreated:  types.StringValue(api.DateCreated),
		DateModified: types.StringValue(api.DateModified),
	}

	if api.EndDate != nil {
		state.EndDate = types.StringValue(*api.EndDate)
	} else {
		state.EndDate = types.StringNull()
	}

	return state
}

func (r *testPlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test_plan"
}

func (r *testPlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix test plan.",
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
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"suites": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of test suite UUIDs.",
			},
			"tags": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("[]"),
			},
			"end_date": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Target completion date (YYYY-MM-DD).",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *testPlanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *testPlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan testPlanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := testPlanStateToPayload(plan)
	var apiResp testPlanAPIModel
	err := r.client.Create(ctx, "/test-plans/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating test plan", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testPlanAPIToState(apiResp))...)
}

func (r *testPlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state testPlanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp testPlanAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/test-plans/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading test plan", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testPlanAPIToState(apiResp))...)
}

func (r *testPlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan testPlanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state testPlanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := testPlanStateToPayload(plan)
	var apiResp testPlanAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/test-plans/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating test plan", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testPlanAPIToState(apiResp))...)
}

func (r *testPlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state testPlanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/test-plans/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting test plan", err.Error())
	}
}

func testPlanStateToPayload(plan testPlanResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	}

	unmarshalJSONFieldIfSet(plan.Suites, "suites", payload)
	unmarshalJSONFieldIfSet(plan.Tags, "tags", payload)

	if !plan.EndDate.IsNull() && !plan.EndDate.IsUnknown() {
		payload["end_date"] = plan.EndDate.ValueString()
	}

	return payload
}

// ===========================================================================
// Test Plan Data Source
// ===========================================================================

var (
	_ datasource.DataSource              = &testPlanDataSource{}
	_ datasource.DataSourceWithConfigure = &testPlanDataSource{}
)

func NewTestPlanDataSource() datasource.DataSource {
	return &testPlanDataSource{}
}

type testPlanDataSource struct {
	client *client.Client
}

func (d *testPlanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test_plan"
}

func (d *testPlanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Fetches a single Argonix test plan by ID.",
		Attributes: map[string]dsschema.Attribute{
			"id":            dsschema.StringAttribute{Required: true},
			"name":          dsschema.StringAttribute{Computed: true},
			"description":   dsschema.StringAttribute{Computed: true},
			"suites":        dsschema.StringAttribute{Computed: true},
			"tags":          dsschema.StringAttribute{Computed: true},
			"end_date":      dsschema.StringAttribute{Computed: true},
			"date_created":  dsschema.StringAttribute{Computed: true},
			"date_modified": dsschema.StringAttribute{Computed: true},
		},
	}
}

func (d *testPlanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *testPlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp testPlanAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/test-plans/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading test plan", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testPlanAPIToState(apiResp))...)
}
