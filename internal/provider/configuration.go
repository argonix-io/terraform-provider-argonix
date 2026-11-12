package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math/big"
	"net/url"
	"reflect"
	"slices"

	"github.com/argonix-io/terraform-provider-argonix/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type fieldKind int

const (
	stringField fieldKind = iota
	boolField
	intField
	jsonObjectField
	jsonArrayField
	idSetField
	floatField
	currencyAmountField
)

type configField struct {
	name, description string
	kind              fieldKind
	required          bool
	readOnly          bool
	replace           bool
	responseName      string
	defaultValue      any
	allowed           []string
	jsonEquivalent    func(actual, expected any) bool
}

type configuration struct {
	name, plural, endpoint, description string
	fields                              []configField
	detailForList                       bool
	readAfterUpdate                     bool
}

func (c configuration) detail(id string) string {
	return c.endpoint + url.PathEscape(id) + "/"
}

func (c configuration) allFields() []configField {
	return append([]configField{output("id", stringField)}, c.fields...)
}

func (f configField) attrType() attr.Type {
	switch f.kind {
	case boolField:
		return types.BoolType
	case intField:
		return types.Int64Type
	case floatField, currencyAmountField:
		return types.Float64Type
	case idSetField:
		return types.SetType{ElemType: types.StringType}
	default:
		return types.StringType
	}
}

func (f configField) resourceAttribute() schema.Attribute {
	optional := !f.required && !f.readOnly
	computed := !f.required
	switch f.kind {
	case boolField:
		a := schema.BoolAttribute{Required: f.required, Optional: optional, Computed: computed, Description: f.description}
		if value, ok := f.defaultValue.(bool); ok {
			a.Default = booldefault.StaticBool(value)
		}
		return a
	case intField:
		a := schema.Int64Attribute{Required: f.required, Optional: optional, Computed: computed, Description: f.description}
		if value, ok := f.defaultValue.(int); ok {
			a.Default = int64default.StaticInt64(int64(value))
		}
		return a
	case floatField, currencyAmountField:
		return schema.Float64Attribute{Required: f.required, Optional: optional, Computed: computed, Description: f.description}
	case idSetField:
		return schema.SetAttribute{Optional: optional, Computed: computed, ElementType: types.StringType, Description: f.description}
	default:
		a := schema.StringAttribute{Required: f.required, Optional: optional, Computed: computed, Description: f.description}
		if value, ok := f.defaultValue.(string); ok {
			a.Default = stringdefault.StaticString(value)
		}
		if f.name == "id" {
			a.PlanModifiers = []planmodifier.String{stringplanmodifier.UseStateForUnknown()}
		}
		if f.replace {
			a.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
		}
		if !f.readOnly && (f.kind == jsonArrayField || f.kind == jsonObjectField) {
			a.Validators = []validator.String{jsonShapeValidator{kind: f.kind}}
		}
		if len(f.allowed) > 0 {
			a.Validators = []validator.String{stringChoiceValidator(f.allowed)}
		}
		return a
	}
}

func (f configField) dataSourceAttribute(required bool) dataschema.Attribute {
	switch f.kind {
	case boolField:
		return dataschema.BoolAttribute{Computed: true, Description: f.description}
	case intField:
		return dataschema.Int64Attribute{Computed: true, Description: f.description}
	case floatField, currencyAmountField:
		return dataschema.Float64Attribute{Computed: true, Description: f.description}
	case idSetField:
		return dataschema.SetAttribute{Computed: true, ElementType: types.StringType, Description: f.description}
	default:
		return dataschema.StringAttribute{Required: required, Computed: !required, Description: f.description}
	}
}

func (c configuration) objectType() types.ObjectType {
	fields := make(map[string]attr.Type, len(c.fields)+1)
	for _, field := range c.allFields() {
		fields[field.name] = field.attrType()
	}
	return types.ObjectType{AttrTypes: fields}
}

func (c configuration) valuesFromTerraform(ctx context.Context, raw tftypes.Value) (map[string]attr.Value, error) {
	value, err := c.objectType().ValueFromTerraform(ctx, raw)
	if err != nil {
		return nil, err
	}
	object, ok := value.(types.Object)
	if !ok || object.IsNull() || object.IsUnknown() {
		return nil, fmt.Errorf("expected a known configuration object")
	}
	return object.Attributes(), nil
}

func configurationID(values map[string]attr.Value) (string, error) {
	id, ok := values["id"].(types.String)
	if !ok || id.ValueString() == "" {
		return "", fmt.Errorf("configuration lacks an ID")
	}
	return id.ValueString(), nil
}

type stringChoiceValidator []string

func (v stringChoiceValidator) Description(context.Context) string {
	return fmt.Sprintf("Must be one of %v.", []string(v))
}

func (v stringChoiceValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v stringChoiceValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if !slices.Contains(v, req.ConfigValue.ValueString()) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid configuration value", v.Description(ctx))
	}
}

type jsonShapeValidator struct{ kind fieldKind }

func (v jsonShapeValidator) Description(context.Context) string {
	if v.kind == jsonObjectField {
		return "Must be a JSON object."
	}
	return "Must be a JSON array."
}

func (v jsonShapeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonShapeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	raw := []byte(req.ConfigValue.ValueString())
	var value any
	err := json.Unmarshal(raw, &value)
	_, object := value.(map[string]any)
	_, array := value.([]any)
	if err != nil || (v.kind == jsonObjectField && !object) || (v.kind == jsonArrayField && !array) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid JSON configuration", v.Description(ctx))
	}
}

type configurationResource struct {
	spec   configuration
	client *client.Client
}

var _ resource.ResourceWithConfigure = &configurationResource{}
var _ resource.ResourceWithImportState = &configurationResource{}

func (r *configurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.spec.name
}

func (r *configurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	fields := make(map[string]schema.Attribute, len(r.spec.fields)+1)
	for _, field := range r.spec.allFields() {
		fields[field.name] = field.resourceAttribute()
	}
	resp.Schema = schema.Schema{Description: r.spec.description, Attributes: fields}
}

func (r *configurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configuredClient(req.ProviderData, &resp.Diagnostics)
}

func configuredClient(value any, diagnostics *diag.Diagnostics) *client.Client {
	if value == nil {
		return nil
	}
	c, ok := value.(*client.Client)
	if !ok {
		diagnostics.AddError("Invalid provider configuration", "Expected an Argonix API client.")
	}
	return c
}

type apiObject map[string]json.RawMessage

func (r *configurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	values, err := r.spec.valuesFromTerraform(ctx, req.Plan.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading resource plan", err.Error())
		return
	}
	payload, err := r.spec.payload(values, nil)
	if err != nil {
		resp.Diagnostics.AddError("Encoding resource plan", err.Error())
		return
	}
	var result apiObject
	if err := r.client.Create(ctx, r.spec.endpoint, payload, &result); err != nil {
		resp.Diagnostics.AddError("Creating "+r.spec.name, err.Error())
		return
	}
	r.spec.writeState(ctx, &resp.State, result, values, &resp.Diagnostics)
}

func (r *configurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	prior, err := r.spec.valuesFromTerraform(ctx, req.State.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading resource state", err.Error())
		return
	}
	id, err := configurationID(prior)
	if err != nil {
		resp.Diagnostics.AddError("Reading resource state", err.Error())
		return
	}
	var result apiObject
	err = r.client.Read(ctx, r.spec.detail(id), &result)
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Reading "+r.spec.name, err.Error())
		return
	}
	r.spec.writeState(ctx, &resp.State, result, prior, &resp.Diagnostics)
}

func (r *configurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	values, err := r.spec.valuesFromTerraform(ctx, req.Plan.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading resource plan", err.Error())
		return
	}
	prior, err := r.spec.valuesFromTerraform(ctx, req.State.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Reading resource state", err.Error())
		return
	}
	id, err := configurationID(prior)
	if err != nil {
		resp.Diagnostics.AddError("Reading resource state", err.Error())
		return
	}
	payload, err := r.spec.payload(values, prior)
	if err != nil {
		resp.Diagnostics.AddError("Encoding resource plan", err.Error())
		return
	}
	var result apiObject
	endpoint := r.spec.detail(id)
	if len(payload) == 0 {
		err = r.client.Read(ctx, endpoint, &result)
	} else {
		err = r.client.Update(ctx, endpoint, payload, &result)
		if err == nil && r.spec.readAfterUpdate {
			result = nil
			err = r.client.Read(ctx, endpoint, &result)
		}
	}
	if err != nil {
		resp.Diagnostics.AddError("Updating "+r.spec.name, err.Error())
		return
	}
	r.spec.writeState(ctx, &resp.State, result, values, &resp.Diagnostics)
}

func (r *configurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, r.spec.detail(id.ValueString())); err != nil {
		resp.Diagnostics.AddError("Deleting "+r.spec.name, err.Error())
	}
}

func (r *configurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (c configuration) payload(values, prior map[string]attr.Value) (map[string]any, error) {
	payload := make(map[string]any)
	for _, field := range c.fields {
		value := values[field.name]
		if field.readOnly || value.IsNull() || value.IsUnknown() {
			continue
		}
		if old := prior[field.name]; old != nil && old.Equal(value) {
			continue
		}
		switch v := value.(type) {
		case types.String:
			if len(field.allowed) > 0 && !slices.Contains(field.allowed, v.ValueString()) {
				return nil, fmt.Errorf("%s must be one of %v", field.name, field.allowed)
			}
			switch field.kind {
			case jsonObjectField, jsonArrayField:
				raw := json.RawMessage(v.ValueString())
				if !json.Valid(raw) {
					return nil, fmt.Errorf("%s contains invalid JSON", field.name)
				}
				payload[field.name] = raw
			default:
				payload[field.name] = v.ValueString()
			}
		case types.Bool:
			payload[field.name] = v.ValueBool()
		case types.Int64:
			payload[field.name] = v.ValueInt64()
		case types.Float64:
			payload[field.name] = v.ValueFloat64()
		case types.Set:
			ids := make([]string, 0, len(v.Elements()))
			for _, element := range v.Elements() {
				id, ok := element.(types.String)
				if !ok || id.IsNull() || id.IsUnknown() {
					return nil, fmt.Errorf("%s contains an unresolved or null ID", field.name)
				}
				ids = append(ids, id.ValueString())
			}
			payload[field.name] = ids
		}
	}
	return payload, nil
}

func (c configuration) writeState(ctx context.Context, state *tfsdk.State, body apiObject, prior map[string]attr.Value, diagnostics *diag.Diagnostics) {
	value, err := c.stateValue(body, prior)
	if err != nil {
		diagnostics.AddError("Decoding "+c.name, err.Error())
		return
	}
	raw, err := value.ToTerraformValue(ctx)
	if err != nil {
		diagnostics.AddError("Saving "+c.name+" state", err.Error())
		return
	}
	state.Raw = raw
}

func (c configuration) stateValue(body apiObject, prior map[string]attr.Value) (types.Object, error) {
	values := make(map[string]attr.Value, len(c.fields)+1)
	for _, field := range c.allFields() {
		key := field.name
		if field.responseName != "" {
			key = field.responseName
		}
		value, err := field.decode(body[key], prior[field.name])
		if err != nil {
			return types.ObjectNull(c.objectType().AttrTypes), fmt.Errorf("%s: %w", key, err)
		}
		values[field.name] = value
	}
	if _, err := configurationID(values); err != nil {
		return types.ObjectNull(c.objectType().AttrTypes), err
	}
	return types.ObjectValueMust(c.objectType().AttrTypes, values), nil
}

func (f configField) decode(raw json.RawMessage, prior attr.Value) (attr.Value, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return f.attrType().ValueFromTerraform(context.Background(), tftypes.NewValue(f.attrType().TerraformType(context.Background()), nil))
	}
	switch f.kind {
	case boolField:
		var value bool
		err := json.Unmarshal(raw, &value)
		return types.BoolValue(value), err
	case intField:
		var value int64
		err := json.Unmarshal(raw, &value)
		return types.Int64Value(value), err
	case floatField:
		var value float64
		err := json.Unmarshal(raw, &value)
		return types.Float64Value(value), err
	case currencyAmountField:
		var micros int64
		err := json.Unmarshal(raw, &micros)
		return types.Float64Value(float64(micros) / 1_000_000), err
	case idSetField:
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, err
		}
		ids := make([]attr.Value, 0, len(items))
		for _, item := range items {
			var id string
			if err := json.Unmarshal(item, &id); err != nil {
				var object struct{ ID string }
				if err := json.Unmarshal(item, &object); err != nil {
					return nil, err
				}
				id = object.ID
			}
			if id == "" {
				return nil, fmt.Errorf("related object lacks an ID")
			}
			ids = append(ids, types.StringValue(id))
		}
		return types.SetValueMust(types.StringType, ids), nil
	case jsonObjectField, jsonArrayField:
		value, err := decodeJSON(raw)
		if err != nil {
			return nil, err
		}
		if old, ok := prior.(types.String); ok && !old.IsNull() && !old.IsUnknown() {
			expected, err := decodeJSON([]byte(old.ValueString()))
			if err == nil && equivalentJSON(value, expected, f.jsonEquivalent) {
				return old, nil
			}
		}
		canonical, err := json.Marshal(value)
		return types.StringValue(string(canonical)), err
	default:
		var value string
		err := json.Unmarshal(raw, &value)
		return types.StringValue(value), err
	}
}

func decodeJSON(raw []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	err := decoder.Decode(&value)
	return value, err
}

func equivalentJSON(actual, expected any, additional func(any, any) bool) bool {
	if equalJSONValue(actual, expected) {
		return true
	}
	return additional != nil && additional(actual, expected)
}

func equalJSONValue(left, right any) bool {
	switch a := left.(type) {
	case map[string]any:
		b, ok := right.(map[string]any)
		return ok && maps.EqualFunc(a, b, equalJSONValue)
	case []any:
		b, ok := right.([]any)
		return ok && slices.EqualFunc(a, b, equalJSONValue)
	case json.Number:
		b, ok := right.(json.Number)
		if !ok {
			return false
		}
		if a == b {
			return true
		}
		// Treat API-normalized decimal spellings as equal without rounding
		// large integers through float64.
		x, okX := new(big.Rat).SetString(a.String())
		y, okY := new(big.Rat).SetString(b.String())
		return okX && okY && x.Cmp(y) == 0
	default:
		return reflect.DeepEqual(left, right)
	}
}

type configurationDataSource struct {
	spec   configuration
	list   bool
	client *client.Client
}

var _ datasource.DataSourceWithConfigure = &configurationDataSource{}

func (d *configurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	name := d.spec.name
	if d.list {
		name = d.spec.plural
	}
	resp.TypeName = req.ProviderTypeName + "_" + name
}

func (d *configurationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	fields := make(map[string]dataschema.Attribute, len(d.spec.fields)+1)
	for _, field := range d.spec.allFields() {
		fields[field.name] = field.dataSourceAttribute(field.name == "id" && !d.list)
	}
	if !d.list {
		resp.Schema = dataschema.Schema{Description: "Reads an Argonix " + d.spec.name + " by ID.", Attributes: fields}
		return
	}
	resp.Schema = dataschema.Schema{
		Description: "Lists Argonix " + d.spec.plural + " in the organization.",
		Attributes: map[string]dataschema.Attribute{
			d.spec.plural: dataschema.ListNestedAttribute{
				Computed:     true,
				NestedObject: dataschema.NestedAttributeObject{Attributes: fields},
			},
		},
	}
}

func (d *configurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configuredClient(req.ProviderData, &resp.Diagnostics)
}

func (d *configurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !d.list {
		var id types.String
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &id)...)
		if resp.Diagnostics.HasError() {
			return
		}
		var result apiObject
		if err := d.client.Read(ctx, d.spec.detail(id.ValueString()), &result); err != nil {
			resp.Diagnostics.AddError("Reading "+d.spec.name, err.Error())
			return
		}
		d.spec.writeState(ctx, &resp.State, result, nil, &resp.Diagnostics)
		return
	}
	var results []apiObject
	if err := d.client.List(ctx, d.spec.endpoint, &results); err != nil {
		resp.Diagnostics.AddError("Listing "+d.spec.plural, err.Error())
		return
	}
	entries := make([]attr.Value, 0, len(results))
	for _, result := range results {
		if d.spec.detailForList {
			var id string
			if err := json.Unmarshal(result["id"], &id); err != nil || id == "" {
				resp.Diagnostics.AddError("Listing "+d.spec.plural, "A list entry lacks its ID.")
				return
			}
			if err := d.client.Read(ctx, d.spec.detail(id), &result); err != nil {
				if client.IsNotFound(err) {
					continue
				}
				resp.Diagnostics.AddError("Reading "+d.spec.name, err.Error())
				return
			}
		}
		entry, err := d.spec.stateValue(result, nil)
		if err != nil {
			resp.Diagnostics.AddError("Decoding "+d.spec.name, err.Error())
			return
		}
		entries = append(entries, entry)
	}
	list := types.ListValueMust(d.spec.objectType(), entries)
	wrapper := types.ObjectValueMust(map[string]attr.Type{d.spec.plural: list.Type(ctx)}, map[string]attr.Value{d.spec.plural: list})
	raw, err := wrapper.ToTerraformValue(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Saving "+d.spec.plural+" state", err.Error())
		return
	}
	resp.State.Raw = raw
}
