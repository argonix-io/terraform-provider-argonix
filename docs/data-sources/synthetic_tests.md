---
page_title: "argonix_synthetic_tests Data Source - Argonix"
description: |-
  Fetches all synthetic tests in the organization.
---

# argonix_synthetic_tests (Data Source)

Fetches all synthetic tests in the organization.

## Example Usage

```terraform
data "argonix_synthetic_tests" "all" {}
```

## Schema

### Read-Only

- `synthetic_tests` (List of Object) — List of all synthetic tests. Each has:
  - `id` (String)
  - `name` (String)
  - `description` (String)
  - `is_active` (Boolean)
  - `test_type` (String)
  - `steps` (String) — JSON-encoded steps.
  - `check_interval` (Number)
  - `timeout` (Number)
  - `tags` (String) — JSON-encoded tags.
  - `locations` (String) — JSON-encoded locations.
  - `current_status` (String)
  - `date_created` (String)
  - `date_modified` (String)
