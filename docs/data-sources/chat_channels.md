---
page_title: "argonix_chat_channels Data Source - Argonix"
description: |-
  Fetches all chat channels in the organization.
---

# argonix_chat_channels (Data Source)

Fetches all chat channels in the organization.

## Example Usage

```terraform
data "argonix_chat_channels" "all" {}
```

## Schema

### Read-Only

- `chat_channels` (List of Object) — List of all chat channels. Each has:
  - `id` (String)
  - `channel_type` (String)
  - `channel_id` (String)
  - `channel_name` (String)
  - `persona_id` (String)
  - `connector_id` (String)
  - `config` (String) — JSON-encoded configuration.
  - `is_active` (Boolean)
  - `date_created` (String)
  - `date_modified` (String)
