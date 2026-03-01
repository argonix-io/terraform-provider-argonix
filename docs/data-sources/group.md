---
page_title: "argonix_group Data Source - Argonix"
description: |-
  Fetches a single Argonix group by ID.
---

# argonix_group (Data Source)

Fetches a single Argonix group by ID.

## Example Usage

```terraform
data "argonix_group" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the group.

### Read-Only

- `name` (String)
- `description` (String)
- `tags` (String) — JSON-encoded tags object.
- `date_created` (String)
- `date_modified` (String)
