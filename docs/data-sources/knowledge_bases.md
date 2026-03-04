---
page_title: "argonix_knowledge_bases Data Source - Argonix"
description: |-
  Fetches all knowledge bases in the organization.
---

# argonix_knowledge_bases (Data Source)

Fetches all knowledge bases in the organization.

## Example Usage

```terraform
data "argonix_knowledge_bases" "all" {}
```

## Schema

### Read-Only

- `knowledge_bases` (List of Object) — List of all knowledge bases. Each has:
  - `id` (String)
  - `name` (String)
  - `source_type` (String)
  - `connector_id` (String)
  - `is_active` (Boolean)
  - `sync_config` (String) — JSON-encoded sync configuration.
  - `last_synced_at` (String)
  - `document_count` (Number)
  - `date_created` (String)
  - `date_modified` (String)
