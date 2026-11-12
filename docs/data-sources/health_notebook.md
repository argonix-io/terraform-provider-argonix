---
page_title: "argonix_health_notebook Data Source - Argonix"
description: |-
  Reads an Argonix health notebook by UUID.
---

# argonix_health_notebook (Data Source)

Reads an Argonix health notebook by UUID.

## Example usage

```terraform
data "argonix_health_notebook" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Name of the notebook.
- `report_kind` (String) — health or security. Changing this value replaces the resource.
- `connector_ids` (Set of String) — Connector IDs in scope.
- `use_applications` (Boolean) — Scope the notebook to applications.
- `perspective` (String) — all, business or platform.
- `application_root_ids` (Set of String) — Selected application IDs.
- `application_label_selector` (String, JSON array) — JSON array of exact key and value pairs for detected applications.
- `include_unattributed` (Boolean) — Include findings without application attribution.
- `minimum_attribution_confidence` (Number) — Minimum attribution confidence from 0 to 100.
- `is_active` (Boolean) — Whether the schedule is active.
- `frequency` (String) — daily, weekly or monthly.
- `schedule_hour` (Number) — UTC hour of execution (0–23).
- `schedule_minute` (Number) — UTC minute of execution (0–59).
- `schedule_days` (String, JSON array) — JSON array of weekday numbers, with Monday as 0.
- `schedule_day_of_month` (Number) — Day of month for monthly notebooks.
- `custom_prompt` (String) — Additional instructions for the report.
- `lookback_hours` (Number) — History window from 1 to 168 hours.
- `include_threats` (Boolean) — Include detected threats.
- `include_findings` (Boolean) — Include security findings.
- `include_monitors` (Boolean) — Include monitor results.
- `include_incidents` (Boolean) — Include incidents.
- `include_cost` (Boolean) — Include cost signals.
- `include_log_anomalies` (Boolean) — Include log anomalies.
- `include_vulnerabilities` (Boolean) — Include vulnerabilities.
- `include_runtime` (Boolean) — Include runtime signals.
- `include_dns` (Boolean) — Include DNS signals.
- `include_network` (Boolean) — Include network signals.
- `include_operational_signals` (Boolean) — Include application operational signals.
- `notification_channel_ids` (Set of String) — Chat channel IDs receiving reports.
- `alert_channel_ids` (Set of String) — Alert channel IDs receiving reports.
- `ticketing_min_severity` (String) — Minimum severity for ticket creation.
- `ticketing_action` (String)
- `latest_score` (Number)
- `latest_severity` (String)
- `latest_evidence_coverage` (String, JSON object)
- `date_created` (String)
- `date_modified` (String)
