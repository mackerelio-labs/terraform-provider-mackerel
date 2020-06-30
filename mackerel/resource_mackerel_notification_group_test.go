package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelNotificationGroup(t *testing.T) {
	resourceName := "mackerel_notification_group.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notification-grouup %s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelNotificationGroupDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelNotificationGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelNotificationGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "all"),
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

func testAccCheckMackerelNotificationGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_notification_group" {
			continue
		}
		groups, err := client.FindNotificationGroups()
		if err != nil {
			return err
		}
		for _, group := range groups {
			if group.ID == r.Primary.ID {
				return fmt.Errorf("notification group still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelNotificationGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("notification group not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no notification group ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		groups, err := client.FindNotificationGroups()
		if err != nil {
			return err
		}

		for _, group := range groups {
			if group.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("notification group not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccMackerelNotificationGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "mackerel_notification_group" "foo" {
  name = "%s"
}
`, name)
}
