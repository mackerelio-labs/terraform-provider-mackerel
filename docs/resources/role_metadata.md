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
  name = "bar"
}

resource "mackerel_role_metadata" "bar" {
  service = mackerel_service.foo.name
  role = mackerel_role.bar.name
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

* `service` - The name of the service.
* `role` - The name of the role.
* `namespace` - Identifier for the metadata
* `metadata_json` - Arbitrary JSON data for the service.

