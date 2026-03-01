---
page_title: "argonix_test_plan Data Source - Argonix"
description: |-
  Fetches a single Argonix test plan by ID.
---

# argonix_test_plan (Data Source)

Fetches a single Argonix test plan by ID.

## Example Usage

```terraform
data "argonix_test_plan" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the test plan.

### Read-Only

- `name` (String)
- `description` (String)
- `suites` (String) — JSON-encoded list of test suite UUIDs.
- `tags` (String) — JSON-encoded tags.
- `end_date` (String) — Target completion date.
- `date_created` (String)
- `date_modified` (String)
