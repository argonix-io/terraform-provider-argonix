---
page_title: "argonix_alert_channel Data Source - Argonix"
description: |-
  Fetches a single Argonix alert channel by ID.
---

# argonix_alert_channel (Data Source)

Fetches a single Argonix alert channel by ID.

## Example Usage

```terraform
data "argonix_alert_channel" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the alert channel.

### Read-Only

- `name` (String)
- `channel_type` (String)
- `is_active` (Boolean)
- `config` (String, Sensitive) — JSON-encoded channel configuration.
- `date_created` (String)
- `date_modified` (String)
