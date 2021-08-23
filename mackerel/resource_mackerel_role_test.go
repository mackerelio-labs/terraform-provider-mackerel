package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelRole(t *testing.T) {
	resourceName := "mackerel_role.bar"
	rand := acctest.RandString(5)
	serviceName := fmt.Sprintf("tf-service-%s", rand)
	name := fmt.Sprintf("tf-role-%s", rand)
	nameUpdated := fmt.Sprintf("tf-rol-%s-updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMackerelRoleDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelRoleConfig(serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", serviceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelRoleConfigUpdated(serviceName, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "service", serviceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", fmt.Sprintf("%s is managed by Terraform", nameUpdated)),
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

func testAccCheckMackerelRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_role" {
			continue
		}

		services, err := client.FindServices()
		if err != nil {
			return err
		}
		for _, service := range services {
			if service.Name != r.Primary.Attributes["service"] {
				continue
			}
			for _, role := range service.Roles {
				if role == r.Primary.Attributes["name"] {
					return fmt.Errorf("mackerel role still exists: %s", r.Primary.ID)
				}
			}
		}
	}
	return nil
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
			return err
		}
		for _, role := range roles {
			if role.Name == rs.Primary.Attributes["name"] {
				return nil
			}
		}

		return fmt.Errorf("role not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccMackerelRoleConfig(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.id
  name = "%s"
}
`, serviceName, name)
}

func testAccMackerelRoleConfigUpdated(serviceName, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_role" "bar" {
  service = mackerel_service.foo.name
  name = "%s"
  memo = "%s is managed by Terraform"
}
`, serviceName, name, name)
}
