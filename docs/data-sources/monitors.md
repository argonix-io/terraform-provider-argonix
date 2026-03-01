---
page_title: "argonix_monitors Data Source - Argonix"
description: |-
  Fetches all monitors in the organization.
---

# argonix_monitors (Data Source)

Fetches all monitors in the organization.

## Example Usage

```terraform
data "argonix_monitors" "all" {}

output "monitor_count" {
  value = length(data.argonix_monitors.all.monitors)
}
```

## Schema

### Read-Only

- `monitors` (List of Object) — List of all monitors. Each monitor has the following attributes:
  - `id` (String)
  - `name` (String)
  - `monitor_type` (String)
  - `is_active` (Boolean)
  - `url` (String)
  - `hostname` (String)
  - `port` (Number)
  - `dns_record_type` (String)
  - `dns_expected` (String)
  - `http_method` (String)
  - `http_headers` (String) — JSON-encoded headers.
  - `http_body` (String)
  - `http_body_content_type` (String)
  - `follow_redirects` (Boolean)
  - `verify_ssl` (Boolean)
  - `http_auth_user` (String)
  - `http_auth_pass` (String, Sensitive)
  - `keyword` (String)
  - `keyword_exists` (Boolean)
  - `check_interval` (Number)
  - `timeout` (Number)
  - `retries` (Number)
  - `remediation_enabled` (Boolean)
  - `remediation_script` (String)
  - `remediation_timeout` (Number)
  - `remediation_wait_seconds` (Number)
  - `heartbeat_token` (String)
  - `heartbeat_grace_seconds` (Number)
  - `multi_step_config` (String) — JSON-encoded.
  - `grpc_service` (String)
  - `grpc_method` (String)
  - `grpc_proto` (String)
  - `grpc_metadata` (String) — JSON-encoded.
  - `grpc_tls` (Boolean)
  - `assertions` (String) — JSON-encoded.
  - `ssl_expiry_warn_days` (Number)
  - `location` (String)
  - `regions` (String) — JSON-encoded.
  - `tags` (String) — JSON-encoded.
  - `group_id` (String)
  - `current_status` (String)
  - `date_created` (String)
  - `date_modified` (String)
