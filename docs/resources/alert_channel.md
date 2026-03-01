---
page_title: "argonix_alert_channel Resource - Argonix"
description: |-
  Manages an Argonix alert channel (notification destination).
---

# argonix_alert_channel (Resource)

Manages an Argonix alert channel. Alert channels are notification destinations used by alert rules (e.g., Slack, email, PagerDuty).

## Example Usage

### Slack Channel

```terraform
resource "argonix_alert_channel" "slack_ops" {
  name         = "Slack Ops Channel"
  channel_type = "slack"
  config = jsonencode({
    webhook_url = "https://hooks.slack.com/services/xxx/yyy/zzz"
  })
}
```

### Email Channel

```terraform
resource "argonix_alert_channel" "email_ops" {
  name         = "Email Ops Team"
  channel_type = "email"
  config = jsonencode({
    addresses = ["ops@example.com", "oncall@example.com"]
  })
}
```

### Webhook Channel

```terraform
resource "argonix_alert_channel" "webhook" {
  name         = "Custom Webhook"
  channel_type = "webhook"
  config = jsonencode({
    url     = "https://example.com/webhook"
    method  = "POST"
    headers = { "X-Custom" = "value" }
  })
}
```

## Schema

### Required

- `name` (String) — Name of the alert channel.
- `channel_type` (String) — Channel type. One of: `email`, `slack`, `webhook`, `pagerduty`, `opsgenie`, `telegram`, `discord`, `teams`, `jira`.
- `config` (String, Sensitive) — JSON-encoded channel configuration. The structure depends on the `channel_type`.

### Optional

- `is_active` (Boolean) — Whether the channel is active. Defaults to `true`.

### Read-Only

- `id` (String) — UUID of the alert channel.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_alert_channel.example <alert-channel-uuid>
```
