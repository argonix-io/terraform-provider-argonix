---
page_title: "argonix_notification_rule Resource - Argonix Provider"
---

# argonix_notification_rule (Resource)

Manages an Argonix notification rule. Notification rules connect monitors or synthetic tests to notification channels. When a trigger condition is met, notifications are dispatched through the configured channels.

## Example Usage

### Monitor notification rule

```hcl
resource "argonix_notification_rule" "down_alert" {
  name              = "Production Down Alert"
  trigger_condition = "goes_down"
  all_monitors      = true
  channels          = jsonencode([argonix_alert_channel.slack.id])
  auto_investigate   = true
}
```

### Synthetic test notification rule

```hcl
resource "argonix_notification_rule" "test_failing" {
  name                = "API Test Failing"
  trigger_condition   = "test_failing"
  all_synthetic_tests = true
  channels            = jsonencode([argonix_alert_channel.email.id])
}
```

### After each test run (CI/CD)

```hcl
resource "argonix_notification_rule" "ci_webhook" {
  name              = "CI/CD Test Result"
  trigger_condition = "test_run_complete"
  synthetic_tests   = jsonencode([argonix_synthetic_test.checkout_flow.id])
  channels          = jsonencode([argonix_alert_channel.webhook.id])
}
```

## Argument Reference

- `name` (String, Required) — Name of the notification rule.
- `trigger_condition` (String, Required) — One of:
  - **Monitor triggers:** `status_change`, `goes_down`, `goes_up`, `degraded`, `ssl_expiry`
  - **Synthetic test triggers:** `test_failing`, `test_passing`, `test_run_complete`
- `channels` (String, Required) — JSON-encoded list of notification channel UUIDs.
- `is_active` (Boolean, Optional) — Whether the rule is active. Default `true`.
- `all_monitors` (Boolean, Optional) — Apply to all monitors. Default `false`.
- `all_synthetic_tests` (Boolean, Optional) — Apply to all synthetic tests. Default `false`.
- `monitors` (String, Optional) — JSON-encoded list of monitor UUIDs.
- `synthetic_tests` (String, Optional) — JSON-encoded list of synthetic test UUIDs.
- `monitor_tags` (String, Optional) — JSON-encoded list of tags to match monitors.
- `consecutive_failures` (Integer, Optional) — Failures before triggering. Default `1`.
- `cooldown_minutes` (Integer, Optional) — Minimum minutes between repeated notifications. Default `5`.
- `auto_investigate` (Boolean, Optional) — When triggered, Argos AI automatically investigates the root cause and posts analysis to channels. Default `false`.

## Attribute Reference

- `id` (String) — UUID of the notification rule.
- `date_created` (String) — ISO 8601 creation timestamp.
- `date_modified` (String) — ISO 8601 last-modified timestamp.
