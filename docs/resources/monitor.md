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
  group_id       = argonix_group.production.id
  tags           = jsonencode(["api", "health"])
  regions        = jsonencode(["eu-france", "us-east"])
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

## Schema

### Required

- `name` (String) — Display name of the monitor.
- `monitor_type` (String) — Type of monitor. One of: `http`, `ping`, `tcp`, `dns`, `ssl`, `keyword`, `grpc`, `heartbeat`, `multi_step_http`.

### Optional

- `is_active` (Boolean) — Whether the monitor is active. Defaults to `true`.
- `url` (String) — URL to monitor (for HTTP/keyword/SSL monitors). Defaults to `""`.
- `hostname` (String) — Hostname to monitor (for Ping/TCP/DNS monitors). Defaults to `""`.
- `port` (Number) — Port to monitor (for TCP monitors).
- `check_interval` (Number) — Check interval in seconds. Defaults to `300`.
- `timeout` (Number) — Request timeout in seconds. Defaults to `30`.
- `retries` (Number) — Number of retries before marking as down. Defaults to `0`.
- `http_method` (String) — HTTP method. Defaults to `"GET"`.
- `group_id` (String) — UUID of the group this monitor belongs to.
- `tags` (String) — JSON-encoded list of tags. Defaults to `"[]"`.
- `regions` (String) — JSON-encoded list of region codes to run checks from. Defaults to `"[]"`.

### Read-Only

- `id` (String) — UUID of the monitor.
- `current_status` (String) — Current status: `up`, `down`, `degraded`, `maintenance`, `unknown`.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

Monitors can be imported using their UUID:

```shell
terraform import argonix_monitor.example <monitor-uuid>
```
