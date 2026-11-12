---
page_title: "argonix_code_profiles Data Source - Argonix"
description: |-
  Lists the Argonix code profiles visible to the provider's API key.
---

# argonix_code_profiles (Data Source)

Lists the Argonix code profiles visible to the provider's API key.

## Example usage

```terraform
data "argonix_code_profiles" "all" {}

output "code_profile_ids" {
  value = [for item in data.argonix_code_profiles.all.code_profiles : item.id]
}
```

## Schema

- `code_profiles` (List of Object) — Objects with the attributes documented for the [single-object data source](code_profile.md).
