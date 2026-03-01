---
page_title: "argonix_alert_rule Resource - Argonix"
description: |-
  Manages an Argonix alert rule.
---

# argonix_alert_rule (Resource)

Manages an Argonix alert rule. Alert rules define the conditions under which alerts are triggered and which channels receive notifications.

## Example Usage

```terraform
resource "argonix_alert_rule" "all_down" {
  name              = "Alert on any monitor down"
  trigger_condition = "goes_down"
  all_monitors      = true
  channels          = jsonencode([
    argonix_alert_channel.slack_ops.id,
    argonix_alert_channel.email_ops.id,
  ])
  cooldown_minutes = 10
}
```

### Alert Rule for Specific Monitors

```terraform
resource "argonix_alert_rule" "api_monitors" {
  name              = "API Monitors Alert"
  trigger_condition = "goes_down"
  monitors = jsonencode([
    argonix_monitor.api_health.id,
  ])
  consecutive_failures = 3
  channels = jsonencode([
    argonix_alert_channel.slack_ops.id,
  ])
}
```

## Schema

### Required

- `name` (String) — Name of the alert rule.
- `trigger_condition` (String) — Trigger condition. One of: `status_change`, `goes_down`, `goes_up`, `degraded`, `ssl_expiry`, `test_failing`, `test_passing`.
- `channels` (String) — JSON-encoded list of alert channel UUIDs to notify.

### Optional

- `is_active` (Boolean) — Whether the rule is active. Defaults to `true`.
- `consecutive_failures` (Number) — Number of consecutive failures before triggering. Defaults to `1`.
- `cooldown_minutes` (Number) — Cooldown period in minutes between repeated alerts. Defaults to `5`.
- `all_monitors` (Boolean) — Apply to all monitors in the organization. Defaults to `false`.
- `all_synthetic_tests` (Boolean) — Apply to all synthetic tests in the organization. Defaults to `false`.
- `monitor_tags` (String) — JSON-encoded list of tags to match monitors. Defaults to `"[]"`.
- `monitors` (String) — JSON-encoded list of specific monitor UUIDs. Defaults to `"[]"`.
- `synthetic_tests` (String) — JSON-encoded list of specific synthetic test UUIDs. Defaults to `"[]"`.

### Read-Only

- `id` (String) — UUID of the alert rule.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_alert_rule.example <alert-rule-uuid>
```
