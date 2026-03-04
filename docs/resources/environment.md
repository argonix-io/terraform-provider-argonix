---
page_title: "argonix_environment Resource - Argonix"
description: |-
  Manages an Argonix test environment.
---

# argonix_environment (Resource)

Manages an Argonix test environment. Environments store key-value variables used in synthetic tests.

## Example Usage

```terraform
resource "argonix_environment" "staging" {
  name      = "Staging"
  variables = jsonencode({
    BASE_URL  = "https://staging.example.com"
    API_TOKEN = var.staging_api_token
  })
}

resource "argonix_environment" "production" {
  name       = "Production"
  is_default = true
  variables  = jsonencode({
    BASE_URL  = "https://api.example.com"
    API_TOKEN = var.prod_api_token
  })
}
```

## Schema

### Required

- `name` (String) — Name of the environment.

### Optional

- `variables` (String, Sensitive) — JSON-encoded key-value variables. Defaults to `"{}"`.
- `is_default` (Boolean) — Whether this is the default environment. Defaults to `false`.

### Read-Only

- `id` (String) — UUID of the environment.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_environment.example <environment-uuid>
```
