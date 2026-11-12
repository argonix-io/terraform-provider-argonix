---
page_title: "argonix_sre_patrol Data Source - Argonix"
description: |-
  Reads an Argonix sre patrol by UUID.
---

# argonix_sre_patrol (Data Source)

Reads an Argonix sre patrol by UUID.

## Example usage

```terraform
data "argonix_sre_patrol" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

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
- `autonomy` (String)
- `last_run_status` (String)
- `last_run_date` (String)
- `date_created` (String)
- `date_modified` (String)
