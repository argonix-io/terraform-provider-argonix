---
page_title: "argonix_manual_test_case Resource - Argonix"
description: |-
  Manages an Argonix manual test case.
---

# argonix_manual_test_case (Resource)

Manages an Argonix manual test case. Manual test cases define step-by-step procedures for manual QA testing.

## Example Usage

```terraform
resource "argonix_manual_test_case" "checkout" {
  title       = "Verify checkout flow"
  description = "End-to-end checkout test"
  priority    = "high"
  steps = jsonencode([
    { description = "Add item to cart",    expected = "Item appears in cart" },
    { description = "Proceed to checkout", expected = "Checkout page loads" },
    { description = "Complete payment",    expected = "Order confirmation shown" }
  ])
}
```

## Schema

### Required

- `title` (String) — Title of the test case.
- `steps` (String) — JSON-encoded ordered list of steps: `[{"description": "...", "expected": "..."}, ...]`.

### Optional

- `description` (String) — Description. Defaults to `""`.
- `preconditions` (String) — Preconditions that must be met before execution. Defaults to `""`.
- `priority` (String) — Priority level: `critical`, `high`, `medium`, `low`. Defaults to `"medium"`.
- `tags` (String) — JSON-encoded list of tags. Defaults to `"[]"`.

### Read-Only

- `id` (String) — UUID of the test case.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_manual_test_case.example <manual-test-case-uuid>
```
