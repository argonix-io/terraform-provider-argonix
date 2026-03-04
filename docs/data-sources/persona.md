---
page_title: "argonix_persona Data Source - Argonix"
description: |-
  Fetches a single Argonix persona by ID.
---

# argonix_persona (Data Source)

Fetches a single Argonix Argos persona by ID.

## Example Usage

```terraform
data "argonix_persona" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the persona.

### Read-Only

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
