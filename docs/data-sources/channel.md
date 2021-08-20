---
page_title: "Mackerel: mackerel_channel"
subcategory: "Notifications"
description: |-
---

# Data Source: mackerel_channel

Use this data source allows access to details of a specific Channel.  
You can get one of the following channels: email, slack or webhook.

## Example Usage

```terraform
data "mackerel_channel" "this" {
  id = "example_id"
}
```

## Argument Reference

* `id` - (Required) The ID of channel.

## Attributes Reference

* `id` - The ID of channel.
* `name` - The name of channel.
* `slack` - The list including `url`, `mentions`, `enabled_graph_image` and `events`.
* `webhook` - The list including `url` and `events`.
* `email` - The list including `emails`, `user_ids` and `events`.
