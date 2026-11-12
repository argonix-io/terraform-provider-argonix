package provider

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestConfigurationSchemas(t *testing.T) {
	ctx := t.Context()
	for _, factory := range configurationResources() {
		r := factory()
		var metadata resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "argonix"}, &metadata)
		t.Run(metadata.TypeName, func(t *testing.T) {
			var resp resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &resp)
			if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
				t.Fatalf("invalid schema: %v", diags)
			}
		})
	}
	for _, factory := range configurationDataSources() {
		d := factory()
		var metadata datasource.MetadataResponse
		d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "argonix"}, &metadata)
		t.Run("data_"+metadata.TypeName, func(t *testing.T) {
			var resp datasource.SchemaResponse
			d.Schema(ctx, datasource.SchemaRequest{}, &resp)
			if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
				t.Fatalf("invalid schema: %v", diags)
			}
		})
	}
}

func TestTerraformProviderPublishesConfigurationSchemas(t *testing.T) {
	server := providerserver.NewProtocol6(New("test")())()
	result, err := server.GetProviderSchema(t.Context(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("provider protocol failed to export schema: %+v", diagnostic)
		}
	}
	for _, spec := range managedConfigurations() {
		name := "argonix_" + spec.name
		if result.ResourceSchemas[name] == nil || result.DataSourceSchemas[name] == nil || result.DataSourceSchemas["argonix_"+spec.plural] == nil {
			t.Fatalf("resource or data sources missing from the provider protocol for %s", name)
		}
	}
}

func TestPeriodicJobReadAndClearRelations(t *testing.T) {
	var response apiObject
	err := json.Unmarshal([]byte(`{
		"id":"job-1","name":"Daily review","kind":"agent","mode":"discussion",
		"connectors":[{"id":"connector-1","name":"Logs"}],"application_root_ids":["app-1"],
		"notification_channels":["chat-1"],"alert_channels":["alert-1"]
	}`), &response)
	if err != nil {
		t.Fatal(err)
	}
	state, err := periodicJobConfiguration.stateValue(response, nil)
	if err != nil {
		t.Fatal(err)
	}
	prior := state.Attributes()
	wantConnectors := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("connector-1")})
	if got := prior["connector_ids"]; !got.Equal(wantConnectors) {
		t.Fatalf("connector IDs = %v", got)
	}
	plan := state.Attributes()
	plan["notification_channel_ids"] = types.SetValueMust(types.StringType, []attr.Value{})
	payload, err := periodicJobConfiguration.payload(plan, prior)
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]any{"notification_channel_ids": []string{}}; !reflect.DeepEqual(payload, want) {
		t.Fatalf("payload = %#v, want %#v", payload, want)
	}
}

func TestSanityJSONRoundTripAndPartialUpdate(t *testing.T) {
	var response apiObject
	err := json.Unmarshal([]byte(`{
		"id":"check-1","name":"Health","test_type":"scenario","schedule_mode":"cron",
		"schedule":"0 8 * * *","schedule_timezone":"Europe/Paris","next_run_at":"2026-11-12T08:00:00Z",
		"steps":[{"id":"health","kind":"http","params":{"url":"https://example.com"}}],
		"input_variables":{"minimum":9007199254740993,"environment":"prod"}
	}`), &response)
	if err != nil {
		t.Fatal(err)
	}
	original := types.StringValue(`{ "environment": "prod", "minimum": 9007199254740993 }`)
	state, err := syntheticTestConfiguration.stateValue(response, map[string]attr.Value{"input_variables": original})
	if err != nil {
		t.Fatal(err)
	}
	prior := state.Attributes()
	if !prior["input_variables"].Equal(original) {
		t.Fatalf("JSON reformatting would produce inconsistent state: %v", prior["input_variables"])
	}
	plan := state.Attributes()
	plan["name"] = types.StringValue("New check name")
	payload, err := syntheticTestConfiguration.payload(plan, prior)
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]any{"name": "New check name"}; !reflect.DeepEqual(payload, want) {
		t.Fatalf("payload = %#v, want %#v; updating the name must preserve the next scheduled run", payload, want)
	}
}

func TestLogSourceNormalizesSnapshotDefaultsWithoutMaskingDrift(t *testing.T) {
	config := types.StringValue(`{"namespace":"production","pod":"api","volume_spike_min_ratio":2}`)
	prior := map[string]attr.Value{"config": config}
	var response apiObject
	if err := json.Unmarshal([]byte(`{
		"id":"log-1","analysis_mode":"snapshot",
		"config":{"namespace":"production","pod":"api","volume_spike_min_ratio":2.0,"analysis_mode":"snapshot","tail_lines":500}
	}`), &response); err != nil {
		t.Fatal(err)
	}
	state, err := logSourceConfiguration.stateValue(response, prior)
	if err != nil {
		t.Fatal(err)
	}
	if !state.Attributes()["config"].Equal(config) {
		t.Fatalf("server defaults or decimal spelling caused state drift: %v", state.Attributes()["config"])
	}
	response["config"] = json.RawMessage(`{"namespace":"production","pod":"different-api","volume_spike_min_ratio":2.0,"analysis_mode":"snapshot","tail_lines":500}`)
	state, err = logSourceConfiguration.stateValue(response, prior)
	if err != nil {
		t.Fatal(err)
	}
	if state.Attributes()["config"].Equal(config) {
		t.Fatal("a change to the log target was hidden by normalization")
	}
}

func TestLegacyAPINullsUseTerraformDefaults(t *testing.T) {
	for _, tc := range []struct {
		spec configuration
		body string
		want map[string]attr.Value
	}{
		{
			spec: workflowConfiguration,
			body: `{"id":"workflow","description":null,"schedule":null,"required_connector_types":null,"is_active":null,"schedule_next_run_at":null}`,
			want: map[string]attr.Value{
				"description":              types.StringValue(""),
				"schedule":                 types.StringValue(""),
				"input_hints":              types.StringValue("{}"),
				"required_connector_types": types.StringValue("[]"),
				"requires_confirmation":    types.BoolValue(false),
				"is_active":                types.BoolValue(true),
				"schedule_next_run_at":     types.StringNull(),
			},
		},
		{
			spec: syntheticTestConfiguration,
			body: `{"id":"check","check_interval":null,"tags":null,"last_duration_ms":null}`,
			want: map[string]attr.Value{
				"check_interval":   types.Int64Value(300),
				"tags":             types.StringValue("[]"),
				"locations":        types.StringValue("[]"),
				"last_duration_ms": types.Float64Null(),
			},
		},
	} {
		t.Run(tc.spec.name, func(t *testing.T) {
			var response apiObject
			if err := json.Unmarshal([]byte(tc.body), &response); err != nil {
				t.Fatal(err)
			}
			prior := map[string]attr.Value{
				"schedule": types.StringValue("0 9 * * *"),
			}
			state, err := tc.spec.stateValue(response, prior)
			if err != nil {
				t.Fatal(err)
			}
			for key, want := range tc.want {
				if got := state.Attributes()[key]; !got.Equal(want) {
					t.Errorf("%s = %v, want %v", key, got, want)
				}
			}
		})
	}
}
