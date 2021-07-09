---
page_title: "Mackerel: mackerel_role"
subcategory: "Role"
description: |-
---

# Resource: mackerel_role

This resource allows creating and management of Role.

## Example Usage
```terraform
resource "mackerel_service" "foo" {
  name = "foo"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.name
  name    = "bar"
  memo    = "foo:bar is managed by Terraform"
}
```

## Argument Reference

* `name` - (Required) The name of role.
* `service` - (Required) The name of service.
* `memo` - Notes related to this role.

## Attributes Reference

No additional attributes are exported.

## Import

Role setting can be imported using their <service_name>:<role_name>, e.g.

```
$ terraform import mackerel_role.foo foo:bar
```