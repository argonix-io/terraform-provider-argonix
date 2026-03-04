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
	_ datasource.DataSource              = &notificationRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationRuleDataSource{}
	_ datasource.DataSource              = &notificationRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationRulesDataSource{}
)

func NewNotificationRuleDataSource() datasource.DataSource {
	return &notificationRuleDataSource{}
}

type notificationRuleDataSource struct {
	client *client.Client
}

func (d *notificationRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rule"
}

func (d *notificationRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix notification rule by ID.",
		Attributes: map[string]schema.Attribute{
			"id":                   schema.StringAttribute{Required: true},
			"name":                 schema.StringAttribute{Computed: true},
			"is_active":            schema.BoolAttribute{Computed: true},
			"trigger_condition":    schema.StringAttribute{Computed: true},
			"consecutive_failures": schema.Int64Attribute{Computed: true},
			"cooldown_minutes":     schema.Int64Attribute{Computed: true},
			"all_monitors":         schema.BoolAttribute{Computed: true},
			"all_synthetic_tests":  schema.BoolAttribute{Computed: true},
			"auto_investigate":     schema.BoolAttribute{Computed: true},
			"auto_remediate":       schema.BoolAttribute{Computed: true},
			"remediation_strategy": schema.StringAttribute{Computed: true},
			"monitor_tags":         schema.StringAttribute{Computed: true},
			"monitors":             schema.StringAttribute{Computed: true},
			"synthetic_tests":      schema.StringAttribute{Computed: true},
			"channels":             schema.StringAttribute{Computed: true},
			"date_created":         schema.StringAttribute{Computed: true},
			"date_modified":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *notificationRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *notificationRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp notificationRuleAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/notification-rules/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading notification rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, notificationRuleAPIToState(apiResp))...)
}

// --- List notification rules data source ---

func NewNotificationRulesDataSource() datasource.DataSource {
	return &notificationRulesDataSource{}
}

type notificationRulesDataSource struct {
	client *client.Client
}

type notificationRulesDataSourceModel struct {
	NotificationRules []notificationRuleResourceModel `tfsdk:"notification_rules"`
}

func (d *notificationRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_rules"
}

func (d *notificationRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all notification rules in the organization.",
		Attributes: map[string]schema.Attribute{
			"notification_rules": schema.ListNestedAttribute{
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
						"auto_investigate":     schema.BoolAttribute{Computed: true},
						"auto_remediate":       schema.BoolAttribute{Computed: true},
						"remediation_strategy": schema.StringAttribute{Computed: true},
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

func (d *notificationRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *notificationRulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []notificationRuleAPIModel
	err := d.client.List(ctx, "/notification-rules/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing notification rules", err.Error())
		return
	}

	state := notificationRulesDataSourceModel{}
	for _, r := range apiResp {
		state.NotificationRules = append(state.NotificationRules, notificationRuleAPIToState(r))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
