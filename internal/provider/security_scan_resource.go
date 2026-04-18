package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &securityScanResource{}
	_ resource.ResourceWithConfigure = &securityScanResource{}
)

func NewSecurityScanResource() resource.Resource {
	return &securityScanResource{}
}

type securityScanResource struct {
	client *client.Client
}

type securityScanResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ScanType      types.String `tfsdk:"scan_type"`
	Connector     types.String `tfsdk:"connector"`
	TriggeredBy   types.String `tfsdk:"triggered_by"`
	Status        types.String `tfsdk:"status"`
	FindingsCount types.Int64  `tfsdk:"findings_count"`
	Summary       types.String `tfsdk:"summary"`
	StartedAt     types.String `tfsdk:"started_at"`
	CompletedAt   types.String `tfsdk:"completed_at"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

type securityScanAPIModel struct {
	ID            string      `json:"id"`
	ScanType      string      `json:"scan_type"`
	Connector     string      `json:"connector"`
	TriggeredBy   string      `json:"triggered_by"`
	Status        string      `json:"status"`
	FindingsCount int64       `json:"findings_count"`
	Summary       interface{} `json:"summary"`
	StartedAt     string      `json:"started_at"`
	CompletedAt   string      `json:"completed_at"`
	DateCreated   string      `json:"date_created"`
	DateModified  string      `json:"date_modified"`
}

func securityScanAPIToState(api securityScanAPIModel) securityScanResourceModel {
	summaryJSON, _ := json.Marshal(api.Summary)

	return securityScanResourceModel{
		ID:            types.StringValue(api.ID),
		ScanType:      types.StringValue(api.ScanType),
		Connector:     types.StringValue(api.Connector),
		TriggeredBy:   types.StringValue(api.TriggeredBy),
		Status:        types.StringValue(api.Status),
		FindingsCount: types.Int64Value(api.FindingsCount),
		Summary:       types.StringValue(string(summaryJSON)),
		StartedAt:     types.StringValue(api.StartedAt),
		CompletedAt:   types.StringValue(api.CompletedAt),
		DateCreated:   types.StringValue(api.DateCreated),
		DateModified:  types.StringValue(api.DateModified),
	}
}

func (r *securityScanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_scan"
}

func (r *securityScanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Triggers an Argonix security scan (CSPM, vulnerability, or secret scan).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scan_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of scan: cspm, vulnerability, or secret.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "UUID of the connector to scan. Empty scans all connectors.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"triggered_by": schema.StringAttribute{
				Computed:    true,
				Description: "How the scan was triggered (manual or scheduled).",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Scan status: pending, running, completed, or failed.",
			},
			"findings_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of findings from the scan.",
			},
			"summary": schema.StringAttribute{
				Computed:    true,
				Description: "JSON-encoded scan summary.",
			},
			"started_at":    schema.StringAttribute{Computed: true},
			"completed_at":  schema.StringAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *securityScanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *securityScanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan securityScanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"scan_type": plan.ScanType.ValueString(),
	}
	if plan.Connector.ValueString() != "" {
		payload["connector"] = plan.Connector.ValueString()
	}

	var apiResp securityScanAPIModel
	err := r.client.Create(ctx, "/security/scans/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating security scan", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityScanAPIToState(apiResp))...)
}

func (r *securityScanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state securityScanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp securityScanAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/security/scans/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading security scan", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, securityScanAPIToState(apiResp))...)
}

func (r *securityScanResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Security scans are immutable — any change forces replacement.
	resp.Diagnostics.AddError("Update not supported", "Security scans are immutable. Changes require replacement.")
}

func (r *securityScanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state securityScanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/security/scans/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting security scan", err.Error())
	}
}
