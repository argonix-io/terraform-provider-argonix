---
page_title: "argonix_connector Resource - Argonix"
description: |-
  Manages an Argonix Argos connector.
---

# argonix_connector (Resource)

Manages an Argonix Argos connector. Connectors link external services (Slack, PagerDuty, Jira, AWS, etc.) to the Argos AI agent.

## Example Usage

```terraform
resource "argonix_connector" "slack" {
  name           = "Slack Production"
  connector_type = "slack"
  config         = jsonencode({
    token      = var.slack_token
    channel_id = "C01234567"
  })
}

resource "argonix_connector" "datadog" {
  name           = "Datadog"
  connector_type = "datadog"
  config         = jsonencode({
    api_key = var.datadog_api_key
    app_key = var.datadog_app_key
  })
}

resource "argonix_connector" "gcp_scoped" {
  name           = "GCP Production"
  connector_type = "gcp"
  config         = jsonencode({
    credentials_json = var.gcp_credentials
  })
  scopes = jsonencode({
    projects = ["my-project-prod", "my-project-staging"]
    regions  = ["europe-west1", "us-central1"]
  })
}
```

## Schema

### Required

- `name` (String) — Name of the connector.
- `connector_type` (String) — Type of connector. One of: `slack`, `teams`, `pagerduty`, `opsgenie`, `jira`, `servicenow`, `github`, `gitlab`, `datadog`, `grafana`, `prometheus`, `cloudwatch`, `elastic`, `splunk`, `sentry`, `new_relic`, `aws`, `gcp`, `azure`, `kubernetes`, `terraform`, `ansible`, `jenkins`, `confluence`, `notion`, `linear`, `zendesk`, `okta`, `custom_webhook`, `mysql`, `postgresql`, `mssql`, `oracle`, `redis`, `memcached`, `bigquery`.

### Optional

- `is_active` (Boolean) — Whether the connector is active. Defaults to `true`.
- `config` (String, Sensitive) — JSON-encoded configuration. Defaults to `"{}"`.
- `capabilities` (String) — JSON-encoded capabilities list. Defaults to `"[]"`.
- `tags` (String) — JSON-encoded tags. Defaults to `"[]"`.
- `scopes` (String) — JSON-encoded scopes object restricting where the connector can operate (e.g. projects, regions, namespaces). Defaults to `"{}"`. See scope keys per connector type via the API.

### Read-Only

- `id` (String) — UUID of the connector.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_connector.example <connector-uuid>
```
