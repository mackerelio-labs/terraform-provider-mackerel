package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelService(t *testing.T) {
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	memo := fmt.Sprintf("%s is managed by Terraform.", name)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelServiceConfig(name, memo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceExists("mackerel_service.foo"),
					resource.TestCheckResourceAttr(
						"mackerel_service.foo", "name", name),
					resource.TestCheckResourceAttr(
						"mackerel_service.foo", "memo", memo),
				),
			},
		},
	})
}

func testAccCheckMackerelServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_service" {
			continue
		}

		services, err := client.FindServices()
		if err != nil {
			return fmt.Errorf("err: %s", err)
		}
		for _, srv := range services {
			if srv.Name == r.Primary.ID {
				return fmt.Errorf("service still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelServiceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("service not found from resources: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no service ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		services, err := client.FindServices()
		if err != nil {
			return fmt.Errorf("err: %s", err)
		}

		var found = false
		for _, srv := range services {
			if srv.Name == rs.Primary.ID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("service not found from mackerel: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckMackerelServiceConfig(name, memo string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
	name = "%s"
	memo = "%s"
}
`, name, memo)
}
