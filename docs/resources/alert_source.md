---
page_title: "argonix_alert_source Resource - Argonix"
description: |-
  Manages an Argonix alert source for ingesting external alerts via webhooks.
---

# argonix_alert_source (Resource)

Manages an Argonix alert source. Alert sources let you ingest alerts from external monitoring tools (Alertmanager, Datadog, Grafana, PagerDuty, OpsGenie) via webhooks. Ingested alerts create incidents and can optionally trigger Argos AI auto-investigation and auto-remediation.

## Example Usage

### Alertmanager Source

```terraform
resource "argonix_alert_source" "alertmanager" {
  name             = "Production Alertmanager"
  source_type      = "alertmanager"
  auto_investigate = true
  channels         = [argonix_alert_channel.slack_ops.id]
}
```

### Datadog Source with Auto-Remediation

```terraform
resource "argonix_alert_source" "datadog" {
  name                 = "Datadog Production"
  source_type          = "datadog"
  connector            = argonix_connector.datadog.id
  auto_investigate     = true
  auto_remediate       = true
  remediation_strategy = "approval_required"
  channels             = [argonix_alert_channel.slack_ops.id, argonix_alert_channel.email_ops.id]
}
```

### Generic Webhook

```terraform
resource "argonix_alert_source" "generic" {
  name        = "Custom Webhook"
  source_type = "generic"
}

output "webhook_url" {
  value     = argonix_alert_source.generic.webhook_url
  sensitive = true
}
```

## Schema

### Required

- `name` (String) — Display name of the alert source.
- `source_type` (String) — Type of external alerting system. One of: `alertmanager`, `datadog`, `grafana`, `pagerduty`, `opsgenie`, `generic`.

### Optional

- `is_active` (Boolean) — Whether the alert source is enabled. Defaults to `true`.
- `connector` (String) — UUID of the connector to use for investigation/remediation context.
- `filters` (String) — JSON-encoded filter configuration for incoming alerts.
- `auto_investigate` (Boolean) — Enable Argos AI auto-investigation for ingested alerts. Defaults to `false`.
- `auto_remediate` (Boolean) — Enable Argos AI auto-remediation. Defaults to `false`.
- `remediation_strategy` (String) — How remediation is executed: `auto` or `approval_required`. Defaults to `approval_required`.
- `channels` (List of String) — List of alert channel UUIDs to notify when alerts are ingested.

### Read-Only

- `id` (String) — UUID of the alert source.
- `webhook_secret` (String, Sensitive) — Secret token for the webhook endpoint.
- `webhook_url` (String) — Full webhook URL for sending alerts to this source.
- `last_received_at` (String) — Timestamp of the last received alert.
- `total_received` (Number) — Total number of alerts received.
- `date_created` (String)
- `date_modified` (String)

## Import

Alert sources can be imported using their UUID:

```shell
terraform import argonix_alert_source.example xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
