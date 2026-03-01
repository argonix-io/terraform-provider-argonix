---
page_title: "argonix_notification_rule Data Source - Argonix Provider"
---

# argonix_notification_rule (Data Source)

Fetches a single Argonix notification rule by ID.

## Example Usage

```hcl
data "argonix_notification_rule" "example" {
  id = "uuid-of-notification-rule"
}
```

## Argument Reference

- `id` (String, Required) — UUID of the notification rule to fetch.

## Attribute Reference

All attributes from the `argonix_notification_rule` resource are available as read-only computed values.
