---
page_title: "argonix_detection_rule Data Source - Argonix Provider"
---

# argonix_detection_rule (Data Source)

Fetches a single Argonix detection rule by ID.

## Example Usage

```hcl
data "argonix_detection_rule" "brute_force" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "rule_type" {
  value = data.argonix_detection_rule.brute_force.rule_type
}
```

## Argument Reference

- `id` (String, Required) — UUID of the detection rule.

## Attribute Reference

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
