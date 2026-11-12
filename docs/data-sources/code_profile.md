---
page_title: "argonix_code_profile Data Source - Argonix"
description: |-
  Reads an Argonix code profile by UUID.
---

# argonix_code_profile (Data Source)

Reads an Argonix code profile by UUID.

Read an Argonix quality profile and its rules. Profile cloning and versioning are managed in Argonix.

## Example usage

```terraform
data "argonix_code_profile" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String)
- `language` (String)
- `version` (Number)
- `is_active` (Boolean)
- `is_default` (Boolean)
- `is_builtin` (Boolean)
- `scanner_defaults` (String, JSON object)
- `rules` (String, JSON array)
- `date_created` (String)
