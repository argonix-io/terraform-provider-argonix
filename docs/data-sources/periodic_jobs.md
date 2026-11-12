---
page_title: "argonix_periodic_jobs Data Source - Argonix"
description: |-
  Lists the Argonix periodic jobs visible to the provider's API key.
---

# argonix_periodic_jobs (Data Source)

Lists the Argonix periodic jobs visible to the provider's API key.

## Example usage

```terraform
data "argonix_periodic_jobs" "all" {}

output "periodic_job_ids" {
  value = [for item in data.argonix_periodic_jobs.all.periodic_jobs : item.id]
}
```

## Schema

- `periodic_jobs` (List of Object) — Objects with the attributes documented for the [single-object data source](periodic_job.md).
