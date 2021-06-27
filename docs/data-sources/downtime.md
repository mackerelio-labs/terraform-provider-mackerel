---
page_title: "Mackerel: mackerel_downtime"
subcategory: "Monitors"
description: |-
---

# Data Source: mackerel_downtime

Use this data source allows access to details of a specific downtime setting.

## Example Usage

```terraform
data "mackerel_downtime" "this" {
  id = "example_id"
}
```

## Argument Reference

* `id` - (Required) The ID of downtime.

## Attributes Reference

* `id` - The ID of downtime.
* `name` - The name of downtime.
* `memo` - Notes for the downtime.
* `duration` - The duration of downtime (in minutes).
* `monitor_scopes` - The set of monitor ids that scope of target monitor configurations.
* `monitor_exclude_scopes` - The set of excluded monitor ids that scope of target monitor configurations.
* `service_scopes` - The set of services that scope of target monitor configurations.
* `service_exclude_scopes` - The set of excluded services that scope of target monitor configurations.
* `role_scopes` - The set of roles that scope of target monitor configurations.
* `role_exclude_scopes` - The set of excluded roles that scope of target monitor configurations.
* `start` - The starting time (in epoch seconds).
* `recurrence` - The configuration for recurrence.
