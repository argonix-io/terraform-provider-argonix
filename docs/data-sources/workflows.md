---
page_title: "argonix_workflows Data Source - Argonix"
description: |-
  Lists the Argonix workflows visible to the provider's API key.
---

# argonix_workflows (Data Source)

Lists the Argonix workflows visible to the provider's API key.

## Example usage

```terraform
data "argonix_workflows" "all" {}

output "workflow_ids" {
  value = [for item in data.argonix_workflows.all.workflows : item.id]
}
```

## Schema

- `workflows` (List of Object) — Objects with the attributes documented for the [single-object data source](workflow.md).
