---
page_title: "argonix_log_sources Data Source - Argonix"
description: |-
  Lists the Argonix log sources visible to the provider's API key.
---

# argonix_log_sources (Data Source)

Lists the Argonix log sources visible to the provider's API key.

## Example usage

```terraform
data "argonix_log_sources" "all" {}

output "log_source_ids" {
  value = [for item in data.argonix_log_sources.all.log_sources : item.id]
}
```

## Schema

- `log_sources` (List of Object) — Objects with the attributes documented for the [single-object data source](log_source.md).
