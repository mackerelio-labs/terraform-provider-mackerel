package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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
					resource.TestCheckResourceAttr(dsName, "url_path", "testpath"),
					resource.TestCheckResourceAttr(dsName, "markdown.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "markdown.0.title", "test markdown"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.markdown", "# h1"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.x", "1"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.y", "2"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.width", "3"),
						resource.TestCheckResourceAttr(dsName, "markdown.0.layout.0.height", "4"),
					),
				),
			},
		},
	})
}

func testAccDataSourceMackerelDashboardConfigMarkdown(rand, title string) string {
	return fmt.Sprintf(`
resource "mackerel_dashboard" "foo" {
  title = "%s"
  memo = "This dashboard is managed by Terraform."
  url_path = "testpath"
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

data "mackerel_dashboard" "foo" {
  id = mackerel_dashboard.foo.id
}
`, title)
}
