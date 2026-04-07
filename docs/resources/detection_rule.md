---
page_title: "argonix_detection_rule Resource - Argonix Provider"
---

# argonix_detection_rule (Resource)

Manages an Argonix threat detection rule. Detection rules define patterns that are evaluated against ingested incidents to automatically detect threats and map them to MITRE ATT&CK tactics.

## Example Usage

### Threshold rule (brute force detection)

```hcl
resource "argonix_detection_rule" "brute_force" {
  name        = "Brute Force Detection"
  rule_type   = "threshold"
  severity    = "high"
  config = jsonencode({
    field          = "severity"
    value          = "critical"
    count          = 5
    window_minutes = 10
  })
  mitre_tactic    = "credential-access"
  mitre_technique = "T1110"
}
```

### Pattern rule (credential stuffing)

```hcl
resource "argonix_detection_rule" "credential_stuffing" {
  name      = "Credential Stuffing"
  rule_type = "pattern"
  severity  = "critical"
  config = jsonencode({
    field = "title"
    regex = "brute.force|credential.stuff|login.fail"
  })
  mitre_tactic    = "credential-access"
  mitre_technique = "T1110.004"
}
```

### Sequence rule (kill chain)

```hcl
resource "argonix_detection_rule" "kill_chain" {
  name      = "Recon-to-Exfil Kill Chain"
  rule_type = "sequence"
  severity  = "critical"
  config = jsonencode({
    tactics        = ["reconnaissance", "initial-access", "exfiltration"]
    window_minutes = 60
  })
  mitre_tactic = "exfiltration"
}
```

## Argument Reference

- `name` (String, Required) — Display name of the detection rule.
- `rule_type` (String, Required) — Type of detection: `threshold`, `pattern`, or `sequence`.
- `config` (String, Required) — JSON-encoded rule configuration. Format depends on `rule_type`:
  - **threshold**: `{"field": "severity", "value": "critical", "count": 5, "window_minutes": 10}`
  - **pattern**: `{"field": "title", "regex": "brute.force|credential.stuff"}`
  - **sequence**: `{"tactics": ["reconnaissance", "initial-access", "exfiltration"], "window_minutes": 60}`
- `description` (String, Optional) — Description of the rule purpose. Default: `""`.
- `severity` (String, Optional) — Severity level: `critical`, `high`, `medium`, or `low`. Default: `"medium"`.
- `is_active` (Boolean, Optional) — Whether the rule is active. Default: `true`.
- `mitre_tactic` (String, Optional) — MITRE ATT&CK tactic ID (e.g. `initial-access`, `lateral-movement`). Default: `""`.
- `mitre_technique` (String, Optional) — MITRE ATT&CK technique ID (e.g. `T1078`, `T1059`). Default: `""`.

## Attribute Reference

- `id` (String) — UUID of the detection rule.
- `detections_count` (Integer) — Number of threat detections triggered by this rule.
- `date_created` (String) — Creation timestamp.
- `date_modified` (String) — Last modification timestamp.

## Import

Detection rules can be imported using their UUID:

```shell
terraform import argonix_detection_rule.brute_force <uuid>
```
