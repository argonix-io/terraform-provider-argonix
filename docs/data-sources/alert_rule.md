---
page_title: "argonix_alert_rule Data Source - Argonix"
description: |-
  Fetches a single Argonix alert rule by ID.
---

# argonix_alert_rule (Data Source)

Fetches a single Argonix alert rule by ID.

## Example Usage

```terraform
data "argonix_alert_rule" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the alert rule.

### Read-Only

- `name` (String)
- `is_active` (Boolean)
- `trigger_condition` (String)
- `consecutive_failures` (Number)
- `cooldown_minutes` (Number)
- `all_monitors` (Boolean)
- `all_synthetic_tests` (Boolean)
- `monitor_tags` (String) — JSON-encoded list of tags.
- `monitors` (String) — JSON-encoded list of monitor UUIDs.
- `synthetic_tests` (String) — JSON-encoded list of synthetic test UUIDs.
- `channels` (String) — JSON-encoded list of channel UUIDs.
- `date_created` (String)
- `date_modified` (String)
