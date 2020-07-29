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
	name := fmt.Sprintf("tf-notification-grouup %s", rand)
	nameUpdated := fmt.Sprintf("tf-notification-group %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelNotificationGroupDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelNotificationGroupConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelNotificationGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "all"),
					resource.TestCheckResourceAttr(resourceName, "child_notification_group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "child_channel_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "service.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelNotificationGroupConfigUpdated(rand, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelNotificationGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "critical"),
					resource.TestCheckResourceAttr(resourceName, "child_notification_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "child_channel_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "monitor.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service.#", "2"),
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

func testAccMackerelNotificationGroupConfigUpdated(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "tf-service-%s"
}

resource "mackerel_service" "bar" {
  name = "tf-service-%s-bar"
}

resource "mackerel_channel" "foo" {
  name = "tf-channel-%s"
  email {}
}

resource "mackerel_monitor" "foo" {
  name = "tf-monitor-%s"
  connectivity {}
}

resource "mackerel_notification_group" "child" {
  name = "tf-notification-group-%s-child"
}

resource "mackerel_notification_group" "foo" {
  name = "%s"
  notification_level = "critical"
  child_notification_group_ids = [
    mackerel_notification_group.child.id]
  child_channel_ids = [
    mackerel_channel.foo.id]
  monitor {
    id = mackerel_monitor.foo.id
    skip_default = false
  }
  // ignore duplicates
  monitor {
    id = mackerel_monitor.foo.id
    skip_default = false
  }
  service {
    name = mackerel_service.foo.name
  }
  // ignore duplicates
  service {
    name = mackerel_service.foo.name
  }
  service {
    name = mackerel_service.bar.name
  }
}
`, rand, rand, rand, rand, rand, name)
}
