package provider
package provider

import (
	"context"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &workflowDataSource{}
	_ datasource.DataSourceWithConfigure = &workflowDataSource{}
	_ datasource.DataSource              = &workflowsDataSource{}
	_ datasource.DataSourceWithConfigure = &workflowsDataSource{}
)

func NewWorkflowDataSource() datasource.DataSource {
	return &workflowDataSource{}
}

type workflowDataSource struct {
	client *client.Client
}

func (d *workflowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow"
}

func (d *workflowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix workflow by ID.",
		Attributes: map[string]schema.Attribute{
			"id":                       schema.StringAttribute{Required: true},
			"name":                     schema.StringAttribute{Computed: true},
			"slug":                     schema.StringAttribute{Computed: true},
			"description":              schema.StringAttribute{Computed: true},
			"category":                 schema.StringAttribute{Computed: true},
			"steps":                    schema.StringAttribute{Computed: true},
			"required_connector_types": schema.StringAttribute{Computed: true},
			"requires_confirmation":    schema.BoolAttribute{Computed: true},
			"schedule":                 schema.StringAttribute{Computed: true},
			"is_active":                schema.BoolAttribute{Computed: true},
			"date_created":             schema.StringAttribute{Computed: true},
			"date_modified":            schema.StringAttribute{Computed: true},
		},
	}
}

func (d *workflowDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *workflowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp workflowAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/argos/workflows/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading workflow", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, workflowAPIToState(apiResp))...)
}

// --- List workflows data source ---

func NewWorkflowsDataSource() datasource.DataSource {
	return &workflowsDataSource{}
}

type workflowsDataSource struct {
	client *client.Client
}

type workflowsDataSourceModel struct {
	Workflows []workflowResourceModel `tfsdk:"workflows"`
}

func (d *workflowsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflows"
}

func (d *workflowsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all workflows in the organization.",
		Attributes: map[string]schema.Attribute{
			"workflows": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                       schema.StringAttribute{Computed: true},
						"name":                     schema.StringAttribute{Computed: true},
						"slug":                     schema.StringAttribute{Computed: true},
						"description":              schema.StringAttribute{Computed: true},
						"category":                 schema.StringAttribute{Computed: true},
						"steps":                    schema.StringAttribute{Computed: true},
						"required_connector_types": schema.StringAttribute{Computed: true},
						"requires_confirmation":    schema.BoolAttribute{Computed: true},
						"schedule":                 schema.StringAttribute{Computed: true},
						"is_active":                schema.BoolAttribute{Computed: true},
						"date_created":             schema.StringAttribute{Computed: true},
						"date_modified":            schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *workflowsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *workflowsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []workflowAPIModel
	err := d.client.List(ctx, "/argos/workflows/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing workflows", err.Error())
		return
	}

	state := workflowsDataSourceModel{}
	for _, item := range apiResp {
		state.Workflows = append(state.Workflows, workflowAPIToState(item))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
