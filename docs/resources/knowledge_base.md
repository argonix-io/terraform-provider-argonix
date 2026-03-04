---
page_title: "argonix_knowledge_base Resource - Argonix"
description: |-
  Manages an Argonix Argos knowledge base.
---

# argonix_knowledge_base (Resource)

Manages an Argonix Argos knowledge base. Knowledge bases provide contextual information to AI personas through manual entries or synced external sources.

## Example Usage

```terraform
resource "argonix_knowledge_base" "runbooks" {
  name        = "Runbooks"
  source_type = "confluence"
  connector_id = argonix_connector.confluence.id
  sync_config = jsonencode({
    space_key    = "OPS"
    sync_interval = "1h"
  })
}

resource "argonix_knowledge_base" "manual" {
  name        = "Internal Procedures"
  source_type = "manual"
}
```

## Schema

### Required

- `name` (String) — Name of the knowledge base.

### Optional

- `source_type` (String) — Source type. One of: `manual`, `confluence`, `notion`, `github`, `gitlab`, `web`. Defaults to `"manual"`.
- `connector_id` (String) — UUID of the connector used for syncing.
- `is_active` (Boolean) — Whether the knowledge base is active. Defaults to `true`.
- `sync_config` (String) — JSON-encoded sync configuration. Defaults to `"{}"`.

### Read-Only

- `id` (String) — UUID of the knowledge base.
- `last_synced_at` (String) — Timestamp of last sync.
- `document_count` (Number) — Number of documents in the knowledge base.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_knowledge_base.example <knowledge-base-uuid>
```
