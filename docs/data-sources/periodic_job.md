---
page_title: "argonix_periodic_job Data Source - Argonix"
description: |-
  Reads an Argonix periodic job by UUID.
---

# argonix_periodic_job (Data Source)

Reads an Argonix periodic job by UUID.

## Example usage

```terraform
data "argonix_periodic_job" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Name of the periodic job.
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
- `mode` (String)
- `timeout_seconds` (Number) — Maximum execution time in seconds.
- `latest_run_status` (String)
- `latest_run_date` (String)
- `date_created` (String)
- `date_modified` (String)
