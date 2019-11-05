package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelRole(t *testing.T) {
	serviceName := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
	roleName := fmt.Sprintf("tf-role-%s", acctest.RandString(5))
	roleMemo := fmt.Sprintf("%s role is managed by Terraform.", roleName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil, // todo
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelRoleConfig(serviceName, roleName, roleMemo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleExists("mackerel_role.bar"),
					resource.TestCheckResourceAttr(
						"mackerel_role.bar", "service", serviceName),
					resource.TestCheckResourceAttr(
						"mackerel_role.bar", "name", roleName),
					resource.TestCheckResourceAttr(
						"mackerel_role.bar", "memo", roleMemo),
				),
			},
		},
	})
}

func testAccCheckMackerelRoleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("role not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no role ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		roles, err := client.FindRoles(rs.Primary.Attributes["service"])
		if err != nil {
			return fmt.Errorf("err: %s", err)
		}

		for _, role := range roles {
			if role.Name == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("role not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccCheckMackerelRoleConfig(serviceName, roleName, roleMemo string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
    name = "%s"
}

resource "mackerel_role" "bar" {
    service = "${mackerel_service.foo.id}"
    name = "%s"
    memo = "%s"
}
`, serviceName, roleName, roleMemo)
}
