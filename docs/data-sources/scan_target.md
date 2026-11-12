---
page_title: "argonix_scan_target Data Source - Argonix"
description: |-
  Reads an Argonix scan target by UUID.
---

# argonix_scan_target (Data Source)

Reads an Argonix scan target by UUID.

## Example usage

```terraform
data "argonix_scan_target" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Display name of the target.
- `identifier` (String) — Repository URL or container image identifier.
- `connector` (String) — Git or registry connector ID.
- `project` (String) — Code project ID.
- `branch` (String) — Branch for a repository target.
- `default_ref` (String) — Tag or digest for an unpinned image.
- `platform` (String) — OCI image platform, such as linux/amd64.
- `enabled_scanners` (String, JSON array) — Nonempty JSON array of scanners supported for this target kind.
- `frequency` (String) — manual, daily or weekly.
- `schedule_hour` (Number) — UTC hour of execution (0–23).
- `schedule_minute` (Number) — UTC minute of execution (0–59).
- `schedule_days` (String, JSON array) — JSON array of weekday numbers, with Monday as 0.
- `is_active` (Boolean) — Whether the target is active.
- `webhook_enabled` (Boolean) — Accept repository webhooks for this target.
- `enforce_merge_requests` (Boolean) — Report scan results to the merge request gate.
- `gate_environment` (String) — Optional environment for the merge request gate.
- `last_dispatched_at` (String)
- `date_created` (String)
- `date_modified` (String)
- `kind` (String) — repository or image.
