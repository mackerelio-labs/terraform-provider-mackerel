---
page_title: "Mackerel: mackerel_ip_ranges"
subcategory: "Metadata"
description: |-
---

# Data Source: mackerel_ip_ranges

Use this data source to get the IP address ranges used by Mackerel, provided by the [Metadata API](https://mackerel.io/docs/entry/spec/metadata).

The ranges are fetched from `https://mackerel.io/meta/ip-ranges.json`, which does not require an API key.

## Example Usage

Allow requests from Mackerel (e.g. external monitoring and Webhook notifications) in an AWS security group.

```terraform
data "mackerel_ip_ranges" "this" {}

resource "aws_security_group_rule" "allow_mackerel" {
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = data.mackerel_ip_ranges.this.outgoing
  ipv6_cidr_blocks  = data.mackerel_ip_ranges.this.outgoing_ipv6
  security_group_id = aws_security_group.example.id
}
```

## Argument Reference

This data source has no arguments.

## Attributes Reference

* `id` - The URL of the Metadata API endpoint.
* `incoming` - List of IPv4 address ranges that Mackerel receives.
* `incoming_ipv6` - List of IPv6 address ranges that Mackerel receives.
* `outgoing` - List of IPv4 address ranges that Mackerel sends from, such as alert notifications (e.g. Webhook).
* `outgoing_ipv6` - List of IPv6 address ranges that Mackerel sends from, such as alert notifications (e.g. Webhook).
