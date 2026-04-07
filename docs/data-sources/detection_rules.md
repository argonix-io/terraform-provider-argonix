---
page_title: "argonix_detection_rules Data Source - Argonix Provider"
---

# argonix_detection_rules (Data Source)

Fetches all detection rules in the organization.

## Example Usage

```hcl
data "argonix_detection_rules" "all" {}

output "active_rules" {
  value = [for r in data.argonix_detection_rules.all.detection_rules : r.name if r.is_active]
}
```

## Attribute Reference

- `detection_rules` (List) — List of detection rule objects, each containing:
  - `id` (String) — UUID.
  - `name` (String) — Display name.
  - `description` (String) — Description.
  - `rule_type` (String) — Type: `threshold`, `pattern`, or `sequence`.
  - `severity` (String) — Severity level.
  - `is_active` (Boolean) — Whether the rule is active.
  - `config` (String) — JSON-encoded rule configuration.
  - `mitre_tactic` (String) — MITRE ATT&CK tactic ID.
  - `mitre_technique` (String) — MITRE ATT&CK technique ID.
  - `detections_count` (Integer) — Number of triggered detections.
  - `date_created` (String) — Creation timestamp.
  - `date_modified` (String) — Last modification timestamp.
