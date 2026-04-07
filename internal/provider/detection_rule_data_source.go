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
	_ datasource.DataSource              = &detectionRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &detectionRuleDataSource{}
	_ datasource.DataSource              = &detectionRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &detectionRulesDataSource{}
)

// --- Singular data source ---

func NewDetectionRuleDataSource() datasource.DataSource {
	return &detectionRuleDataSource{}
}

type detectionRuleDataSource struct {
	client *client.Client
}

func (d *detectionRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detection_rule"
}

func (d *detectionRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix detection rule by ID.",
		Attributes: map[string]schema.Attribute{
			"id":               schema.StringAttribute{Required: true},
			"name":             schema.StringAttribute{Computed: true},
			"description":      schema.StringAttribute{Computed: true},
			"rule_type":        schema.StringAttribute{Computed: true},
			"severity":         schema.StringAttribute{Computed: true},
			"is_active":        schema.BoolAttribute{Computed: true},
			"config":           schema.StringAttribute{Computed: true},
			"mitre_tactic":     schema.StringAttribute{Computed: true},
			"mitre_technique":  schema.StringAttribute{Computed: true},
			"detections_count": schema.Int64Attribute{Computed: true},
			"date_created":     schema.StringAttribute{Computed: true},
			"date_modified":    schema.StringAttribute{Computed: true},
		},
	}
}

func (d *detectionRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *detectionRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp detectionRuleAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/security/detection-rules/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading detection rule", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, detectionRuleAPIToState(apiResp))...)
}

// --- Plural data source ---

func NewDetectionRulesDataSource() datasource.DataSource {
	return &detectionRulesDataSource{}
}

type detectionRulesDataSource struct {
	client *client.Client
}

type detectionRulesDataSourceModel struct {
	DetectionRules []detectionRuleResourceModel `tfsdk:"detection_rules"`
}

func (d *detectionRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detection_rules"
}

func (d *detectionRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all detection rules in the organization.",
		Attributes: map[string]schema.Attribute{
			"detection_rules": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":               schema.StringAttribute{Computed: true},
						"name":             schema.StringAttribute{Computed: true},
						"description":      schema.StringAttribute{Computed: true},
						"rule_type":        schema.StringAttribute{Computed: true},
						"severity":         schema.StringAttribute{Computed: true},
						"is_active":        schema.BoolAttribute{Computed: true},
						"config":           schema.StringAttribute{Computed: true},
						"mitre_tactic":     schema.StringAttribute{Computed: true},
						"mitre_technique":  schema.StringAttribute{Computed: true},
						"detections_count": schema.Int64Attribute{Computed: true},
						"date_created":     schema.StringAttribute{Computed: true},
						"date_modified":    schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *detectionRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *detectionRulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []detectionRuleAPIModel
	err := d.client.List(ctx, "/security/detection-rules/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing detection rules", err.Error())
		return
	}

	state := detectionRulesDataSourceModel{}
	for _, r := range apiResp {
		state.DetectionRules = append(state.DetectionRules, detectionRuleAPIToState(r))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
