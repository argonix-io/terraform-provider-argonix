---
page_title: "argonix_workflows Data Source - Argonix"
description: |-
  Fetches all workflows in the organization.
---

# argonix_workflows (Data Source)

Fetches all workflows in the organization.

## Example Usage

```terraform
data "argonix_workflows" "all" {}
```

## Schema

### Read-Only

- `workflows` (List of Object) — List of all workflows. Each has:
  - `id` (String)
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
