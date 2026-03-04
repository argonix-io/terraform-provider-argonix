---
page_title: "argonix_knowledge_base Data Source - Argonix"
description: |-
  Fetches a single Argonix knowledge base by ID.
---

# argonix_knowledge_base (Data Source)

Fetches a single Argonix knowledge base by ID.

## Example Usage

```terraform
data "argonix_knowledge_base" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the knowledge base.

### Read-Only

- `name` (String)
- `source_type` (String)
- `connector_id` (String)
- `is_active` (Boolean)
- `sync_config` (String) — JSON-encoded sync configuration.
- `last_synced_at` (String)
- `document_count` (Number)
- `date_created` (String)
- `date_modified` (String)
