---
page_title: "Mackerel: mackerel_dashboard"
subcategory: "Dashboard"
description: |-
---

# Data Source: mackerel_dashboard

Use this data source allows access to details of a specific dashboard.  

## Example Usage

```terraform
data "mackerel_dashboard" "this" {
  id = "example_id"
}
```

## Argument Reference

* `id` - (Required) The ID of dashboard.

## Attributes Reference

* `id` - The ID of dashboard.
* `title` - The title of dashboard.
* `memo` - The memo of dashboard.
* `url_path` - The URL path of dashboard.
* `created_at` - Creation time (epoch seconds).
* `updated_at` - Last update time (epoch seconds).
* `graph` - The graph widget.
  * `title` - The title of graph widget.
  * `host` - The host graph.
    * `host_id` - The ID of host.
    * `name` - The name of graph (e.g., "loadavg")
  * `role` - The role graph.
    * `role_fullname` - The service name and role name concatenated by `:`.
    * `name` - The name of graph (e.g., "loadav5").
    * `is_stacked` - Whether the graph is a stacked or line chart. If true, it will be a stacked graph.
  * `service` - The service graph.
    * `service_name` - The name of service.
    * `name` - The name of graph.
  * `expression` - The expression graph.
    * `expression` - The expression for graphs.
  * `query` - The query graph.
    * `query` - The PromQL-style query.
    * `legend` - The query legend.
  * `range` - The display period for graphs. If unspecified, it will be variable and the display period can be changed from the controller displayed at the top of the dashboard.
    * `relative` - （The period from (current time + `offset` - `period`) to (current time + `offset`) is displayed. Negative values for `offset` can be used to display graphs for a specified period in the past.
      * `period` - Duration (seconds).
      * `offset` - Difference from the current time (seconds).
    * `absolute` - The period from start to end is displayed.
      * `start` - Start time (epoch seconds).
      * `end` - End time (epoch seconds).
  * `layout` - The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and the y axis in the down direction as the positive direction.
    * `x` - The x coordinate of widget.
    * `y` - The y coordinate of widget.
    * `width` - The width of widget.
    * `height` - The height of widget.
* `value` - The value widget.
  * `title` - The title of value widget.
  * `metric` - The metric of value widget.
    * `host` - The host metric.
      * `host_id` - The ID of host.
      * `name` - The name of metric (e.g., "loadavg5").
    * `service` - The service metric.
      * `service_name` - The name of service.
      * `name` - The name of metric.
    * `expression` - The expression metric.
      * `expression` - The expression for metric.
    * `query` - The query metric.
      * `query` - The PromQL-style query.
      * `legend` - The query legend.
  * `fraction_size` - Number of decimal places to display (0-16).
  * `suffix` - Units to be displayed after the numerical value.
  * `layout` - The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and they axis in the down direction as the positive direction.
    * `x` - The x coordinate of widget.
    * `y` - The y coordinate of widget.
    * `width` - The width of widget.
    * `height` - The height of widget.
* `markdown` - The markdown widget.
  * `title` - The title of markdown widget.
  * `markdown` - String in Markdown format.
  * `layout` - The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and they axis in the down direction as the positive direction.
    * `x` - The x coordinate of widget.
    * `y` - The y coordinate of widget.
    * `width` - The width of widget.
    * `height` - The height of widget.
* `alert_status` - The alertStatus widget.
  * `title` - The title of alertStatus widget.
  * `role_fullname` - The service name and role name concatenated by `:`.
  * `layout` - The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and they axis in the down direction as the positive direction.
      * `x` - The x coordinate of widget.
      * `y` - The y coordinate of widget.
      * `width` - The width of widget.
      * `height` - The height of widget.
