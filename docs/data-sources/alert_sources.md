---
page_title: "argonix_alert_sources Data Source - Argonix"
description: |-
  Fetches all alert sources in the organization.
---

# argonix_alert_sources (Data Source)

Fetches all alert sources in the organization.

## Example Usage

```terraform
data "argonix_alert_sources" "all" {}

output "active_sources" {
  value = [for s in data.argonix_alert_sources.all.alert_sources : s.name if s.is_active]
}
```

## Schema

### Read-Only

- `alert_sources` (List of Object) — List of all alert sources. Each has:
  - `id` (String)
  - `name` (String)
  - `source_type` (String)
  - `is_active` (Boolean)
  - `connector` (String)
  - `filters` (String)
  - `auto_investigate` (Boolean)
  - `auto_remediate` (Boolean)
  - `channels` (List of String)
  - `webhook_secret` (String, Sensitive)
  - `webhook_url` (String)
  - `last_received_at` (String)
  - `total_received` (Number)
  - `date_created` (String)
  - `date_modified` (String)
