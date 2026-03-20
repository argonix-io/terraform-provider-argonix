---
page_title: "argonix_alert_source Data Source - Argonix"
description: |-
  Fetches a single Argonix alert source by ID.
---

# argonix_alert_source (Data Source)

Fetches a single Argonix alert source by ID.

## Example Usage

```terraform
data "argonix_alert_source" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "webhook_url" {
  value = data.argonix_alert_source.example.webhook_url
}
```

## Schema

### Required

- `id` (String) — UUID of the alert source.

### Read-Only

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
