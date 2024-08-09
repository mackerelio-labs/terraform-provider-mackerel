package mackerel

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMackerelDowntime(t *testing.T) {
	dsName := "data.mackerel_downtime.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-downtime-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelDowntimeConfig(rand, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This downtime is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "start", "1735707600"),
					resource.TestCheckResourceAttr(dsName, "duration", "3600"),
					resource.TestCheckResourceAttr(dsName, "recurrence.#", "1"),
					resource.TestCheckResourceAttr(dsName, "recurrence.0.type", "weekly"),
					resource.TestCheckResourceAttr(dsName, "recurrence.0.interval", "2"),
					resource.TestCheckResourceAttr(dsName, "recurrence.0.weekdays.#", "5"),
					resource.TestCheckResourceAttr(dsName, "recurrence.0.until", "1767193199"),
					resource.TestCheckResourceAttr(dsName, "service_scopes.#", "1"),
					resource.TestCheckResourceAttr(dsName, "service_exclude_scopes.#", "1"),
					resource.TestCheckResourceAttr(dsName, "role_scopes.#", "1"),
					resource.TestCheckResourceAttr(dsName, "role_exclude_scopes.#", "1"),
					resource.TestCheckResourceAttr(dsName, "monitor_scopes.#", "0"),
					resource.TestCheckResourceAttr(dsName, "monitor_exclude_scopes.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceMackerelDowntimeConfig(rand, name string) string {
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

data "mackerel_downtime" "foo" {
  id = mackerel_downtime.foo.id
}
`, rand, rand, rand, rand, name)
}

func TestAccDataSourceMackerelDowntimeNotMatchAnyDowntime(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `data "mackerel_downtime" "foo" { id = "not-found" }`,
				ExpectError: regexp.MustCompile(`the ID 'not-found' does not match any downtime in mackerel\.io`),
			},
		},
	})
}
