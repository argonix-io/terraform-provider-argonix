---
page_title: "argonix_cost_budgets Data Source - Argonix"
description: |-
  Lists the Argonix cost budgets visible to the provider's API key.
---

# argonix_cost_budgets (Data Source)

Lists the Argonix cost budgets visible to the provider's API key.

## Example usage

```terraform
data "argonix_cost_budgets" "all" {}

output "cost_budget_ids" {
  value = [for item in data.argonix_cost_budgets.all.cost_budgets : item.id]
}
```

## Schema

- `cost_budgets` (List of Object) — Objects with the attributes documented for the [single-object data source](cost_budget.md).
