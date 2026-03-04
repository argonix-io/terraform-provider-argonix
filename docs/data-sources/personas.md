---
page_title: "argonix_personas Data Source - Argonix"
description: |-
  Fetches all personas in the organization.
---

# argonix_personas (Data Source)

Fetches all Argos personas in the organization.

## Example Usage

```terraform
data "argonix_personas" "all" {}
```

## Schema

### Read-Only

- `personas` (List of Object) — List of all personas. Each has:
  - `id` (String)
  - `name` (String)
  - `description` (String)
  - `template` (String)
  - `system_prompt` (String)
  - `is_active` (Boolean)
  - `connector_ids` (String) — JSON-encoded list of connector UUIDs.
  - `allowed_tools` (String) — JSON-encoded list of allowed tools.
  - `approval_rules` (String) — JSON-encoded approval rules.
  - `date_created` (String)
  - `date_modified` (String)
