---
page_title: "argonix_synthetic_test Resource - Argonix"
description: |-
  Manages an Argonix synthetic test or Sanity check.
---

# argonix_synthetic_test (Resource)

Manages an Argonix synthetic test or Sanity check.

Use `test_type = "scenario"` for a Sanity check. Existing API and browser tests continue to use `api` and `browser`. The default inputs apply to scheduled runs. Launch individual runs from Argonix.

## Example usage

```terraform
resource "argonix_synthetic_test" "sanity" {
  name              = "Production health check"
  test_type         = "scenario"
  schedule_mode     = "cron"
  schedule          = "*/15 * * * *"
  schedule_timezone = "Europe/Paris"
  run_timeout       = 120
  input_variables   = jsonencode({ endpoint = "https://example.com/health" })

  steps = jsonencode([{
    id     = "health"
    kind   = "http"
    params = { url = "{{input.endpoint}}" }
    assertions = [{ path = "status_code", operator = "eq", value = 200 }]
  }])
}
```

## Schema

### Required

- `name` (String) — Name of the check.
- `test_type` (String) — api, browser or scenario (Sanity check).
- `steps` (String, JSON array) — JSON array of check steps.

### Optional

- `description` (String) — Description of the check. Defaults to `""`.
- `is_active` (Boolean) — Whether the check is active. Defaults to `true`.
- `check_interval` (Number) — Interval between runs in seconds. Defaults to `300`.
- `timeout` (Number) — Per-step timeout in seconds. Defaults to `30`.
- `tags` (String, JSON array) — JSON array of tags. Defaults to `"[]"`.
- `locations` (String, JSON array) — JSON array of region codes. Defaults to `"[]"`.
- `input_variables` (String, JSON object) — JSON object of default Sanity check input variables.
- `run_timeout` (Number) — Maximum duration of a Sanity check in seconds.
- `schedule_mode` (String) — manual, interval or cron.
- `schedule` (String) — Cron expression when schedule_mode is cron.
- `schedule_timezone` (String) — IANA timezone for the cron expression.
- `application_key` (String) — Application UUID associated with the check.
- `resource_id` (String) — Inventory resource UUID associated with the check.

### Read-only

- `id` (String) — Argonix UUID.
- `allow_mutating_actions` (Boolean)
- `steps_count` (Number)
- `result_is_stale` (Boolean)
- `current_status` (String)
- `next_run_at` (String)
- `last_run_at` (String)
- `last_duration_ms` (Number)
- `last_warning_count` (Number)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_synthetic_test.example <uuid>
```
