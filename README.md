# Terraform Provider for Argonix

The Argonix Terraform provider allows you to manage your [Argonix](https://argonix.io) monitoring infrastructure as code.

## Features

- **Monitors** — HTTP, TCP, Ping, DNS, SSL, Keyword, gRPC, Heartbeat, Multi-step HTTP with full config (assertions, remediation, HTTP auth, custom headers, etc.)
- **Synthetic Tests and Sanity Checks** — API and browser tests, parameterized scenarios, application links, and cron schedules
- **Notification Rules** — Flexible alerting on monitor and synthetic test events (including after each test run for CI/CD)
- **Alert Channels** — Slack, email, webhook, PagerDuty, and more
- **Status Pages** — Public or private status pages with custom branding
- **Groups** — Organize monitors in folders and manage organization groups
- **Discovery and Logs** — Discover monitors from Kubernetes, OpenShift and cloud inventory, and configure log analysis sources
- **Test Management** — Test suites, manual test cases, and test plans
- **Argos** — Workflows, connectors, personas, knowledge bases, chat channels, periodic jobs, health notebooks, and SRE patrols
- **FinOps** — Tenants, cost budgets, and saved cost views
- **Security** — Security policies, detection rules, code analysis projects, and repository or image scanner targets. Read quality gates and code profiles as data sources

## Requirements

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- [Go](https://go.dev/dl/) >= 1.25.5 (to build the provider)

## Installation

### Terraform Registry

```terraform
terraform {
  required_providers {
    argonix = {
      source  = "argonix-io/argonix"
      version = "~> 1.2.0"
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

Create an API key from the Argonix dashboard under **Settings → API Keys**. The provider discovers the organization from the key. Grant `iac:read` for configuration lookups and `iac:write` to manage resources. Code profiles and quality gates require `security:read`. The key inherits the organization role of its creator.

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
resource "argonix_monitor_group" "production" {
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
  group_id       = argonix_monitor_group.production.id
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
| `argonix_group` | Organization groups |
| `argonix_monitor_group` | Monitor folders |
| `argonix_monitor_discovery_source` | Kubernetes, OpenShift and cloud monitor discovery |
| `argonix_log_source` | Scheduled analysis of logs from an existing connector |
| `argonix_notification_rule` | Notification rules |
| `argonix_alert_channel` | Notification channels |
| `argonix_status_page` | Public/private status pages |
| `argonix_test_suite` | Test suites |
| `argonix_manual_test_case` | Manual test cases |
| `argonix_test_plan` | Test plans |
| `argonix_workflow` | Workflows with input hints, scheduled inputs, and timezones |
| `argonix_periodic_job` | Scheduled Argos jobs |
| `argonix_health_notebook` | Health reports and security briefs |
| `argonix_sre_patrol` | Operational observation and reports |
| `argonix_code_project` | Repository configuration for static analysis |
| `argonix_scan_target` | Repository and container image scanner targets |
| `argonix_finops_tenant` | FinOps tenants for cost allocation |
| `argonix_cost_budget` | Cost budgets and alert thresholds |
| `argonix_cost_perspective` | Saved cost views |

## Data Sources

| Data Source | Description |
|------------|-------------|
| `argonix_monitor` / `argonix_monitors` | Read monitors |
| `argonix_synthetic_test` / `argonix_synthetic_tests` | Read synthetic tests |
| `argonix_group` / `argonix_groups` | Read groups |
| `argonix_notification_rule` / `argonix_notification_rules` | Read notification rules |
| `argonix_alert_channel` / `argonix_alert_channels` | Read alert channels |
| `argonix_status_page` / `argonix_status_pages` | Read status pages |
| `argonix_test_suite` | Read test suites |
| `argonix_manual_test_case` | Read manual test cases |
| `argonix_workflow` / `argonix_workflows` | Read workflow definitions |
| `argonix_periodic_job` / `argonix_periodic_jobs` | Read periodic jobs |
| `argonix_health_notebook` / `argonix_health_notebooks` | Read health notebooks and security briefs |
| `argonix_sre_patrol` / `argonix_sre_patrols` | Read SRE patrols |
| `argonix_code_project` / `argonix_code_projects` | Read code projects |
| `argonix_scan_target` / `argonix_scan_targets` | Read scanner targets |
| `argonix_code_quality_gate` / `argonix_code_quality_gates` | Read active quality gates |
| `argonix_code_profile` / `argonix_code_profiles` | Read code analysis profiles |
| `argonix_monitor_group` / `argonix_monitor_groups` | Read monitor folders |
| `argonix_monitor_discovery_source` / `argonix_monitor_discovery_sources` | Read monitor discovery sources |
| `argonix_log_source` / `argonix_log_sources` | Read log analysis sources |
| `argonix_finops_tenant` / `argonix_finops_tenants` | Read FinOps tenants |
| `argonix_cost_budget` / `argonix_cost_budgets` | Read cost budgets |
| `argonix_cost_perspective` / `argonix_cost_perspectives` | Read cost views |

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
