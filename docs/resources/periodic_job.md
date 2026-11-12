---
page_title: "argonix_periodic_job Resource - Argonix"
description: |-
  Manages an Argonix periodic job.
---

# argonix_periodic_job (Resource)

Manages an Argonix periodic job.

Agent jobs require a prompt and use the API's default discussion mode. `application_reconcile` jobs are inactive by default. Set `is_active = true` to schedule an application reconciliation.

## Example usage

```terraform
resource "argonix_periodic_job" "summary" {
  name          = "Daily operational summary"
  kind          = "agent"
  prompt        = "Summarize recent service alerts and flag recurring issues."
  connector_ids = [var.observability_connector_id]
  schedule_hour = 8
  schedule_days = jsonencode([0, 1, 2, 3, 4])
}
```

## Schema

### Required

- `name` (String) — Name of the periodic job.

### Optional

- `kind` (String) — agent or application_reconcile. Changing this value replaces the resource.
- `prompt` (String) — Prompt required for an agent job.
- `connector_ids` (Set of String) — Connector IDs available to the job.
- `application_root_ids` (Set of String) — Application IDs available to the job.
- `schedule_hour` (Number) — UTC hour of execution (0–23).
- `schedule_minute` (Number) — UTC minute of execution (0–59).
- `schedule_days` (String, JSON array) — JSON array of weekday numbers, with Monday as 0. Empty means every day.
- `notification_channel_ids` (Set of String) — Chat channel IDs receiving notifications.
- `alert_channel_ids` (Set of String) — Alert channel IDs receiving results.
- `auto_create_ticket` (Boolean) — Create a ticket through a configured matching connector and alert channel.
- `deliverable_retention_days` (Number) — Number of days to retain deliverables.
- `deliverable_retention_action` (String) — delete or move expired deliverables.
- `deliverable_retention_destination` (String) — Destination when retention action is move.
- `is_active` (Boolean) — Whether the schedule is active.
- `timeout_seconds` (Number) — Maximum execution time in seconds.

### Read-only

- `id` (String) — Argonix UUID.
- `mode` (String)
- `latest_run_status` (String)
- `latest_run_date` (String)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_periodic_job.example <uuid>
```
