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
	_ datasource.DataSource              = &personaDataSource{}
	_ datasource.DataSourceWithConfigure = &personaDataSource{}
	_ datasource.DataSource              = &personasDataSource{}
	_ datasource.DataSourceWithConfigure = &personasDataSource{}
)

func NewPersonaDataSource() datasource.DataSource {
	return &personaDataSource{}
}

type personaDataSource struct {
	client *client.Client
}

func (d *personaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_persona"
}

func (d *personaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix Argos persona by ID.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true},
			"name":           schema.StringAttribute{Computed: true},
			"description":    schema.StringAttribute{Computed: true},
			"template":       schema.StringAttribute{Computed: true},
			"system_prompt":  schema.StringAttribute{Computed: true},
			"is_active":      schema.BoolAttribute{Computed: true},
			"connector_ids":  schema.StringAttribute{Computed: true},
			"allowed_tools":  schema.StringAttribute{Computed: true},
			"approval_rules": schema.StringAttribute{Computed: true},
			"date_created":   schema.StringAttribute{Computed: true},
			"date_modified":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *personaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *personaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp personaAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/argos/personas/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading persona", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, personaAPIToState(apiResp))...)
}

// --- List personas data source ---

func NewPersonasDataSource() datasource.DataSource {
	return &personasDataSource{}
}

type personasDataSource struct {
	client *client.Client
}

type personasDataSourceModel struct {
	Personas []personaResourceModel `tfsdk:"personas"`
}

func (d *personasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personas"
}

func (d *personasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all Argos personas in the organization.",
		Attributes: map[string]schema.Attribute{
			"personas": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":             schema.StringAttribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"description":    schema.StringAttribute{Computed: true},
						"template":       schema.StringAttribute{Computed: true},
						"system_prompt":  schema.StringAttribute{Computed: true},
						"is_active":      schema.BoolAttribute{Computed: true},
						"connector_ids":  schema.StringAttribute{Computed: true},
						"allowed_tools":  schema.StringAttribute{Computed: true},
						"approval_rules": schema.StringAttribute{Computed: true},
						"date_created":   schema.StringAttribute{Computed: true},
						"date_modified":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *personasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *personasDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []personaAPIModel
	err := d.client.List(ctx, "/argos/personas/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing personas", err.Error())
		return
	}

	state := personasDataSourceModel{}
	for _, item := range apiResp {
		state.Personas = append(state.Personas, personaAPIToState(item))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
