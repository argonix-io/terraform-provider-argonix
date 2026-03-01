---
page_title: "argonix_monitors Data Source - Argonix"
description: |-
  Fetches all monitors in the organization.
---

# argonix_monitors (Data Source)

Fetches all monitors in the organization.

## Example Usage

```terraform
data "argonix_monitors" "all" {}

output "monitor_count" {
  value = length(data.argonix_monitors.all.monitors)
}
```

## Schema

### Read-Only

- `monitors` (List of Object) — List of all monitors. Each monitor has the following attributes:
  - `id` (String)
  - `name` (String)
  - `monitor_type` (String)
  - `is_active` (Boolean)
  - `url` (String)
  - `hostname` (String)
  - `port` (Number)
  - `check_interval` (Number)
  - `timeout` (Number)
  - `retries` (Number)
  - `http_method` (String)
  - `current_status` (String)
  - `group_id` (String)
  - `tags` (String) — JSON-encoded tags.
  - `regions` (String) — JSON-encoded regions.
  - `date_created` (String)
  - `date_modified` (String)
