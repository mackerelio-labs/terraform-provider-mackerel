---
page_title: "Mackerel: mackerel_service_metadata"
subcategory: "Service"
description: |-
---

# Data Source: mackerel_service_metadata

Use this data source allows access to details of a specific Service Metadata.  

## Example Usage

```terraform
resource "mackerel_service" "foo" {
  name = "foo"
}

resource "mackerel_service_metadata" "foo" {
  service = mackerel_service.foo.name
  namespace = "foo"
  metadata_json = jsonencode({
    id = 1
  })
}

data "mackerel_service_metadata" "foo" {
  service = mackerel_service_metadata.foo.service
  namespace = mackerel_service_metadata.foo.namespace
}
```

## Argument Reference

* `service` - (Required) The name of the service.
* `namespace` - (Required) Identifier for the metadata

## Attributes Reference

* `service` - The name of the service.
* `namespace` - Identifier for the metadata
* `metadata_json` - Arbitrary JSON data for the service.
