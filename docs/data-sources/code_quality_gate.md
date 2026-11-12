---
page_title: "argonix_code_quality_gate Data Source - Argonix"
description: |-
  Reads an Argonix code quality gate by UUID.
---

# argonix_code_quality_gate (Data Source)

Reads an Argonix code quality gate by UUID.

Read an active Argonix quality gate and use its ID in an `argonix_code_project` resource. Gate versioning is managed in Argonix.

## Example usage

```terraform
data "argonix_code_quality_gate" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `organization` (String)
- `name` (String)
- `description` (String)
- `version` (Number)
- `is_default` (Boolean)
- `is_active` (Boolean)
- `conditions` (String, JSON array)
- `date_created` (String)
- `date_modified` (String)
