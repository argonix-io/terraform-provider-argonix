---
page_title: "argonix_code_projects Data Source - Argonix"
description: |-
  Lists the Argonix code projects visible to the provider's API key.
---

# argonix_code_projects (Data Source)

Lists the Argonix code projects visible to the provider's API key.

## Example usage

```terraform
data "argonix_code_projects" "all" {}

output "code_project_ids" {
  value = [for item in data.argonix_code_projects.all.code_projects : item.id]
}
```

## Schema

- `code_projects` (List of Object) — Objects with the attributes documented for the [single-object data source](code_project.md).
