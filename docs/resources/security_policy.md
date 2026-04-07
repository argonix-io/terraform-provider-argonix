---
page_title: "argonix_security_policy Resource - Argonix Provider"
---

# argonix_security_policy (Resource)

Manages an Argonix security policy (deployment gate). Security policies define quality gates that must pass before a deployment is allowed, based on security scan results.

## Example Usage

### Production deployment gate

```hcl
resource "argonix_security_policy" "production" {
  name        = "Production Gate"
  environment = "production"
  rules = jsonencode({
    max_critical         = 0
    max_high             = 0
    scan_max_age_hours   = 24
    required_scan_types  = ["trivy", "gitleaks", "semgrep"]
  })
}
```

### Staging (more permissive)

```hcl
resource "argonix_security_policy" "staging" {
  name        = "Staging Gate"
  environment = "staging"
  rules = jsonencode({
    max_critical        = 0
    max_high            = 10
    scan_max_age_hours  = 72
    required_scan_types = ["trivy"]
  })
}
```

## Argument Reference

- `name` (String, Required) — Display name of the security policy.
- `rules` (String, Required) — JSON-encoded policy rules object. Supported keys:
  - `max_critical` — Maximum number of critical findings allowed (integer).
  - `max_high` — Maximum number of high findings allowed (integer).
  - `scan_max_age_hours` — Maximum age of the latest scan in hours (integer).
  - `required_scan_types` — List of scan tools that must have run (e.g. `["trivy", "gitleaks"]`).
- `description` (String, Optional) — Description of the policy purpose. Default: `""`.
- `environment` (String, Optional) — Target environment (e.g. `production`, `staging`). Empty applies to all environments. Default: `""`.
- `is_active` (Boolean, Optional) — Whether the policy is active. Default: `true`.

## Attribute Reference

- `id` (String) — UUID of the security policy.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

Security policies can be imported using their UUID:

```shell
terraform import argonix_security_policy.production <uuid>
```
