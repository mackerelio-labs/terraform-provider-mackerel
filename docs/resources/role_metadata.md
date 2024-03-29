---
page_title: "Mackerel: mackerel_role_metadata"
subcategory: "Role"
description: |-
---

# Resource: mackerel_role_metadata

This resource allows creating and management of Role Metadata.

## Example Usage
```terraform
resource "mackerel_service" "foo" {
  name = "foo"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.id
  name    = "bar"
}

resource "mackerel_role_metadata" "bar" {
  service   = mackerel_service.foo.name
  role      = mackerel_role.bar.name
  namespace = "bar"

  metadata_json = jsonencode({
    id = 1
  })
}
```

## Argument Reference

* `service` - (Required) The name of the service.
* `role` - (Required) The name of role.
* `namespace` - (Required) Identifier for the metadata
* `metadata_json` - (Required) Arbitrary JSON data for the service.

## Attributes Reference

No additional attributes are exported.

## Import

Role metadata can be imported using their <service_name>:<role_name>/<metadata>, e.g.

```
$ terraform import mackerel_role.foo foo:bar/bar
```