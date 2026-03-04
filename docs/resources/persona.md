---
page_title: "argonix_persona Resource - Argonix"
description: |-
  Manages an Argonix Argos persona.
---

# argonix_persona (Resource)

Manages an Argonix Argos persona. Personas define AI agent behaviour templates for different operational contexts.

## Example Usage

```terraform
resource "argonix_persona" "devops" {
  name          = "DevOps Agent"
  description   = "Handles infrastructure incidents"
  template      = "devops"
  system_prompt = "You are a DevOps engineer. Focus on infrastructure issues."
  connector_ids = jsonencode([argonix_connector.slack.id, argonix_connector.datadog.id])
}

resource "argonix_persona" "security" {
  name     = "Security Agent"
  template = "security"
  approval_rules = jsonencode({
    require_approval_for = ["write", "delete"]
  })
}
```

## Schema

### Required

- `name` (String) — Name of the persona.

### Optional

- `description` (String) — Description of the persona. Defaults to `""`.
- `template` (String) — Template preset. One of: `devops`, `it_support`, `hr`, `security`, `custom`. Defaults to `"custom"`.
- `system_prompt` (String) — Custom system prompt for the AI agent. Defaults to `""`.
- `is_active` (Boolean) — Whether the persona is active. Defaults to `true`.
- `connector_ids` (String) — JSON-encoded list of connector UUIDs. Defaults to `"[]"`.
- `allowed_tools` (String) — JSON-encoded list of allowed tool names. Defaults to `"[]"`.
- `approval_rules` (String) — JSON-encoded approval rules. Defaults to `"{}"`.

### Read-Only

- `id` (String) — UUID of the persona.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_persona.example <persona-uuid>
```
