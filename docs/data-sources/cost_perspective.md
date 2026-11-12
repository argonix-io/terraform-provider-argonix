---
page_title: "argonix_cost_perspective Data Source - Argonix"
description: |-
  Reads an Argonix cost perspective by UUID.
---

# argonix_cost_perspective (Data Source)

Reads an Argonix cost perspective by UUID.

## Example usage

```terraform
data "argonix_cost_perspective" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — View name.
- `filters` (String, JSON object) — JSON object of saved cost filters.
- `date_created` (String)
- `date_updated` (String)
