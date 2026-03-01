---
page_title: "argonix_status_page Data Source - Argonix"
description: |-
  Fetches a single Argonix status page by ID.
---

# argonix_status_page (Data Source)

Fetches a single Argonix status page by ID.

## Example Usage

```terraform
data "argonix_status_page" "example" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- `id` (String) — UUID of the status page.

### Read-Only

- `name` (String)
- `slug` (String)
- `custom_domain` (String)
- `visibility` (String)
- `logo_url` (String)
- `favicon_url` (String)
- `accent_color` (String)
- `custom_css` (String)
- `header_text` (String)
- `footer_text` (String)
- `meta_title` (String)
- `meta_description` (String)
- `show_health_graph` (Boolean)
- `is_active` (Boolean)
- `date_created` (String)
- `date_modified` (String)
