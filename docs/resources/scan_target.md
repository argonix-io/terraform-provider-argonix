---
page_title: "argonix_scan_target Resource - Argonix"
description: |-
  Manages an Argonix repository or container image scanner target.
---

# argonix_scan_target (Resource)

Manages an Argonix repository or container image scanner target.

This resource manages repository and image scanning schedules. Link a target to a code project using `project`.

## Example usage

```terraform
resource "argonix_scan_target" "service_image" {
  name             = "Service image"
  kind             = "image"
  identifier       = "ghcr.io/example/service:1.0.0"
  project          = argonix_code_project.service.id
  enabled_scanners = jsonencode(["trivy"])
  frequency        = "daily"
  schedule_hour    = 3
}
```

## Schema

### Required

- `identifier` (String) — Repository URL or container image identifier.
- `enabled_scanners` (String, JSON array) — Nonempty JSON array of scanners supported for this target kind.
- `kind` (String) — repository or image.

### Optional

- `name` (String) — Display name of the target.
- `connector` (String) — Git or registry connector ID.
- `project` (String) — Code project ID.
- `branch` (String) — Branch for a repository target.
- `default_ref` (String) — Tag or digest for an unpinned image.
- `platform` (String) — OCI image platform, such as linux/amd64.
- `frequency` (String) — manual, daily or weekly.
- `schedule_hour` (Number) — UTC hour of execution (0–23).
- `schedule_minute` (Number) — UTC minute of execution (0–59).
- `schedule_days` (String, JSON array) — JSON array of weekday numbers, with Monday as 0.
- `is_active` (Boolean) — Whether the target is active.
- `webhook_enabled` (Boolean) — Accept repository webhooks for this target.
- `enforce_merge_requests` (Boolean) — Report scan results to the merge request gate.
- `gate_environment` (String) — Optional environment for the merge request gate.

### Read-only

- `id` (String) — Argonix UUID.
- `last_dispatched_at` (String)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_scan_target.example <uuid>
```
