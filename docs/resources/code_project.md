---
page_title: "argonix_code_project Resource - Argonix"
description: |-
  Manages a Git repository for defensive code analysis.
---

# argonix_code_project (Resource)

Manages a Git repository for defensive code analysis.

A code project configures a Git repository and its analysis policy. Use a canonical repository URL with no trailing slash or `.git` suffix. Analysis runs, credentials and manual review operations are managed in Argonix. Deleting a project requires the organization admin role and removes its analyses and findings.

## Example usage

```terraform
resource "argonix_code_project" "service" {
  name                  = "Service"
  repository_url        = "https://github.com/example/service"
  connector             = var.git_connector_id
  default_branch        = "main"
  enabled_engines       = jsonencode(["opengrep", "gitleaks"])
  webhook_analysis_mode = "standard"
}
```

## Schema

### Required

- `repository_url` (String) — Canonical repository URL without a .git suffix or trailing slash.

### Optional

- `name` (String) — Project name; defaults to the repository name.
- `key` (String) — Unique project key; generated from the name if omitted.
- `connector` (String) — GitHub or GitLab connector ID.
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
- `webhook_analysis_mode` (String) — Static analysis mode for repository webhooks: none or standard.

### Read-only

- `id` (String) — Argonix UUID.
- `repository_provider` (String)
- `entry_kind` (String)
- `latest_analysis` (String)
- `date_created` (String)
- `date_modified` (String)

## Import

```shell
terraform import argonix_code_project.example <uuid>
```
