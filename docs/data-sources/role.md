---
page_title: "Mackerel: mackerel_role"
subcategory: "Role"
description: |-
---

# Data Source: mackerel_role

Use this data source allows access to details of a specific Role.  

## Example Usage

```terraform
data "mackerel_role" "bar" {
  service = "foo"
  name    = "bar"
}
```

## Argument Reference

* `service` - (Required) The name of the service.
* `name` - (Required) The name of the role.

## Attributes Reference

* `service` - The name of the service.
* `name` - The name of the role.
* `memo` - Notes related to this role.
