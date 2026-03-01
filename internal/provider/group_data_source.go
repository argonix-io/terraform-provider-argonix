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
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
	_ datasource.DataSource              = &groupsDataSource{}
	_ datasource.DataSourceWithConfigure = &groupsDataSource{}
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	client *client.Client
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix group by ID.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"name":          schema.StringAttribute{Computed: true},
			"description":   schema.StringAttribute{Computed: true},
			"tags":          schema.StringAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp groupAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/groups/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading group", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, groupAPIToState(apiResp))...)
}

// --- List groups data source ---

func NewGroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

type groupsDataSource struct {
	client *client.Client
}

type groupsDataSourceModel struct {
	Groups []groupResourceModel `tfsdk:"groups"`
}

func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all groups in the organization.",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"name":          schema.StringAttribute{Computed: true},
						"description":   schema.StringAttribute{Computed: true},
						"tags":          schema.StringAttribute{Computed: true},
						"date_created":  schema.StringAttribute{Computed: true},
						"date_modified": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *groupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []groupAPIModel
	err := d.client.List(ctx, "/groups/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing groups", err.Error())
		return
	}

	state := groupsDataSourceModel{}
	for _, g := range apiResp {
		state.Groups = append(state.Groups, groupAPIToState(g))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
