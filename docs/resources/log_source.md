---
page_title: "argonix_log_source Resource - Argonix"
description: |-
  Manages an Argonix log analysis source on an existing connector.
---

# argonix_log_source (Resource)

Manages an Argonix log analysis source on an existing connector.

Use a connector that supports log analysis. Kubernetes and OpenShift snapshot sources require `namespace` and `pod` in `config`. Store connector credentials with the connector.

## Example usage

```terraform
resource "argonix_log_source" "service" {
  name             = "Service logs"
  connector        = var.loki_connector_id
  query            = "{namespace=\"production\",app=\"api\"}"
  schedule_minutes = 15
  sample_limit     = 10000
  config = jsonencode({
    new_pattern_min_count = 5
  })
}
```

## Schema

### Required

- `name` (String) — Log source name.
- `connector` (String) — Log connector UUID in the organization.

### Optional

- `query` (String) — Provider-specific log filter.
- `field_mapping` (String, JSON object) — JSON object mapping the message, timestamp, status, service and host fields.
- `is_active` (Boolean) — Enable periodic log analysis.
- `schedule_minutes` (Number) — Number of minutes between analyses (5–1440).
- `sample_limit` (Number) — Maximum sample lines per analysis.
- `health_intelligence_mode` (String) — disabled, shadow or active.
- `config` (String, JSON object) — Provider settings and optional detection and retention overrides as a JSON object.

### Read-only

- `id` (String) — Argonix UUID.
- `connector_name` (String)
- `connector_type` (String)
- `analysis_mode` (String)
- `detection_settings` (String, JSON object)
- `retention_settings` (String, JSON object)
- `last_analyzed_at` (String)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_log_source.example <uuid>
```
