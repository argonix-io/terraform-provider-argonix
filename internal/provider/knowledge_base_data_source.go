package provider
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
	_ datasource.DataSource              = &knowledgeBaseDataSource{}
	_ datasource.DataSourceWithConfigure = &knowledgeBaseDataSource{}
	_ datasource.DataSource              = &knowledgeBasesDataSource{}
	_ datasource.DataSourceWithConfigure = &knowledgeBasesDataSource{}
)

func NewKnowledgeBaseDataSource() datasource.DataSource {
	return &knowledgeBaseDataSource{}
}

type knowledgeBaseDataSource struct {
	client *client.Client
}

func (d *knowledgeBaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge_base"
}

func (d *knowledgeBaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single Argonix knowledge base by ID.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Required: true},
			"name":           schema.StringAttribute{Computed: true},
			"source_type":    schema.StringAttribute{Computed: true},
			"connector_id":   schema.StringAttribute{Computed: true},
			"is_active":      schema.BoolAttribute{Computed: true},
			"sync_config":    schema.StringAttribute{Computed: true},
			"last_synced_at": schema.StringAttribute{Computed: true},
			"document_count": schema.Int64Attribute{Computed: true},
			"date_created":   schema.StringAttribute{Computed: true},
			"date_modified":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *knowledgeBaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *knowledgeBaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config struct {
		ID types.String `tfsdk:"id"`
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp knowledgeBaseAPIModel
	err := d.client.Read(ctx, fmt.Sprintf("/argos/knowledge-bases/%s/", config.ID.ValueString()), &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading knowledge base", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, knowledgeBaseAPIToState(apiResp))...)
}

// --- List knowledge bases data source ---

func NewKnowledgeBasesDataSource() datasource.DataSource {
	return &knowledgeBasesDataSource{}
}

type knowledgeBasesDataSource struct {
	client *client.Client
}

type knowledgeBasesDataSourceModel struct {
	KnowledgeBases []knowledgeBaseResourceModel `tfsdk:"knowledge_bases"`
}

func (d *knowledgeBasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge_bases"
}

func (d *knowledgeBasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all knowledge bases in the organization.",
		Attributes: map[string]schema.Attribute{
			"knowledge_bases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":             schema.StringAttribute{Computed: true},
						"name":           schema.StringAttribute{Computed: true},
						"source_type":    schema.StringAttribute{Computed: true},
						"connector_id":   schema.StringAttribute{Computed: true},
						"is_active":      schema.BoolAttribute{Computed: true},
						"sync_config":    schema.StringAttribute{Computed: true},
						"last_synced_at": schema.StringAttribute{Computed: true},
						"document_count": schema.Int64Attribute{Computed: true},
						"date_created":   schema.StringAttribute{Computed: true},
						"date_modified":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *knowledgeBasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *knowledgeBasesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiResp []knowledgeBaseAPIModel
	err := d.client.List(ctx, "/argos/knowledge-bases/", &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error listing knowledge bases", err.Error())
		return
	}

	state := knowledgeBasesDataSourceModel{}
	for _, item := range apiResp {
		state.KnowledgeBases = append(state.KnowledgeBases, knowledgeBaseAPIToState(item))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
