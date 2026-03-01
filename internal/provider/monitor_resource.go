package provider

import (
	"context"
	"encoding/json"
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
	_ resource.Resource              = &monitorResource{}
	_ resource.ResourceWithConfigure = &monitorResource{}
)

func NewMonitorResource() resource.Resource {
	return &monitorResource{}
}

type monitorResource struct {
	client *client.Client
}

// monitorResourceModel is the Terraform state model for a monitor.
type monitorResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	MonitorType   types.String `tfsdk:"monitor_type"`
	IsActive      types.Bool   `tfsdk:"is_active"`
	URL           types.String `tfsdk:"url"`
	Hostname      types.String `tfsdk:"hostname"`
	Port          types.Int64  `tfsdk:"port"`
	CheckInterval types.Int64  `tfsdk:"check_interval"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	Retries       types.Int64  `tfsdk:"retries"`
	HTTPMethod    types.String `tfsdk:"http_method"`
	CurrentStatus types.String `tfsdk:"current_status"`
	GroupID       types.String `tfsdk:"group_id"`
	Tags          types.String `tfsdk:"tags"`
	Regions       types.String `tfsdk:"regions"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

// monitorAPIModel is the JSON shape returned by the Argonix API.
type monitorAPIModel struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	MonitorType   string      `json:"monitor_type"`
	IsActive      bool        `json:"is_active"`
	URL           string      `json:"url"`
	Hostname      string      `json:"hostname"`
	Port          *int64      `json:"port"`
	CheckInterval int64       `json:"check_interval"`
	Timeout       int64       `json:"timeout"`
	Retries       int64       `json:"retries"`
	HTTPMethod    string      `json:"http_method"`
	CurrentStatus string      `json:"current_status"`
	GroupID       *string     `json:"group"`
	Tags          interface{} `json:"tags"`
	Regions       interface{} `json:"regions"`
	DateCreated   string      `json:"date_created"`
	DateModified  string      `json:"date_modified"`
}

func monitorAPIToState(api monitorAPIModel) monitorResourceModel {
	state := monitorResourceModel{
		ID:            types.StringValue(api.ID),
		Name:          types.StringValue(api.Name),
		MonitorType:   types.StringValue(api.MonitorType),
		IsActive:      types.BoolValue(api.IsActive),
		URL:           types.StringValue(api.URL),
		Hostname:      types.StringValue(api.Hostname),
		CheckInterval: types.Int64Value(api.CheckInterval),
		Timeout:       types.Int64Value(api.Timeout),
		Retries:       types.Int64Value(api.Retries),
		HTTPMethod:    types.StringValue(api.HTTPMethod),
		CurrentStatus: types.StringValue(api.CurrentStatus),
		DateCreated:   types.StringValue(api.DateCreated),
		DateModified:  types.StringValue(api.DateModified),
	}

	if api.Port != nil {
		state.Port = types.Int64Value(*api.Port)
	} else {
		state.Port = types.Int64Null()
	}

	if api.GroupID != nil {
		state.GroupID = types.StringValue(*api.GroupID)
	} else {
		state.GroupID = types.StringNull()
	}

	tagsJSON, _ := json.Marshal(api.Tags)
	state.Tags = types.StringValue(string(tagsJSON))
	regionsJSON, _ := json.Marshal(api.Regions)
	state.Regions = types.StringValue(string(regionsJSON))

	return state
}

func (r *monitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *monitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "UUID of the monitor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the monitor.",
			},
			"monitor_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of monitor: http, ping, tcp, dns, ssl, keyword, grpc, heartbeat, multi_step_http.",
			},
			"is_active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the monitor is active.",
			},
			"url": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "URL to monitor (for HTTP/keyword/SSL monitors).",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Hostname to monitor (for ping/TCP/DNS monitors).",
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Port to monitor (for TCP monitors).",
			},
			"check_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(300),
				Description: "Check interval in seconds.",
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(30),
				Description: "Request timeout in seconds.",
			},
			"retries": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Number of retries before marking as down.",
			},
			"http_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("GET"),
				Description: "HTTP method (GET, POST, PUT, etc.).",
			},
			"current_status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status: up, down, degraded, maintenance, unknown.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "UUID of the group this monitor belongs to.",
			},
			"tags": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of tags.",
			},
			"regions": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of region codes.",
			},
			"date_created": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp.",
			},
			"date_modified": schema.StringAttribute{
				Computed:    true,
				Description: "Last modification timestamp.",
			},
		},
	}
}

func (r *monitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *monitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan monitorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := monitorStateToPayload(plan)
	var apiResp monitorAPIModel
	err := r.client.Create(ctx, "/monitors/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating monitor", err.Error())
		return
	}

	state := monitorAPIToState(apiResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *monitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state monitorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp monitorAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/monitors/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading monitor", err.Error())
		return
	}

	newState := monitorAPIToState(apiResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *monitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan monitorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state monitorResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := monitorStateToPayload(plan)
	var apiResp monitorAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/monitors/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating monitor", err.Error())
		return
	}

	newState := monitorAPIToState(apiResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state monitorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/monitors/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting monitor", err.Error())
		return
	}
}

func monitorStateToPayload(plan monitorResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":           plan.Name.ValueString(),
		"monitor_type":   plan.MonitorType.ValueString(),
		"is_active":      plan.IsActive.ValueBool(),
		"url":            plan.URL.ValueString(),
		"hostname":       plan.Hostname.ValueString(),
		"check_interval": plan.CheckInterval.ValueInt64(),
		"timeout":        plan.Timeout.ValueInt64(),
		"retries":        plan.Retries.ValueInt64(),
		"http_method":    plan.HTTPMethod.ValueString(),
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		payload["port"] = plan.Port.ValueInt64()
	}

	if !plan.GroupID.IsNull() && !plan.GroupID.IsUnknown() {
		payload["group"] = plan.GroupID.ValueString()
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags interface{}
		json.Unmarshal([]byte(plan.Tags.ValueString()), &tags)
		payload["tags"] = tags
	}

	if !plan.Regions.IsNull() && !plan.Regions.IsUnknown() {
		var regions interface{}
		json.Unmarshal([]byte(plan.Regions.ValueString()), &regions)
		payload["regions"] = regions
	}

	return payload
}
