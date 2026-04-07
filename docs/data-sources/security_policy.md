---
page_title: "argonix_security_policy Data Source - Argonix Provider"
---

# argonix_security_policy (Data Source)

Fetches a single Argonix security policy by ID.

## Example Usage

```hcl
data "argonix_security_policy" "production" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "policy_rules" {
  value = data.argonix_security_policy.production.rules
}
```

## Argument Reference

- `id` (String, Required) — UUID of the security policy.

## Attribute Reference

- `name` (String) — Display name.
- `description` (String) — Description.
- `rules` (String) — JSON-encoded policy rules.
- `environment` (String) — Target environment.
- `is_active` (Boolean) — Whether the policy is active.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.
