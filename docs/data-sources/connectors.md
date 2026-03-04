---
page_title: "argonix_connectors Data Source - Argonix"
description: |-
  Fetches all connectors in the organization.
---

# argonix_connectors (Data Source)

Fetches all connectors in the organization.

## Example Usage

```terraform
data "argonix_connectors" "all" {}
```

## Schema

### Read-Only

- `connectors` (List of Object) — List of all connectors. Each has:
  - `id` (String)
  - `name` (String)
  - `connector_type` (String)
  - `is_active` (Boolean)
  - `config` (String, Sensitive) — JSON-encoded configuration.
  - `capabilities` (String) — JSON-encoded capabilities.
  - `tags` (String) — JSON-encoded tags.
  - `date_created` (String)
  - `date_modified` (String)
