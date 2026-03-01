package provider

import (
	"context"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// ── Single monitor data source ──────────────────────────────────────────

var (
	_ datasource.DataSource              = &monitorDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorDataSource{}
)

func NewMonitorDataSource() datasource.DataSource {
	return &monitorDataSource{}
}

type monitorDataSource struct {
	client *client.Client
}

func (d *monitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func monitorDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":                       schema.StringAttribute{Required: true, Description: "UUID of the monitor."},
		"name":                     schema.StringAttribute{Computed: true, Description: "Display name."},
		"monitor_type":             schema.StringAttribute{Computed: true, Description: "Type: http, ping, tcp, dns, ssl, keyword, grpc, heartbeat, multi_step_http."},
		"is_active":                schema.BoolAttribute{Computed: true, Description: "Whether the monitor is active."},
		"url":                      schema.StringAttribute{Computed: true, Description: "URL to monitor."},
		"hostname":                 schema.StringAttribute{Computed: true, Description: "Hostname to monitor."},
		"port":                     schema.Int64Attribute{Computed: true, Description: "Port to monitor."},
		"dns_record_type":          schema.StringAttribute{Computed: true, Description: "DNS record type."},
		"dns_expected":             schema.StringAttribute{Computed: true, Description: "Expected DNS response."},
		"http_method":              schema.StringAttribute{Computed: true, Description: "HTTP method."},
		"http_headers":             schema.StringAttribute{Computed: true, Description: "JSON-encoded HTTP headers."},
		"http_body":                schema.StringAttribute{Computed: true, Description: "HTTP request body."},
		"http_body_content_type":   schema.StringAttribute{Computed: true, Description: "Content-Type for HTTP body."},
		"follow_redirects":         schema.BoolAttribute{Computed: true, Description: "Whether to follow redirects."},
		"verify_ssl":               schema.BoolAttribute{Computed: true, Description: "Whether to verify SSL."},
		"http_auth_user":           schema.StringAttribute{Computed: true, Description: "HTTP Basic Auth username."},
		"http_auth_pass":           schema.StringAttribute{Computed: true, Sensitive: true, Description: "HTTP Basic Auth password."},
		"keyword":                  schema.StringAttribute{Computed: true, Description: "Keyword to search for."},
		"keyword_exists":           schema.BoolAttribute{Computed: true, Description: "Alert when keyword is missing (true) or found (false)."},
		"check_interval":           schema.Int64Attribute{Computed: true, Description: "Check interval in seconds."},
		"timeout":                  schema.Int64Attribute{Computed: true, Description: "Request timeout in seconds."},
		"retries":                  schema.Int64Attribute{Computed: true, Description: "Retries before marking down."},
		"remediation_enabled":      schema.BoolAttribute{Computed: true, Description: "Automatic remediation enabled."},
		"remediation_script":       schema.StringAttribute{Computed: true, Description: "Shell script for remediation."},
		"remediation_timeout":      schema.Int64Attribute{Computed: true, Description: "Remediation script timeout in seconds."},
		"remediation_wait_seconds": schema.Int64Attribute{Computed: true, Description: "Seconds to wait after remediation."},
		"heartbeat_token":          schema.StringAttribute{Computed: true, Description: "Auto-generated heartbeat token."},
		"heartbeat_grace_seconds":  schema.Int64Attribute{Computed: true, Description: "Grace period for heartbeat."},
		"multi_step_config":        schema.StringAttribute{Computed: true, Description: "JSON-encoded multi-step config."},
		"grpc_service":             schema.StringAttribute{Computed: true, Description: "gRPC service name."},
		"grpc_method":              schema.StringAttribute{Computed: true, Description: "gRPC method."},
		"grpc_proto":               schema.StringAttribute{Computed: true, Description: "gRPC protobuf definition."},
		"grpc_metadata":            schema.StringAttribute{Computed: true, Description: "JSON-encoded gRPC metadata."},
		"grpc_tls":                 schema.BoolAttribute{Computed: true, Description: "Use TLS for gRPC."},
		"assertions":               schema.StringAttribute{Computed: true, Description: "JSON-encoded assertions."},
		"ssl_expiry_warn_days":     schema.Int64Attribute{Computed: true, Description: "Days before SSL expiry warning."},
		"location":                 schema.StringAttribute{Computed: true, Description: "Primary check location."},
		"regions":                  schema.StringAttribute{Computed: true, Description: "JSON-encoded region list."},
		"tags":                     schema.StringAttribute{Computed: true, Description: "JSON-encoded tags."},
		"group_id":                 schema.StringAttribute{Computed: true, Description: "UUID of the group."},
		"current_status":           schema.StringAttribute{Computed: true, Description: "Current status."},
		"date_created":             schema.StringAttribute{Computed: true, Description: "Creation timestamp."},
		"date_modified":            schema.StringAttribute{Computed: true, Description: "Last modification timestamp."},
	}
}

func (d *monitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an Argonix monitor by ID.",
		Attributes:  monitorDataSourceAttributes(),
	}
}

func (d *monitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *monitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg monitorResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp monitorAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/monitors/%s/", cfg.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading monitor", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, monitorAPIToState(apiResp))...)
}

// ── List monitors data source ───────────────────────────────────────────

var (
	_ datasource.DataSource              = &monitorsDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorsDataSource{}
)

func NewMonitorsDataSource() datasource.DataSource {
	return &monitorsDataSource{}
}

type monitorsDataSource struct {
	client *client.Client
}

type monitorsDataSourceModel struct {
	Monitors []monitorResourceModel `tfsdk:"monitors"`
}

func (d *monitorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitors"
}

func (d *monitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all Argonix monitors for the organization.",
		Attributes: map[string]schema.Attribute{
			"monitors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: monitorDataSourceAttributes(),
				},
			},
		},
	}
}

func (d *monitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *monitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []monitorAPIModel
	err := d.client.List(ctx, "/monitors/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing monitors", err.Error())
		return
	}

	state := monitorsDataSourceModel{}
	for _, m := range apiResp {
		state.Monitors = append(state.Monitors, monitorAPIToState(m))
	}
	if state.Monitors == nil {
		state.Monitors = []monitorResourceModel{}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// In the single data source, override "id" to be Required instead of Computed.
func monitorDataSourceSingleAttributes() map[string]schema.Attribute {
	attrs := monitorDataSourceAttributes()
	attrs["id"] = schema.StringAttribute{Required: true, Description: "UUID of the monitor."}
	return attrs
}
