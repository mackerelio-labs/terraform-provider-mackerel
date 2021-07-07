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

* `name` - (Required) The name of role.
* `memo` - Notes related to this service.

## Attributes Reference

* `name` - The name of service.
* `memo` - Notes related to this service.

