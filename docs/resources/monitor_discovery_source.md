---
page_title: "argonix_monitor_discovery_source Resource - Argonix"
description: |-
  Manages discovery of monitors from Kubernetes, OpenShift or cloud inventory.
---

# argonix_monitor_discovery_source (Resource)

Manages discovery of monitors from Kubernetes, OpenShift or cloud inventory.

The connector determines which resource kinds are supported. Managed monitors and their latest discovery statistics are visible in Argonix.

## Example usage

```terraform
resource "argonix_monitor_discovery_source" "production" {
  name                  = "Production ingress"
  connector             = var.kubernetes_connector_id
  source_type           = "kubernetes"
  resource_kinds        = jsonencode(["ingress"])
  namespace_allow       = jsonencode(["production"])
  default_group         = var.monitor_group_id
  sync_interval_minutes = 15
}
```

## Schema

### Required

- `name` (String) — Discovery source name.
- `connector` (String) — Connector UUID in the organization.

### Optional

- `source_type` (String) — kubernetes or cloud_inventory.
- `enabled` (Boolean) — Enable periodic discovery.
- `resource_kinds` (String, JSON array) — JSON array of resource kinds supported by the connector.
- `include_private_targets` (Boolean) — Include internal hostnames and addresses when the probe can reach them.
- `namespace_allow` (String, JSON array) — JSON array of namespaces to include.
- `namespace_deny` (String, JSON array) — JSON array of namespaces to exclude.
- `host_include_regex` (String) — Only discover hosts matching this regular expression.
- `host_exclude_regex` (String) — Skip hosts matching this regular expression.
- `ingress_classes` (String, JSON array) — JSON array of Ingress classes to include.
- `label_selector` (String, JSON object) — JSON object of labels that discovered resources must match.
- `label_tag_keys` (String, JSON array) — JSON array of label keys used as tags on created monitors.
- `expand_paths` (Boolean) — Create monitors for individual Ingress or HTTPRoute paths.
- `default_check_interval` (Number) — Check interval in seconds for newly created monitors.
- `default_group` (String) — Monitor group UUID in the organization.
- `default_auto_investigate` (Boolean) — Enable investigation on monitors created from this source.
- `sync_interval_minutes` (Number) — Minimum number of minutes between discovery runs.
- `missing_grace_syncs` (Number) — Consecutive missing observations before the managed monitor is removed.

### Read-only

- `id` (String) — Argonix UUID.
- `connector_name` (String)
- `connector_type` (String)
- `default_group_name` (String)
- `default_auto_remediate` (Boolean)
- `last_sync_status` (String)
- `last_sync_at` (String)
- `last_sync_message` (String)
- `last_created_count` (Number)
- `last_updated_count` (Number)
- `last_deleted_count` (Number)
- `last_skipped_count` (Number)
- `monitors_count` (Number)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_monitor_discovery_source.example <uuid>
```
