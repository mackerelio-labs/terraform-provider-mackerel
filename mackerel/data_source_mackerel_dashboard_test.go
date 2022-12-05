package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMackerelDashboardGraph(t *testing.T) {
	dsName := "data.mackerel_dashboard.foo"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelDashboardConfigGraph(rand, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "title", title),
					resource.TestCheckResourceAttr(dsName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "url_path", rand),
					resource.TestCheckResourceAttr(dsName, "graph.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "graph.0.title", "test graph role"),
						resource.TestCheckResourceAttr(dsName, "graph.0.graph.0.role.0.role_fullname", fmt.Sprintf("tf-service-%s-include:tf-role-%s-include", rand, rand)),
						resource.TestCheckResourceAttr(dsName, "graph.0.graph.0.role.0.name", "loadavg5"),
						resource.TestCheckResourceAttr(dsName, "graph.0.range.0.relative.0.period", "3600"),
						resource.TestCheckResourceAttr(dsName, "graph.0.range.0.relative.0.offset", "1800"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.x", "1"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.y", "2"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(dsName, "graph.0.layout.0.height", "4"),
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
  // graph {
  //   title = "test graph host"
	// 	graph {
	// 		host {
	// 			host_id = "<host_id>"
	// 			name = "loadavg"
	// 		}
	// 	}
	// 	range {
	// 		relative {
	// 			period = 3600
	// 			offset = 1800
	// 		}
	// 	}
  //   layout {
	// 		x = 1
	// 		y = 2
	// 		width = 3
	// 		height = 4
	// 	}
	// }
	graph {
		title = "test graph role"
		graph {
			role {
				role_fullname = "${mackerel_service.include.name}:${mackerel_role.include.name}"
				name = "loadavg5"
				is_stacked = true
			}
		}
		range {
			relative {
				period = 3600
				offset = 1800
			}
		}
		layout {
			x = 1
			y = 2
			width = 3
			height = 4
		}
	}
	// TODO
	graph {
		title = "test graph service"
		graph {
			service {
				service_name = "hatena-mac"
				name = "Sample.*"
			}
		}
		range {
			absolute {
				start = 1669794987
				end = 1669798555
			}
		}
		layout {
			x = 2
			y = 10
			width = 4
			height = 5
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
