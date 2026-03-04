---
page_title: "argonix_chat_channel Resource - Argonix"
description: |-
  Manages an Argonix Argos chat channel.
---

# argonix_chat_channel (Resource)

Manages an Argonix Argos chat channel. Chat channels connect AI personas to messaging platforms (Slack, Teams, Jira) for interactive incident management.

## Example Usage

```terraform
resource "argonix_chat_channel" "ops_slack" {
  channel_type = "slack"
  channel_id   = "C01234567"
  channel_name = "#ops-incidents"
  persona_id   = argonix_persona.devops.id
  connector_id = argonix_connector.slack.id
}

resource "argonix_chat_channel" "teams" {
  channel_type = "teams"
  channel_id   = "19:abc123@thread.tacv2"
  channel_name = "Incidents"
  persona_id   = argonix_persona.devops.id
  connector_id = argonix_connector.teams.id
}
```

## Schema

### Required

- `channel_type` (String) — Type of channel. One of: `slack`, `teams`, `jira`.
- `channel_id` (String) — External channel identifier.
- `connector_id` (String) — UUID of the connector for this channel.

### Optional

- `channel_name` (String) — Display name for the channel. Defaults to `""`.
- `persona_id` (String) — UUID of the persona assigned to this channel. Defaults to `""`.
- `config` (String) — JSON-encoded additional configuration. Defaults to `"{}"`.
- `is_active` (Boolean) — Whether the channel is active. Defaults to `true`.

### Read-Only

- `id` (String) — UUID of the chat channel.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_chat_channel.example <chat-channel-uuid>
```
