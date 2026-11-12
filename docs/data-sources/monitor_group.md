---
page_title: "argonix_monitor_group Data Source - Argonix"
description: |-
  Reads an Argonix monitor group by UUID.
---

# argonix_monitor_group (Data Source)

Reads an Argonix monitor group by UUID.

## Example usage

```terraform
data "argonix_monitor_group" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Monitor folder name.
- `description` (String) — Monitor folder description.
- `is_managed` (Boolean)
- `monitors_count` (Number)
- `date_created` (String)
- `date_modified` (String)
