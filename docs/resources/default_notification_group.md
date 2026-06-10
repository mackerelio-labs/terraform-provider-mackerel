---
page_title: "Mackerel: mackerel_default_notification_group"
subcategory: "Notifications"
description: |-

---

# Resource: mackerel_default_notification_group

This resource manages the default notification group settings.

The default notification group is an existing Mackerel notification group identified by the API type `group-default`.
This resource does not create or delete the default notification group.

## Example Usage

```terraform
resource "mackerel_channel" "email" {
  name = "email"

  email {
    emails = ["alice@example.com"]
    events = ["alert", "alertGroup"]
  }
}

resource "mackerel_notification_group" "example" {
  name = "example"
}

resource "mackerel_default_notification_group" "default" {
  notification_level = "critical"

  child_notification_group_ids = [
    mackerel_notification_group.example.id,
  ]

  child_channel_ids = [
    mackerel_channel.email.id,
  ]
}
```

To manage no child notification groups or channels, set the corresponding attribute to an empty list (`[]`).

## Argument Reference

* `notification_level` - (Optional) The level of notification ("all" or "critical". Default "all").
* `child_notification_group_ids` - (Required) A set of notification group IDs.
* `child_channel_ids` - (Required) A set of notification channel IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the default notification group.
