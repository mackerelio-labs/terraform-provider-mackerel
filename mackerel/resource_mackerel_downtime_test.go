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
	name := fmt.Sprintf("tf-downtime-%s", rand)
	nameUpdated := fmt.Sprintf("tf-downtime-%s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMackerelDowntimeDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelDowntimeConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDowntimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "start", "1735707600"),
					resource.TestCheckResourceAttr(resourceName, "duration", "3600"),
					resource.TestCheckNoResourceAttr(resourceName, "recurrence.#"),
					resource.TestCheckNoResourceAttr(resourceName, "service_scopes.#"),
					resource.TestCheckNoResourceAttr(resourceName, "service_exclude_scopes.#"),
					resource.TestCheckNoResourceAttr(resourceName, "role_scopes.#"),
					resource.TestCheckNoResourceAttr(resourceName, "role_exclude_scopes.#"),
					resource.TestCheckNoResourceAttr(resourceName, "monitor_scopes.#"),
					resource.TestCheckNoResourceAttr(resourceName, "monitor_exclude_scopes.#"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelDowntimeConfigUpdated(rand, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDowntimeExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This downtime is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "start", "1735707600"),
					resource.TestCheckResourceAttr(resourceName, "duration", "3600"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.type", "weekly"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.interval", "2"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.weekdays.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "recurrence.0.until", "1767193199"),
					resource.TestCheckResourceAttr(resourceName, "service_scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "service_exclude_scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_exclude_scopes.#", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "monitor_scopes.#"),
					resource.TestCheckNoResourceAttr(resourceName, "monitor_exclude_scopes.#"),
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
  name = "%s"
  start = 1735707600
  duration = 3600
}
`, name)
}

func testAccMackerelDowntimeConfigUpdated(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "include" {
  name = "tf-service-%s-include"
}

resource "mackerel_role" "exclude" {
  service = mackerel_service.include.name
  name = "tf-role-%s"
}

resource "mackerel_service" "exclude" {
  name = "tf-service-%s-exclude"
}

resource "mackerel_role" "include" {
  service = mackerel_service.exclude.name
  name = "tf-role-%s"
}

resource "mackerel_downtime" "foo" {
  name = "%s"
  memo = "This downtime is managed by Terraform."
  start = 1735707600
  duration = 3600

  recurrence {
    type = "weekly"
    interval = 2
    weekdays = [
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday"]
    until = 1767193199
  }

  service_scopes = [
    mackerel_service.include.name]
  service_exclude_scopes = [
    mackerel_service.exclude.name]

  role_scopes = [
    "${mackerel_role.include.service}: ${mackerel_role.include.name}"]
  role_exclude_scopes = [
    "${mackerel_role.exclude.service}: ${mackerel_role.exclude.name}"]
}

`, rand, rand, rand, rand, name)
}