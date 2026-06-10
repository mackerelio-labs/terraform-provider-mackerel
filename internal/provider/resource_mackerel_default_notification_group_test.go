package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
	mackerelclient "github.com/mackerelio/mackerel-client-go"
)

func Test_MackerelDefaultNotificationGroupResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	provider.NewMackerelDefaultNotificationGroupResource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %+v", resp.Diagnostics)
	}

	if _, ok := resp.Schema.Blocks["monitor"]; ok {
		t.Fatal("default notification group resource must not expose monitor block")
	}
	if _, ok := resp.Schema.Blocks["service"]; ok {
		t.Fatal("default notification group resource must not expose service block")
	}
	if _, ok := resp.Schema.Attributes["id"]; !ok {
		t.Fatal("default notification group resource must expose id attribute")
	}
	if _, ok := resp.Schema.Attributes["name"]; ok {
		t.Fatal("default notification group resource must not expose name attribute")
	}
	if _, ok := resp.Schema.Attributes["notification_level"]; !ok {
		t.Fatal("default notification group resource must expose notification_level attribute")
	}
	if _, ok := resp.Schema.Attributes["child_notification_group_ids"]; !ok {
		t.Fatal("default notification group resource must expose child_notification_group_ids attribute")
	}
	assertRequiredSetAttribute(t, resp.Schema.Attributes["child_notification_group_ids"])
	assertRequiredSetAttribute(t, resp.Schema.Attributes["child_channel_ids"])

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func assertRequiredSetAttribute(t *testing.T, attr schema.Attribute) {
	t.Helper()

	setAttr, ok := attr.(schema.SetAttribute)
	if !ok {
		t.Fatalf("attribute type = %T, want schema.SetAttribute", attr)
	}
	if !setAttr.IsRequired() {
		t.Fatal("attribute must be required")
	}
	if setAttr.IsOptional() {
		t.Fatal("attribute must not be optional")
	}
	if setAttr.IsComputed() {
		t.Fatal("attribute must not be computed")
	}
}

func TestAccMackerelDefaultNotificationGroup(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	resourceName := "mackerel_default_notification_group.default"
	rand := acctest.RandString(5)

	preCheck(t)

	client := mackerelClient()
	original, err := testAccFindDefaultNotificationGroup(client)
	if err != nil {
		t.Fatalf("failed to find default notification group: %+v", err)
	}
	t.Cleanup(func() {
		if err := testAccRestoreDefaultNotificationGroup(client, original); err != nil {
			t.Errorf("failed to restore default notification group: %+v", err)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMackerelDefaultNotificationGroupConfigEmpty(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDefaultNotificationGroup(resourceName, "all", 0, 0),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "all"),
					resource.TestCheckResourceAttr(resourceName, "child_notification_group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "child_channel_ids.#", "0"),
				),
			},
			{
				Config: testAccMackerelDefaultNotificationGroupConfigWithChildren(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDefaultNotificationGroup(resourceName, "critical", 1, 1),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "critical"),
					resource.TestCheckResourceAttr(resourceName, "child_notification_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "child_channel_ids.#", "1"),
				),
			},
			{
				Config: testAccMackerelDefaultNotificationGroupConfigEmptyWithChildren(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelDefaultNotificationGroup(resourceName, "all", 0, 0),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notification_level", "all"),
					resource.TestCheckResourceAttr(resourceName, "child_notification_group_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "child_channel_ids.#", "0"),
				),
			},
		},
	})
}

func testAccFindDefaultNotificationGroup(client *mackerelclient.Client) (*mackerelclient.NotificationGroup, error) {
	groups, err := client.FindNotificationGroups()
	if err != nil {
		return nil, err
	}

	var defaultGroup *mackerelclient.NotificationGroup
	for _, group := range groups {
		if group.Type != mackerelclient.NotificationGroupTypeGroupDefault {
			continue
		}
		if defaultGroup != nil {
			return nil, fmt.Errorf("multiple default notification groups found")
		}
		defaultGroup = group
	}
	if defaultGroup == nil {
		return nil, fmt.Errorf("default notification group is not found")
	}
	return defaultGroup, nil
}

func testAccRestoreDefaultNotificationGroup(client *mackerelclient.Client, group *mackerelclient.NotificationGroup) error {
	param := testAccDefaultNotificationGroupUpdateParam(group)
	_, err := client.UpdateNotificationGroup(group.ID, &param)
	return err
}

func testAccDefaultNotificationGroupUpdateParam(group *mackerelclient.NotificationGroup) mackerelclient.NotificationGroup {
	param := *group
	param.Type = ""
	if param.ChildNotificationGroupIDs == nil {
		param.ChildNotificationGroupIDs = []string{}
	}
	if param.ChildChannelIDs == nil {
		param.ChildChannelIDs = []string{}
	}
	return param
}

func testAccCheckMackerelDefaultNotificationGroup(resourceName, notificationLevel string, childNotificationGroupCount, childChannelCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("default notification group not found from resources: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no default notification group ID is set")
		}

		group, err := testAccFindDefaultNotificationGroup(mackerelClient())
		if err != nil {
			return err
		}
		if group.ID != rs.Primary.ID {
			return fmt.Errorf("default notification group ID = %q, want %q", group.ID, rs.Primary.ID)
		}
		if string(group.NotificationLevel) != notificationLevel {
			return fmt.Errorf("default notification group notification level = %q, want %q", group.NotificationLevel, notificationLevel)
		}
		if got := len(group.ChildNotificationGroupIDs); got != childNotificationGroupCount {
			return fmt.Errorf("default notification group child notification group count = %d, want %d", got, childNotificationGroupCount)
		}
		if got := len(group.ChildChannelIDs); got != childChannelCount {
			return fmt.Errorf("default notification group child channel count = %d, want %d", got, childChannelCount)
		}
		return nil
	}
}

func testAccMackerelDefaultNotificationGroupConfigEmpty() string {
	return `
resource "mackerel_default_notification_group" "default" {
  notification_level = "all"

  child_notification_group_ids = []
  child_channel_ids = []
}
`
}

func testAccMackerelDefaultNotificationGroupConfigWithChildren(rand string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "tf-channel-%s"
  email {}
}

resource "mackerel_notification_group" "child" {
  name = "tf-notification-group-%s-child"
}

resource "mackerel_default_notification_group" "default" {
  notification_level = "critical"

  child_notification_group_ids = [
    mackerel_notification_group.child.id]
  child_channel_ids = [
    mackerel_channel.foo.id]
}
`, rand, rand)
}

func testAccMackerelDefaultNotificationGroupConfigEmptyWithChildren(rand string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "tf-channel-%s"
  email {}
}

resource "mackerel_notification_group" "child" {
  name = "tf-notification-group-%s-child"
}

resource "mackerel_default_notification_group" "default" {
  notification_level = "all"

  child_notification_group_ids = []
  child_channel_ids = []
}
`, rand, rand)
}
