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
	_ resource.Resource              = &statusPageResource{}
	_ resource.ResourceWithConfigure = &statusPageResource{}
)

func NewStatusPageResource() resource.Resource {
	return &statusPageResource{}
}

type statusPageResource struct {
	client *client.Client
}

type statusPageResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Slug            types.String `tfsdk:"slug"`
	CustomDomain    types.String `tfsdk:"custom_domain"`
	Visibility      types.String `tfsdk:"visibility"`
	LogoURL         types.String `tfsdk:"logo_url"`
	FaviconURL      types.String `tfsdk:"favicon_url"`
	AccentColor     types.String `tfsdk:"accent_color"`
	CustomCSS       types.String `tfsdk:"custom_css"`
	HeaderText      types.String `tfsdk:"header_text"`
	FooterText      types.String `tfsdk:"footer_text"`
	MetaTitle       types.String `tfsdk:"meta_title"`
	MetaDescription types.String `tfsdk:"meta_description"`
	ShowHealthGraph types.Bool   `tfsdk:"show_health_graph"`
	IsActive        types.Bool   `tfsdk:"is_active"`
	DateCreated     types.String `tfsdk:"date_created"`
	DateModified    types.String `tfsdk:"date_modified"`
}

type statusPageAPIModel struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	CustomDomain    string `json:"custom_domain"`
	Visibility      string `json:"visibility"`
	LogoURL         string `json:"logo_url"`
	FaviconURL      string `json:"favicon_url"`
	AccentColor     string `json:"accent_color"`
	CustomCSS       string `json:"custom_css"`
	HeaderText      string `json:"header_text"`
	FooterText      string `json:"footer_text"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	ShowHealthGraph bool   `json:"show_health_graph"`
	IsActive        bool   `json:"is_active"`
	DateCreated     string `json:"date_created"`
	DateModified    string `json:"date_modified"`
}

func statusPageAPIToState(api statusPageAPIModel) statusPageResourceModel {
	return statusPageResourceModel{
		ID:              types.StringValue(api.ID),
		Name:            types.StringValue(api.Name),
		Slug:            types.StringValue(api.Slug),
		CustomDomain:    types.StringValue(api.CustomDomain),
		Visibility:      types.StringValue(api.Visibility),
		LogoURL:         types.StringValue(api.LogoURL),
		FaviconURL:      types.StringValue(api.FaviconURL),
		AccentColor:     types.StringValue(api.AccentColor),
		CustomCSS:       types.StringValue(api.CustomCSS),
		HeaderText:      types.StringValue(api.HeaderText),
		FooterText:      types.StringValue(api.FooterText),
		MetaTitle:       types.StringValue(api.MetaTitle),
		MetaDescription: types.StringValue(api.MetaDescription),
		ShowHealthGraph: types.BoolValue(api.ShowHealthGraph),
		IsActive:        types.BoolValue(api.IsActive),
		DateCreated:     types.StringValue(api.DateCreated),
		DateModified:    types.StringValue(api.DateModified),
	}
}

func (r *statusPageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page"
}

func (r *statusPageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Argonix status page.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"slug": schema.StringAttribute{
				Required:    true,
				Description: "URL slug for the status page (must be unique).",
			},
			"custom_domain": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"visibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("public"),
				Description: "Visibility: public or private.",
			},
			"logo_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"favicon_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"accent_color": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("#3B82F6"),
			},
			"custom_css": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"header_text": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"footer_text": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"meta_title": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"meta_description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"show_health_graph": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"is_active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *statusPageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *statusPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan statusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := statusPageStateToPayload(plan)
	var apiResp statusPageAPIModel
	err := r.client.Create(ctx, "/status-pages/", payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating status page", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, statusPageAPIToState(apiResp))...)
}

func (r *statusPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state statusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apiResp statusPageAPIModel
	err := r.client.Read(ctx, fmt.Sprintf("/status-pages/%s/", state.ID.ValueString()), &apiResp)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading status page", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, statusPageAPIToState(apiResp))...)
}

func (r *statusPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan statusPageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state statusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := statusPageStateToPayload(plan)
	var apiResp statusPageAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/status-pages/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating status page", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, statusPageAPIToState(apiResp))...)
}

func (r *statusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state statusPageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, fmt.Sprintf("/status-pages/%s/", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error deleting status page", err.Error())
	}
}

func statusPageStateToPayload(plan statusPageResourceModel) map[string]interface{} {
	return map[string]interface{}{
		"name":              plan.Name.ValueString(),
		"slug":              plan.Slug.ValueString(),
		"custom_domain":     plan.CustomDomain.ValueString(),
		"visibility":        plan.Visibility.ValueString(),
		"logo_url":          plan.LogoURL.ValueString(),
		"favicon_url":       plan.FaviconURL.ValueString(),
		"accent_color":      plan.AccentColor.ValueString(),
		"custom_css":        plan.CustomCSS.ValueString(),
		"header_text":       plan.HeaderText.ValueString(),
		"footer_text":       plan.FooterText.ValueString(),
		"meta_title":        plan.MetaTitle.ValueString(),
		"meta_description":  plan.MetaDescription.ValueString(),
		"show_health_graph": plan.ShowHealthGraph.ValueBool(),
		"is_active":         plan.IsActive.ValueBool(),
	}
}
