---
page_title: "argonix_synthetic_test Resource - Argonix"
description: |-
  Manages an Argonix synthetic test.
---

# argonix_synthetic_test (Resource)

Manages an Argonix synthetic test. Synthetic tests execute multi-step API or browser workflows on a schedule to verify end-to-end functionality.

## Example Usage

```terraform
resource "argonix_synthetic_test" "login_flow" {
  name      = "Login Flow"
  test_type = "api"
  steps = jsonencode([
    {
      name   = "Get CSRF Token"
      method = "GET"
      url    = "https://api.example.com/auth/csrf"
    },
    {
      name   = "Login"
      method = "POST"
      url    = "https://api.example.com/auth/login"
      body   = "{\"email\": \"test@example.com\", \"password\": \"secret\"}"
    }
  ])
  check_interval = 300
  locations      = jsonencode(["eu-france"])
}
```

## Schema

### Required

- `name` (String) — Name of the synthetic test.
- `test_type` (String) — Type of synthetic test: `api` or `browser`.
- `steps` (String) — JSON-encoded array of step objects defining the test workflow.

### Optional

- `description` (String) — Description of the test. Defaults to `""`.
- `is_active` (Boolean) — Whether the test is active. Defaults to `true`.
- `check_interval` (Number) — Check interval in seconds. Defaults to `300`.
- `timeout` (Number) — Request timeout in seconds. Defaults to `30`.
- `tags` (String) — JSON-encoded list of tags. Defaults to `"[]"`.
- `locations` (String) — JSON-encoded list of region codes to run the test from. Defaults to `"[]"`.

### Read-Only

- `id` (String) — UUID of the synthetic test.
- `current_status` (String) — Current status of the test.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_synthetic_test.example <synthetic-test-uuid>
```
