---
page_title: "argonix_cost_perspective Resource - Argonix"
description: |-
  Manages a saved FinOps view.
---

# argonix_cost_perspective (Resource)

Manages a saved FinOps view.

Perspectives save a named set of filters for the FinOps dashboard.

## Example usage

```terraform
resource "argonix_cost_perspective" "platform" {
  name    = "Platform costs"
  filters = jsonencode({ tenant = argonix_finops_tenant.platform.slug })
}
```

## Schema

### Required

- `name` (String) — View name.

### Optional

- `filters` (String, JSON object) — JSON object of saved cost filters.

### Read-only

- `id` (String) — Argonix UUID.
- `date_created` (String)
- `date_updated` (String)

## Import

```shell
terraform import argonix_cost_perspective.example <uuid>
```
