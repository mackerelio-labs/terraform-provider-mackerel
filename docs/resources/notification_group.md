---
page_title: "Mackerel: mackerel_notification_group"
subcategory: "Notifications"
description: |-

---

# Resource: mackerel_notification_group

This resource allows creating and management of notification_group.

## Example Usage

```terraform
resource "mackerel_notification_group" "example" {
  name               = "Example notification group"
  notification_level = "critical"

  child_notification_group_ids = []
  child_channel_ids = ["2vh7AZ21abc"]

  monitor {
    id           = "2qtozU21abc"
    skip_default = false
  }

  service {
    name = "Example-Service-1"
  }
  service {
    name = "Example-Service-2"
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) The name of the notification group.
* `notification_level` - The level of notification ("all" or "critical". Default "all").
* `child_notification_group_ids` - A set of notification group IDs.
* `child_channel_ids` -  A set of notification channel IDs.
* `monitor` - Configuration block(s) with monitor rules. See [Monitor](#monitor) below for details.
* `service` - Configuration block(s) with services. See [Service](#service) below for details.

### Monitor

* `id` - (Required) The monitor rule ID.
* `skip_default` - If true, send notifications to this notification group only.

### Service

* `name` - (Required) The name of the service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of notification group.

## Import

Notification group setting can be imported using their ID, e.g.

```
$ terraform import mackerel_notification_group.this notification_group_id
```

