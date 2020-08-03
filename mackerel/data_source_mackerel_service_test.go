package mackerel

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMackerelService(t *testing.T) {
	name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(`data "mackerel_service" "foo" { name = "%s" }`, name),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`the name '%s' does not match any service in mackerel\.io`, name)),
			},
			{
				Config: testAccDataSourceMackerelServiceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mackerel_service.foo", "id", name),
					resource.TestCheckResourceAttr("data.mackerel_service.foo", "name", name),
					resource.TestCheckResourceAttr("data.mackerel_service.foo", "memo", "This service is managed by Terraform."),
				),
			},
		},
	})
}

func testAccDataSourceMackerelServiceConfig(name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
  memo = "This service is managed by Terraform."
}

data "mackerel_service" "foo" {
  name = mackerel_service.foo.name
}
`, name)
}
