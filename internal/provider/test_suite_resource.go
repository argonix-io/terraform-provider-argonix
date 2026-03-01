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
// Test Suite Resource
// ===========================================================================

var (
	_ resource.Resource              = &testSuiteResource{}
	_ resource.ResourceWithConfigure = &testSuiteResource{}
)

func NewTestSuiteResource() resource.Resource {
	return &testSuiteResource{}
}

type testSuiteResource struct {
	client *client.Client
}

type testSuiteResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Tags            types.String `tfsdk:"tags"`
	SyntheticTests  types.String `tfsdk:"synthetic_tests"`
	ManualTestCases types.String `tfsdk:"manual_test_cases"`
	DateCreated     types.String `tfsdk:"date_created"`
	DateModified    types.String `tfsdk:"date_modified"`
}

type testSuiteAPIModel struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	Tags            interface{} `json:"tags"`
	SyntheticTests  interface{} `json:"synthetic_tests"`
	ManualTestCases interface{} `json:"manual_test_cases"`
	DateCreated     string      `json:"date_created"`
	DateModified    string      `json:"date_modified"`
}

func testSuiteAPIToState(api testSuiteAPIModel) testSuiteResourceModel {
	tagsJSON, _ := json.Marshal(api.Tags)
	stJSON, _ := json.Marshal(api.SyntheticTests)
	mtcJSON, _ := json.Marshal(api.ManualTestCases)

	state := testSuiteResourceModel{
		ID:              types.StringValue(api.ID),
		Name:            types.StringValue(api.Name),
		Description:     types.StringValue(api.Description),
		Tags:            types.StringValue(string(tagsJSON)),
		SyntheticTests:  types.StringValue(string(stJSON)),
		ManualTestCases: types.StringValue(string(mtcJSON)),
		DateCreated:     types.StringValue(api.DateCreated),
		DateModified:    types.StringValue(api.DateModified),
	}
	return state
}

func (r *testSuiteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test_suite"
}

func (r *testSuiteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix test suite.",
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
			"tags": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("[]"),
			},
			"synthetic_tests": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of synthetic test UUIDs.",
			},
			"manual_test_cases": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of manual test case UUIDs.",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *testSuiteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *testSuiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan testSuiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := testSuiteStateToPayload(plan)
	var apiResp testSuiteAPIModel
	err := r.client.Create(ctx, "/test-suites/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating test suite", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testSuiteAPIToState(apiResp))...)
}

func (r *testSuiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state testSuiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp testSuiteAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/test-suites/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading test suite", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testSuiteAPIToState(apiResp))...)
}

func (r *testSuiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan testSuiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state testSuiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := testSuiteStateToPayload(plan)
	var apiResp testSuiteAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/test-suites/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating test suite", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testSuiteAPIToState(apiResp))...)
}

func (r *testSuiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state testSuiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/test-suites/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting test suite", err.Error())
	}
}

func testSuiteStateToPayload(plan testSuiteResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	}

	unmarshalJSONFieldIfSet(plan.Tags, "tags", payload)
	unmarshalJSONFieldIfSet(plan.SyntheticTests, "synthetic_tests", payload)
	unmarshalJSONFieldIfSet(plan.ManualTestCases, "manual_test_cases", payload)

	return payload
}

func unmarshalJSONFieldIfSet(field types.String, key string, payload map[string]interface{}) {
	if !field.IsNull() && !field.IsUnknown() {
		var v interface{}
		if err := json.Unmarshal([]byte(field.ValueString()), &v); err == nil {
			payload[key] = v
		}
	}
}

// ===========================================================================
// Test Suite Data Source
// ===========================================================================

var (
	_ datasource.DataSource              = &testSuiteDataSource{}
	_ datasource.DataSourceWithConfigure = &testSuiteDataSource{}
)

func NewTestSuiteDataSource() datasource.DataSource {
	return &testSuiteDataSource{}
}

type testSuiteDataSource struct {
	client *client.Client
}

func (d *testSuiteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_test_suite"
}

func (d *testSuiteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Fetches a single Argonix test suite by ID.",
		Attributes: map[string]dsschema.Attribute{
			"id":                dsschema.StringAttribute{Required: true},
			"name":              dsschema.StringAttribute{Computed: true},
			"description":       dsschema.StringAttribute{Computed: true},
			"tags":              dsschema.StringAttribute{Computed: true},
			"synthetic_tests":   dsschema.StringAttribute{Computed: true},
			"manual_test_cases": dsschema.StringAttribute{Computed: true},
			"date_created":      dsschema.StringAttribute{Computed: true},
			"date_modified":     dsschema.StringAttribute{Computed: true},
		},
	}
}

func (d *testSuiteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *testSuiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp testSuiteAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/test-suites/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading test suite", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, testSuiteAPIToState(apiResp))...)
}
