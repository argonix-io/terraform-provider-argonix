---
page_title: "argonix_workflow Resource - Argonix"
description: |-
  Manages an Argonix Argos workflow.
---

# argonix_workflow (Resource)

Manages an Argonix Argos workflow.

Set a unique `slug` when creating a workflow. Use `schedule_inputs` for required inputs referenced by scheduled steps. `input_hints` describes the inputs discovered from `steps`; it does not declare them.

## Example usage

```terraform
resource "argonix_workflow" "weekday_review" {
  name                = "Weekday review"
  slug                = "weekday-review"
  category            = "general"
  requires_confirmation = true
  schedule            = "0 9 * * 1-5"
  schedule_timezone   = "Europe/Paris"
  steps               = jsonencode(var.review_steps)
  schedule_inputs     = jsonencode({ environment = "production" })
  input_hints = jsonencode({
    environment = { description = "Environment to inspect", example = "production" }
  })
}
```

## Schema

### Required

- `name` (String) — Display name of the workflow.
- `category` (String) — identity, incident, onboarding, devops, security or general.
- `steps` (String, JSON array) — JSON array of workflow step definitions.

### Optional

- `slug` (String) — Unique URL-friendly identifier.
- `description` (String) — Description of the workflow. Defaults to `""`.
- `input_hints` (String, JSON object) — JSON object describing inputs referenced in the steps. Defaults to `"{}"`.
- `required_connector_types` (String, JSON array) — JSON array of required connector types. Defaults to `"[]"`.
- `requires_confirmation` (Boolean) — Require human confirmation before execution. Defaults to `false`.
- `extract_inputs_with_ai` (Boolean) — Allow AI extraction of inputs from the request.
- `schedule` (String) — Cron expression, or an empty string for no schedule. Defaults to `"", "Cron expression`.
- `schedule_timezone` (String) — IANA timezone for the cron expression.
- `schedule_inputs` (String, JSON object) — JSON object of inputs supplied to scheduled runs.
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

### Read-only

- `id` (String) — Argonix UUID.
- `inputs` (String, JSON array)
- `schedule_last_run_at` (String)
- `schedule_next_run_at` (String)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_workflow.example <uuid>
```
