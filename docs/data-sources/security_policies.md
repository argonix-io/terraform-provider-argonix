---
page_title: "argonix_security_policies Data Source - Argonix Provider"
---

# argonix_security_policies (Data Source)

Fetches all security policies in the organization.

## Example Usage

```hcl
data "argonix_security_policies" "all" {}

output "active_policies" {
  value = [for p in data.argonix_security_policies.all.security_policies : p.name if p.is_active]
}
```

## Attribute Reference

- `security_policies` (List) — List of security policy objects, each containing:
  - `id` (String) — UUID.
  - `name` (String) — Display name.
  - `description` (String) — Description.
  - `rules` (String) — JSON-encoded policy rules.
  - `environment` (String) — Target environment.
  - `is_active` (Boolean) — Whether the policy is active.
  - `date_created` (String) — Creation timestamp.
  - `date_modified` (String) — Last modification timestamp.
