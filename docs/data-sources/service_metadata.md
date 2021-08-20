---
page_title: "Mackerel: mackerel_service_metadata"
subcategory: "Service"
description: |-
---

# Data Source: mackerel_service_metadata

Use this data source allows access to details of a specific Service Metadata.  

## Example Usage

```terraform
data "mackerel_service_metadata" "foo" {
  service   = "foo"
  namespace = "bar"
}
```

## Argument Reference

* `service` - (Required) The name of the service.
* `namespace` - (Required) Identifier for the metadata

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `metadata_json` - Arbitrary JSON data for the service.
