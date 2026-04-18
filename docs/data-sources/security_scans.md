---
page_title: "argonix_security_scans Data Source - Argonix"
subcategory: "Security"
description: |-
  Fetches all security scans in the organization.
---

# argonix_security_scans (Data Source)

Fetches all security scans in the organization.

## Example Usage

```terraform
data "argonix_security_scans" "all" {}
```

## Attribute Reference

- `security_scans` - List of security scans. Each scan has:
  - `id` - UUID of the scan.
  - `scan_type` - Type of scan.
  - `connector` - UUID of the scanned connector.
  - `triggered_by` - How the scan was triggered.
  - `status` - Scan status.
  - `findings_count` - Number of findings.
  - `summary` - JSON-encoded scan summary.
  - `started_at` - Timestamp when the scan started.
  - `completed_at` - Timestamp when the scan completed.
  - `date_created` - Creation timestamp.
  - `date_modified` - Last modification timestamp.
