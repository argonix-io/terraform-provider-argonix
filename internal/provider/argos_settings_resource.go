package provider

import (
	"context"
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
	_ resource.Resource              = &argosSettingsResource{}
	_ resource.ResourceWithConfigure = &argosSettingsResource{}
)

func NewArgosSettingsResource() resource.Resource {
	return &argosSettingsResource{}
}

type argosSettingsResource struct {
	client *client.Client
}

type argosSettingsResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	LLMProvider        types.String `tfsdk:"llm_provider"`
	LLMModel           types.String `tfsdk:"llm_model"`
	LLMApiKey          types.String `tfsdk:"llm_api_key"`
	LLMBaseURL         types.String `tfsdk:"llm_base_url"`
	CustomInstructions types.String `tfsdk:"custom_instructions"`
	DemoMode           types.Bool   `tfsdk:"demo_mode"`
	DateModified       types.String `tfsdk:"date_modified"`
}

type argosSettingsAPIModel struct {
	ID                 string `json:"id"`
	LLMProvider        string `json:"llm_provider"`
	LLMModel           string `json:"llm_model"`
	LLMApiKeySet       bool   `json:"llm_api_key_set"`
	LLMBaseURL         string `json:"llm_base_url"`
	CustomInstructions string `json:"custom_instructions"`
	DemoMode           bool   `json:"demo_mode"`
	DateModified       string `json:"date_modified"`
}

func argosSettingsAPIToState(api argosSettingsAPIModel, currentApiKey types.String) argosSettingsResourceModel {
	return argosSettingsResourceModel{
		ID:                 types.StringValue(api.ID),
		LLMProvider:        types.StringValue(api.LLMProvider),
		LLMModel:           types.StringValue(api.LLMModel),
		LLMApiKey:          currentApiKey, // preserve — API never returns the key
		LLMBaseURL:         types.StringValue(api.LLMBaseURL),
		CustomInstructions: types.StringValue(api.CustomInstructions),
		DemoMode:           types.BoolValue(api.DemoMode),
		DateModified:       types.StringValue(api.DateModified),
	}
}

func (r *argosSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_argos_settings"
}

func (r *argosSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Argos AI agent settings for the organization. This is a singleton resource — each organization has exactly one.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"llm_provider": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("google"),
				Description: "LLM provider: `local`, `google`, `anthropic`, or `openai`.",
			},
			"llm_model": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Model override (blank = default for the provider). Any model name is accepted when `llm_base_url` is set.",
			},
			"llm_api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Custom API key — bring your own key for the LLM provider. Leave empty to use the platform default.",
			},
			"llm_base_url": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Custom base URL for the LLM endpoint (e.g. vLLM, Azure OpenAI, self-hosted Ollama).",
			},
			"custom_instructions": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Custom system prompt instructions prepended to every Argos conversation.",
			},
			"demo_mode": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "When enabled, Argos returns scripted demo responses instead of calling the LLM.",
			},
			"date_modified": schema.StringAttribute{Computed: true},
		},
	}
}

func (r *argosSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// settingsEndpoint is the singleton list endpoint; the API auto-creates when accessed.
const settingsEndpoint = "/argos/settings/"

func (r *argosSettingsResource) readSettings(ctx context.Context) (*argosSettingsAPIModel, error) {
	var apiResp argosSettingsAPIModel
	// The API returns the singleton directly on GET /argos/settings/
	if err := r.client.Read(ctx, settingsEndpoint, &apiResp); err != nil {
		return nil, err
	}
	return &apiResp, nil
}

func (r *argosSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan argosSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the auto-created singleton first to get its ID
	existing, err := r.readSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading argos settings", err.Error())
		return
	}

	// PATCH with desired state
	payload := argosSettingsStateToPayload(plan)
	var apiResp argosSettingsAPIModel
	err = r.client.Update(ctx, fmt.Sprintf("/argos/settings/%s/", existing.ID), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating argos settings", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, argosSettingsAPIToState(apiResp, plan.LLMApiKey))...)
}

func (r *argosSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state argosSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.readSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading argos settings", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, argosSettingsAPIToState(*apiResp, state.LLMApiKey))...)
}

func (r *argosSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan argosSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var state argosSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := argosSettingsStateToPayload(plan)
	var apiResp argosSettingsAPIModel
	err := r.client.Update(ctx, fmt.Sprintf("/argos/settings/%s/", state.ID.ValueString()), payload, &apiResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating argos settings", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, argosSettingsAPIToState(apiResp, plan.LLMApiKey))...)
}

func (r *argosSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Singleton — reset to defaults instead of deleting
	var state argosSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"llm_provider":        "google",
		"llm_model":           "",
		"llm_api_key":         "",
		"llm_base_url":        "",
		"custom_instructions": "",
		"demo_mode":           false,
	}
	var apiResp argosSettingsAPIModel
	_ = r.client.Update(ctx, fmt.Sprintf("/argos/settings/%s/", state.ID.ValueString()), payload, &apiResp)
}

func argosSettingsStateToPayload(plan argosSettingsResourceModel) map[string]interface{} {
	payload := map[string]interface{}{
		"llm_provider":        plan.LLMProvider.ValueString(),
		"llm_model":           plan.LLMModel.ValueString(),
		"llm_base_url":        plan.LLMBaseURL.ValueString(),
		"custom_instructions": plan.CustomInstructions.ValueString(),
		"demo_mode":           plan.DemoMode.ValueBool(),
	}
	// Only send api_key when explicitly set (sensitive, write-only)
	if !plan.LLMApiKey.IsNull() && !plan.LLMApiKey.IsUnknown() && plan.LLMApiKey.ValueString() != "" {
		payload["llm_api_key"] = plan.LLMApiKey.ValueString()
	}
	return payload
}
