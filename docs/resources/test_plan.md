---
page_title: "argonix_test_plan Resource - Argonix"
description: |-
  Manages an Argonix test plan.
---

# argonix_test_plan (Resource)

Manages an Argonix test plan. Test plans organize test suites with a target completion date for release planning.

## Example Usage

```terraform
resource "argonix_test_plan" "release_v2" {
  name        = "Release v2.0 Test Plan"
  description = "All tests for v2.0 release"
  suites      = jsonencode([argonix_test_suite.smoke.id])
  end_date    = "2026-04-01"
}
```

## Schema

### Required

- `name` (String) — Name of the test plan.

### Optional

- `description` (String) — Description. Defaults to `""`.
- `suites` (String) — JSON-encoded list of test suite UUIDs. Defaults to `"[]"`.
- `tags` (String) — JSON-encoded list of tags. Defaults to `"[]"`.
- `end_date` (String) — Target completion date in `YYYY-MM-DD` format.

### Read-Only

- `id` (String) — UUID of the test plan.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_test_plan.example <test-plan-uuid>
```
