---
page_title: "argonix_maintenance_window Data Source - Argonix"
description: |-
  Fetches a single Argonix maintenance window by ID.
---

# argonix_maintenance_window (Data Source)

Fetches a single Argonix maintenance window by ID.

## Example Usage

```terraform
data "argonix_maintenance_window" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the maintenance window.

### Read-Only

- `name` (String)
- `group_id` (String)
- `starts_at` (String)
- `ends_at` (String)
- `repeat` (String)
- `time_from` (String)
- `time_to` (String)
- `weekdays` (String)
- `day_of_month` (Number)
- `cron_expression` (String)
- `is_active` (Boolean)
- `schedule_summary` (String)
- `date_created` (String)
- `date_modified` (String)
