---
page_title: "argonix_connector Data Source - Argonix"
description: |-
  Fetches a single Argonix connector by ID.
---

# argonix_connector (Data Source)

Fetches a single Argonix connector by ID.

## Example Usage

```terraform
data "argonix_connector" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the connector.

### Read-Only

- `name` (String)
- `connector_type` (String)
- `is_active` (Boolean)
- `config` (String, Sensitive) — JSON-encoded configuration.
- `capabilities` (String) — JSON-encoded capabilities.
- `tags` (String) — JSON-encoded tags.
- `date_created` (String)
- `date_modified` (String)
