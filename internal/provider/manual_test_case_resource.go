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
// Manual Test Case Resource
// ===========================================================================

var (
	_ resource.Resource              = &manualTestCaseResource{}
	_ resource.ResourceWithConfigure = &manualTestCaseResource{}
)

func NewManualTestCaseResource() resource.Resource {
	return &manualTestCaseResource{}
}

type manualTestCaseResource struct {
	client *client.Client
}

type manualTestCaseResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Title         types.String `tfsdk:"title"`
	Description   types.String `tfsdk:"description"`
	Preconditions types.String `tfsdk:"preconditions"`
	Steps         types.String `tfsdk:"steps"`
	Priority      types.String `tfsdk:"priority"`
	Tags          types.String `tfsdk:"tags"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

type manualTestCaseAPIModel struct {
	ID            string      `json:"id"`
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Preconditions string      `json:"preconditions"`
	Steps         interface{} `json:"steps"`
	Priority      string      `json:"priority"`
	Tags          interface{} `json:"tags"`
	DateCreated   string      `json:"date_created"`
	DateModified  string      `json:"date_modified"`
}

func manualTestCaseAPIToState(api manualTestCaseAPIModel) manualTestCaseResourceModel {
	stepsJSON, _ := json.Marshal(api.Steps)
	tagsJSON, _ := json.Marshal(api.Tags)
	return manualTestCaseResourceModel{
		ID:            types.StringValue(api.ID),
		Title:         types.StringValue(api.Title),
		Description:   types.StringValue(api.Description),
		Preconditions: types.StringValue(api.Preconditions),
		Steps:         types.StringValue(string(stepsJSON)),
		Priority:      types.StringValue(api.Priority),
		Tags:          types.StringValue(string(tagsJSON)),
		DateCreated:   types.StringValue(api.DateCreated),
		DateModified:  types.StringValue(api.DateModified),
	}
}

func (r *manualTestCaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manual_test_case"
}

func (r *manualTestCaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix manual test case.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"preconditions": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"steps": schema.StringAttribute{
				Required:    true,
				Description: `JSON-encoded ordered list: [{"description": "...", "expected": "..."}, ...]`,
			},
			"priority": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("medium"),
				Description: "Priority: critical, high, medium, low.",
			},
			"tags": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("[]"),
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *manualTestCaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *manualTestCaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan manualTestCaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := manualTestCaseStateToPayload(plan)
	var apiResp manualTestCaseAPIModel
	err := r.client.Create(ctx, "/manual-test-cases/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating manual test case", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, manualTestCaseAPIToState(apiResp))...)
}

func (r *manualTestCaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state manualTestCaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp manualTestCaseAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/manual-test-cases/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading manual test case", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, manualTestCaseAPIToState(apiResp))...)
}

func (r *manualTestCaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan manualTestCaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state manualTestCaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := manualTestCaseStateToPayload(plan)
	var apiResp manualTestCaseAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/manual-test-cases/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating manual test case", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, manualTestCaseAPIToState(apiResp))...)
}

func (r *manualTestCaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state manualTestCaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/manual-test-cases/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting manual test case", err.Error())
	}
}

func manualTestCaseStateToPayload(plan manualTestCaseResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"title":         plan.Title.ValueString(),
		"description":   plan.Description.ValueString(),
		"preconditions": plan.Preconditions.ValueString(),
		"priority":      plan.Priority.ValueString(),
	}

	var steps interface{}
	json.Unmarshal([]byte(plan.Steps.ValueString()), &steps)
	payload["steps"] = steps

	unmarshalJSONFieldIfSet(plan.Tags, "tags", payload)

	return payload
}

// ===========================================================================
// Manual Test Case Data Source
// ===========================================================================

var (
	_ datasource.DataSource              = &manualTestCaseDataSource{}
	_ datasource.DataSourceWithConfigure = &manualTestCaseDataSource{}
)

func NewManualTestCaseDataSource() datasource.DataSource {
	return &manualTestCaseDataSource{}
}

type manualTestCaseDataSource struct {
	client *client.Client
}

func (d *manualTestCaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manual_test_case"
}

func (d *manualTestCaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Fetches a single Argonix manual test case by ID.",
		Attributes: map[string]dsschema.Attribute{
			"id":            dsschema.StringAttribute{Required: true},
			"title":         dsschema.StringAttribute{Computed: true},
			"description":   dsschema.StringAttribute{Computed: true},
			"preconditions": dsschema.StringAttribute{Computed: true},
			"steps":         dsschema.StringAttribute{Computed: true},
			"priority":      dsschema.StringAttribute{Computed: true},
			"tags":          dsschema.StringAttribute{Computed: true},
			"date_created":  dsschema.StringAttribute{Computed: true},
			"date_modified": dsschema.StringAttribute{Computed: true},
		},
	}
}

func (d *manualTestCaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *manualTestCaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp manualTestCaseAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/manual-test-cases/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading manual test case", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, manualTestCaseAPIToState(apiResp))...)
}
