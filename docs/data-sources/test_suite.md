---
page_title: "argonix_test_suite Data Source - Argonix"
description: |-
  Fetches a single Argonix test suite by ID.
---

# argonix_test_suite (Data Source)

Fetches a single Argonix test suite by ID.

## Example Usage

```terraform
data "argonix_test_suite" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the test suite.

### Read-Only

- `name` (String)
- `description` (String)
- `tags` (String) — JSON-encoded tags.
- `synthetic_tests` (String) — JSON-encoded list of synthetic test UUIDs.
- `manual_test_cases` (String) — JSON-encoded list of manual test case UUIDs.
- `date_created` (String)
- `date_modified` (String)
