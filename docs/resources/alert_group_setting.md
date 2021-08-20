---
page_title: "Mackerel: mackerel_alert_group_setting"
subcategory: "Alerts"
description: |-

---

# Resource: mackerel_alert_group_setting

This resource allows creating and management of alert group setting.

## Example Usage

### Create Alert Group Setting

```terraform
resource "mackerel_alert_group_setting" "production" {
  name = "production alert group"

  service_scopes = [
    "web"
  ]
}
```

## Argument Reference

* `name` - The name of the alert group setting.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of alert group setting.
* `memo` - Notes related to the alert group setting.
* `monitor_scopes` - An array of monitor IDs.
* `notification_interval` - The time interval (in minutes) for resending notifications.
* `role_scopes` - An array of the role's fullnames.
* `service_scopes` - An array of service names

## Import

Alert group setting can be imported using their ID, e.g.

```
$ terraform import mackerel_alert_group_setting.production ABCDEFG
```
