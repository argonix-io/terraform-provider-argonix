---
page_title: "argonix_workflow Data Source - Argonix"
description: |-
  Fetches a single Argonix workflow by ID.
---

# argonix_workflow (Data Source)

Fetches a single Argonix workflow by ID.

## Example Usage

```terraform
data "argonix_workflow" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the workflow.

### Read-Only

- `name` (String)
- `slug` (String)
- `description` (String)
- `category` (String)
- `steps` (String) — JSON-encoded workflow steps.
- `required_connector_types` (String) — JSON-encoded list.
- `requires_confirmation` (Boolean)
- `schedule` (String)
- `is_active` (Boolean)
- `date_created` (String)
- `date_modified` (String)
