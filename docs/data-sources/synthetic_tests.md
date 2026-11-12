---
page_title: "argonix_synthetic_tests Data Source - Argonix"
description: |-
  Lists the Argonix synthetic tests visible to the provider's API key.
---

# argonix_synthetic_tests (Data Source)

Lists the Argonix synthetic tests visible to the provider's API key.

## Example usage

```terraform
data "argonix_synthetic_tests" "all" {}

output "synthetic_test_ids" {
  value = [for item in data.argonix_synthetic_tests.all.synthetic_tests : item.id]
}
```

## Schema

- `synthetic_tests` (List of Object) — Objects with the attributes documented for the [single-object data source](synthetic_test.md).
