---
page_title: "argonix_finops_tenant Data Source - Argonix"
description: |-
  Reads an Argonix finops tenant by UUID.
---

# argonix_finops_tenant (Data Source)

Reads an Argonix finops tenant by UUID.

## Example usage

```terraform
data "argonix_finops_tenant" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Tenant name.
- `color` (String) — Display color.
- `parent` (String) — Parent tenant UUID.
- `monthly_budget_micros` (Number) — Monthly budget in millionths of a currency unit.
- `contact_email` (String) — Tenant contact email address.
- `description` (String) — Tenant description.
- `slug` (String)
- `is_system` (Boolean)
- `rules` (String, JSON array)
- `date_created` (String)
- `date_modified` (String)
- `kind` (String) — business_unit, client, team, product or environment.
