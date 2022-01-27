---
page_title: "Mackerel: mackerel_monitor"
subcategory: "Monitors"
description: |-

---

# Resource: mackerel_monitor

This resource allows creating and management of monitor.

## Example Usage

### host_metric

```terraform
resource "mackerel_monitor" "cpu_high" {
  name                  = "cpu % is high"
  is_mute               = false
  notification_interval = 10

  host_metric {
    metric   = "cpu%"
    operator = ">"
    warning  = 80
    critical = 90
    duration = 3

    scopes = ["app", "proxy:ec2"]
  }
}
```

### connectivity

```terraform
resource "mackerel_monitor" "connectivity" {
  name = "connectivity"

  connectivity {
    scopes = ["app", "proxy:ec2"]
  }
}
```

### service_metric

```terraform
resource "mackerel_monitor" "service_metric" {
  name = "error rate is high"

  service_metric {
    service  = "app"
    metric   = "execution.error.rate"
    operator = ">"
    duration = 3
    warning  = 2
  }
}
```

### external

```terraform
resource "mackerel_monitor" "external" {
  name                  = "Example Domain"
  notification_interval = 10

  external {
    method                 = "GET"
    url                    = "https://example.com"
    service                = "app"
    response_time_critical = 10000
    response_time_warning  = 5000
    response_time_duration = 3
    headers                = { Cache-Control = "no-cache" }
  }
}
```

### expression

```terraform
resource "mackerel_monitor" "role_avg" {
  name = "role average"

  expression {
    expression = "avg(roleSlots(service:role,loadavg5))"
    operator   = ">"
    warning    = 5
    critical   = 10
  }
}
```

### anomaly_detection

```terraform
resource "mackerel_monitor" "anomaly_detection" {
  name                  = "anomaly detection"
  memo                  = "my anomaly detection for roles"
  notification_interval = 10

  anomaly_detection {
    scopes              = ["myService: myRole"]
    warning_sensitivity = "insensitive"
    maxCheckAttempts    = 3
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) The name of the monitor.
* `memo` - The notes for the monitoring configuration.
* `is_mute` - Whether monitoring is muted or not. Valid values are `true` and `false`.
* `notification_interval` - The time interval for re-sending notifications in minutes. If empty, notifications will not be re-sent. Default is `0`.

### host_metric

* `metric` - (Required) The name of the host metric targeted by monitoring.
* `operator` - (Required) The comparison operator to determines the conditions that state whether the designated variable is either big or small. The observed value is on the left of the operator and the designated value is on the right. Valid values are `>` and `<`.
* `duration` - (Required) The duration of the monitor.
* `warning` - (Required, at least one of `warning` or `critical`) The threshold that generates a warning alert.
* `critical` - (Required, at least one of `warning` or `critical`) The threshold that generates a critical alert.
* `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made. Default is `1`.
* `scopes` - The set of monitoring target’s service name or role name.
* `exclude_scopes` - The set of monitoring exclusion target’s service name or role name.

### connectivity

* `scopes` - The set of monitoring target’s service name or role name.
* `exclude_scopes` - The set of monitoring exclusion target’s service name or role name.

### service_metric

* `service` - (Required) Name of the service targeted by monitoring.
* `metric` - (Required) The name of the host metric targeted by monitoring.
* `operator` - (Required) The comparison operator to determines the conditions that state whether the designated variable is either big or small. The observed value is on the left of the operator and the designated value is on the right. Valid values are `>` and `<`.
* `duration` - (Required) The duration of the monitor.
* `warning` - (Required, at least one of `warning` or `critical`) The threshold that generates a warning alert.
* `critical` - (Required, at least one of `warning` or `critical`) The threshold that generates a critical alert.
* `missing_duration_warning` - The threshold in minutes to generate a warning alert for interruption monitoring.
* `missing_duration_critical` - The threshold in minutes to generate a critical alert for interruption monitoring.
* `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made. Default is `1`.

### external

* `method` - (Required) Request method. Valid values are `GET`, `POST`, `PUT` or `DELETE`.
* `url` - (Required) Monitoring target URL.
* `service` - Service name. When response time is monitored, it will be graphed as the service metrics of this.
* `response_time_warning` - The response time threshold for warning alerts in milliseconds. Required with `service`.
* `response_time_critical` - The response time threshold for critical alerts in milliseconds. Required with `service`.
* `response_time_duration` - The duration to monitor the average of response time. Required with `service`.
* `request_body` - HTTP request body.
* `contains_string` - String which should be contained by the response body.
* `certification_expiration_warning` - Certification expiration date monitor’s “Warning” threshold. number of days remaining until expiration.
* `certification_expiration_critical` - Certification expiration date monitor’s “Critical” threshold. number of days remaining until expiration.
* `skip_certificate_verification` - Whether verify the certificate when monitoring a server with a self-signed certificate or not. Valid values are `true` and `false`.
* `headers` - The values configured as the HTTP request header.
* `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made. Default is `1`.
* `follow_redirect` - Evaluates the response of the redirector as a result. Valid values are `true` and `false`. Default is `false`.

### expression

* `expression` - (Required)
* `operator` - (Required) The comparison operator to determines the conditions that state whether the designated variable is either big or small. The observed value is on the left of the operator and the designated value is on the right. Valid values are `>` and `<`.
* `warning` - (Required, at least one of `warning` or `critical`) The threshold that generates a warning alert.
* `critical` - (Required, at least one of `warning` or `critical`) The threshold that generates a critical alert.

### anomaly_detection

* `scopes` - (Required) Expression of the monitoring target. Only valid for graph sequences that become one line.
* `warning_sensitivity` - (Required, at least one of `warning_sensitivity` or `critical_sensitivity`) The sensitivity to generates warning alerts. Valid values are `insensitive`, `normal` and `sensitive`.
* `critical_sensitivity` - (Required, at least one of `warning_sensitivity` or `critical_sensitivity`) The sensitivity to generates warning critical. Valid values are `insensitive`, `normal` and `sensitive`.
* `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made. Default is `1`.
* `training_period_from` - Epoch seconds. Anomaly detection use metric data starting from the specified time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of monitor.

## Import

Monitor setting can be imported using their ID, e.g.

```
$ terraform import mackerel_monitor.this monitor_id
```
