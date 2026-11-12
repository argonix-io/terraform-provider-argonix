---
page_title: "argonix_finops_tenants Data Source - Argonix"
description: |-
  Lists the Argonix finops tenants visible to the provider's API key.
---

# argonix_finops_tenants (Data Source)

Lists the Argonix finops tenants visible to the provider's API key.

## Example usage

```terraform
data "argonix_finops_tenants" "all" {}

output "finops_tenant_ids" {
  value = [for item in data.argonix_finops_tenants.all.finops_tenants : item.id]
}
```

## Schema

- `finops_tenants` (List of Object) — Objects with the attributes documented for the [single-object data source](finops_tenant.md).
