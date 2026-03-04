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
	_ datasource.DataSource              = &environmentDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentDataSource{}
	_ datasource.DataSource              = &environmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentsDataSource{}
)

func NewEnvironmentDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

type environmentDataSource struct {
	client *client.Client
}

func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (d *environmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix environment by ID.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"name":          schema.StringAttribute{Computed: true},
			"variables":     schema.StringAttribute{Computed: true, Sensitive: true},
			"is_default":    schema.BoolAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp environmentAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/environments/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, environmentAPIToState(apiResp))...)
}

// --- List environments data source ---

func NewEnvironmentsDataSource() datasource.DataSource {
	return &environmentsDataSource{}
}

type environmentsDataSource struct {
	client *client.Client
}

type environmentsDataSourceModel struct {
	Environments []environmentResourceModel `tfsdk:"environments"`
}

func (d *environmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

func (d *environmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all environments in the organization.",
		Attributes: map[string]schema.Attribute{
			"environments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"name":          schema.StringAttribute{Computed: true},
						"variables":     schema.StringAttribute{Computed: true, Sensitive: true},
						"is_default":    schema.BoolAttribute{Computed: true},
						"date_created":  schema.StringAttribute{Computed: true},
						"date_modified": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *environmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *environmentsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []environmentAPIModel
	err := d.client.List(ctx, "/environments/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing environments", err.Error())
		return
	}

	state := environmentsDataSourceModel{}
	for _, item := range apiResp {
		state.Environments = append(state.Environments, environmentAPIToState(item))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
