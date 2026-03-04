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
	_ datasource.DataSource              = &maintenanceWindowDataSource{}
	_ datasource.DataSourceWithConfigure = &maintenanceWindowDataSource{}
	_ datasource.DataSource              = &maintenanceWindowsDataSource{}
	_ datasource.DataSourceWithConfigure = &maintenanceWindowsDataSource{}
)

func NewMaintenanceWindowDataSource() datasource.DataSource {
	return &maintenanceWindowDataSource{}
}

type maintenanceWindowDataSource struct {
	client *client.Client
}

func (d *maintenanceWindowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maintenance_window"
}

func (d *maintenanceWindowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix maintenance window by ID.",
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Required: true},
			"name":             schema.StringAttribute{Computed: true},
			"group_id":         schema.StringAttribute{Computed: true},
			"starts_at":        schema.StringAttribute{Computed: true},
			"ends_at":          schema.StringAttribute{Computed: true},
			"repeat":           schema.StringAttribute{Computed: true},
			"time_from":        schema.StringAttribute{Computed: true},
			"time_to":          schema.StringAttribute{Computed: true},
			"weekdays":         schema.StringAttribute{Computed: true},
			"day_of_month":     schema.Int64Attribute{Computed: true},
			"cron_expression":  schema.StringAttribute{Computed: true},
			"is_active":        schema.BoolAttribute{Computed: true},
			"schedule_summary": schema.StringAttribute{Computed: true},
			"date_created":     schema.StringAttribute{Computed: true},
			"date_modified":    schema.StringAttribute{Computed: true},
		},
	}
}

func (d *maintenanceWindowDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *maintenanceWindowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp maintenanceWindowAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/maintenance-windows/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading maintenance window", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, maintenanceWindowAPIToState(apiResp))...)
}

// --- List maintenance windows data source ---

func NewMaintenanceWindowsDataSource() datasource.DataSource {
	return &maintenanceWindowsDataSource{}
}

type maintenanceWindowsDataSource struct {
	client *client.Client
}

type maintenanceWindowsDataSourceModel struct {
	MaintenanceWindows []maintenanceWindowResourceModel `tfsdk:"maintenance_windows"`
}

func (d *maintenanceWindowsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maintenance_windows"
}

func (d *maintenanceWindowsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all maintenance windows in the organization.",
		Attributes: map[string]schema.Attribute{
			"maintenance_windows": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":               schema.StringAttribute{Computed: true},
						"name":             schema.StringAttribute{Computed: true},
						"group_id":         schema.StringAttribute{Computed: true},
						"starts_at":        schema.StringAttribute{Computed: true},
						"ends_at":          schema.StringAttribute{Computed: true},
						"repeat":           schema.StringAttribute{Computed: true},
						"time_from":        schema.StringAttribute{Computed: true},
						"time_to":          schema.StringAttribute{Computed: true},
						"weekdays":         schema.StringAttribute{Computed: true},
						"day_of_month":     schema.Int64Attribute{Computed: true},
						"cron_expression":  schema.StringAttribute{Computed: true},
						"is_active":        schema.BoolAttribute{Computed: true},
						"schedule_summary": schema.StringAttribute{Computed: true},
						"date_created":     schema.StringAttribute{Computed: true},
						"date_modified":    schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *maintenanceWindowsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *maintenanceWindowsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []maintenanceWindowAPIModel
	err := d.client.List(ctx, "/maintenance-windows/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing maintenance windows", err.Error())
		return
	}

	state := maintenanceWindowsDataSourceModel{}
	for _, item := range apiResp {
		state.MaintenanceWindows = append(state.MaintenanceWindows, maintenanceWindowAPIToState(item))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
