---
page_title: "argonix_finops_tenant Resource - Argonix"
description: |-
  Manages a FinOps tenant for cost allocation.
---

# argonix_finops_tenant (Resource)

Manages a FinOps tenant for cost allocation.

Creating, editing, or deleting a tenant requires an organization admin. Read the tenant's allocation rules in `rules`.

## Example usage

```terraform
resource "argonix_finops_tenant" "platform" {
  name = "Platform"
  kind = "team"
}
```

## Schema

### Required

- `name` (String) — Tenant name.

### Optional

- `color` (String) — Display color.
- `parent` (String) — Parent tenant UUID.
- `monthly_budget_micros` (Number) — Monthly budget in millionths of a currency unit.
- `contact_email` (String) — Tenant contact email address.
- `description` (String) — Tenant description.
- `kind` (String) — business_unit, client, team, product or environment.

### Read-only

- `id` (String) — Argonix UUID.
- `slug` (String)
- `is_system` (Boolean)
- `rules` (String, JSON array)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_finops_tenant.example <uuid>
```
