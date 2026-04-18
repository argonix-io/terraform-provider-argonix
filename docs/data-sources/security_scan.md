---
page_title: "argonix_security_scan Data Source - Argonix"
subcategory: "Security"
description: |-
  Fetches a single Argonix security scan by ID.
---

# argonix_security_scan (Data Source)

Fetches a single security scan by ID.

## Example Usage

```terraform
data "argonix_security_scan" "latest" {
  id = "scan-uuid-here"
}
```

## Argument Reference

- `id` - (Required) UUID of the security scan.

## Attribute Reference

- `scan_type` - Type of scan: `cspm`, `vulnerability`, or `secret`.
- `connector` - UUID of the scanned connector.
- `triggered_by` - How the scan was triggered.
- `status` - Scan status.
- `findings_count` - Number of findings.
- `summary` - JSON-encoded scan summary.
- `started_at` - Timestamp when the scan started.
- `completed_at` - Timestamp when the scan completed.
- `date_created` - Creation timestamp.
- `date_modified` - Last modification timestamp.
