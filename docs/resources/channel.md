---
page_title: "Mackerel: mackerel_channel"
subcategory: "Notifications"
description: |-

---

# Resource: mackerel_channel

This resource allows creating and management of channel, which manages either email, slack or webhook.

## Example Usage

### Channel of email

```terraform
resource "mackerel_channel" "email" {
  name = "email"

  email {
    emails = ["alice@example.com"]
    events = ["alert", "alertGroup"]
  }
}
```

### Channel of slack

```terraform
resource "mackerel_channel" "slack" {
  name = "slack"

  slack {
    url                 = "https://hooks.slack.com/services/ABCD/12345"
    enabled_graph_image = true
    mentions = {
      critical = "@here"
    }
  }
}
```

### Channel of webhook

```terraform
resource "mackerel_channel" "webhook" {
  name = "webhook"

  webhook {
    url = "https://webhook.com/AAAAAAAA"
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - The name of the channel.

### email

* `emails` - A set of email addresses to receive notifications.
* `user_ids` - A set of user IDs to receive notifications.
* `events` - A set of notification events. Valid values are `alert` or `alertGroup`.

### slack

* `url` - Incoming Webhook URL for Slack.
* `mentions` - A map of mentions. Valid values are `ok`, `warning`, or `critical`.
* `enabled_graph_image` - A boolean value whether to post the corresponding graph. Default `falsse`.
* `events` - A set of notification events. Valid values are `alert` or `alertGroup`.

### webhook

* `url` - URL to receive HTTP request.
* `events` - A set of notification events. Valid values are `alert` or `alertGroup`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of channel.
* `name` - The name of channel.

## Import

Channel setting can be imported using their ID, e.g.

```
$ terraform import mackerel_channel.email ABCDEFG
```
