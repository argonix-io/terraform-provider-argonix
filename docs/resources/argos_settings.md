---
page_title: "argonix_argos_settings Resource - Argonix"
description: |-
  Manages Argos AI agent settings for the organization.
---

# argonix_argos_settings (Resource)

Manages the Argos AI agent settings for your organization. This is a **singleton** resource — each organization has exactly one.

Use this resource to configure which LLM provider Argos uses, bring your own API key, point to a custom endpoint (vLLM, Azure OpenAI, self-hosted Ollama), and set custom system prompt instructions.

## Example Usage

### Use a cloud provider

```terraform
resource "argonix_argos_settings" "main" {
  llm_provider = "google"
  llm_model    = "gemini-2.5-flash"

  custom_instructions = "Always respond in French. Prioritize Kubernetes investigation."
}
```

### Bring Your Own Model (self-hosted vLLM)

```terraform
resource "argonix_argos_settings" "main" {
  llm_provider = "openai"  # vLLM exposes an OpenAI-compatible API
  llm_model    = "meta-llama/Llama-3.3-70B-Instruct"
  llm_api_key  = var.vllm_api_key
  llm_base_url = "https://vllm.internal.company.com/v1"

  custom_instructions = "Follow our runbook procedures. Never delete resources without confirmation."
}
```

### Bring Your Own Model (self-hosted Ollama)

```terraform
resource "argonix_argos_settings" "main" {
  llm_provider = "local"
  llm_model    = "qwen3.5:32b"
  llm_base_url = "http://ollama.prod.internal:11434/v1"
}
```

### Bring Your Own Key (cloud provider)

```terraform
resource "argonix_argos_settings" "main" {
  llm_provider = "anthropic"
  llm_model    = "claude-sonnet-4-6-20260120"
  llm_api_key  = var.anthropic_api_key

  custom_instructions = "Use structured output. Always include remediation steps."
}
```

## Schema

### Optional

- `llm_provider` (String) — LLM provider. One of: `local`, `google`, `anthropic`, `openai`. Defaults to `"google"`.
- `llm_model` (String) — Model override (blank = default for the provider). Any model name is accepted when `llm_base_url` is set.
- `llm_api_key` (String, Sensitive) — Custom API key for the LLM provider. Leave empty to use the platform default.
- `llm_base_url` (String) — Custom base URL for the LLM endpoint (e.g. vLLM, Azure OpenAI, self-hosted Ollama).
- `custom_instructions` (String) — Custom system prompt instructions prepended to every Argos conversation.
- `demo_mode` (Boolean) — When enabled, Argos returns scripted demo responses instead of calling the LLM.

### Read-Only

- `id` (String) — UUID of the settings object.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_argos_settings.main <settings-uuid>
```

~> **Note:** Since this is a singleton, you can retrieve the settings UUID via the API:
```bash
curl -s https://api.argonix.io/api/0.1/organizations/{org_id}/argos/settings/ \
  -H "Authorization: Api-Key ax_YOUR_KEY" | jq -r '.id'
```
