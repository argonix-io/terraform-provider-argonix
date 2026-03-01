---
page_title: "argonix_status_page Resource - Argonix"
description: |-
  Manages an Argonix status page.
---

# argonix_status_page (Resource)

Manages an Argonix status page. Status pages provide a public or private view of your service health for your users.

## Example Usage

```terraform
resource "argonix_status_page" "public" {
  name         = "Example Status"
  slug         = "example-status"
  visibility   = "public"
  accent_color = "#10B981"
  header_text  = "Service Status"
}
```

### Status Page with Custom Domain

```terraform
resource "argonix_status_page" "branded" {
  name            = "Acme Corp Status"
  slug            = "acme-status"
  custom_domain   = "status.acme.com"
  visibility      = "public"
  logo_url        = "https://cdn.acme.com/logo.png"
  favicon_url     = "https://cdn.acme.com/favicon.ico"
  accent_color    = "#FF6600"
  header_text     = "Acme Service Status"
  footer_text     = "© 2026 Acme Corp"
  meta_title      = "Acme Status"
  meta_description = "Real-time status of all Acme services."
  show_health_graph = true
}
```

## Schema

### Required

- `name` (String) — Name of the status page.
- `slug` (String) — URL slug for the status page. Must be unique.

### Optional

- `custom_domain` (String) — Custom domain for the status page. Defaults to `""`.
- `visibility` (String) — Visibility: `public` or `private`. Defaults to `"public"`.
- `logo_url` (String) — URL of the logo image. Defaults to `""`.
- `favicon_url` (String) — URL of the favicon. Defaults to `""`.
- `accent_color` (String) — Hex accent color. Defaults to `"#3B82F6"`.
- `custom_css` (String) — Custom CSS for the status page. Defaults to `""`.
- `header_text` (String) — Header text displayed on the page. Defaults to `""`.
- `footer_text` (String) — Footer text displayed on the page. Defaults to `""`.
- `meta_title` (String) — HTML meta title for SEO. Defaults to `""`.
- `meta_description` (String) — HTML meta description for SEO. Defaults to `""`.
- `show_health_graph` (Boolean) — Whether to show a health graph. Defaults to `false`.
- `is_active` (Boolean) — Whether the page is active. Defaults to `true`.

### Read-Only

- `id` (String) — UUID of the status page.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

```shell
terraform import argonix_status_page.example <status-page-uuid>
```
