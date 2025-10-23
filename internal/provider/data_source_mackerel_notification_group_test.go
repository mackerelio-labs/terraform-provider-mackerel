package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelNotificationGroupDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := &fwdatasource.SchemaResponse{}
	provider.NewMackerelNotificationGroupDataSource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccDataSourceMackerelNotificationGroup(t *testing.T) {
	dsName := "data.mackerel_notification_group.foo"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-notification-group-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelNotificationGroupConfig(rand, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "notification_level", "critical"),
					resource.TestCheckResourceAttr(dsName, "child_notification_group_ids.#", "1"),
					resource.TestCheckResourceAttr(dsName, "child_channel_ids.#", "1"),
					resource.TestCheckResourceAttr(dsName, "monitor.#", "1"),
					resource.TestCheckResourceAttr(dsName, "service.#", "2"),
				),
			},
		},
	})
}

func testAccDataSourceMackerelNotificationGroupConfig(rand, name string) string {
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

data "mackerel_notification_group" "foo" {
  id = mackerel_notification_group.foo.id
}
`, rand, rand, rand, rand, rand, name)
}
