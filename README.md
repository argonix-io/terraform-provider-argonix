# Terraform Provider for Argonix

The Argonix Terraform provider allows you to manage your [Argonix](https://argonix.io) monitoring infrastructure as code.

## Features

- **Monitors** — HTTP, TCP, Ping, DNS, SSL, Keyword, gRPC, Heartbeat, Multi-step HTTP with full config (assertions, remediation, HTTP auth, custom headers, etc.)
- **Synthetic Tests** — API and browser test flows with multi-step scenarios
- **Notification Rules** — Flexible alerting on monitor and synthetic test events (including after each test run for CI/CD)
- **Alert Channels** — Slack, email, webhook, PagerDuty, and more
- **Status Pages** — Public or private status pages with custom branding
- **Groups** — Organize monitors into logical groups
- **Test Management** — Test suites, manual test cases, and test plans

## Requirements

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- [Go](https://go.dev/dl/) >= 1.24 (to build the provider)

## Installation

### Terraform Registry

```terraform
terraform {
  required_providers {
    argonix = {
      source = "argonix-io/argonix"
    }
  }
}
```

### Local Development

```bash
git clone https://github.com/argonix-io/terraform-provider-argonix.git
cd terraform-provider-argonix
go install .
```

## Authentication

Create an API key from the Argonix dashboard under **Settings → API Keys**. The organization is automatically determined from the key — no need to specify it separately.

```terraform
provider "argonix" {
  api_key = var.argonix_api_key
}

variable "argonix_api_key" {
  type      = string
  sensitive = true
}
```

Or via environment variables:

```bash
export ARGONIX_API_KEY="ax_..."
```

## Quick Start

```terraform
# Create a group
resource "argonix_group" "production" {
  name        = "Production"
  description = "Production monitors"
}

# HTTP monitor
resource "argonix_monitor" "api_health" {
  name           = "API Health Check"
  monitor_type   = "http"
  url            = "https://api.example.com/health"
  check_interval = 60
  timeout        = 10
  group_id       = argonix_group.production.id
  tags           = jsonencode(["api", "health"])
  regions        = jsonencode(["eu-france"])
  assertions = jsonencode([
    { type = "status_code", operator = "equals", value = "200" }
  ])
}

# Synthetic API test
resource "argonix_synthetic_test" "login_flow" {
  name      = "Login Flow"
  test_type = "api"
  steps = jsonencode([
    { name = "Get CSRF", method = "GET", url = "https://api.example.com/csrf" },
    { name = "Login", method = "POST", url = "https://api.example.com/login" }
  ])
}

# Alert channel
resource "argonix_alert_channel" "slack" {
  name         = "Slack Ops"
  channel_type = "slack"
  config = jsonencode({
    webhook_url = "https://hooks.slack.com/services/xxx/yyy/zzz"
  })
}

# Notification rule
resource "argonix_notification_rule" "all_down" {
  name              = "Alert on any monitor down"
  trigger_condition = "goes_down"
  all_monitors      = true
  channels          = jsonencode([argonix_alert_channel.slack.id])
}

# Public status page
resource "argonix_status_page" "public" {
  name       = "Service Status"
  slug       = "status"
  visibility = "public"
}
```

## Resources

| Resource | Description |
|----------|-------------|
| `argonix_monitor` | Uptime monitors (HTTP, TCP, Ping, DNS, SSL, Keyword, gRPC, Heartbeat, Multi-step) |
| `argonix_synthetic_test` | API and browser synthetic tests |
| `argonix_group` | Monitor groups |
| `argonix_notification_rule` | Notification rules |
| `argonix_alert_channel` | Notification channels |
| `argonix_status_page` | Public/private status pages |
| `argonix_test_suite` | Test suites |
| `argonix_manual_test_case` | Manual test cases |
| `argonix_test_plan` | Test plans |

## Data Sources

| Data Source | Description |
|------------|-------------|
| `argonix_monitor` / `argonix_monitors` | Read monitors |
| `argonix_synthetic_test` / `argonix_synthetic_tests` | Read synthetic tests |
| `argonix_group` / `argonix_groups` | Read groups |
| `argonix_notification_rule` / `argonix_notification_rules` | Read notification rules |
| `argonix_alert_channel` / `argonix_alert_channels` | Read alert channels |
| `argonix_status_page` / `argonix_status_pages` | Read status pages |
| `argonix_test_suite` / `argonix_test_suites` | Read test suites |
| `argonix_manual_test_case` / `argonix_manual_test_cases` | Read manual test cases |

## Development

```bash
# Build
go build ./...

# Run tests
go test ./...

# Install locally
go install .
```

## Documentation

Full documentation for each resource and data source is available in the [`docs/`](docs/) directory and on the [Terraform Registry](https://registry.terraform.io/providers/argonix-io/argonix/latest/docs).

## License

[Mozilla Public License 2.0](LICENSE)
