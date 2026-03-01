---
page_title: "argonix_status_pages Data Source - Argonix"
description: |-
  Fetches all status pages in the organization.
---

# argonix_status_pages (Data Source)

Fetches all status pages in the organization.

## Example Usage

```terraform
data "argonix_status_pages" "all" {}
```

## Schema

### Read-Only

- `status_pages` (List of Object) — List of all status pages. Each has:
  - `id` (String)
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
