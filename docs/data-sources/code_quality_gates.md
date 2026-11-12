---
page_title: "argonix_code_quality_gates Data Source - Argonix"
description: |-
  Lists the Argonix code quality gates visible to the provider's API key.
---

# argonix_code_quality_gates (Data Source)

Lists the Argonix code quality gates visible to the provider's API key.

## Example usage

```terraform
data "argonix_code_quality_gates" "all" {}

output "code_quality_gate_ids" {
  value = [for item in data.argonix_code_quality_gates.all.code_quality_gates : item.id]
}
```

## Schema

- `code_quality_gates` (List of Object) — Objects with the attributes documented for the [single-object data source](code_quality_gate.md).
