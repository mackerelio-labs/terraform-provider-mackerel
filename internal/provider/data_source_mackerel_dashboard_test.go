package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelDashboardDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := fwdatasource.SchemaResponse{}
	if provider.NewMackerelDashboardDataSource().Schema(ctx, req, &resp); resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccDataSourceMackerelDashboardGraph(t *testing.T) {
	dsName := "data.mackerel_dashboard.foo"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelDashboardConfigGraph(rand, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "title", title),
					resource.TestCheckResourceAttr(dsName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "url_path", rand),
					resource.TestCheckResourceAttr(dsName, "graph.#", "3"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "graph.0.title", "test graph role"),
						resource.TestCheckResourceAttr(dsName, "graph.0.role.0.role_fullname", fmt.Sprintf("tf-service-%s-include:tf-role-%s-include", rand, rand)),
						resource.TestCheckResourceAttr(dsName, "graph.0.role.0.name", "loadavg5"),
						resource.TestCheckResourceAttr(dsName, "graph.0.range.0.relative.0.period", "3600"),
						resource.TestCheckResourceAttr(dsName, "graph.0.range.0.relative.0.offset", "1800"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.x", "2"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.y", "12"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.width", "10"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.height", "8"),
						resource.TestCheckResourceAttr(dsName, "graph.1.title", "test graph expression"),
						resource.TestCheckResourceAttr(dsName, "graph.1.expression.0.expression", fmt.Sprintf("role(tf-service-%s-include:tf-role-%s-include, loadavg5)", rand, rand)),
						resource.TestCheckResourceAttr(dsName, "graph.1.range.0.absolute.0.start", "1667275734"),
						resource.TestCheckResourceAttr(dsName, "graph.1.range.0.absolute.0.end", "1672546734"),
						resource.TestCheckResourceAttr(dsName, "graph.1.layout.0.x", "4"),
						resource.TestCheckResourceAttr(dsName, "graph.1.layout.0.y", "32"),
						resource.TestCheckResourceAttr(dsName, "graph.1.layout.0.width", "10"),
						resource.TestCheckResourceAttr(dsName, "graph.1.layout.0.height", "8"),
						resource.TestCheckResourceAttr(dsName, "graph.2.title", "test graph query"),
						resource.TestCheckResourceAttr(dsName, "graph.2.query.0.query", "container.cpu.utilization{k8s.deployment.name=\"httpbin\"}"),
						resource.TestCheckResourceAttr(dsName, "graph.2.query.0.legend", "{{k8s.node.name}}"),
						resource.TestCheckResourceAttr(dsName, "graph.2.range.0.relative.0.period", "3600"),
						resource.TestCheckResourceAttr(dsName, "graph.2.range.0.relative.0.offset", "1800"),
						resource.TestCheckResourceAttr(dsName, "graph.2.layout.0.x", "0"),
						resource.TestCheckResourceAttr(dsName, "graph.2.layout.0.y", "52"),
						resource.TestCheckResourceAttr(dsName, "graph.2.layout.0.width", "10"),
						resource.TestCheckResourceAttr(dsName, "graph.2.layout.0.height", "8"),
					),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelDashboardValue(t *testing.T) {
	dsName := "data.mackerel_dashboard.foo"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelDashboardConfigValue(rand, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "title", title),
					resource.TestCheckResourceAttr(dsName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "url_path", rand),
					resource.TestCheckResourceAttr(dsName, "value.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "value.0.title", "test value expression"),
						resource.TestCheckResourceAttr(dsName, "value.0.metric.0.expression.0.expression", fmt.Sprintf("role(tf-service-%s-include:tf-role-%s-include, loadavg5)", rand, rand)),
						resource.TestCheckResourceAttr(dsName, "value.0.fraction_size", "5"),
						resource.TestCheckResourceAttr(dsName, "value.0.suffix", "test suffix"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.x", "3"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.y", "15"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.height", "4"),
					),
				),
			},
			{
				Config: testAccDataSourceMackerelDashboardConfigValue_Query(rand, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "title", title),
					resource.TestCheckResourceAttr(dsName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "url_path", rand),
					resource.TestCheckResourceAttr(dsName, "value.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "value.0.title", "test value query"),
						resource.TestCheckResourceAttr(dsName, "value.0.metric.0.query.0.query", "avg(avg_over_time(container.cpu.utilization{k8s.deployment.name=\"httpbin\"}[1h]))"),
						resource.TestCheckResourceAttr(dsName, "value.0.metric.0.query.0.legend", "average utilization"),
						resource.TestCheckResourceAttr(dsName, "value.0.fraction_size", "10"),
						resource.TestCheckResourceAttr(dsName, "value.0.suffix", "test query suffix"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.x", "3"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.y", "15"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(dsName, "value.0.layout.0.height", "4"),
					),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelDashboardMarkdown(t *testing.T) {
	dsName := "data.mackerel_dashboard.foo"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelDashboardConfigMarkdown(rand, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "title", title),
					resource.TestCheckResourceAttr(dsName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "url_path", rand),
					resource.TestCheckResourceAttr(dsName, "markdown.#", "2"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "markdown.0.title", "test markdown"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.markdown", "# h1"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.x", "1"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.y", "2"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.height", "4"),
						resource.TestCheckResourceAttr(dsName, "markdown.1.title", "test markdown 2"),
						resource.TestCheckResourceAttr(dsName, "markdown.1.markdown", "# h2"),
						resource.TestCheckResourceAttr(dsName, "markdown.1.layout.0.x", "2"),
						resource.TestCheckResourceAttr(dsName, "markdown.1.layout.0.y", "10"),
						resource.TestCheckResourceAttr(dsName, "markdown.1.layout.0.width", "3"),
						resource.TestCheckResourceAttr(dsName, "markdown.1.layout.0.height", "4"),
					),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelDashboardAlertStatus(t *testing.T) {
	dsName := "data.mackerel_dashboard.foo"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelDashboardConfigAlertStatus(rand, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "title", title),
					resource.TestCheckResourceAttr(dsName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "url_path", rand),
					resource.TestCheckResourceAttr(dsName, "alert_status.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "alert_status.0.title", "test alertStatus"),
						resource.TestCheckResourceAttr(dsName, "alert_status.0.role_fullname", fmt.Sprintf("tf-service-%s-include:tf-role-%s-include", rand, rand)),
						resource.TestCheckResourceAttr(dsName, "alert_status.0.layout.0.x", "1"),
						resource.TestCheckResourceAttr(dsName, "alert_status.0.layout.0.y", "2"),
						resource.TestCheckResourceAttr(dsName, "alert_status.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(dsName, "alert_status.0.layout.0.height", "4"),
					),
				),
			},
		},
	})
}

func testAccDataSourceMackerelDashboardConfigGraph(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "foo" {
  title = "%s"
  memo = "This dashboard is managed by Terraform."
  url_path = "%s"
  graph {
    title = "test graph role"
    role {
      role_fullname = "${mackerel_service.include.name}:${mackerel_role.include.name}"
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
  graph {
    title = "test graph expression"
    expression {
      expression = "role(${mackerel_service.include.name}:${mackerel_role.include.name}, loadavg5)"
    }
    range {
      absolute {
        start = 1667275734
        end = 1672546734
      }
    }
    layout {
      x = 4
      y = 32
      width = 10
      height = 8
    }
  }
  graph {
    title = "test graph query"
    query {
      query = "container.cpu.utilization{k8s.deployment.name=\"httpbin\"}"
      legend = "{{k8s.node.name}}"
    }
    range {
      relative {
        period = 3600
        offset = 1800
      }
    }
    layout {
      x = 0
      y = 52
      width = 10
      height = 8
    }
  }
}


data "mackerel_dashboard" "foo" {
  id = mackerel_dashboard.foo.id
}
`, rand, rand, title, rand)
}

func testAccDataSourceMackerelDashboardConfigValue(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "foo" {
  title = "%s"
  memo = "This dashboard is managed by Terraform."
  url_path = "%s"
  value {
    title = "test value expression"
    metric {
      expression {
        expression = "role(${mackerel_service.include.name}:${mackerel_role.include.name}, loadavg5)"
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

data "mackerel_dashboard" "foo" {
  id = mackerel_dashboard.foo.id
}
`, rand, rand, title, rand)
}

func testAccDataSourceMackerelDashboardConfigValue_Query(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "foo" {
  title = "%s"
  memo = "This dashboard is managed by Terraform."
  url_path = "%s"
  value {
    title = "test value query"
    metric {
      query {
        query = "avg(avg_over_time(container.cpu.utilization{k8s.deployment.name=\"httpbin\"}[1h]))"
        legend = "average utilization"
      }
    }
    fraction_size = 10
    suffix = "test query suffix"
    layout {
      x = 3
      y = 15
      width = 3
      height = 4
    }
  }
}

data "mackerel_dashboard" "foo" {
  id = mackerel_dashboard.foo.id
}
`, rand, rand, title, rand)
}

func testAccDataSourceMackerelDashboardConfigMarkdown(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_dashboard" "foo" {
  title = "%s"
  memo = "This dashboard is managed by Terraform."
  url_path = "%s"
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
  markdown {
    title = "test markdown 2"
    markdown = "# h2"
    layout {
      x = 2
      y = 10
      width = 3
      height = 4
    }
  }
}

data "mackerel_dashboard" "foo" {
  id = mackerel_dashboard.foo.id
}
`, title, rand)
}

func testAccDataSourceMackerelDashboardConfigAlertStatus(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "foo" {
  title = "%s"
  memo = "This dashboard is managed by Terraform."
  url_path = "%s"
  alert_status {
    title = "test alertStatus"
    role_fullname = "${mackerel_service.include.name}:${mackerel_role.include.name}"
    layout {
      x = 1
      y = 2
      width = 3
      height = 4
    }
  }
}

data "mackerel_dashboard" "foo" {
  id = mackerel_dashboard.foo.id
}
`, rand, rand, title, rand)
}
