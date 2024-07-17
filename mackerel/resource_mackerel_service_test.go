package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelService_withMemo(t *testing.T) {
	resourceName := "mackerel_service.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-%s", rand)
	nameUpdated := fmt.Sprintf("tf-updated-%s", rand)
	memo := fmt.Sprintf("%s is managed by Terraform.", name)
	memoUpdated := fmt.Sprintf("%s is managed by Terraform.", nameUpdated)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckMackerelServiceDestroy,
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

func TestAccMackerelService_noMemo(t *testing.T) {
	resourceName := "mackerel_service.foo"
	rand := acctest.RandString(5)
	name := "tf-" + rand
	nameUpdated := "tf-updated-" + rand

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckMackerelServiceDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelServiceConfig_noMemo(name),
				Check:  testAccCheckMackerelServiceExists(resourceName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact("")),
				},
			},
			// Test: Update
			{
				Config: testAccMackerelServiceConfig_noMemo(nameUpdated),
				Check:  testAccCheckMackerelServiceExists(resourceName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(nameUpdated)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(nameUpdated)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("memo"), knownvalue.StringExact("")),
				},
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

func testAccMackerelServiceConfig_noMemo(name string) string {
	return `
resource "mackerel_service" "foo" {
  name = "` + name + `"
}`
}
