---
page_title: "Mackerel: mackerel_service_metric_names"
subcategory: "Service"
description: |-
---

# Data Source: mackerel_service_metric_names

Use this data source allows access to details of a specific Service metric names.

## Example Usage

All of the service metric names.

```terraform
data "mackerel_service_metric_names" "foo" {
  name = "foo"  // service name
}
```

Filter by prefix for the meric names.

```terraform
data "mackerel_service_metric_names" "foo-xxx" {
  name   = "foo"        // service name
  prefix = "custom.xxx" // prefix of the service metric names
}
```

## Argument Reference

* `name` - (Required) The name of the service.
* `prefix` - Prefix of the metric names.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `metric_names` - Set of the service metric names.
