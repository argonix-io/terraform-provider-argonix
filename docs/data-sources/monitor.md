---
page_title: "argonix_monitor Data Source - Argonix"
description: |-
  Fetches a single Argonix monitor by ID.
---

# argonix_monitor (Data Source)

Fetches a single Argonix monitor by ID.

## Example Usage

```terraform
data "argonix_monitor" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "monitor_name" {
  value = data.argonix_monitor.example.name
}
```

## Schema

### Required

- `id` (String) — UUID of the monitor.

### Read-Only

- `name` (String) — Display name of the monitor.
- `monitor_type` (String) — Type of monitor.
- `is_active` (Boolean) — Whether the monitor is active.
- `url` (String) — URL being monitored.
- `hostname` (String) — Hostname being monitored.
- `port` (Number) — Port being monitored.
- `check_interval` (Number) — Check interval in seconds.
- `timeout` (Number) — Request timeout in seconds.
- `retries` (Number) — Number of retries.
- `http_method` (String) — HTTP method.
- `current_status` (String) — Current status.
- `group_id` (String) — UUID of the group.
- `tags` (String) — JSON-encoded tags.
- `regions` (String) — JSON-encoded regions.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.
