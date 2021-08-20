---
page_title: "Mackerel: mackerel_notification_group"
subcategory: "Notifications"
description: |-
---

# Data Source: mackerel_notification_group

Use this data source allows access to details of a specific notification group setting.

## Example Usage

```terraform
data "mackerel_notification_group" "this" {
  id = "example_id"
}
```

## Argument Reference

* `id` - (Required) The ID of notification group.

## Attributes Reference

* `id` - The ID of notification group.
* `name` - The name of notification group.
* `notification_level` - The level of notification ("all" or "critical").
* `child_notification_group_ids` - A set of notification group IDs.
* `child_channel_ids` -  A set of notification channel IDs.
* `monitor` - A set of notification target monitor rules.
* `service` - A set of notification target services.
