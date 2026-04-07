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
	_ datasource.DataSource              = &securityPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &securityPolicyDataSource{}
	_ datasource.DataSource              = &securityPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &securityPoliciesDataSource{}
)

// --- Singular data source ---

func NewSecurityPolicyDataSource() datasource.DataSource {
	return &securityPolicyDataSource{}
}

type securityPolicyDataSource struct {
	client *client.Client
}

func (d *securityPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_policy"
}

func (d *securityPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix security policy by ID.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"name":          schema.StringAttribute{Computed: true},
			"description":   schema.StringAttribute{Computed: true},
			"rules":         schema.StringAttribute{Computed: true},
			"environment":   schema.StringAttribute{Computed: true},
			"is_active":     schema.BoolAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (d *securityPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *securityPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp securityPolicyAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/security/policies/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading security policy", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityPolicyAPIToState(apiResp))...)
}

// --- Plural data source ---

func NewSecurityPoliciesDataSource() datasource.DataSource {
	return &securityPoliciesDataSource{}
}

type securityPoliciesDataSource struct {
	client *client.Client
}

type securityPoliciesDataSourceModel struct {
	SecurityPolicies []securityPolicyResourceModel `tfsdk:"security_policies"`
}

func (d *securityPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_policies"
}

func (d *securityPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all security policies in the organization.",
		Attributes: map[string]schema.Attribute{
			"security_policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true},
						"name":          schema.StringAttribute{Computed: true},
						"description":   schema.StringAttribute{Computed: true},
						"rules":         schema.StringAttribute{Computed: true},
						"environment":   schema.StringAttribute{Computed: true},
						"is_active":     schema.BoolAttribute{Computed: true},
						"date_created":  schema.StringAttribute{Computed: true},
						"date_modified": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *securityPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *securityPoliciesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []securityPolicyAPIModel
	err := d.client.List(ctx, "/security/policies/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing security policies", err.Error())
		return
	}

	state := securityPoliciesDataSourceModel{}
	for _, r := range apiResp {
		state.SecurityPolicies = append(state.SecurityPolicies, securityPolicyAPIToState(r))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
