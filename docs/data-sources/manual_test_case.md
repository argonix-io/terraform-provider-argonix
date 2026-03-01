---
page_title: "argonix_manual_test_case Data Source - Argonix"
description: |-
  Fetches a single Argonix manual test case by ID.
---

# argonix_manual_test_case (Data Source)

Fetches a single Argonix manual test case by ID.

## Example Usage

```terraform
data "argonix_manual_test_case" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the manual test case.

### Read-Only

- `title` (String)
- `description` (String)
- `preconditions` (String)
- `steps` (String) — JSON-encoded steps.
- `priority` (String)
- `tags` (String) — JSON-encoded tags.
- `date_created` (String)
- `date_modified` (String)
