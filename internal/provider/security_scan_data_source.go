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
	_ datasource.DataSource              = &securityScanDataSource{}
	_ datasource.DataSourceWithConfigure = &securityScanDataSource{}
	_ datasource.DataSource              = &securityScansDataSource{}
	_ datasource.DataSourceWithConfigure = &securityScansDataSource{}
)

// --- Singular data source ---

func NewSecurityScanDataSource() datasource.DataSource {
	return &securityScanDataSource{}
}

type securityScanDataSource struct {
	client *client.Client
}

func (d *securityScanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_scan"
}

func (d *securityScanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix security scan by ID.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Required: true},
			"scan_type":       schema.StringAttribute{Computed: true},
			"connector":       schema.StringAttribute{Computed: true},
			"triggered_by":    schema.StringAttribute{Computed: true},
			"status":          schema.StringAttribute{Computed: true},
			"findings_count":  schema.Int64Attribute{Computed: true},
			"summary":         schema.StringAttribute{Computed: true},
			"started_at":      schema.StringAttribute{Computed: true},
			"completed_at":    schema.StringAttribute{Computed: true},
			"date_created":    schema.StringAttribute{Computed: true},
			"date_modified":   schema.StringAttribute{Computed: true},
		},
	}
}

func (d *securityScanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *securityScanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp securityScanAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/security/scans/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading security scan", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityScanAPIToState(apiResp))...)
}

// --- Plural data source ---

func NewSecurityScansDataSource() datasource.DataSource {
	return &securityScansDataSource{}
}

type securityScansDataSource struct {
	client *client.Client
}

type securityScansDataSourceModel struct {
	SecurityScans []securityScanResourceModel `tfsdk:"security_scans"`
}

func (d *securityScansDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_scans"
}

func (d *securityScansDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all security scans in the organization.",
		Attributes: map[string]schema.Attribute{
			"security_scans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":              schema.StringAttribute{Computed: true},
						"scan_type":       schema.StringAttribute{Computed: true},
						"connector":       schema.StringAttribute{Computed: true},
						"triggered_by":    schema.StringAttribute{Computed: true},
						"status":          schema.StringAttribute{Computed: true},
						"findings_count":  schema.Int64Attribute{Computed: true},
						"summary":         schema.StringAttribute{Computed: true},
						"started_at":      schema.StringAttribute{Computed: true},
						"completed_at":    schema.StringAttribute{Computed: true},
						"date_created":    schema.StringAttribute{Computed: true},
						"date_modified":   schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *securityScansDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *securityScansDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []securityScanAPIModel
	err := d.client.List(ctx, "/security/scans/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing security scans", err.Error())
		return
	}

	state := securityScansDataSourceModel{}
	for _, r := range apiResp {
		state.SecurityScans = append(state.SecurityScans, securityScanAPIToState(r))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
