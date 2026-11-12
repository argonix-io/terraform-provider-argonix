---
page_title: "argonix_scan_targets Data Source - Argonix"
description: |-
  Lists the Argonix scan targets visible to the provider's API key.
---

# argonix_scan_targets (Data Source)

Lists the Argonix scan targets visible to the provider's API key.

## Example usage

```terraform
data "argonix_scan_targets" "all" {}

output "scan_target_ids" {
  value = [for item in data.argonix_scan_targets.all.scan_targets : item.id]
}
```

## Schema

- `scan_targets` (List of Object) — Objects with the attributes documented for the [single-object data source](scan_target.md).
