---
page_title: "argonix_workflow Data Source - Argonix"
description: |-
  Reads an Argonix workflow by UUID.
---

# argonix_workflow (Data Source)

Reads an Argonix workflow by UUID.

## Example usage

```terraform
data "argonix_workflow" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Display name of the workflow.
- `slug` (String) — Unique URL-friendly identifier.
- `description` (String) — Description of the workflow. Defaults to `""`.
- `category` (String) — identity, incident, onboarding, devops, security or general.
- `steps` (String, JSON array) — JSON array of workflow step definitions.
- `input_hints` (String, JSON object) — JSON object describing inputs referenced in the steps. Defaults to `"{}"`.
- `inputs` (String, JSON array)
- `required_connector_types` (String, JSON array) — JSON array of required connector types. Defaults to `"[]"`.
- `requires_confirmation` (Boolean) — Require human confirmation before execution. Defaults to `false`.
- `extract_inputs_with_ai` (Boolean) — Allow AI extraction of inputs from the request.
- `schedule` (String) — Cron expression, or an empty string for no schedule. Defaults to `"", "Cron expression`.
- `schedule_timezone` (String) — IANA timezone for the cron expression.
- `schedule_inputs` (String, JSON object) — JSON object of inputs supplied to scheduled runs.
- `schedule_last_run_at` (String)
- `schedule_next_run_at` (String)
- `security_trigger` (String) — Event that triggers the workflow for defensive triage.
- `security_min_severity` (String) — Minimum security event severity.
- `security_detection_rule` (String) — Optional detection rule ID.
- `incident_trigger` (String) — Optional incident event for investigation.
- `incident_min_severity` (String) — Minimum incident severity.
- `incident_min_confidence` (Number) — Minimum confidence for matching incidents.
- `incident_max_affected_resources` (Number) — Maximum affected resource count for matching incidents.
- `incident_allowed_environments` (String, JSON array) — JSON array of environments allowed for incident investigation.
- `incident_timeout_minutes` (Number) — Incident workflow timeout in minutes.
- `is_active` (Boolean) — Whether the workflow is active. Defaults to `true`.
- `date_created` (String)
- `date_modified` (String)
