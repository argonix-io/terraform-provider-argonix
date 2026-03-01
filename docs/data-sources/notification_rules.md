---
page_title: "argonix_notification_rules Data Source - Argonix Provider"
---

# argonix_notification_rules (Data Source)

Fetches all notification rules in the organization.

## Example Usage

```hcl
data "argonix_notification_rules" "all" {}

output "rule_names" {
  value = [for r in data.argonix_notification_rules.all.notification_rules : r.name]
}
```

## Attribute Reference

- `notification_rules` (List) — List of notification rule objects. Each has the same attributes as `argonix_notification_rule`.
