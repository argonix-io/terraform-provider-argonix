---
page_title: "argonix_alert_channels Data Source - Argonix"
description: |-
  Fetches all alert channels in the organization.
---

# argonix_alert_channels (Data Source)

Fetches all alert channels in the organization.

## Example Usage

```terraform
data "argonix_alert_channels" "all" {}
```

## Schema

### Read-Only

- `alert_channels` (List of Object) — List of all alert channels. Each has:
  - `id` (String)
  - `name` (String)
  - `channel_type` (String)
  - `is_active` (Boolean)
  - `config` (String, Sensitive) — JSON-encoded configuration.
  - `date_created` (String)
  - `date_modified` (String)
