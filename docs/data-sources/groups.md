---
page_title: "argonix_groups Data Source - Argonix"
description: |-
  Fetches all groups in the organization.
---

# argonix_groups (Data Source)

Fetches all groups in the organization.

## Example Usage

```terraform
data "argonix_groups" "all" {}
```

## Schema

### Read-Only

- `groups` (List of Object) — List of all groups. Each has:
  - `id` (String)
  - `name` (String)
  - `description` (String)
  - `tags` (String) — JSON-encoded tags object.
  - `date_created` (String)
  - `date_modified` (String)
