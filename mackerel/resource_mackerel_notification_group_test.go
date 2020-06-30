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
	rNameUpdated := fmt.Sprintf("tf-notification-gruop %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelNotificationGroupDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelNotificationGroupConfig(rand, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelNotificationGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "all"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelNotificationGroupConfigUpdate(rand, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelNotificationGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "critical"),
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

func testAccMackerelNotificationGroupConfig(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "tf-service-%s"
}

resource "mackerel_channel" "foo" {
  name = "tf-channel-%s"
  email { }
}

resource "mackerel_notification_group" "foo" {
  name = "%s"
  child_channel_ids = [mackerel_channel.foo.id]
  service {
    name = mackerel_service.foo.id
  }
}
`, rand, rand, name)
}

func testAccMackerelNotificationGroupConfigUpdate(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "tf-service-%s"
}

resource "mackerel_channel" "foo" {
  name = "tf-channel-%s"
  email { }
}

resource "mackerel_notification_group" "foo" {
  name = "%s"
  notification_level = "critical"
}`, rand, rand, name)
}
