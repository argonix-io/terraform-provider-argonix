---
page_title: "argonix_monitor Resource - Argonix"
description: |-
  Manages an Argonix monitor.
---

# argonix_monitor (Resource)

Manages an Argonix monitor. Supports HTTP, TCP, Ping, DNS, SSL, Keyword, gRPC, Heartbeat, and Multi-step HTTP monitor types.

## Example Usage

### HTTP Monitor

```terraform
resource "argonix_monitor" "api_health" {
  name           = "API Health Check"
  monitor_type   = "http"
  url            = "https://api.example.com/health"
  check_interval = 60
  timeout        = 10
  http_method    = "GET"
  verify_ssl     = true
  group_id       = argonix_group.production.id
  tags           = jsonencode(["api", "health"])
  regions        = jsonencode(["eu-france"])
  assertions = jsonencode([
    { type = "status_code", operator = "equals", value = "200" }
  ])
}
```

### TCP Monitor

```terraform
resource "argonix_monitor" "database_ping" {
  name           = "Database Ping"
  monitor_type   = "tcp"
  hostname       = "db.example.com"
  port           = 5432
  check_interval = 120
}
```

### DNS Monitor

```terraform
resource "argonix_monitor" "dns_check" {
  name            = "DNS A Record"
  monitor_type    = "dns"
  hostname        = "example.com"
  dns_record_type = "A"
  dns_expected    = "93.184.216.34"
}
```

### Keyword Monitor

```terraform
resource "argonix_monitor" "keyword_check" {
  name           = "Homepage Keyword"
  monitor_type   = "keyword"
  url            = "https://example.com"
  keyword        = "Welcome"
  keyword_exists = true
}
```

### gRPC Monitor

```terraform
resource "argonix_monitor" "grpc_health" {
  name         = "gRPC Health"
  monitor_type = "grpc"
  hostname     = "grpc.example.com"
  port         = 443
  grpc_service = "grpc.health.v1.Health"
  grpc_method  = "Check"
  grpc_tls     = true
}
```

### Heartbeat Monitor

```terraform
resource "argonix_monitor" "cron_heartbeat" {
  name                    = "Cron Job Heartbeat"
  monitor_type            = "heartbeat"
  check_interval          = 3600
  heartbeat_grace_seconds = 300
}
```

### Monitor with Remediation

```terraform
resource "argonix_monitor" "auto_heal" {
  name                     = "Auto-Healing API"
  monitor_type             = "http"
  url                      = "https://api.example.com/health"
  remediation_enabled      = true
  remediation_script       = "#!/bin/bash\nsystemctl restart myapp"
  remediation_timeout      = 120
  remediation_wait_seconds = 60
  auto_investigate         = true
  auto_remediate           = true
  remediation_strategy     = "approval_required"
}
```

## Schema

### Required

- `name` (String) — Display name of the monitor.
- `monitor_type` (String) — Type of monitor. One of: `http`, `ping`, `tcp`, `dns`, `ssl`, `keyword`, `grpc`, `heartbeat`, `multi_step_http`.

### Optional

**General**

- `is_active` (Boolean) — Whether the monitor is active. Defaults to `true`.
- `group_id` (String) — UUID of the group this monitor belongs to.
- `tags` (String) — JSON-encoded list of tags. Defaults to `"[]"`.
- `regions` (String) — JSON-encoded list of region codes to run checks from. Defaults to `"[]"`.
- `location` (String) — Primary check location. Defaults to `"eu-france"`.

**Target**

- `url` (String) — URL to monitor (for HTTP/keyword/SSL monitors). Defaults to `""`.
- `hostname` (String) — Hostname to monitor (for Ping/TCP/DNS monitors). Defaults to `""`.
- `port` (Number) — Port to monitor (for TCP monitors).

**Scheduling**

- `check_interval` (Number) — Check interval in seconds. Defaults to `300`.
- `timeout` (Number) — Request timeout in seconds. Defaults to `30`.
- `retries` (Number) — Number of retries before marking as down. Defaults to `0`.

**DNS**

- `dns_record_type` (String) — DNS record type to query: `A`, `AAAA`, `CNAME`, `MX`, `NS`, `TXT`, etc. Defaults to `"A"`.
- `dns_expected` (String) — Expected DNS response value. Defaults to `""`.

**HTTP Options**

- `http_method` (String) — HTTP method: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, `OPTIONS`. Defaults to `"GET"`.
- `http_headers` (String) — JSON-encoded object of custom HTTP headers. Defaults to `"{}"`.
- `http_body` (String) — HTTP request body (for POST/PUT/PATCH). Defaults to `""`.
- `http_body_content_type` (String) — Content-Type header for the HTTP body. Defaults to `"application/json"`.
- `follow_redirects` (Boolean) — Whether to follow HTTP redirects. Defaults to `true`.
- `verify_ssl` (Boolean) — Whether to verify SSL certificates. Defaults to `true`.
- `http_auth_user` (String) — HTTP Basic Auth username. Defaults to `""`.
- `http_auth_pass` (String, Sensitive) — HTTP Basic Auth password. Defaults to `""`.

**Keyword**

- `keyword` (String) — Keyword to search for in the response body. Defaults to `""`.
- `keyword_exists` (Boolean) — If `true`, alert when keyword is NOT found. If `false`, alert when keyword IS found. Defaults to `true`.

**Remediation**

- `remediation_enabled` (Boolean) — Enable automatic remediation when the monitor goes down. Defaults to `false`.
- `remediation_script` (String) — Shell script to execute when remediation is triggered. Defaults to `""`.
- `remediation_timeout` (Number) — Maximum execution time for the remediation script in seconds. Defaults to `60`.
- `remediation_wait_seconds` (Number) — Seconds to wait after remediation before rechecking. Defaults to `30`.

**Argos AI**

- `auto_investigate` (Boolean) — Enable Argos AI auto-investigation when the monitor goes down. Defaults to `false`.
- `auto_remediate` (Boolean) — Enable Argos AI auto-remediation after investigation. Defaults to `false`.
- `remediation_strategy` (String) — Remediation strategy: `auto` or `approval_required`. Defaults to `"approval_required"`.

**Heartbeat**

- `heartbeat_grace_seconds` (Number) — Grace period in seconds before a missed heartbeat triggers an alert. Defaults to `0`.

**Multi-step HTTP**

- `multi_step_config` (String) — JSON-encoded multi-step HTTP configuration. Defaults to `"[]"`.

**gRPC**

- `grpc_service` (String) — gRPC service name. Defaults to `""`.
- `grpc_method` (String) — gRPC method to call. Defaults to `""`.
- `grpc_proto` (String) — Protobuf definition for the gRPC service. Defaults to `""`.
- `grpc_metadata` (String) — JSON-encoded gRPC metadata key-value pairs. Defaults to `"{}"`.
- `grpc_tls` (Boolean) — Whether to use TLS for gRPC connections. Defaults to `true`.

**Assertions & SSL**

- `assertions` (String) — JSON-encoded list of assertion objects to validate responses. Defaults to `"[]"`.
- `ssl_expiry_warn_days` (Number) — Days before SSL expiry to trigger a warning. Defaults to `30`.

### Read-Only

- `id` (String) — UUID of the monitor.
- `heartbeat_token` (String) — Auto-generated token for heartbeat push URL.
- `current_status` (String) — Current status: `up`, `down`, `degraded`, `maintenance`, `unknown`.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

Monitors can be imported using their UUID:

```shell
terraform import argonix_monitor.example <monitor-uuid>
```
