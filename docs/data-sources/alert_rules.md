---
page_title: "argonix_alert_rules Data Source - Argonix"
description: |-
  Fetches all alert rules in the organization.
---

# argonix_alert_rules (Data Source)

Fetches all alert rules in the organization.

## Example Usage

```terraform
data "argonix_alert_rules" "all" {}
```

## Schema

### Read-Only

- `alert_rules` (List of Object) — List of all alert rules. Each has:
  - `id` (String)
  - `name` (String)
  - `is_active` (Boolean)
  - `trigger_condition` (String)
  - `consecutive_failures` (Number)
  - `cooldown_minutes` (Number)
  - `all_monitors` (Boolean)
  - `all_synthetic_tests` (Boolean)
  - `monitor_tags` (String) — JSON-encoded.
  - `monitors` (String) — JSON-encoded.
  - `synthetic_tests` (String) — JSON-encoded.
  - `channels` (String) — JSON-encoded.
  - `date_created` (String)
  - `date_modified` (String)
