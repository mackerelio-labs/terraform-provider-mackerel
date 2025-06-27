---
page_title: "Mackerel: mackerel_service"
subcategory: "Service"
description: |-
---

# Resource: mackerel_service

This resource allows creating and management of Service.

## Example Usage
```terraform
resource "mackerel_service" "foo" {
  name = "foo"
  memo = "Notes related to this service."
}
```

## Argument Reference

* `name` - (Required) The name of service.
* `memo` - Notes related to this service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `roles` - Set of roles in the service. This is a computed field and will be populated after the service is created.

## Import

Service setting can be imported using their name, e.g.

```
$ terraform import mackerel_service.foo name
```
