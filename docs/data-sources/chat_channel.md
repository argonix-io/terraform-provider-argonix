---
page_title: "argonix_chat_channel Data Source - Argonix"
description: |-
  Fetches a single Argonix chat channel by ID.
---

# argonix_chat_channel (Data Source)

Fetches a single Argonix chat channel by ID.

## Example Usage

```terraform
data "argonix_chat_channel" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the chat channel.

### Read-Only

- `channel_type` (String)
- `channel_id` (String)
- `channel_name` (String)
- `persona_id` (String)
- `connector_id` (String)
- `config` (String) — JSON-encoded configuration.
- `is_active` (Boolean)
- `date_created` (String)
- `date_modified` (String)
