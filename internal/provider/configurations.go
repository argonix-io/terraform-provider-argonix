package provider

import (
	"encoding/json"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func required(name string, kind fieldKind, description string) configField {
	return configField{name: name, kind: kind, description: description, required: true}
}

func setting(name string, kind fieldKind, description string) configField {
	return configField{name: name, kind: kind, description: description}
}

func defaulted(name string, kind fieldKind, value any, description string) configField {
	return configField{name: name, kind: kind, defaultValue: value, description: description}
}

func output(name string, kind fieldKind) configField {
	return configField{name: name, kind: kind, readOnly: true}
}

func ids(name, response, description string) configField {
	return configField{name: name, kind: idSetField, responseName: response, description: description}
}

func immutable(name string, kind fieldKind, description string) configField {
	return configField{name: name, kind: kind, description: description, replace: true}
}

var workflowConfiguration = configuration{
	name: "workflow", plural: "workflows", endpoint: "/argos/workflows/",
	description: "Manages an Argonix Argos workflow.", detailForList: true,
	fields: []configField{
		required("name", stringField, "Display name of the workflow."),
		setting("slug", stringField, "Unique URL-friendly identifier."),
		defaulted("description", stringField, "", "Description of the workflow."),
		required("category", stringField, "identity, incident, onboarding, devops, security or general."),
		required("steps", jsonArrayField, "JSON array of workflow step definitions."),
		defaulted("input_hints", jsonObjectField, "{}", "JSON object describing inputs referenced in the steps."),
		output("inputs", jsonArrayField),
		defaulted("required_connector_types", jsonArrayField, "[]", "JSON array of required connector types."),
		defaulted("requires_confirmation", boolField, false, "Require human confirmation before execution."),
		setting("extract_inputs_with_ai", boolField, "Allow AI extraction of inputs from the request."),
		defaulted("schedule", stringField, "", "Cron expression, or an empty string for no schedule."),
		setting("schedule_timezone", stringField, "IANA timezone for the cron expression."),
		setting("schedule_inputs", jsonObjectField, "JSON object of inputs supplied to scheduled runs."),
		output("schedule_last_run_at", stringField),
		output("schedule_next_run_at", stringField),
		setting("security_trigger", stringField, "Event that triggers the workflow for defensive triage."),
		setting("security_min_severity", stringField, "Minimum security event severity."),
		setting("security_detection_rule", stringField, "Optional detection rule ID."),
		setting("incident_trigger", stringField, "Optional incident event for investigation."),
		setting("incident_min_severity", stringField, "Minimum incident severity."),
		setting("incident_min_confidence", intField, "Minimum confidence for matching incidents."),
		setting("incident_max_affected_resources", intField, "Maximum affected resource count for matching incidents."),
		setting("incident_allowed_environments", jsonArrayField, "JSON array of environments allowed for incident investigation."),
		setting("incident_timeout_minutes", intField, "Incident workflow timeout in minutes."),
		defaulted("is_active", boolField, true, "Whether the workflow is active."),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var syntheticTestConfiguration = configuration{
	name: "synthetic_test", plural: "synthetic_tests", endpoint: "/synthetic-tests/",
	description: "Manages an Argonix synthetic test or Sanity check.", detailForList: true,
	fields: []configField{
		required("name", stringField, "Name of the check."),
		defaulted("description", stringField, "", "Description of the check."),
		defaulted("is_active", boolField, true, "Whether the check is active."),
		required("test_type", stringField, "api, browser or scenario (Sanity check)."),
		required("steps", jsonArrayField, "JSON array of check steps."),
		defaulted("check_interval", intField, 300, "Interval between runs in seconds."),
		defaulted("timeout", intField, 30, "Per-step timeout in seconds."),
		defaulted("tags", jsonArrayField, "[]", "JSON array of tags."),
		defaulted("locations", jsonArrayField, "[]", "JSON array of region codes."),
		setting("input_variables", jsonObjectField, "JSON object of default Sanity check input variables."),
		setting("run_timeout", intField, "Maximum duration of a Sanity check in seconds."),
		setting("schedule_mode", stringField, "manual, interval or cron."),
		setting("schedule", stringField, "Cron expression when schedule_mode is cron."),
		setting("schedule_timezone", stringField, "IANA timezone for the cron expression."),
		setting("application_key", stringField, "Application UUID associated with the check."),
		{name: "resource_id", kind: stringField, responseName: "resolved_resource", description: "Inventory resource UUID associated with the check."},
		output("allow_mutating_actions", boolField),
		output("steps_count", intField),
		output("result_is_stale", boolField),
		output("current_status", stringField),
		output("next_run_at", stringField),
		output("last_run_at", stringField),
		output("last_duration_ms", floatField),
		output("last_warning_count", intField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var periodicJobConfiguration = configuration{
	name: "periodic_job", plural: "periodic_jobs", endpoint: "/argos/jobs/",
	description: "Manages an Argonix periodic job.", detailForList: true,
	fields: []configField{
		required("name", stringField, "Name of the periodic job."),
		immutable("kind", stringField, "agent or application_reconcile."),
		setting("prompt", stringField, "Prompt required for an agent job."),
		ids("connector_ids", "connectors", "Connector IDs available to the job."),
		ids("application_root_ids", "application_root_ids", "Application IDs available to the job."),
		setting("schedule_hour", intField, "UTC hour of execution (0–23)."),
		setting("schedule_minute", intField, "UTC minute of execution (0–59)."),
		setting("schedule_days", jsonArrayField, "JSON array of weekday numbers, with Monday as 0. Empty means every day."),
		ids("notification_channel_ids", "notification_channels", "Chat channel IDs receiving notifications."),
		ids("alert_channel_ids", "alert_channels", "Alert channel IDs receiving results."),
		setting("auto_create_ticket", boolField, "Create a ticket through a configured matching connector and alert channel."),
		setting("deliverable_retention_days", intField, "Number of days to retain deliverables."),
		setting("deliverable_retention_action", stringField, "delete or move expired deliverables."),
		setting("deliverable_retention_destination", stringField, "Destination when retention action is move."),
		setting("is_active", boolField, "Whether the schedule is active."),
		output("mode", stringField),
		setting("timeout_seconds", intField, "Maximum execution time in seconds."),
		output("latest_run_status", stringField),
		output("latest_run_date", stringField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var healthNotebookConfiguration = configuration{
	name: "health_notebook", plural: "health_notebooks", endpoint: "/argos/health-notebooks/",
	description: "Manages an Argonix health notebook or security brief.", detailForList: true,
	fields: []configField{
		setting("name", stringField, "Name of the notebook."),
		immutable("report_kind", stringField, "health or security."),
		ids("connector_ids", "connectors", "Connector IDs in scope."),
		setting("use_applications", boolField, "Scope the notebook to applications."),
		setting("perspective", stringField, "all, business or platform."),
		ids("application_root_ids", "application_root_ids", "Selected application IDs."),
		setting("application_label_selector", jsonArrayField, "JSON array of exact key and value pairs for detected applications."),
		setting("include_unattributed", boolField, "Include findings without application attribution."),
		setting("minimum_attribution_confidence", intField, "Minimum attribution confidence from 0 to 100."),
		setting("is_active", boolField, "Whether the schedule is active."),
		setting("frequency", stringField, "daily, weekly or monthly."),
		setting("schedule_hour", intField, "UTC hour of execution (0–23)."),
		setting("schedule_minute", intField, "UTC minute of execution (0–59)."),
		setting("schedule_days", jsonArrayField, "JSON array of weekday numbers, with Monday as 0."),
		setting("schedule_day_of_month", intField, "Day of month for monthly notebooks."),
		setting("custom_prompt", stringField, "Additional instructions for the report."),
		setting("lookback_hours", intField, "History window from 1 to 168 hours."),
		setting("include_threats", boolField, "Include detected threats."),
		setting("include_findings", boolField, "Include security findings."),
		setting("include_monitors", boolField, "Include monitor results."),
		setting("include_incidents", boolField, "Include incidents."),
		setting("include_cost", boolField, "Include cost signals."),
		setting("include_log_anomalies", boolField, "Include log anomalies."),
		setting("include_vulnerabilities", boolField, "Include vulnerabilities."),
		setting("include_runtime", boolField, "Include runtime signals."),
		setting("include_dns", boolField, "Include DNS signals."),
		setting("include_network", boolField, "Include network signals."),
		setting("include_operational_signals", boolField, "Include application operational signals."),
		ids("notification_channel_ids", "notification_channels", "Chat channel IDs receiving reports."),
		ids("alert_channel_ids", "alert_channels", "Alert channel IDs receiving reports."),
		setting("ticketing_min_severity", stringField, "Minimum severity for ticket creation."),
		output("ticketing_action", stringField),
		output("latest_score", intField),
		output("latest_severity", stringField),
		output("latest_evidence_coverage", jsonObjectField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var srePatrolConfiguration = configuration{
	name: "sre_patrol", plural: "sre_patrols", endpoint: "/argos/patrols/",
	description: "Manages an Argonix SRE patrol for operational observation.",
	fields: []configField{
		setting("name", stringField, "Name of the patrol."),
		ids("connector_ids", "connectors", "Connector IDs in scope."),
		setting("perspective", stringField, "all, business or platform."),
		setting("use_applications", boolField, "Scope the patrol to applications."),
		ids("application_root_ids", "application_root_ids", "Selected application IDs."),
		setting("include_unattributed", boolField, "Include signals without application attribution."),
		setting("minimum_attribution_confidence", intField, "Minimum attribution confidence from 0 to 100."),
		ids("notification_channel_ids", "notification_channels", "Chat channel IDs receiving results."),
		ids("alert_channel_ids", "alert_channels", "Alert channel IDs receiving alerts."),
		ids("delivery_channel_ids", "delivery_channels", "Alert channel IDs receiving patrol reports."),
		setting("is_active", boolField, "Whether the patrol is active."),
		setting("cadence_minutes", intField, "Minutes between sweeps, from 5 to 60."),
		setting("min_confidence", intField, "Minimum confidence from 0 to 100."),
		setting("ticketing_min_severity", stringField, "Minimum severity for ticket creation."),
		setting("custom_prompt", stringField, "Additional instructions for interpreting operational signals."),
		output("autonomy", stringField),
		output("last_run_status", stringField),
		output("last_run_date", stringField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var codeProjectConfiguration = configuration{
	name: "code_project", plural: "code_projects", endpoint: "/security/code-projects/",
	description: "Manages a Git repository for defensive code analysis.",
	fields: []configField{
		setting("name", stringField, "Project name; defaults to the repository name."),
		setting("key", stringField, "Unique project key; generated from the name if omitted."),
		required("repository_url", stringField, "Canonical repository URL without a .git suffix or trailing slash."),
		setting("connector", stringField, "GitHub or GitLab connector ID."),
		output("repository_provider", stringField),
		output("entry_kind", stringField),
		setting("provider_project_id", stringField, "Project ID at the Git provider."),
		setting("default_branch", stringField, "Default branch."),
		setting("visibility", stringField, "private, internal or public."),
		setting("quality_gate", stringField, "Quality gate ID."),
		setting("enabled_engines", jsonArrayField, "JSON array of static analysis engines."),
		setting("publish_quality_gate", boolField, "Publish quality gate status to the Git provider."),
		setting("publish_review_feedback", boolField, "Publish review feedback to the Git provider."),
		setting("excluded_paths", jsonArrayField, "JSON array of excluded paths or glob patterns."),
		{name: "webhook_analysis_mode", kind: stringField, allowed: []string{"none", "standard"}, description: "Static analysis mode for repository webhooks: none or standard."},
		setting("webhook_standard_scope", stringField, "Standard analysis scope on the default branch."),
		setting("webhook_standard_non_default_scope", stringField, "Standard analysis scope on non-default branches."),
		setting("webhook_standard_long_lived_scope", stringField, "Standard analysis scope on long-lived branches."),
		setting("webhook_long_lived_patterns", jsonArrayField, "JSON array of long-lived branch patterns."),
		setting("webhook_default_gate_scope", stringField, "all or new_code for the default branch."),
		setting("webhook_non_default_gate_scope", stringField, "all or new_code for non-default branches."),
		setting("webhook_long_lived_gate_scope", stringField, "all or new_code for long-lived branches."),
		setting("long_lived_quality_gate", stringField, "Quality gate ID for long-lived branches."),
		setting("sql_dialect", stringField, "SQL dialect required if SQLFluff is enabled."),
		setting("auto_assign", boolField, "Assign findings automatically where supported."),
		setting("is_active", boolField, "Whether the project is active."),
		output("latest_analysis", stringField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var scanTargetConfiguration = configuration{
	name: "scan_target", plural: "scan_targets", endpoint: "/security/scan-targets/",
	description: "Manages an Argonix repository or container image scanner target.",
	fields: []configField{
		setting("name", stringField, "Display name of the target."),
		{name: "kind", kind: stringField, required: true, allowed: []string{"repository", "image"}, description: "repository or image."},
		required("identifier", stringField, "Repository URL or container image identifier."),
		setting("connector", stringField, "Git or registry connector ID."),
		setting("project", stringField, "Code project ID."),
		setting("branch", stringField, "Branch for a repository target."),
		setting("default_ref", stringField, "Tag or digest for an unpinned image."),
		setting("platform", stringField, "OCI image platform, such as linux/amd64."),
		required("enabled_scanners", jsonArrayField, "Nonempty JSON array of scanners supported for this target kind."),
		setting("frequency", stringField, "manual, daily or weekly."),
		setting("schedule_hour", intField, "UTC hour of execution (0–23)."),
		setting("schedule_minute", intField, "UTC minute of execution (0–59)."),
		setting("schedule_days", jsonArrayField, "JSON array of weekday numbers, with Monday as 0."),
		setting("is_active", boolField, "Whether the target is active."),
		setting("webhook_enabled", boolField, "Accept repository webhooks for this target."),
		setting("enforce_merge_requests", boolField, "Report scan results to the merge request gate."),
		setting("gate_environment", stringField, "Optional environment for the merge request gate."),
		output("last_dispatched_at", stringField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var monitorGroupConfiguration = configuration{
	name: "monitor_group", plural: "monitor_groups", endpoint: "/monitor-groups/",
	description: "Manages a folder for Argonix monitors.",
	fields: []configField{
		required("name", stringField, "Monitor folder name."),
		setting("description", stringField, "Monitor folder description."),
		output("is_managed", boolField), output("monitors_count", intField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var finopsTenantConfiguration = configuration{
	name: "finops_tenant", plural: "finops_tenants", endpoint: "/finops/tenants/",
	description: "Manages a FinOps tenant for cost allocation.", readAfterUpdate: true,
	fields: []configField{
		required("name", stringField, "Tenant name."),
		{name: "kind", kind: stringField, allowed: []string{"business_unit", "client", "team", "product", "environment"}, description: "business_unit, client, team, product or environment."},
		setting("color", stringField, "Display color."),
		setting("parent", stringField, "Parent tenant UUID."),
		setting("monthly_budget_micros", intField, "Monthly budget in millionths of a currency unit."),
		setting("contact_email", stringField, "Tenant contact email address."),
		setting("description", stringField, "Tenant description."),
		output("slug", stringField), output("is_system", boolField), output("rules", jsonArrayField),
		output("date_created", stringField), output("date_modified", stringField),
	},
}

var costBudgetConfiguration = configuration{
	name: "cost_budget", plural: "cost_budgets", endpoint: "/finops/budgets/",
	description: "Manages a FinOps cost budget.",
	fields: []configField{
		required("name", stringField, "Budget name."),
		{name: "amount", kind: currencyAmountField, required: true, responseName: "amount_micros", description: "Budget amount in currency units, with up to six decimal places."},
		setting("period", stringField, "monthly or quarterly."),
		setting("currency", stringField, "Three-letter uppercase currency code."),
		setting("filters", jsonObjectField, "JSON object of cost filters. An empty object covers the organization."),
		setting("alert_thresholds_pct", jsonArrayField, "JSON array of unique percentages in ascending order."),
		ids("alert_channels", "alert_channels", "Alert channel IDs receiving budget alerts."),
		setting("is_active", boolField, "Whether the budget is active."),
		output("recent_alerts", jsonArrayField), output("date_created", stringField), output("date_updated", stringField),
	},
}

var costPerspectiveConfiguration = configuration{
	name: "cost_perspective", plural: "cost_perspectives", endpoint: "/finops/perspectives/",
	description: "Manages a saved FinOps view.",
	fields: []configField{
		required("name", stringField, "View name."),
		setting("filters", jsonObjectField, "JSON object of saved cost filters."),
		output("date_created", stringField), output("date_updated", stringField),
	},
}

var monitorDiscoverySourceConfiguration = configuration{
	name: "monitor_discovery_source", plural: "monitor_discovery_sources", endpoint: "/monitor-discovery-sources/",
	description: "Manages discovery of monitors from Kubernetes, OpenShift or cloud inventory.",
	fields: []configField{
		required("name", stringField, "Discovery source name."),
		required("connector", stringField, "Connector UUID in the organization."),
		setting("source_type", stringField, "kubernetes or cloud_inventory."),
		setting("enabled", boolField, "Enable periodic discovery."),
		setting("resource_kinds", jsonArrayField, "JSON array of resource kinds supported by the connector."),
		setting("include_private_targets", boolField, "Include internal hostnames and addresses when the probe can reach them."),
		setting("namespace_allow", jsonArrayField, "JSON array of namespaces to include."),
		setting("namespace_deny", jsonArrayField, "JSON array of namespaces to exclude."),
		setting("host_include_regex", stringField, "Only discover hosts matching this regular expression."),
		setting("host_exclude_regex", stringField, "Skip hosts matching this regular expression."),
		setting("ingress_classes", jsonArrayField, "JSON array of Ingress classes to include."),
		setting("label_selector", jsonObjectField, "JSON object of labels that discovered resources must match."),
		setting("label_tag_keys", jsonArrayField, "JSON array of label keys used as tags on created monitors."),
		setting("expand_paths", boolField, "Create monitors for individual Ingress or HTTPRoute paths."),
		setting("default_check_interval", intField, "Check interval in seconds for newly created monitors."),
		setting("default_group", stringField, "Monitor group UUID in the organization."),
		setting("default_auto_investigate", boolField, "Enable investigation on monitors created from this source."),
		setting("sync_interval_minutes", intField, "Minimum number of minutes between discovery runs."),
		setting("missing_grace_syncs", intField, "Consecutive missing observations before the managed monitor is removed."),
		output("connector_name", stringField), output("connector_type", stringField), output("default_group_name", stringField),
		output("default_auto_remediate", boolField), output("last_sync_status", stringField),
		output("last_sync_at", stringField), output("last_sync_message", stringField),
		output("last_created_count", intField), output("last_updated_count", intField),
		output("last_deleted_count", intField), output("last_skipped_count", intField),
		output("monitors_count", intField), output("date_created", stringField), output("date_modified", stringField),
	},
}

var logSourceConfiguration = configuration{
	name: "log_source", plural: "log_sources", endpoint: "/log-sources/",
	description: "Manages an Argonix log analysis source on an existing connector.",
	fields: []configField{
		required("name", stringField, "Log source name."),
		required("connector", stringField, "Log connector UUID in the organization."),
		setting("query", stringField, "Provider-specific log filter."),
		{name: "config", kind: jsonObjectField, jsonEquivalent: equivalentLogSourceConfig, description: "Provider settings and optional detection and retention overrides as a JSON object."},
		setting("field_mapping", jsonObjectField, "JSON object mapping the message, timestamp, status, service and host fields."),
		setting("is_active", boolField, "Enable periodic log analysis."),
		setting("schedule_minutes", intField, "Number of minutes between analyses (5–1440)."),
		setting("sample_limit", intField, "Maximum sample lines per analysis."),
		setting("health_intelligence_mode", stringField, "disabled, shadow or active."),
		output("connector_name", stringField), output("connector_type", stringField), output("analysis_mode", stringField),
		output("detection_settings", jsonObjectField), output("retention_settings", jsonObjectField),
		output("last_analyzed_at", stringField), output("date_created", stringField), output("date_modified", stringField),
	},
}

func equivalentLogSourceConfig(actual, expected any) bool {
	result, resultOK := actual.(map[string]any)
	config, configOK := expected.(map[string]any)
	if !resultOK || !configOK || result["analysis_mode"] != "snapshot" {
		return false
	}
	// Snapshot sources get server defaults. Ignore these keys only when the
	// user left them out and their returned values match the API defaults.
	result = maps.Clone(result)
	if _, explicit := config["analysis_mode"]; !explicit {
		delete(result, "analysis_mode")
	}
	if _, explicit := config["tail_lines"]; !explicit && equalJSONValue(result["tail_lines"], json.Number("500")) {
		delete(result, "tail_lines")
	}
	return equalJSONValue(result, config)
}

var codeQualityGateConfiguration = configuration{
	name: "code_quality_gate", plural: "code_quality_gates", endpoint: "/security/code-quality-gates/",
	fields: []configField{
		output("organization", stringField), output("name", stringField), output("description", stringField),
		output("version", intField), output("is_default", boolField), output("is_active", boolField),
		output("conditions", jsonArrayField), output("date_created", stringField), output("date_modified", stringField),
	},
}

var codeProfileConfiguration = configuration{
	name: "code_profile", plural: "code_profiles", endpoint: "/security/code-profiles/",
	fields: []configField{
		output("name", stringField), output("language", stringField), output("version", intField),
		output("is_active", boolField), output("is_default", boolField), output("is_builtin", boolField),
		output("scanner_defaults", jsonObjectField), output("rules", jsonArrayField), output("date_created", stringField),
	},
}

func managedConfigurations() []configuration {
	return []configuration{
		workflowConfiguration, syntheticTestConfiguration, periodicJobConfiguration,
		healthNotebookConfiguration, srePatrolConfiguration, codeProjectConfiguration, scanTargetConfiguration,
		monitorGroupConfiguration, finopsTenantConfiguration, costBudgetConfiguration, costPerspectiveConfiguration,
		monitorDiscoverySourceConfiguration, logSourceConfiguration,
	}
}

func configurationResources() []func() resource.Resource {
	resources := make([]func() resource.Resource, 0)
	for _, spec := range managedConfigurations() {
		resources = append(resources, func() resource.Resource { return &configurationResource{spec: spec} })
	}
	return resources
}

func configurationDataSources() []func() datasource.DataSource {
	sources := make([]func() datasource.DataSource, 0)
	configurations := append(managedConfigurations(), codeQualityGateConfiguration, codeProfileConfiguration)
	for _, spec := range configurations {
		sources = append(sources,
			func() datasource.DataSource { return &configurationDataSource{spec: spec} },
			func() datasource.DataSource { return &configurationDataSource{spec: spec, list: true} },
		)
	}
	return sources
}
