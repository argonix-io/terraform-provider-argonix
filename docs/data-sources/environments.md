---
page_title: "argonix_environments Data Source - Argonix"
description: |-
  Fetches all environments in the organization.
---

# argonix_environments (Data Source)

Fetches all environments in the organization.

## Example Usage

```terraform
data "argonix_environments" "all" {}
```

## Schema

### Read-Only

- `environments` (List of Object) — List of all environments. Each has:
  - `id` (String)
  - `name` (String)
  - `variables` (String, Sensitive) — JSON-encoded key-value variables.
  - `is_default` (Boolean)
  - `date_created` (String)
  - `date_modified` (String)
