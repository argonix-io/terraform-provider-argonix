---
page_title: "argonix_workflow Resource - Argonix"
description: |-
  Manages an Argonix Argos workflow.
---

# argonix_workflow (Resource)

Manages an Argonix Argos workflow. Workflows define automated sequences of AI-driven actions for incident response, onboarding, security, and more.

## Example Usage

```terraform
resource "argonix_workflow" "restart_pod" {
  name     = "Restart Failing Pod"
  slug     = "restart-failing-pod"
  category = "devops"
  steps    = jsonencode([
    { action = "identify_pod", connector_type = "kubernetes" },
    { action = "restart_pod", connector_type = "kubernetes" },
    { action = "notify", connector_type = "slack" }
  ])
  required_connector_types = jsonencode(["kubernetes", "slack"])
  requires_confirmation    = true
}

resource "argonix_workflow" "restart_pod_documented" {
  name     = "Restart Failing Pod"
  slug     = "restart-failing-pod-doc"
  category = "devops"
  steps    = jsonencode([
    { id = "restart", tool = "k8s.delete_pod", params = { namespace = "{{input.namespace}}", name = "{{input.pod_name}}" } }
  ])
  # input_hints only annotates the {{input.x}} referenced in steps (the input
  # SET is always derived from steps — hints never define it).
  input_hints = jsonencode({
    namespace = { description = "Kubernetes namespace of the pod", example = "prod" }
    pod_name  = { description = "Name of the pod to restart", example = "api-gateway-7d9f" }
  })
  required_connector_types = jsonencode(["kubernetes"])
}

resource "argonix_workflow" "daily_report" {
  name     = "Daily Status Report"
  slug     = "daily-status-report"
  category = "general"
  schedule = "0 9 * * *"
}
```

## Schema

### Required

- `name` (String) — Name of the workflow.

### Optional

- `slug` (String) — URL-friendly slug. Defaults to `""`.
- `description` (String) — Description of the workflow. Defaults to `""`.
- `category` (String) — Category. One of: `identity`, `incident`, `onboarding`, `devops`, `security`, `general`. Defaults to `"general"`.
- `steps` (String) — JSON-encoded workflow steps. Defaults to `"[]"`.
- `input_hints` (String) — JSON-encoded map annotating each `{{input.x}}` referenced in `steps` with an optional `description` and `example`. It does **not** define which inputs exist (that set is always derived from the `{{input.x}}` references in `steps`); orphan hints are ignored. Defaults to `"{}"`.
- `required_connector_types` (String) — JSON-encoded list of required connector types. Defaults to `"[]"`.
- `requires_confirmation` (Boolean) — Whether human confirmation is required. Defaults to `true`.
- `schedule` (String) — Cron expression for scheduled execution. Defaults to `""`.
- `is_active` (Boolean) — Whether the workflow is active. Defaults to `true`.

### Read-Only

- `id` (String) — UUID of the workflow.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_workflow.example <workflow-uuid>
```
