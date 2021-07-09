---
page_title: "Mackerel: mackerel_role_metadata"
subcategory: "Role"
description: |-
---

# Data Source: mackerel_role_metadata

Use this data source allows access to details of a specific Role Metadata.  

## Example Usage

```terraform
data "mackerel_role_metadata" "bar" {
  service   = "foo"
  role      = "bar"
  namespace = "foo"
}
```

## Argument Reference

* `service` - (Required) The name of the service.
* `role` - (Required) The name of the role.
* `namespace` - (Required) Identifier for the metadata

## Attributes Reference

* `metadata_json` - Arbitrary JSON data for the role.
