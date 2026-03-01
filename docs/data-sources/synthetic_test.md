---
page_title: "argonix_synthetic_test Data Source - Argonix"
description: |-
  Fetches a single Argonix synthetic test by ID.
---

# argonix_synthetic_test (Data Source)

Fetches a single Argonix synthetic test by ID.

## Example Usage

```terraform
data "argonix_synthetic_test" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the synthetic test.

### Read-Only

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
