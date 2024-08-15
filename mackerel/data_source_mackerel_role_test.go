package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMackerelRole(t *testing.T) {
	dsName := "data.mackerel_role.foo"
	rand := acctest.RandString(5)
	serviceName := fmt.Sprintf("tf-service-%s", rand)
	name := fmt.Sprintf("tf-role-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelRoleConfig(serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "id", fmt.Sprintf("%s:%s", serviceName, name)),
					resource.TestCheckResourceAttr(dsName, "service", serviceName),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This role is managed by Terraform."),
				),
			},
		},
	})
}

func testAccDataSourceMackerelRoleConfig(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "foo" {
  service = mackerel_service.foo.name
  name = "%s"
  memo = "This role is managed by Terraform."
}

data "mackerel_role" "foo" {
  depends_on = [mackerel_role.foo]
  service = mackerel_role.foo.service
  name = mackerel_role.foo.name
}
`, serviceName, name)
}
