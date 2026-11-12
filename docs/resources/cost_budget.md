---
page_title: "argonix_cost_budget Resource - Argonix"
description: |-
  Manages a FinOps cost budget.
---

# argonix_cost_budget (Resource)

Manages a FinOps cost budget.

Set a budget in currency units and scope it with `filters`. Alert thresholds must be sorted and unique.

## Example usage

```terraform
resource "argonix_cost_budget" "platform" {
  name                 = "Platform budget"
  amount               = 1500
  currency             = "EUR"
  period               = "monthly"
  filters              = jsonencode({ tenant = argonix_finops_tenant.platform.slug })
  alert_thresholds_pct = jsonencode([80, 100])
}
```

## Schema

### Required

- `name` (String) — Budget name.
- `amount` (Number) — Budget amount in currency units, with up to six decimal places.

### Optional

- `period` (String) — monthly or quarterly.
- `currency` (String) — Three-letter uppercase currency code.
- `filters` (String, JSON object) — JSON object of cost filters. An empty object covers the organization.
- `alert_thresholds_pct` (String, JSON array) — JSON array of unique percentages in ascending order.
- `alert_channels` (Set of String) — Alert channel IDs receiving budget alerts.
- `is_active` (Boolean) — Whether the budget is active.

### Read-only

- `id` (String) — Argonix UUID.
- `recent_alerts` (String, JSON array)
- `date_created` (String)
- `date_updated` (String)

## Import

```shell
terraform import argonix_cost_budget.example <uuid>
```
