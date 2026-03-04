---
page_title: "argonix_environment Data Source - Argonix"
description: |-
  Fetches a single Argonix environment by ID.
---

# argonix_environment (Data Source)

Fetches a single Argonix environment by ID.

## Example Usage

```terraform
data "argonix_environment" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the environment.

### Read-Only

- `name` (String)
- `variables` (String, Sensitive) — JSON-encoded key-value variables.
- `is_default` (Boolean)
- `date_created` (String)
- `date_modified` (String)
