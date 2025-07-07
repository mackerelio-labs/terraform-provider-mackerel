---
page_title: "Mackerel: mackerel_aws_integration"
subcategory: "Integrations"
description: |-

---

# Resource: mackerel_aws_integration

This resource allows creating and management of AWS Integration.

## Example Usage

```terraform
resource "mackerel_service" "foo" {
  name = "foo"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.name
  name    = "bar"
}

resource "mackerel_aws_integration" "baz" {
  name          = "baz"
  memo          = "This aws integration is managed by Terraform."
  key           = ""
  secret_key    = ""
  role_arn      = "arn:aws:iam::123456789012:role/mackerel-integration-role"
  external_id   = "jCymhqx4Xy88SrTpDMtoXo65Tj5vd2vcRiJiWfd9KUuM"
  region        = "ap-northeast-1"
  included_tags = "Name:staging-server,Environment:staging"
  excluded_tags = "Name:develop-server,Environment:develop"

  ec2 {
    enable               = true
    role                 = "${mackerel_service.foo.name}: ${mackerel_role.bar.name}"
    excluded_metrics     = []
    retire_automatically = true
  }

  alb {
    enable           = true
    role             = "${mackerel_service.foo.name}: ${mackerel_role.bar.name}"
    excluded_metrics = ["alb.request.count", "alb.bytes.processed"]
  }

  rds {
    role             = "${mackerel_service.foo.name}: ${mackerel_role.bar.name}"
    excluded_metrics = ["rds.cpu.used"]
  }

  nlb {
    enable = false
  }
}
```

## Argument Reference

* `name` - (Required) The name of aws integration.
* `memo` - Notes related to this aws integration.
* `key` - The AWS IAM user access key used for integration settings.
* `secret_key` - The AWS IAM user secret key used for integration settings.
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
* `retire_automatically` - (Services that support automatic retirement only) Whether automatic retirement is enabled.

## Attributes Reference

In addition to the above arguments except for the secret key, the following attributes are exported:

* `id` - The ID of aws integration setting.

## Import

AWS Integration setting can be imported using their ID, e.g.

```
$ terraform import mackerel_aws_integration.foo ABCDEFG
```
