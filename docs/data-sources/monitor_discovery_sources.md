---
page_title: "argonix_monitor_discovery_sources Data Source - Argonix"
description: |-
  Lists the Argonix monitor discovery sources visible to the provider's API key.
---

# argonix_monitor_discovery_sources (Data Source)

Lists the Argonix monitor discovery sources visible to the provider's API key.

## Example usage

```terraform
data "argonix_monitor_discovery_sources" "all" {}

output "monitor_discovery_source_ids" {
  value = [for item in data.argonix_monitor_discovery_sources.all.monitor_discovery_sources : item.id]
}
```

## Schema

- `monitor_discovery_sources` (List of Object) — Objects with the attributes documented for the [single-object data source](monitor_discovery_source.md).
