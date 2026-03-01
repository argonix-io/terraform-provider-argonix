---
page_title: "Argonix Provider"
description: |-
  The Argonix provider is used to manage monitoring, alerting, status pages, and test management resources on the Argonix platform.
---

# Argonix Provider

The Argonix provider allows you to manage your [Argonix](https://argonix.io) monitoring infrastructure as code. It supports monitors, synthetic tests, alert rules, alert channels, status pages, and test management resources.

The organization is automatically determined from the API key — no need to specify it separately.

## Authentication

Create an API key from the Argonix dashboard under **Settings → API Keys**. The key must start with `ax_`.

You can provide the key directly in the provider configuration or via the `ARGONIX_API_KEY` environment variable.

## Example Usage

```terraform
terraform {
  required_providers {
    argonix = {
      source = "argonix-io/argonix"
    }
  }
}

provider "argonix" {
  api_key = var.argonix_api_key
}

variable "argonix_api_key" {
  type      = string
  sensitive = true
}
```

Using environment variables:

```bash
export ARGONIX_API_KEY="ax_..."
export ARGONIX_URL="https://api.argonix.io"  # optional
```

## Schema

### Required

- `api_key` (String, Sensitive) — API key for authenticating with Argonix. The organization is automatically determined from the key. Can also be set via the `ARGONIX_API_KEY` environment variable.

### Optional

- `url` (String) — Base URL of the Argonix API. Defaults to `https://api.argonix.io`. Can also be set via the `ARGONIX_URL` environment variable.
