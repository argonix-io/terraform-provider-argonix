---
page_title: "argonix_log_source Data Source - Argonix"
description: |-
  Reads an Argonix log source by UUID.
---

# argonix_log_source (Data Source)

Reads an Argonix log source by UUID.

## Example usage

```terraform
data "argonix_log_source" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Log source name.
- `connector` (String) — Log connector UUID in the organization.
- `query` (String) — Provider-specific log filter.
- `field_mapping` (String, JSON object) — JSON object mapping the message, timestamp, status, service and host fields.
- `is_active` (Boolean) — Enable periodic log analysis.
- `schedule_minutes` (Number) — Number of minutes between analyses (5–1440).
- `sample_limit` (Number) — Maximum sample lines per analysis.
- `health_intelligence_mode` (String) — disabled, shadow or active.
- `connector_name` (String)
- `connector_type` (String)
- `analysis_mode` (String)
- `detection_settings` (String, JSON object)
- `retention_settings` (String, JSON object)
- `last_analyzed_at` (String)
- `date_created` (String)
- `date_modified` (String)
- `config` (String, JSON object) — Provider settings and optional detection and retention overrides as a JSON object.
