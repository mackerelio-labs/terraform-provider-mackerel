---
page_title: "Mackerel: mackerel_alert_group_setting"
subcategory: "Alerts"
description: |-
---

# Data Source: mackerel_alert_group_setting

Use this data source allows access to details of a specific Alert Group.

## Example Usage

```terraform
data "mackerel_alert_group_setting" "this" {
  id = "example_id"
}
```

## Argument Reference

* `id` - The ID of alert group setting.

## Attributes Reference

* `memo` - Notes related to the alert group setting.
* `monitor_scopes` - An array of monitor IDs.
* `name` - The name of the alert group setting.
* `notification_interval` - The time interval (in minutes) for resending notifications.
* `role_scopes` - An array of the role's fullnames.
* `service_scopes` - An array of service names
