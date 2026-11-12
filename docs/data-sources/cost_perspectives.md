---
page_title: "argonix_cost_perspectives Data Source - Argonix"
description: |-
  Lists the Argonix cost perspectives visible to the provider's API key.
---

# argonix_cost_perspectives (Data Source)

Lists the Argonix cost perspectives visible to the provider's API key.

## Example usage

```terraform
data "argonix_cost_perspectives" "all" {}

output "cost_perspective_ids" {
  value = [for item in data.argonix_cost_perspectives.all.cost_perspectives : item.id]
}
```

## Schema

- `cost_perspectives` (List of Object) — Objects with the attributes documented for the [single-object data source](cost_perspective.md).
