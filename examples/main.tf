# Example Terraform configuration for the Argonix provider.

terraform {
  required_providers {
    argonix = {
      source = "hashicorp.com/argonix-io/argonix"
    }
  }
}

provider "argonix" {
  # url = "https://api.argonix.io"  # optional, this is the default
  api_key = var.argonix_api_key
}

variable "argonix_api_key" {
  type      = string
  sensitive = true
}

# --- Group ---
resource "argonix_group" "production" {
  name        = "Production"
  description = "Production monitors"
  tags        = jsonencode({ env = "prod" })
}

# --- Monitor ---
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

resource "argonix_monitor" "database_ping" {
  name         = "Database Ping"
  monitor_type = "tcp"
  hostname     = "db.example.com"
  port         = 5432
  check_interval = 120
}

# --- Synthetic Test ---
resource "argonix_synthetic_test" "login_flow" {
  name      = "Login Flow"
  test_type = "api"
  steps = jsonencode([
    {
      name    = "Get CSRF Token"
      method  = "GET"
      url     = "https://api.example.com/auth/csrf"
    },
    {
      name    = "Login"
      method  = "POST"
      url     = "https://api.example.com/auth/login"
      body    = "{\"email\": \"test@example.com\", \"password\": \"secret\"}"
    }
  ])
  check_interval = 300
  locations      = jsonencode(["eu-france"])
}

# --- Alert Channel ---
resource "argonix_alert_channel" "slack_ops" {
  name         = "Slack Ops Channel"
  channel_type = "slack"
  config = jsonencode({
    webhook_url = "https://hooks.slack.com/services/xxx/yyy/zzz"
  })
}

resource "argonix_alert_channel" "email_ops" {
  name         = "Email Ops Team"
  channel_type = "email"
  config = jsonencode({
    addresses = ["ops@example.com", "oncall@example.com"]
  })
}

# --- Notification Rule ---
resource "argonix_notification_rule" "all_down" {
  name              = "Alert on any monitor down"
  trigger_condition = "goes_down"
  all_monitors      = true
  channels          = jsonencode([argonix_alert_channel.slack_ops.id, argonix_alert_channel.email_ops.id])
  cooldown_minutes  = 10
}

# --- Status Page ---
resource "argonix_status_page" "public" {
  name         = "Example Status"
  slug         = "example-status"
  visibility   = "public"
  accent_color = "#10B981"
  header_text  = "Service Status"
}


# --- Test Suite ---
resource "argonix_test_suite" "smoke" {
  name           = "Smoke Tests"
  description    = "Quick verification suite"
  synthetic_tests = jsonencode([argonix_synthetic_test.login_flow.id])
}

# --- Manual Test Case ---
resource "argonix_manual_test_case" "checkout" {
  title       = "Verify checkout flow"
  description = "End-to-end checkout test"
  priority    = "high"
  steps = jsonencode([
    { description = "Add item to cart",   expected = "Item appears in cart" },
    { description = "Proceed to checkout", expected = "Checkout page loads" },
    { description = "Complete payment",    expected = "Order confirmation shown" }
  ])
}

# --- Test Plan ---
resource "argonix_test_plan" "release_v2" {
  name        = "Release v2.0 Test Plan"
  description = "All tests for v2.0 release"
  suites      = jsonencode([argonix_test_suite.smoke.id])
  end_date    = "2026-04-01"
}

# --- Data Sources ---
data "argonix_monitors" "all" {}

data "argonix_groups" "all" {}

data "argonix_alert_channels" "all" {}

output "monitor_count" {
  value = length(data.argonix_monitors.all.monitors)
}
