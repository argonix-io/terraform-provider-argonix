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
	_ datasource.DataSource              = &statusPageDataSource{}
	_ datasource.DataSourceWithConfigure = &statusPageDataSource{}
	_ datasource.DataSource              = &statusPagesDataSource{}
	_ datasource.DataSourceWithConfigure = &statusPagesDataSource{}
)

func NewStatusPageDataSource() datasource.DataSource {
	return &statusPageDataSource{}
}

type statusPageDataSource struct {
	client *client.Client
}

func (d *statusPageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (d *statusPageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix status page by ID.",
		Attributes:  statusPageDataSourceAttributes(),
	}
}

func statusPageDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":               schema.StringAttribute{Required: true},
		"name":             schema.StringAttribute{Computed: true},
		"slug":             schema.StringAttribute{Computed: true},
		"custom_domain":    schema.StringAttribute{Computed: true},
		"visibility":       schema.StringAttribute{Computed: true},
		"logo_url":         schema.StringAttribute{Computed: true},
		"favicon_url":      schema.StringAttribute{Computed: true},
		"accent_color":     schema.StringAttribute{Computed: true},
		"custom_css":       schema.StringAttribute{Computed: true},
		"header_text":      schema.StringAttribute{Computed: true},
		"footer_text":      schema.StringAttribute{Computed: true},
		"meta_title":       schema.StringAttribute{Computed: true},
		"meta_description": schema.StringAttribute{Computed: true},
		"show_health_graph": schema.BoolAttribute{Computed: true},
		"is_active":        schema.BoolAttribute{Computed: true},
		"date_created":     schema.StringAttribute{Computed: true},
		"date_modified":    schema.StringAttribute{Computed: true},
	}
}

func (d *statusPageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *statusPageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp statusPageAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/status-pages/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading status page", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, statusPageAPIToState(apiResp))...)
}

// --- List status pages data source ---

func NewStatusPagesDataSource() datasource.DataSource {
	return &statusPagesDataSource{}
}

type statusPagesDataSource struct {
	client *client.Client
}

type statusPagesDataSourceModel struct {
	StatusPages []statusPageResourceModel `tfsdk:"status_pages"`
}

func (d *statusPagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_pages"
}

func (d *statusPagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all status pages in the organization.",
		Attributes: map[string]schema.Attribute{
			"status_pages": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: statusPageDataSourceAttributes(),
				},
			},
		},
	}
}

func (d *statusPagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *statusPagesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []statusPageAPIModel
	err := d.client.List(ctx, "/status-pages/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing status pages", err.Error())
		return
	}

	state := statusPagesDataSourceModel{}
	for _, sp := range apiResp {
		state.StatusPages = append(state.StatusPages, statusPageAPIToState(sp))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
