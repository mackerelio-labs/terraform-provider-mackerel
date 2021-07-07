---
page_title: "Mackerel: mackerel_role_metadata"
subcategory: "Role"
description: |-
---

# Data Source: mackerel_role_metadata

Use this data source allows access to details of a specific Role Metadata.  

## Example Usage

```terraform
resource "mackerel_service" "foo" {
  name = "foo"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.name
  name = "bar"
}

resource "mackerel_role_metadata" "bar" {
  service = mackerel_role.foo.service
  role = mackerel_role.bar.name
  namespace = "bar"
  metadata_json = jsonencode({
    id = 1
  })
}

data "mackerel_role_metadata" "bar" {
  service = mackerel_role_metadata.foo.service
  role = mackerel_role_metadata.bar.role
  namespace = mackerel_role_metadata.bar.namespace
}
```

## Argument Reference

* `service` - (Required) The name of the service.
* `role` - (Required) The name of the role.
* `namespace` - (Required) Identifier for the metadata

## Attributes Reference

* `service` - The name of the service.
* `role` - The name of the role.
* `namespace` - Identifier for the metadata
* `metadata_json` - Arbitrary JSON data for the role.
