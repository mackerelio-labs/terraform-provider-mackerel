---
page_title: "Mackerel: mackerel_service_metadata"
subcategory: "Service"
description: |-
---

# Resource: mackerel_service_metadata

This resource allows creating and management of Service Metadata.

## Example Usage
```terraform
resource "mackerel_service" "foo" {
  name = "foo"
}

resource "mackerel_service_metadata" "foo" {
  service   = mackerel_service.foo.id
  namespace = "bar"

  metadata_json = jsonencode({
    id = 1
  })
}
```

## Argument Reference

* `service` - (Required) The name of the service.
* `namespace` - (Required) Identifier for the metadata
* `metadata_json` - (Required) Arbitrary JSON data for the service.

## Import

Service metadata can be imported using their <service_name>/<namespace>, e.g.

```
$ terraform import mackerel_service_metadata.foo foo:bar
```