---
page_title: "argonix_monitor_group Resource - Argonix"
description: |-
  Manages a folder for Argonix monitors.
---

# argonix_monitor_group (Resource)

Manages a folder for Argonix monitors.

Use this resource for monitor folders and assign its ID to the monitor's `group_id`.

## Example usage

```terraform
resource "argonix_monitor_group" "production" {
  name        = "Production"
  description = "Production monitors"
}
```

## Schema

### Required

- `name` (String) — Monitor folder name.

### Optional

- `description` (String) — Monitor folder description.

### Read-only

- `id` (String) — Argonix UUID.
- `is_managed` (Boolean)
- `monitors_count` (Number)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_monitor_group.example <uuid>
```
