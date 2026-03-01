---
page_title: "argonix_monitor Data Source - Argonix"
description: |-
  Fetches a single Argonix monitor by ID.
---

# argonix_monitor (Data Source)

Fetches a single Argonix monitor by ID.

## Example Usage

```terraform
data "argonix_monitor" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "monitor_name" {
  value = data.argonix_monitor.example.name
}

output "is_ssl_valid" {
  value = data.argonix_monitor.example.current_status
}
```

## Schema

### Required

- `id` (String) — UUID of the monitor.

### Read-Only

- `name` (String) — Display name of the monitor.
- `monitor_type` (String) — Type of monitor: `http`, `ping`, `tcp`, `dns`, `ssl`, `keyword`, `grpc`, `heartbeat`, `multi_step_http`.
- `is_active` (Boolean) — Whether the monitor is active.
- `url` (String) — URL being monitored.
- `hostname` (String) — Hostname being monitored.
- `port` (Number) — Port being monitored.
- `dns_record_type` (String) — DNS record type.
- `dns_expected` (String) — Expected DNS response.
- `http_method` (String) — HTTP method.
- `http_headers` (String) — JSON-encoded HTTP headers.
- `http_body` (String) — HTTP request body.
- `http_body_content_type` (String) — Content-Type for HTTP body.
- `follow_redirects` (Boolean) — Whether redirects are followed.
- `verify_ssl` (Boolean) — Whether SSL is verified.
- `http_auth_user` (String) — HTTP Basic Auth username.
- `http_auth_pass` (String, Sensitive) — HTTP Basic Auth password.
- `keyword` (String) — Keyword searched for.
- `keyword_exists` (Boolean) — Alert when keyword is missing (true) or found (false).
- `check_interval` (Number) — Check interval in seconds.
- `timeout` (Number) — Request timeout in seconds.
- `retries` (Number) — Number of retries.
- `remediation_enabled` (Boolean) — Automatic remediation enabled.
- `remediation_script` (String) — Remediation shell script.
- `remediation_timeout` (Number) — Remediation script timeout in seconds.
- `remediation_wait_seconds` (Number) — Seconds to wait after remediation.
- `heartbeat_token` (String) — Auto-generated heartbeat token.
- `heartbeat_grace_seconds` (Number) — Grace period for heartbeat.
- `multi_step_config` (String) — JSON-encoded multi-step config.
- `grpc_service` (String) — gRPC service name.
- `grpc_method` (String) — gRPC method.
- `grpc_proto` (String) — gRPC protobuf definition.
- `grpc_metadata` (String) — JSON-encoded gRPC metadata.
- `grpc_tls` (Boolean) — Whether TLS is used for gRPC.
- `assertions` (String) — JSON-encoded assertions.
- `ssl_expiry_warn_days` (Number) — Days before SSL expiry warning.
- `location` (String) — Primary check location.
- `regions` (String) — JSON-encoded regions.
- `tags` (String) — JSON-encoded tags.
- `group_id` (String) — UUID of the group.
- `current_status` (String) — Current status.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.
