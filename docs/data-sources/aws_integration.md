---
page_title: "Mackerel: mackerel_aws_integration"
subcategory: "Integrations"
description: |-

---

# Data Source: mackerel_aws_integration

Use this data source allows access to details of a specific aws integration setting.

## Example Usage

```terraform
data "mackerel_aws_integration" "foo" {
  id = "example_id"
}
```

## Argument Reference

* `id` - (Required) The ID of aws integration setting.

## Attributes Reference

* `name` - The name of aws integration.
* `memo` - Notes related to this aws integration.
* `key` - The AWS IAM user access key used for integration settings.
* `role_arn` - The AWS IAM role used for integration settings.
* `external_id` - This is an external ID used during integration configuration using the AWS IAM role.
* `region` - The region in which the integration will be enabled.
* `included_tags` - A list of tags to be included in the integration.
* `excluded_tags` - A list of tags to be removed from the integration.

### AWS Services

[See available AWS service identifiers(mackerel documentation) for details.](https://mackerel.io/api-docs/entry/aws-integration#awsServiceNames)

* `enable` - Whether integration settings are enabled. Default is `true`.
* `role` - The set of monitoring target’s service name or role name.
* `excluded_metrics` - 	Metrics to exclude from integration.
* `retire_automatically` - (EC2, RDS, ElastiCache and Lambda only) Whether automatic retirement is enabled.
