package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelDashboardGraph(t *testing.T) {
	resourceName := "mackerel_dashboard.graph"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard graph %s", rand)
	titleUpdated := fmt.Sprintf("tf-dashboard graph %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMackerelDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelDashboardConfigGraph(rand, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "graph.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.title", "test graph role"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.role.0.role_fullname", fmt.Sprintf("tf-service-%s-include:tf-role-%s-include", rand, rand)),
					resource.TestCheckResourceAttr(resourceName, "graph.0.role.0.name", "loadavg5"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.range.0.relative.0.period", "3600"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.range.0.relative.0.offset", "1800"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.x", "2"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.y", "12"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.width", "10"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.height", "8"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelDashboardConfigGraphUpdated(rand, titleUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", titleUpdated),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "graph.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.title", "test graph role"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.role.0.role_fullname", fmt.Sprintf("tf-service-%s-include:tf-role-%s-include", rand, rand)),
					resource.TestCheckResourceAttr(resourceName, "graph.0.role.0.name", "loadavg5"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.range.0.relative.0.period", "3600"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.range.0.relative.0.offset", "1800"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.x", "2"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.y", "12"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.width", "10"),
					resource.TestCheckResourceAttr(resourceName, "graph.0.layout.0.height", "8"),
					resource.TestCheckResourceAttr(resourceName, "graph.1.title", "test graph expression"),
					resource.TestCheckResourceAttr(resourceName, "graph.1.expression.0.expression", fmt.Sprintf("role(tf-service-%s-include:tf-role-%s-include, loadavg5)", rand, rand)),
					resource.TestCheckResourceAttr(resourceName, "graph.1.range.0.absolute.0.start", "1667275734"),
					resource.TestCheckResourceAttr(resourceName, "graph.1.range.0.absolute.0.end", "1672546734"),
					resource.TestCheckResourceAttr(resourceName, "graph.1.layout.0.x", "4"),
					resource.TestCheckResourceAttr(resourceName, "graph.1.layout.0.y", "32"),
					resource.TestCheckResourceAttr(resourceName, "graph.1.layout.0.width", "10"),
					resource.TestCheckResourceAttr(resourceName, "graph.1.layout.0.height", "8"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMackerelDashboardValue(t *testing.T) {
	resourceName := "mackerel_dashboard.value"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard value %s", rand)
	titleUpdated := fmt.Sprintf("tf-dashboard value %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMackerelDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelDashboardConfigValue(rand, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "value.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "value.0.title", "test value expression"),
						resource.TestCheckResourceAttr(resourceName, "value.0.metric.0.expression.0.expression", fmt.Sprintf("role(tf-service-%s-include:tf-role-%s-include, loadavg5)", rand, rand)),
						resource.TestCheckResourceAttr(resourceName, "value.0.fraction_size", "5"),
						resource.TestCheckResourceAttr(resourceName, "value.0.suffix", "test suffix"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.x", "3"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.y", "15"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.height", "4"),
					),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelDashboardConfigValueUpdated(rand, titleUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", titleUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "value.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "value.0.title", "test value expression"),
						resource.TestCheckResourceAttr(resourceName, "value.0.metric.0.expression.0.expression", fmt.Sprintf("role(tf-service-%s-include:tf-role-%s-include, loadavg5)", rand, rand)),
						resource.TestCheckResourceAttr(resourceName, "value.0.fraction_size", "10"),
						resource.TestCheckResourceAttr(resourceName, "value.0.suffix", "test suffix"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.x", "6"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.y", "15"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.width", "10"),
						resource.TestCheckResourceAttr(resourceName, "value.0.layout.0.height", "7"),
					),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMackerelDashboardMarkdown(t *testing.T) {
	resourceName := "mackerel_dashboard.markdown"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard markdown %s", rand)
	titleUpdated := fmt.Sprintf("tf-dashboard markdown %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMackerelDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelDashboardConfigMarkdown(rand, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "markdown.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "markdown.0.title", "test markdown"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.markdown", "# h1"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.x", "1"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.y", "2"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.height", "4"),
					),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelDashboardConfigMarkdownUpdated(rand, titleUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", titleUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "markdown.#", "2"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "markdown.0.title", "test markdown"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.markdown", "# h1"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.x", "1"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.y", "2"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(resourceName, "markdown.0.layout.0.height", "4"),
						resource.TestCheckResourceAttr(resourceName, "markdown.1.title", "test markdown 2"),
						resource.TestCheckResourceAttr(resourceName, "markdown.1.markdown", "# h2"),
						resource.TestCheckResourceAttr(resourceName, "markdown.1.layout.0.x", "2"),
						resource.TestCheckResourceAttr(resourceName, "markdown.1.layout.0.y", "10"),
						resource.TestCheckResourceAttr(resourceName, "markdown.1.layout.0.width", "3"),
						resource.TestCheckResourceAttr(resourceName, "markdown.1.layout.0.height", "4"),
					),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMackerelDashboardAlertStatus(t *testing.T) {
	resourceName := "mackerel_dashboard.alertstatus"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard alertstatus %s", rand)
	titleUpdated := fmt.Sprintf("tf-dashboard alertstatus %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMackerelDashboardDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelDashboardConfigAlertStatus(rand, title),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "alert_status.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.title", "test alertStatus"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.role_fullname", fmt.Sprintf("tf-service-%s-include:tf-role-%s-include", rand, rand)),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.x", "1"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.y", "2"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.height", "4"),
					),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelDashboardConfigAlertStatusUpdated(rand, titleUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", titleUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This dashboard is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "alert_status.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.title", "test alertStatus"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.role_fullname", fmt.Sprintf("tf-service-%s-include:tf-role-%s-include", rand, rand)),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.x", "5"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.y", "7"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(resourceName, "alert_status.0.layout.0.height", "4"),
					),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMackerelDashboardDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_dashboard" {
			continue
		}
		dashboards, err := client.FindDashboards()
		if err != nil {
			return err
		}
		for _, dashboard := range dashboards {
			if dashboard.ID == r.Primary.ID {
				return fmt.Errorf("dashboard still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelDashboardExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("dashboard not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no dashboard ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		dashboards, err := client.FindDashboards()
		if err != nil {
			return err
		}

		for _, dashboard := range dashboards {
			if dashboard.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("dashboard not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccMackerelDashboardConfigGraph(rand string, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}
	
resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "graph" {
  title = "%s"
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
}
`, rand, rand, title, rand)
}

func testAccMackerelDashboardConfigGraphUpdated(rand string, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "graph" {
  title = "%s"
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
}
`, rand, rand, title, rand)
}

func testAccMackerelDashboardConfigValue(rand string, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}
	
resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "value" {
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
`, rand, rand, title, rand)
}

func testAccMackerelDashboardConfigValueUpdated(rand string, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}
	
resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "value" {
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
    fraction_size = 10
    suffix = "test suffix"
    layout {
      x = 6
      y = 15
      width = 10
      height = 7
    }
  }
}
`, rand, rand, title, rand)
}

func testAccMackerelDashboardConfigMarkdown(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_dashboard" "markdown" {
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
}
`, title, rand)
}

func testAccMackerelDashboardConfigMarkdownUpdated(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_dashboard" "markdown" {
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
`, title, rand)
}

func testAccMackerelDashboardConfigAlertStatus(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}
	
resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "alertstatus" {
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
`, rand, rand, title, rand)
}

func testAccMackerelDashboardConfigAlertStatusUpdated(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}
	
resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-%s-include"
}

resource "mackerel_dashboard" "alertstatus" {
  title = "%s"
  memo = "This dashboard is managed by Terraform."
  url_path = "%s"
  alert_status {
    title = "test alertStatus"
    role_fullname = "${mackerel_service.include.name}:${mackerel_role.include.name}"
    layout {
      x = 5
      y = 7
      width = 3
      height = 4
    }
  }
}
`, rand, rand, title, rand)
}
