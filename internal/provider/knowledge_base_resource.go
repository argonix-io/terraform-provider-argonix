package provider

import (
	"context"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &knowledgeBaseResource{}
	_ resource.ResourceWithConfigure = &knowledgeBaseResource{}
)

func NewKnowledgeBaseResource() resource.Resource {
	return &knowledgeBaseResource{}
}

type knowledgeBaseResource struct {
	client *client.Client
}

type knowledgeBaseResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	SourceType    types.String `tfsdk:"source_type"`
	ConnectorID   types.String `tfsdk:"connector_id"`
	IsActive      types.Bool   `tfsdk:"is_active"`
	SyncConfig    types.String `tfsdk:"sync_config"`
	LastSyncedAt  types.String `tfsdk:"last_synced_at"`
	DocumentCount types.Int64  `tfsdk:"document_count"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

type knowledgeBaseAPIModel struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	SourceType    string  `json:"source_type"`
	Connector     *string `json:"connector"`
	IsActive      bool    `json:"is_active"`
	SyncConfig    *string `json:"sync_config"`
	LastSyncedAt  *string `json:"last_synced_at"`
	DocumentCount int64   `json:"document_count"`
	DateCreated   string  `json:"date_created"`
	DateModified  string  `json:"date_modified"`
}

func knowledgeBaseAPIToState(api knowledgeBaseAPIModel) knowledgeBaseResourceModel {
	state := knowledgeBaseResourceModel{
		ID:            types.StringValue(api.ID),
		Name:          types.StringValue(api.Name),
		SourceType:    types.StringValue(api.SourceType),
		IsActive:      types.BoolValue(api.IsActive),
		DocumentCount: types.Int64Value(api.DocumentCount),
		DateCreated:   types.StringValue(api.DateCreated),
		DateModified:  types.StringValue(api.DateModified),
	}
	if api.Connector != nil {
		state.ConnectorID = types.StringValue(*api.Connector)
	} else {
		state.ConnectorID = types.StringValue("")
	}
	if api.SyncConfig != nil {
		state.SyncConfig = types.StringValue(*api.SyncConfig)
	} else {
		state.SyncConfig = types.StringValue("")
	}
	if api.LastSyncedAt != nil {
		state.LastSyncedAt = types.StringValue(*api.LastSyncedAt)
	} else {
		state.LastSyncedAt = types.StringValue("")
	}
	return state
}

func (r *knowledgeBaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_knowledge_base"
}

func (r *knowledgeBaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix Argos AI knowledge base.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the knowledge base.",
			},
			"source_type": schema.StringAttribute{
				Required:    true,
				Description: "Source type: manual, confluence, notion, gitlab_wiki, github_wiki, jira, slack, google_drive, sharepoint.",
			},
			"connector_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "UUID of the connector to use for syncing.",
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"sync_config": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "JSON-encoded sync configuration.",
			},
			"last_synced_at": schema.StringAttribute{Computed: true},
			"document_count": schema.Int64Attribute{
				Computed: true,
				Default:  int64default.StaticInt64(0),
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *knowledgeBaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *client.Client")
		return
	}
	r.client = c
}

func (r *knowledgeBaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan knowledgeBaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := knowledgeBaseStateToPayload(plan)
	var apiResp knowledgeBaseAPIModel
	err := r.client.Create(ctx, "/argos/knowledge-bases/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating knowledge base", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, knowledgeBaseAPIToState(apiResp))...)
}

func (r *knowledgeBaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state knowledgeBaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp knowledgeBaseAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/argos/knowledge-bases/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading knowledge base", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, knowledgeBaseAPIToState(apiResp))...)
}

func (r *knowledgeBaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan knowledgeBaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state knowledgeBaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := knowledgeBaseStateToPayload(plan)
	var apiResp knowledgeBaseAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/argos/knowledge-bases/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating knowledge base", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, knowledgeBaseAPIToState(apiResp))...)
}

func (r *knowledgeBaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state knowledgeBaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/argos/knowledge-bases/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting knowledge base", err.Error())
	}
}

func knowledgeBaseStateToPayload(plan knowledgeBaseResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":        plan.Name.ValueString(),
		"source_type": plan.SourceType.ValueString(),
		"is_active":   plan.IsActive.ValueBool(),
	}
	if v := plan.ConnectorID.ValueString(); v != "" {
		payload["connector"] = v
	}
	if v := plan.SyncConfig.ValueString(); v != "" {
		payload["sync_config"] = v
	}
	return payload
}
