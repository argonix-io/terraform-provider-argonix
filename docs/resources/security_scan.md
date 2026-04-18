---
page_title: "argonix_security_scan Resource - Argonix"
subcategory: "Security"
description: |-
  Triggers an Argonix security scan (CSPM, vulnerability, or secret scan).
---

# argonix_security_scan (Resource)

Triggers an Argonix security scan. Scans are immutable — any configuration change forces a new scan.

## Example Usage

```terraform
# Scan all connectors
resource "argonix_security_scan" "cspm_scan" {
  scan_type = "cspm"
}

# Scan a specific connector
resource "argonix_security_scan" "vuln_scan" {
  scan_type = "vulnerability"
  connector = argonix_connector.aws.id
}
```

## Argument Reference

- `scan_type` - (Required) Type of scan. One of: `cspm`, `vulnerability`, `secret`.
- `connector` - (Optional) UUID of the connector to scan. When omitted, all connectors are scanned.

## Attribute Reference

- `id` - UUID of the scan.
- `status` - Scan status: `pending`, `running`, `completed`, or `failed`.
- `triggered_by` - How the scan was triggered (`manual` or `scheduled`).
- `findings_count` - Number of findings from the scan.
- `summary` - JSON-encoded scan summary.
- `started_at` - Timestamp when the scan started.
- `completed_at` - Timestamp when the scan completed.
- `date_created` - Creation timestamp.
- `date_modified` - Last modification timestamp.
