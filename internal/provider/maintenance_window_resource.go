package provider

import (
	"context"
	"fmt"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &maintenanceWindowResource{}
	_ resource.ResourceWithConfigure = &maintenanceWindowResource{}
)

func NewMaintenanceWindowResource() resource.Resource {
	return &maintenanceWindowResource{}
}

type maintenanceWindowResource struct {
	client *client.Client
}

type maintenanceWindowResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	GroupID         types.String `tfsdk:"group_id"`
	StartsAt        types.String `tfsdk:"starts_at"`
	EndsAt          types.String `tfsdk:"ends_at"`
	Repeat          types.String `tfsdk:"repeat"`
	TimeFrom        types.String `tfsdk:"time_from"`
	TimeTo          types.String `tfsdk:"time_to"`
	Weekdays        types.String `tfsdk:"weekdays"`
	DayOfMonth      types.Int64  `tfsdk:"day_of_month"`
	CronExpression  types.String `tfsdk:"cron_expression"`
	IsActive        types.Bool   `tfsdk:"is_active"`
	ScheduleSummary types.String `tfsdk:"schedule_summary"`
	DateCreated     types.String `tfsdk:"date_created"`
	DateModified    types.String `tfsdk:"date_modified"`
}

type maintenanceWindowAPIModel struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Group           *string `json:"group"`
	StartsAt        *string `json:"starts_at"`
	EndsAt          *string `json:"ends_at"`
	Repeat          string  `json:"repeat"`
	TimeFrom        *string `json:"time_from"`
	TimeTo          *string `json:"time_to"`
	Weekdays        *string `json:"weekdays"`
	DayOfMonth      *int64  `json:"day_of_month"`
	CronExpression  *string `json:"cron_expression"`
	IsActive        bool    `json:"is_active"`
	ScheduleSummary string  `json:"schedule_summary"`
	DateCreated     string  `json:"date_created"`
	DateModified    string  `json:"date_modified"`
}

func maintenanceWindowAPIToState(api maintenanceWindowAPIModel) maintenanceWindowResourceModel {
	state := maintenanceWindowResourceModel{
		ID:              types.StringValue(api.ID),
		Name:            types.StringValue(api.Name),
		Repeat:          types.StringValue(api.Repeat),
		IsActive:        types.BoolValue(api.IsActive),
		ScheduleSummary: types.StringValue(api.ScheduleSummary),
		DateCreated:     types.StringValue(api.DateCreated),
		DateModified:    types.StringValue(api.DateModified),
	}
	state.GroupID = optStrToTF(api.Group)
	state.StartsAt = optStrToTF(api.StartsAt)
	state.EndsAt = optStrToTF(api.EndsAt)
	state.TimeFrom = optStrToTF(api.TimeFrom)
	state.TimeTo = optStrToTF(api.TimeTo)
	state.Weekdays = optStrToTF(api.Weekdays)
	state.CronExpression = optStrToTF(api.CronExpression)
	if api.DayOfMonth != nil {
		state.DayOfMonth = types.Int64Value(*api.DayOfMonth)
	} else {
		state.DayOfMonth = types.Int64Value(0)
	}
	return state
}

func optStrToTF(s *string) types.String {
	if s != nil {
		return types.StringValue(*s)
	}
	return types.StringValue("")
}

func (r *maintenanceWindowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maintenance_window"
}

func (r *maintenanceWindowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix maintenance window.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the maintenance window.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "UUID of the group to apply the window to. Empty = all monitors.",
			},
			"starts_at": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Start datetime (ISO 8601) for one-time or range windows.",
			},
			"ends_at": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "End datetime (ISO 8601).",
			},
			"repeat": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("once"),
				Description: "Repeat mode: once, daily, weekly, monthly, cron.",
			},
			"time_from": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Daily start time (HH:MM) for recurring windows.",
			},
			"time_to": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Daily end time (HH:MM) for recurring windows.",
			},
			"weekdays": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Comma-separated weekday numbers (0=Mon, 6=Sun) for weekly repeat.",
			},
			"day_of_month": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Day of month (1-31) for monthly repeat.",
			},
			"cron_expression": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Cron expression for cron repeat mode.",
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"schedule_summary": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable schedule summary (read-only).",
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *maintenanceWindowResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *maintenanceWindowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan maintenanceWindowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := maintenanceWindowStateToPayload(plan)
	var apiResp maintenanceWindowAPIModel
	err := r.client.Create(ctx, "/maintenance-windows/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating maintenance window", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, maintenanceWindowAPIToState(apiResp))...)
}

func (r *maintenanceWindowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state maintenanceWindowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp maintenanceWindowAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/maintenance-windows/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading maintenance window", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, maintenanceWindowAPIToState(apiResp))...)
}

func (r *maintenanceWindowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan maintenanceWindowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state maintenanceWindowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := maintenanceWindowStateToPayload(plan)
	var apiResp maintenanceWindowAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/maintenance-windows/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating maintenance window", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, maintenanceWindowAPIToState(apiResp))...)
}

func (r *maintenanceWindowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state maintenanceWindowResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/maintenance-windows/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting maintenance window", err.Error())
	}
}

func maintenanceWindowStateToPayload(plan maintenanceWindowResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"name":      plan.Name.ValueString(),
		"repeat":    plan.Repeat.ValueString(),
		"is_active": plan.IsActive.ValueBool(),
	}
	if v := plan.GroupID.ValueString(); v != "" {
		payload["group"] = v
	}
	if v := plan.StartsAt.ValueString(); v != "" {
		payload["starts_at"] = v
	}
	if v := plan.EndsAt.ValueString(); v != "" {
		payload["ends_at"] = v
	}
	if v := plan.TimeFrom.ValueString(); v != "" {
		payload["time_from"] = v
	}
	if v := plan.TimeTo.ValueString(); v != "" {
		payload["time_to"] = v
	}
	if v := plan.Weekdays.ValueString(); v != "" {
		payload["weekdays"] = v
	}
	if !plan.DayOfMonth.IsNull() && !plan.DayOfMonth.IsUnknown() && plan.DayOfMonth.ValueInt64() > 0 {
		payload["day_of_month"] = plan.DayOfMonth.ValueInt64()
	}
	if v := plan.CronExpression.ValueString(); v != "" {
		payload["cron_expression"] = v
	}
	return payload
}
