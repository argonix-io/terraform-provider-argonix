---
page_title: "argonix_health_notebooks Data Source - Argonix"
description: |-
  Lists the Argonix health notebooks visible to the provider's API key.
---

# argonix_health_notebooks (Data Source)

Lists the Argonix health notebooks visible to the provider's API key.

## Example usage

```terraform
data "argonix_health_notebooks" "all" {}

output "health_notebook_ids" {
  value = [for item in data.argonix_health_notebooks.all.health_notebooks : item.id]
}
```

## Schema

- `health_notebooks` (List of Object) — Objects with the attributes documented for the [single-object data source](health_notebook.md).
