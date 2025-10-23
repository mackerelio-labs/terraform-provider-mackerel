package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelAlertGroupSettingResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	provider.NewMackerelAlertGroupSettingResource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostica: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostica: %+v", diags)
	}
}

func TestAccCompat_MackerelAlertGroupSettingResource(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		config func(name string) string
	}{
		"minimal": {
			config: func(name string) string {
				return `
resource "mackerel_alert_group_setting" "alert_group" {
  name = "` + name + `"
}`
			},
		},
		"empty": {
			config: func(name string) string {
				return `
resource "mackerel_alert_group_setting" "alert_group" {
  name = "` + name + `"
  service_scopes = []
  role_scopes = []
  monitor_scopes = []
}`
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			name := acctest.RandomWithPrefix("tf-alert-group-setting-compat")
			config := tt.config(name)

			resource.Test(t, resource.TestCase{
				PreCheck: func() { preCheck(t) },

				Steps: []resource.TestStep{
					{
						ProtoV5ProviderFactories: protoV5SDKProviderFactories,
						Config:                   config,
					},
					stepNoPlanInFramework(config),
					{
						ProtoV5ProviderFactories: protoV5SDKProviderFactories,
						Config:                   config,
					},
					stepNoPlanInFramework(config),
				},
			})
		})
	}
}

func TestAccMackerelAlertGroupSetting(t *testing.T) {
	resourceName := "mackerel_alert_group_setting.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-alert-group-setting-%s", rand)
	nameUpdated := fmt.Sprintf("tf-alert-group-setting-%s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelAlertGroupSettingDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelAlertGroupSettingConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAlertGroupSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "memo", ""),
					resource.TestCheckResourceAttr(resourceName, "service_scopes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "role_scopes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "monitor_scopes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelAlertGroupSettingConfigUpdated(rand, nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelAlertGroupSettingExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "memo", "This alert group setting is managed by Terraform."),
					resource.TestCheckResourceAttr(resourceName, "service_scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "role_scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "monitor_scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "notification_interval", "60"),
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

func testAccCheckMackerelAlertGroupSettingDestroy(s *terraform.State) error {
	client := mackerelClient()
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_alert_group_setting" {
			continue
		}
		settings, err := client.FindAlertGroupSettings()
		if err != nil {
			return err
		}
		for _, s := range settings {
			if s.ID == r.Primary.ID {
				return fmt.Errorf("alert group setting still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelAlertGroupSettingExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("alert group setting not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no alert group setting ID is set")
		}

		client := mackerelClient()
		if _, err := client.GetAlertGroupSetting(rs.Primary.ID); err != nil {
			return err
		}

		return nil
	}
}

func testAccMackerelAlertGroupSettingConfig(name string) string {
	return fmt.Sprintf(`
resource "mackerel_alert_group_setting" "foo" {
  name = "%s"
}
`, name)
}

func testAccMackerelAlertGroupSettingConfigUpdated(rand, name string) string {
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
`, rand, rand, rand, name)
}
