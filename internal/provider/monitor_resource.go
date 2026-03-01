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
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// Type
	MonitorType types.String `tfsdk:"monitor_type"`
	IsActive    types.Bool   `tfsdk:"is_active"`
	// Target
	URL      types.String `tfsdk:"url"`
	Hostname types.String `tfsdk:"hostname"`
	Port     types.Int64  `tfsdk:"port"`
	// DNS
	DNSRecordType types.String `tfsdk:"dns_record_type"`
	DNSExpected   types.String `tfsdk:"dns_expected"`
	// HTTP options
	HTTPMethod          types.String `tfsdk:"http_method"`
	HTTPHeaders         types.String `tfsdk:"http_headers"`
	HTTPBody            types.String `tfsdk:"http_body"`
	HTTPBodyContentType types.String `tfsdk:"http_body_content_type"`
	FollowRedirects     types.Bool   `tfsdk:"follow_redirects"`
	VerifySSL           types.Bool   `tfsdk:"verify_ssl"`
	HTTPAuthUser        types.String `tfsdk:"http_auth_user"`
	HTTPAuthPass        types.String `tfsdk:"http_auth_pass"`
	// Keyword
	Keyword       types.String `tfsdk:"keyword"`
	KeywordExists types.Bool   `tfsdk:"keyword_exists"`
	// Scheduling
	CheckInterval types.Int64 `tfsdk:"check_interval"`
	Timeout       types.Int64 `tfsdk:"timeout"`
	Retries       types.Int64 `tfsdk:"retries"`
	// Remediation
	RemediationEnabled     types.Bool   `tfsdk:"remediation_enabled"`
	RemediationScript      types.String `tfsdk:"remediation_script"`
	RemediationTimeout     types.Int64  `tfsdk:"remediation_timeout"`
	RemediationWaitSeconds types.Int64  `tfsdk:"remediation_wait_seconds"`
	// Heartbeat
	HeartbeatToken        types.String `tfsdk:"heartbeat_token"`
	HeartbeatGraceSeconds types.Int64  `tfsdk:"heartbeat_grace_seconds"`
	// Multi-step HTTP
	MultiStepConfig types.String `tfsdk:"multi_step_config"`
	// gRPC
	GRPCService  types.String `tfsdk:"grpc_service"`
	GRPCMethod   types.String `tfsdk:"grpc_method"`
	GRPCProto    types.String `tfsdk:"grpc_proto"`
	GRPCMetadata types.String `tfsdk:"grpc_metadata"`
	GRPCTLS      types.Bool   `tfsdk:"grpc_tls"`
	// Assertions
	Assertions types.String `tfsdk:"assertions"`
	// SSL
	SSLExpiryWarnDays types.Int64 `tfsdk:"ssl_expiry_warn_days"`
	// Location / Regions / Tags
	Location types.String `tfsdk:"location"`
	Regions  types.String `tfsdk:"regions"`
	Tags     types.String `tfsdk:"tags"`
	GroupID  types.String `tfsdk:"group_id"`
	// State (read-only)
	CurrentStatus types.String `tfsdk:"current_status"`
	DateCreated   types.String `tfsdk:"date_created"`
	DateModified  types.String `tfsdk:"date_modified"`
}

// monitorAPIModel is the JSON shape returned by the Argonix API.
type monitorAPIModel struct {
	ID                     string      `json:"id"`
	Name                   string      `json:"name"`
	MonitorType            string      `json:"monitor_type"`
	IsActive               bool        `json:"is_active"`
	URL                    string      `json:"url"`
	Hostname               string      `json:"hostname"`
	Port                   *int64      `json:"port"`
	DNSRecordType          string      `json:"dns_record_type"`
	DNSExpected            string      `json:"dns_expected"`
	HTTPMethod             string      `json:"http_method"`
	HTTPHeaders            interface{} `json:"http_headers"`
	HTTPBody               string      `json:"http_body"`
	HTTPBodyContentType    string      `json:"http_body_content_type"`
	FollowRedirects        bool        `json:"follow_redirects"`
	VerifySSL              bool        `json:"verify_ssl"`
	HTTPAuthUser           string      `json:"http_auth_user"`
	HTTPAuthPass           string      `json:"http_auth_pass"`
	Keyword                string      `json:"keyword"`
	KeywordExists          bool        `json:"keyword_exists"`
	CheckInterval          int64       `json:"check_interval"`
	Timeout                int64       `json:"timeout"`
	Retries                int64       `json:"retries"`
	RemediationEnabled     bool        `json:"remediation_enabled"`
	RemediationScript      string      `json:"remediation_script"`
	RemediationTimeout     int64       `json:"remediation_timeout"`
	RemediationWaitSeconds int64       `json:"remediation_wait_seconds"`
	HeartbeatToken         string      `json:"heartbeat_token"`
	HeartbeatGraceSeconds  int64       `json:"heartbeat_grace_seconds"`
	MultiStepConfig        interface{} `json:"multi_step_config"`
	GRPCService            string      `json:"grpc_service"`
	GRPCMethod             string      `json:"grpc_method"`
	GRPCProto              string      `json:"grpc_proto"`
	GRPCMetadata           interface{} `json:"grpc_metadata"`
	GRPCTLS                bool        `json:"grpc_tls"`
	Assertions             interface{} `json:"assertions"`
	SSLExpiryWarnDays      int64       `json:"ssl_expiry_warn_days"`
	Location               string      `json:"location"`
	Regions                interface{} `json:"regions"`
	Tags                   interface{} `json:"tags"`
	GroupID                *string     `json:"group"`
	CurrentStatus          string      `json:"current_status"`
	DateCreated            string      `json:"date_created"`
	DateModified           string      `json:"date_modified"`
}

func marshalJSONField(v interface{}) types.String {
	b, _ := json.Marshal(v)
	return types.StringValue(string(b))
}

func monitorAPIToState(api monitorAPIModel) monitorResourceModel {
	state := monitorResourceModel{
		ID:                     types.StringValue(api.ID),
		Name:                   types.StringValue(api.Name),
		MonitorType:            types.StringValue(api.MonitorType),
		IsActive:               types.BoolValue(api.IsActive),
		URL:                    types.StringValue(api.URL),
		Hostname:               types.StringValue(api.Hostname),
		DNSRecordType:          types.StringValue(api.DNSRecordType),
		DNSExpected:            types.StringValue(api.DNSExpected),
		HTTPMethod:             types.StringValue(api.HTTPMethod),
		HTTPHeaders:            marshalJSONField(api.HTTPHeaders),
		HTTPBody:               types.StringValue(api.HTTPBody),
		HTTPBodyContentType:    types.StringValue(api.HTTPBodyContentType),
		FollowRedirects:        types.BoolValue(api.FollowRedirects),
		VerifySSL:              types.BoolValue(api.VerifySSL),
		HTTPAuthUser:           types.StringValue(api.HTTPAuthUser),
		HTTPAuthPass:           types.StringValue(api.HTTPAuthPass),
		Keyword:                types.StringValue(api.Keyword),
		KeywordExists:          types.BoolValue(api.KeywordExists),
		CheckInterval:          types.Int64Value(api.CheckInterval),
		Timeout:                types.Int64Value(api.Timeout),
		Retries:                types.Int64Value(api.Retries),
		RemediationEnabled:     types.BoolValue(api.RemediationEnabled),
		RemediationScript:      types.StringValue(api.RemediationScript),
		RemediationTimeout:     types.Int64Value(api.RemediationTimeout),
		RemediationWaitSeconds: types.Int64Value(api.RemediationWaitSeconds),
		HeartbeatToken:         types.StringValue(api.HeartbeatToken),
		HeartbeatGraceSeconds:  types.Int64Value(api.HeartbeatGraceSeconds),
		MultiStepConfig:        marshalJSONField(api.MultiStepConfig),
		GRPCService:            types.StringValue(api.GRPCService),
		GRPCMethod:             types.StringValue(api.GRPCMethod),
		GRPCProto:              types.StringValue(api.GRPCProto),
		GRPCMetadata:           marshalJSONField(api.GRPCMetadata),
		GRPCTLS:                types.BoolValue(api.GRPCTLS),
		Assertions:             marshalJSONField(api.Assertions),
		SSLExpiryWarnDays:      types.Int64Value(api.SSLExpiryWarnDays),
		Location:               types.StringValue(api.Location),
		Regions:                marshalJSONField(api.Regions),
		Tags:                   marshalJSONField(api.Tags),
		CurrentStatus:          types.StringValue(api.CurrentStatus),
		DateCreated:            types.StringValue(api.DateCreated),
		DateModified:           types.StringValue(api.DateModified),
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

	return state
}

func (r *monitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *monitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix monitor.",
		Attributes: map[string]schema.Attribute{
			// Identity
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
			// Target
			"url": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "URL to monitor (for http, keyword, ssl monitors).",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Hostname to monitor (for ping, tcp, dns monitors).",
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Port to monitor (for tcp monitors).",
			},
			// DNS
			"dns_record_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("A"),
				Description: "DNS record type to query: A, AAAA, CNAME, MX, NS, TXT, etc.",
			},
			"dns_expected": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Expected DNS response value.",
			},
			// HTTP options
			"http_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("GET"),
				Description: "HTTP method: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS.",
			},
			"http_headers": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("{}"),
				Description: "JSON-encoded object of custom HTTP headers.",
			},
			"http_body": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "HTTP request body (for POST/PUT/PATCH).",
			},
			"http_body_content_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("application/json"),
				Description: "Content-Type header for the HTTP body.",
			},
			"follow_redirects": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to follow HTTP redirects.",
			},
			"verify_ssl": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to verify SSL certificates.",
			},
			"http_auth_user": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "HTTP Basic Auth username.",
			},
			"http_auth_pass": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Default:     stringdefault.StaticString(""),
				Description: "HTTP Basic Auth password.",
			},
			// Keyword
			"keyword": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Keyword to search for in the response body (for keyword monitors).",
			},
			"keyword_exists": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "If true, alert when keyword is NOT found. If false, alert when keyword IS found.",
			},
			// Scheduling
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
			// Remediation
			"remediation_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable automatic remediation when the monitor goes down.",
			},
			"remediation_script": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Shell script to execute when remediation is triggered.",
			},
			"remediation_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(60),
				Description: "Maximum execution time for the remediation script in seconds.",
			},
			"remediation_wait_seconds": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(30),
				Description: "Seconds to wait after remediation before rechecking.",
			},
			// Heartbeat
			"heartbeat_token": schema.StringAttribute{
				Computed:    true,
				Description: "Auto-generated token for heartbeat push URL (read-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"heartbeat_grace_seconds": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Grace period in seconds before a missed heartbeat triggers an alert.",
			},
			// Multi-step HTTP
			"multi_step_config": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded multi-step HTTP configuration.",
			},
			// gRPC
			"grpc_service": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "gRPC service name.",
			},
			"grpc_method": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "gRPC method to call.",
			},
			"grpc_proto": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Protobuf definition for the gRPC service.",
			},
			"grpc_metadata": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("{}"),
				Description: "JSON-encoded gRPC metadata key-value pairs.",
			},
			"grpc_tls": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to use TLS for gRPC connections.",
			},
			// Assertions
			"assertions": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of assertion objects to validate responses.",
			},
			// SSL
			"ssl_expiry_warn_days": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(30),
				Description: "Days before SSL expiry to trigger a warning.",
			},
			// Location / Regions / Tags / Group
			"location": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("eu-france"),
				Description: "Primary check location.",
			},
			"regions": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of region codes to run checks from.",
			},
			"tags": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("[]"),
				Description: "JSON-encoded list of tags.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "UUID of the group this monitor belongs to.",
			},
			// State (read-only)
			"current_status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status: up, down, degraded, maintenance, unknown.",
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
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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
	resp.Diagnostics.Append(resp.State.Set(ctx, monitorAPIToState(apiResp))...)
}

func (r *monitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state monitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
	resp.Diagnostics.Append(resp.State.Set(ctx, monitorAPIToState(apiResp))...)
}

func (r *monitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan monitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state monitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
	resp.Diagnostics.Append(resp.State.Set(ctx, monitorAPIToState(apiResp))...)
}

func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state monitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Delete(ctx, fmt.Sprintf("/monitors/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting monitor", err.Error())
	}
}

func unmarshalJSONPayload(s types.String) interface{} {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	var v interface{}
	json.Unmarshal([]byte(s.ValueString()), &v)
	return v
}

func monitorStateToPayload(plan monitorResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":                    plan.Name.ValueString(),
		"monitor_type":            plan.MonitorType.ValueString(),
		"is_active":               plan.IsActive.ValueBool(),
		"url":                     plan.URL.ValueString(),
		"hostname":                plan.Hostname.ValueString(),
		"dns_record_type":         plan.DNSRecordType.ValueString(),
		"dns_expected":            plan.DNSExpected.ValueString(),
		"http_method":             plan.HTTPMethod.ValueString(),
		"http_body":               plan.HTTPBody.ValueString(),
		"http_body_content_type":  plan.HTTPBodyContentType.ValueString(),
		"follow_redirects":        plan.FollowRedirects.ValueBool(),
		"verify_ssl":              plan.VerifySSL.ValueBool(),
		"http_auth_user":          plan.HTTPAuthUser.ValueString(),
		"http_auth_pass":          plan.HTTPAuthPass.ValueString(),
		"keyword":                 plan.Keyword.ValueString(),
		"keyword_exists":          plan.KeywordExists.ValueBool(),
		"check_interval":          plan.CheckInterval.ValueInt64(),
		"timeout":                 plan.Timeout.ValueInt64(),
		"retries":                 plan.Retries.ValueInt64(),
		"remediation_enabled":     plan.RemediationEnabled.ValueBool(),
		"remediation_script":      plan.RemediationScript.ValueString(),
		"remediation_timeout":     plan.RemediationTimeout.ValueInt64(),
		"remediation_wait_seconds": plan.RemediationWaitSeconds.ValueInt64(),
		"heartbeat_grace_seconds": plan.HeartbeatGraceSeconds.ValueInt64(),
		"grpc_service":            plan.GRPCService.ValueString(),
		"grpc_method":             plan.GRPCMethod.ValueString(),
		"grpc_proto":              plan.GRPCProto.ValueString(),
		"grpc_tls":                plan.GRPCTLS.ValueBool(),
		"ssl_expiry_warn_days":    plan.SSLExpiryWarnDays.ValueInt64(),
		"location":                plan.Location.ValueString(),
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		payload["port"] = plan.Port.ValueInt64()
	}
	if !plan.GroupID.IsNull() && !plan.GroupID.IsUnknown() {
		payload["group"] = plan.GroupID.ValueString()
	}

	payload["http_headers"] = unmarshalJSONPayload(plan.HTTPHeaders)
	payload["multi_step_config"] = unmarshalJSONPayload(plan.MultiStepConfig)
	payload["grpc_metadata"] = unmarshalJSONPayload(plan.GRPCMetadata)
	payload["assertions"] = unmarshalJSONPayload(plan.Assertions)
	payload["regions"] = unmarshalJSONPayload(plan.Regions)
	payload["tags"] = unmarshalJSONPayload(plan.Tags)

	return payload
}
