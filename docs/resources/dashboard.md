---
page_title: "Mackerel: mackerel_dashboard"
subcategory: "Dashboard"
description: |-

---

# Resource: mackerel_dashboard

This resource allows creating and management of dashboard.

## Example Usage

### Graph Widget

```terraform
resource "mackerel_service" "foo" {
	name = "tf-service-foo"
}
	
resource "mackerel_role" "foo" {
	service = mackerel_service.foo.name
	name    = "tf-role-foo"
}

resource "mackerel_dashboard" "graph" {
  title = "foo"
  memo = "This dashboard is managed by Terraform."
  url_path = "bar"
	graph {
		title = "graph role"
		role {
			role_fullname = "${mackerel_service.foo.name}:${mackerel_role.foo.name}"
			name = "loadavg5"
			is_stacked = true
		}
		range {
			relative {
				period = 3600
				offset = 1800
			}
		}
		layout {
			x = 2
			y = 12
			width = 10
			height = 8
		}
	}
}
```

### Value Widget

```terraform
resource "mackerel_service" "foo" {
	name = "tf-service-foo"
}
	
resource "mackerel_role" "foo" {
	service = mackerel_service.foo.name
	name    = "tf-role-foo"
}

resource "mackerel_dashboard" "value" {
  title = "foo"
  memo = "This dashboard is managed by Terraform."
  url_path = "bar"
	value {
    title = "test value expression"
    metric {
			expression {
				expression = "role(${mackerel_service.foo.name}:${mackerel_role.foo.name}, loadavg5)"
			}
		}
		fraction_size = 5
		suffix = "test suffix"
    layout {
			x = 3
			y = 15
			width = 3
			height = 4
		}
  }
}
```

### Markdown Widget

```terraform
resource "mackerel_dashboard" "markdown" {
  title = "foo"
  memo = "This dashboard is managed by Terraform."
  url_path = "bar"
  markdown {
    title = "test markdown"
    markdown = "# h1"
    layout {
			x = 1
			y = 2
			width = 3
			height = 4
		}
  }
}
```

### Expression Widget

```terraform
resource "mackerel_service" "foo" {
	name = "tf-service-foo"
}
	
resource "mackerel_role" "foo" {
	service = mackerel_service.foo.name
	name    = "tf-role-foo"
}

resource "mackerel_dashboard" "expression" {
  title = "foo"
  memo = "This dashboard is managed by Terraform."
  url_path = "bar"
  alert_status {
    title = "test alertStatus"
    role_fullname = "${mackerel_service.foo.name}:${mackerel_role.foo.name}"
    layout {
			x = 1
			y = 2
			width = 3
			height = 4
		}
  }
}
```

## Argument Reference

* `title` - (Required) The title of dashboard.
* `memo` - The memo of dashboard.
* `url_path` - The URL path of dashboard.

### graph

* `title` - The title of graph widget.
* `host` - The host graph.
  * `host_id` - (Required) The ID of host.
  * `name` - (Required) The name of graph (e.g., "loadavg")
* `role` - The role graph.
  * `role_fullname` - (Required) The service name and role name concatenated by `:`.
  * `name` - (Required) The name of graph (e.g., "loadav5").
  * `is_stacked` - Whether the graph is a stacked or line chart. If true, it will be astacked graph.
* `service` - The service graph.
  * `service_name` - (Required) The name of service.
  * `name` - (Required) The name of graph.
* `expression` - The expression graph.
  * `expression` - (Required) The expression for graphs.
* `range` - The display period for graphs. If unspecified, it will be variable and thedisplay period can be changed from the controller displayed at the top of the dashboard.
  * `relative` - （The period from (current time + `offset` - `period`) to (current time + `offset`) is displayed. Negative values for `offset` can be used to display graphs for a specified period in the past.
    * `period` - (Required) Duration (seconds).
    * `offset` - (Required) Difference from the current time (seconds).
  * `absolute` - The period from start to end is displayed.
    * `start` - (Required) Start time (epoch seconds).
    * `end` - (Required) End time (epoch seconds).
* `layout` - (Required) The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and the y axis in the down direction as the positive direction.
  * `x` - (Required) The x coordinate of widget.
  * `y` - (Required) The y coordinate of widget.
  * `width` - (Required) The width of widget.
  * `height` - (Required) The height of widget.

### value

* `title` - The title of value widget.
* `metric` - The metric of value widget.
  * `host` - The host metric.
    * `host_id` - (Required) The ID of host.
    * `name` - (Required) The name of metric (e.g., "loadavg5").
  * `service` - The service metric.
    * `service_name` - (Required) The name of service.
    * `name` - (Required) The name of metric.
  * `expression` - The expression metric.
    * `expression` - (Required) The expression for metric.
* `fraction_size` - Number of decimal places to display (0-16).
* `suffix` - Units to be displayed after the numerical value.
* `layout` - (Required) The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and they axis in the down direction as the positive direction.
  * `x` - (Required) The x coordinate of widget.
  * `y` - (Required) The y coordinate of widget.
  * `width` - (Required) The width of widget.
  * `height` - (Required) The height of widget.

### markdown

* `title` - The title of markdown widget.
* `markdown` - (Required) String in Markdown format.
* `layout` - (Required) The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and they axis in the down direction as the positive direction.
  * `x` - (Required) The x coordinate of widget.
  * `y` - (Required) The y coordinate of widget.
  * `width` - (Required) The width of widget.
  * `height` - (Required) The height of widget.

### alert_status

* `title` - The title of alertStatus widget.
* `role_fullname` - (Required) The service name and role name concatenated by `:`.
* `layout` - (Required) The coordinates are specified with the upper left corner of the widget display area as the origin (x = 0, y = 0), with the x axis in the right direction and they axis in the down direction as the positive direction.
  * `x` - (Required) The x coordinate of widget.
  * `y` - (Required) The y coordinate of widget.
  * `width` - (Required) The width of widget.
  * `height` - (Required) The height of widget.
