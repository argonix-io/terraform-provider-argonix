---
page_title: "argonix_sre_patrols Data Source - Argonix"
description: |-
  Lists the Argonix sre patrols visible to the provider's API key.
---

# argonix_sre_patrols (Data Source)

Lists the Argonix sre patrols visible to the provider's API key.

## Example usage

```terraform
data "argonix_sre_patrols" "all" {}

output "sre_patrol_ids" {
  value = [for item in data.argonix_sre_patrols.all.sre_patrols : item.id]
}
```

## Schema

- `sre_patrols` (List of Object) — Objects with the attributes documented for the [single-object data source](sre_patrol.md).
