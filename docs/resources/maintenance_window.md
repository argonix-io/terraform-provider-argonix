---
page_title: "argonix_maintenance_window Resource - Argonix"
description: |-
  Manages an Argonix maintenance window.
---

# argonix_maintenance_window (Resource)

Manages an Argonix maintenance window. Maintenance windows suppress alerts during planned maintenance periods.

## Example Usage

```terraform
resource "argonix_maintenance_window" "weekly" {
  name      = "Weekly Deploy Window"
  group_id  = argonix_group.production.id
  repeat    = "weekly"
  time_from = "02:00"
  time_to   = "04:00"
  weekdays  = "1,3,5"
}

resource "argonix_maintenance_window" "one_time" {
  name      = "Database Migration"
  group_id  = argonix_group.production.id
  repeat    = "once"
  starts_at = "2025-01-15T02:00:00Z"
  ends_at   = "2025-01-15T06:00:00Z"
}

resource "argonix_maintenance_window" "cron" {
  name            = "Nightly Backup"
  group_id        = argonix_group.production.id
  repeat          = "cron"
  cron_expression = "0 3 * * *"
}
```

## Schema

### Required

- `name` (String) — Name of the maintenance window.
- `group_id` (String) — UUID of the group to apply the window to.

### Optional

- `starts_at` (String) — Start time (ISO 8601). Used with `once` repeat.
- `ends_at` (String) — End time (ISO 8601). Used with `once` repeat.
- `repeat` (String) — Recurrence type: `once`, `daily`, `weekly`, `monthly`, `cron`. Defaults to `"once"`.
- `time_from` (String) — Daily start time (HH:MM). Used with recurring schedules.
- `time_to` (String) — Daily end time (HH:MM). Used with recurring schedules.
- `weekdays` (String) — Comma-separated weekday numbers (1=Mon, 7=Sun). Used with `weekly`.
- `day_of_month` (Number) — Day of month (1-31). Used with `monthly`.
- `cron_expression` (String) — Cron expression. Used with `cron` repeat.
- `is_active` (Boolean) — Whether the window is active. Defaults to `true`.

### Read-Only

- `id` (String) — UUID of the maintenance window.
- `schedule_summary` (String) — Human-readable schedule summary.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_maintenance_window.example <maintenance-window-uuid>
```
