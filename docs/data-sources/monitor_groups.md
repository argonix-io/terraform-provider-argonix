---
page_title: "argonix_monitor_groups Data Source - Argonix"
description: |-
  Lists the Argonix monitor groups visible to the provider's API key.
---

# argonix_monitor_groups (Data Source)

Lists the Argonix monitor groups visible to the provider's API key.

## Example usage

```terraform
data "argonix_monitor_groups" "all" {}

output "monitor_group_ids" {
  value = [for item in data.argonix_monitor_groups.all.monitor_groups : item.id]
}
```

## Schema

- `monitor_groups` (List of Object) — Objects with the attributes documented for the [single-object data source](monitor_group.md).
