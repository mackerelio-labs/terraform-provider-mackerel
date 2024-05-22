package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMackerelAlertGroupSetting(t *testing.T) {
	dsName := "data.mackerel_alert_group_setting.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-alert-group-setting-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelAlertGroupSettingConfig(rand, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "memo", "This alert group setting is managed by Terraform."),
					resource.TestCheckResourceAttr(dsName, "service_scopes.#", "1"),
					resource.TestCheckResourceAttr(dsName, "role_scopes.#", "1"),
					resource.TestCheckResourceAttr(dsName, "monitor_scopes.#", "1"),
					resource.TestCheckResourceAttr(dsName, "notification_interval", "60"),
				),
			},
		},
	})
}

func testAccDataSourceMackerelAlertGroupSettingConfig(rand, name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "tf-service-%s"
}

resource "mackerel_role" "foo" {
  service = mackerel_service.foo.id
  name = "tf-role-%s"
}

resource "mackerel_monitor" "foo" {
  name = "tf-monitor-%s"
  connectivity {}
}

resource "mackerel_alert_group_setting" "foo" {
  name = "%s"
  memo = "This alert group setting is managed by Terraform."
  service_scopes = [mackerel_service.foo.id]
  role_scopes = [mackerel_role.foo.id]
  monitor_scopes = [mackerel_monitor.foo.id]
  notification_interval = 60
}

data "mackerel_alert_group_setting" "foo" {
  id = mackerel_alert_group_setting.foo.id
}
`, rand, rand, rand, name)
}
