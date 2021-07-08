---
page_title: "Mackerel: mackerel_monitor"
subcategory: "Monitors"
description: |-
---

# Data Source: mackerel_monitor

Use this data source allows access to details of a specific monitor.

## Example Usage

```terraform
data "mackerel_monitor" "this" {
  id = "example_id"
}
```

## Argument Reference

* `id` - The ID of monitor.

## Attributes Reference

* `name` - The name of the monitor.
* `memo` - The notes for the monitoring configuration.
* `is_mute` - Whether monitoring is muted or not.
* `notification_interval` - The time interval for re-sending notifications in minutes.
* `host_metric` - The settings for the monitor of host metric.
  * `metric` - The name of the host metric targeted by monitoring.
  * `operator` - The comparison operator to determines the conditions that state whether the designated variable is either big or small. The observed value is on the left of the operator and the designated value is on the right.
  * `duration` - The duration of the monitor.
  * `warning` - The threshold that generates a warning alert.
  * `critical` - The threshold that generates a critical alert.
  * `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made.
  * `scopes` - The set of monitoring target’s service name or role name.
  * `exclude_scopes` - The set of monitoring exclusion target’s service name or role name.
* `connectivity` - The settings for the monitor of connectivity.
  * `scopes` - The set of monitoring target’s service name or role name.
  * `exclude_scopes` - The set of monitoring exclusion target’s service name or role name.
* `service_metric` - The settings for the monitor of service metric.
  * `service` - Name of the service targeted by monitoring.
  * `metric` - The name of the host metric targeted by monitoring.
  * `operator` - The comparison operator to determines the conditions that state whether the designated variable is either big or small. The observed value is on the left of the operator and the designated value is on the right.
  * `duration` - The duration of the monitor.
  * `warning` - The threshold that generates a warning alert.
  * `critical` - The threshold that generates a critical alert.
  * `missing_duration_warning` - The threshold in minutes to generate a warning alert for interruption monitoring.
  * `missing_duration_critical` - The threshold in minutes to generate a critical alert for interruption monitoring.
  * `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made.
* `external` - The settings for the monitor of external URL monitoring.
  * `method` - Request method.
  * `url` - Monitoring target URL.
  * `service` - Service name. When response time is monitored, it will be graphed as the service metrics of this.
  * `response_time_warning` - The response time threshold for warning alerts in milliseconds.
  * `response_time_critical` - The response time threshold for critical alerts in milliseconds.
  * `response_time_duration` - The duration to monitor the average of response time.
  * `request_body` - HTTP request body.
  * `contains_string` - String which should be contained by the response body.
  * `certification_expiration_warning` - Certification expiration date monitor’s “Warning” threshold. number of days remaining until expiration.
  * `certification_expiration_critical` - Certification expiration date monitor’s “Critical” threshold. number of days remaining until expiration.
  * `skip_certificate_verification` - Whether verify the certificate when monitoring a server with a self-signed certificate or not.
  * `headers` - The values configured as the HTTP request header.
  * `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made.
* `expression` -  The settings for the monitor of expression monitoring.
  * `expression` - Expression of the monitoring target.
  * `operator` - The comparison operator to determines the conditions that state whether the designated variable is either big or small. The observed value is on the left of the operator and the designated value is on the right.
  * `warning` - The threshold that generates a warning alert.
  * `critical` - The threshold that generates a critical alert.
* `anomaly_detection` - The settings for the monitor of  Anomaly Detection for roles.
  * `scopes` - Expression of the monitoring target. Only valid for graph sequences that become one line.
  * `warning_sensitivity` - The sensitivity to generates warning alerts.
  * `critical_sensitivity` - The sensitivity to generates warning critical.
  * `max_check_attempts` - Number of consecutive Warning/Critical counts before an alert is made.
  * `training_period_from` - Epoch seconds. Anomaly detection use metric data starting from the specified time.
