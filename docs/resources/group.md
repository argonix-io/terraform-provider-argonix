---
page_title: "argonix_group Resource - Argonix"
description: |-
  Manages an Argonix group.
---

# argonix_group (Resource)

Manages an Argonix group. Use this resource for organization groups. To put monitors in folders, use [`argonix_monitor_group`](monitor_group.md).

## Example Usage

```terraform
resource "argonix_group" "production" {
  name        = "Production"
  description = "Production monitors"
  tags        = jsonencode({ env = "prod" })
}
```

## Schema

### Required

- `name` (String) — Name of the group.

### Optional

- `description` (String) — Description of the group. Defaults to `""`.
- `tags` (String) — JSON-encoded tags object. Defaults to `"{}"`.

### Read-Only

- `id` (String) — UUID of the group.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_group.example <group-uuid>
```
