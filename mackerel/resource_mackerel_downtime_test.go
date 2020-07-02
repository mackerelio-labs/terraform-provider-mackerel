package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mackerelio/mackerel-client-go"
)

func TestAccMackerelDowntime(t *testing.T) {
	resourceName := "mackerel_downtime.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-downtime-%s", rand)
	rNameUpdated := fmt.Sprintf("tf-downtime-%s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelDowntimeDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelDowntimeConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDowntimeExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelDowntimeConfigUpdated(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDowntimeExists(resourceName),
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

func testAccCheckMackerelDowntimeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mackerel.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_downtime" {
			continue
		}

		downtimes, err := client.FindDowntimes()
		if err != nil {
			return err
		}
		for _, dt := range downtimes {
			if dt.ID == r.Primary.ID {
				return fmt.Errorf("downtime still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelDowntimeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("downtime not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no downtime ID is set")
		}

		client := testAccProvider.Meta().(*mackerel.Client)
		downtimes, err := client.FindDowntimes()
		if err != nil {
			return err
		}
		for _, dt := range downtimes {
			if dt.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("downtime not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccMackerelDowntimeConfig(name string) string {
	return fmt.Sprintf(`
resource "mackerel_downtime" "foo" {
  name     = "%s"
  memo     = "Planned maintenance"
  start    = 1735707600
  duration = 3600
  recurrence {
    type = "daily"
    interval = 2
    until = 1738332000
  }
}
`, name)
}

func testAccMackerelDowntimeConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "mackerel_downtime" "foo" {
  name     = "%s"
  memo     = "Planned maintenance"
  start    = 1735707600
  duration = 3600
  recurrence {
    type = "weekly"
    interval = 2
    weekdays = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
    until = 1738332000
  }
}
`, name)
}
