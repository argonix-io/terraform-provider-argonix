---
page_title: "argonix_synthetic_test Data Source - Argonix"
description: |-
  Reads an Argonix synthetic test by UUID.
---

# argonix_synthetic_test (Data Source)

Reads an Argonix synthetic test by UUID.

## Example usage

```terraform
data "argonix_synthetic_test" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Name of the check.
- `description` (String) — Description of the check. Defaults to `""`.
- `is_active` (Boolean) — Whether the check is active. Defaults to `true`.
- `test_type` (String) — api, browser or scenario (Sanity check).
- `steps` (String, JSON array) — JSON array of check steps.
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
- `resource_id` (String) — Inventory resource UUID associated with the check.
