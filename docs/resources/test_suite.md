---
page_title: "argonix_test_suite Resource - Argonix"
description: |-
  Manages an Argonix test suite.
---

# argonix_test_suite (Resource)

Manages an Argonix test suite. Test suites group synthetic tests and manual test cases together for organized execution.

## Example Usage

```terraform
resource "argonix_test_suite" "smoke" {
  name            = "Smoke Tests"
  description     = "Quick verification suite"
  synthetic_tests = jsonencode([argonix_synthetic_test.login_flow.id])
  manual_test_cases = jsonencode([argonix_manual_test_case.checkout.id])
}
```

## Schema

### Required

- `name` (String) — Name of the test suite.

### Optional

- `description` (String) — Description. Defaults to `""`.
- `tags` (String) — JSON-encoded list of tags. Defaults to `"[]"`.
- `synthetic_tests` (String) — JSON-encoded list of synthetic test UUIDs. Defaults to `"[]"`.
- `manual_test_cases` (String) — JSON-encoded list of manual test case UUIDs. Defaults to `"[]"`.

### Read-Only

- `id` (String) — UUID of the test suite.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_test_suite.example <test-suite-uuid>
```
