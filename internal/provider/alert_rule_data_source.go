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
	_ datasource.DataSource              = &alertRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &alertRuleDataSource{}
	_ datasource.DataSource              = &alertRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &alertRulesDataSource{}
)

func NewAlertRuleDataSource() datasource.DataSource {
	return &alertRuleDataSource{}
}

type alertRuleDataSource struct {
	client *client.Client
}

func (d *alertRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_rule"
}

func (d *alertRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix alert rule by ID.",
		Attributes: map[string]schema.Attribute{
			"id":                   schema.StringAttribute{Required: true},
			"name":                 schema.StringAttribute{Computed: true},
			"is_active":            schema.BoolAttribute{Computed: true},
			"trigger_condition":    schema.StringAttribute{Computed: true},
			"consecutive_failures": schema.Int64Attribute{Computed: true},
			"cooldown_minutes":     schema.Int64Attribute{Computed: true},
			"all_monitors":         schema.BoolAttribute{Computed: true},
			"all_synthetic_tests":  schema.BoolAttribute{Computed: true},
			"monitor_tags":         schema.StringAttribute{Computed: true},
			"monitors":             schema.StringAttribute{Computed: true},
			"synthetic_tests":      schema.StringAttribute{Computed: true},
			"channels":             schema.StringAttribute{Computed: true},
			"date_created":         schema.StringAttribute{Computed: true},
			"date_modified":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *alertRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp alertRuleAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/alert-rules/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading alert rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, alertRuleAPIToState(apiResp))...)
}

// --- List alert rules data source ---

func NewAlertRulesDataSource() datasource.DataSource {
	return &alertRulesDataSource{}
}

type alertRulesDataSource struct {
	client *client.Client
}

type alertRulesDataSourceModel struct {
	AlertRules []alertRuleResourceModel `tfsdk:"alert_rules"`
}

func (d *alertRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_rules"
}

func (d *alertRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all alert rules in the organization.",
		Attributes: map[string]schema.Attribute{
			"alert_rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                   schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"is_active":            schema.BoolAttribute{Computed: true},
						"trigger_condition":    schema.StringAttribute{Computed: true},
						"consecutive_failures": schema.Int64Attribute{Computed: true},
						"cooldown_minutes":     schema.Int64Attribute{Computed: true},
						"all_monitors":         schema.BoolAttribute{Computed: true},
						"all_synthetic_tests":  schema.BoolAttribute{Computed: true},
						"monitor_tags":         schema.StringAttribute{Computed: true},
						"monitors":             schema.StringAttribute{Computed: true},
						"synthetic_tests":      schema.StringAttribute{Computed: true},
						"channels":             schema.StringAttribute{Computed: true},
						"date_created":         schema.StringAttribute{Computed: true},
						"date_modified":        schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *alertRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alertRulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []alertRuleAPIModel
	err := d.client.List(ctx, "/alert-rules/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing alert rules", err.Error())
		return
	}

	state := alertRulesDataSourceModel{}
	for _, r := range apiResp {
		state.AlertRules = append(state.AlertRules, alertRuleAPIToState(r))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
