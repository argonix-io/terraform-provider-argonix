---
page_title: "argonix_code_project Data Source - Argonix"
description: |-
  Reads an Argonix code project by UUID.
---

# argonix_code_project (Data Source)

Reads an Argonix code project by UUID.

## Example usage

```terraform
data "argonix_code_project" "example" {
  id = var.argonix_object_id
}
```

## Schema

### Required

- `id` (String) — Argonix UUID.

### Read-only

- `name` (String) — Project name; defaults to the repository name.
- `key` (String) — Unique project key; generated from the name if omitted.
- `repository_url` (String) — Canonical repository URL without a .git suffix or trailing slash.
- `connector` (String) — GitHub or GitLab connector ID.
- `repository_provider` (String)
- `entry_kind` (String)
- `provider_project_id` (String) — Project ID at the Git provider.
- `default_branch` (String) — Default branch.
- `visibility` (String) — private, internal or public.
- `quality_gate` (String) — Quality gate ID.
- `enabled_engines` (String, JSON array) — JSON array of static analysis engines.
- `publish_quality_gate` (Boolean) — Publish quality gate status to the Git provider.
- `publish_review_feedback` (Boolean) — Publish review feedback to the Git provider.
- `excluded_paths` (String, JSON array) — JSON array of excluded paths or glob patterns.
- `webhook_standard_scope` (String) — Standard analysis scope on the default branch.
- `webhook_standard_non_default_scope` (String) — Standard analysis scope on non-default branches.
- `webhook_standard_long_lived_scope` (String) — Standard analysis scope on long-lived branches.
- `webhook_long_lived_patterns` (String, JSON array) — JSON array of long-lived branch patterns.
- `webhook_default_gate_scope` (String) — all or new_code for the default branch.
- `webhook_non_default_gate_scope` (String) — all or new_code for non-default branches.
- `webhook_long_lived_gate_scope` (String) — all or new_code for long-lived branches.
- `long_lived_quality_gate` (String) — Quality gate ID for long-lived branches.
- `sql_dialect` (String) — SQL dialect required if SQLFluff is enabled.
- `auto_assign` (Boolean) — Assign findings automatically where supported.
- `is_active` (Boolean) — Whether the project is active.
- `latest_analysis` (String)
- `date_created` (String)
- `date_modified` (String)
- `webhook_analysis_mode` (String) — Static analysis mode for repository webhooks: none or standard.
