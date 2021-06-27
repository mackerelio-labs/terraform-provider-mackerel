---
page_title: "Mackerel: mackerel_downtime"
subcategory: "Monitors"
description: |-

---

# Resource: mackerel_downtime

This resource allows creating and management of downtime.

## Example Usage

```terraform
resource "mackerel_downtime" "maintenance" {
  name     = "maintenance"
  start    = "1624814027"
  duration = 20

  service_scopes = ["web"]

  recurrence {
    interval = 1
    type     = "weekly"
    weekdays = ["Saturday", "Sunday"]
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) The name of the channel.
* `start` - (Required) The starting time in epoch seconds.
* `duration` - (Required) The duration of downtime in minutes.
* `memo` - Notes for the downtime.
* `monitor_scopes` - An array of monitor ids that scope of target monitor configurations.
* `monitor_exclude_scopes` - An array of excluded monitor ids that scope of target monitor configurations.
* `service_scopes` - An array of services that scope of target monitor configurations.
* `service_exclude_scopes` - An array of excluded services that scope of target monitor configurations.
* `role_scopes` - An array of roles that scope of target monitor configurations.
* `role_exclude_scopes` - An array of excluded roles that scope of target monitor configurations.
* `recurrence` - The configuration for recurrence. See [Recurrence](#recurrence) below for details.

### Recurrence

* `interval` - (Required) Recurrence interval.
* `type` - (Required) Recurrence options. Valid values are `hourly`, `daily`, `weekly`, `monthly` or `yearly`.
* `until` - The time at which recurrence ends in epoch seconds.
* `weekdays` - Configuration for the day of the week. Valid values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` or `Saturday`. Only available when the type is set to `weekly`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of downtime.

## Import

Downtime setting can be imported using their ID, e.g.

```
$ terraform import mackerel_downtime.this downtime_id
```
