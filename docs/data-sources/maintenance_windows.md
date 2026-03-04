---
page_title: "argonix_maintenance_windows Data Source - Argonix"
description: |-
  Fetches all maintenance windows in the organization.
---

# argonix_maintenance_windows (Data Source)

Fetches all maintenance windows in the organization.

## Example Usage

```terraform
data "argonix_maintenance_windows" "all" {}
```

## Schema

### Read-Only

- `maintenance_windows` (List of Object) — List of all maintenance windows. Each has:
  - `id` (String)
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
