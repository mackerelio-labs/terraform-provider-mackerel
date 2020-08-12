package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelService(t *testing.T) {
	resourceName := "mackerel_service.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-%s", rand)
	nameUpdated := fmt.Sprintf("tf-updated-%s", rand)
	memo := fmt.Sprintf("%s is managed by Terraform.", name)
	memoUpdated := fmt.Sprintf("%s is managed by Terraform.", nameUpdated)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelServiceDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelServiceConfig(name, memo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", memo),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelServiceConfig(nameUpdated, memoUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", memoUpdated),
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

func testAccCheckMackerelServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_service" {
			continue
		}

		services, err := client.FindServices()
		if err != nil {
			return err
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
			return err
		}

		for _, srv := range services {
			if srv.Name == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("service not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccMackerelServiceConfig(name, memo string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
  memo = "%s"
}
`, name, memo)
}
