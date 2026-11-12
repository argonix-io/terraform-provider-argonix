---
page_title: "argonix_cost_budget Data Source - Argonix"
description: |-
  Reads an Argonix cost budget by UUID.
---

# argonix_cost_budget (Data Source)

Reads an Argonix cost budget by UUID.

## Example usage

```terraform
data "argonix_cost_budget" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Budget name.
- `period` (String) — monthly or quarterly.
- `currency` (String) — Three-letter uppercase currency code.
- `filters` (String, JSON object) — JSON object of cost filters. An empty object covers the organization.
- `alert_thresholds_pct` (String, JSON array) — JSON array of unique percentages in ascending order.
- `alert_channels` (Set of String) — Alert channel IDs receiving budget alerts.
- `is_active` (Boolean) — Whether the budget is active.
- `recent_alerts` (String, JSON array)
- `date_created` (String)
- `date_updated` (String)
- `amount` (Number) — Budget amount in currency units, with up to six decimal places.
