---
page_title: "argonix_sre_patrol Resource - Argonix"
description: |-
  Manages an Argonix SRE patrol for operational observation.
---

# argonix_sre_patrol (Resource)

Manages an Argonix SRE patrol for operational observation.

New patrols use the API's default observation mode.

## Example usage

```terraform
resource "argonix_sre_patrol" "production" {
  name                 = "Production SRE patrol"
  use_applications     = true
  application_root_ids = [var.application_id]
  cadence_minutes      = 15
}
```

## Schema

### Optional

- `name` (String) — Name of the patrol.
- `connector_ids` (Set of String) — Connector IDs in scope.
- `perspective` (String) — all, business or platform.
- `use_applications` (Boolean) — Scope the patrol to applications.
- `application_root_ids` (Set of String) — Selected application IDs.
- `include_unattributed` (Boolean) — Include signals without application attribution.
- `minimum_attribution_confidence` (Number) — Minimum attribution confidence from 0 to 100.
- `notification_channel_ids` (Set of String) — Chat channel IDs receiving results.
- `alert_channel_ids` (Set of String) — Alert channel IDs receiving alerts.
- `delivery_channel_ids` (Set of String) — Alert channel IDs receiving patrol reports.
- `is_active` (Boolean) — Whether the patrol is active.
- `cadence_minutes` (Number) — Minutes between sweeps, from 5 to 60.
- `min_confidence` (Number) — Minimum confidence from 0 to 100.
- `ticketing_min_severity` (String) — Minimum severity for ticket creation.
- `custom_prompt` (String) — Additional instructions for interpreting operational signals.

### Read-only

- `id` (String) — Argonix UUID.
- `autonomy` (String)
- `last_run_status` (String)
- `last_run_date` (String)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_sre_patrol.example <uuid>
```
