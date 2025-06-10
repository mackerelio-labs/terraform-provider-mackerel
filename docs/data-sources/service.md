---
page_title: "Mackerel: mackerel_service"
subcategory: "Service"
description: |-
---

# Data Source: mackerel_service

Use this data source allows access to details of a specific Service.  

## Example Usage

```terraform
data "mackerel_service" "foo" {
  name = "foo"
}
```

## Argument Reference

* `name` - (Required) The name of the service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `memo` - Notes related to this service.
* `roles` - List of roles in the service.
