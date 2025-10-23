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

func Test_MackerelChannelResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	provider.NewMackerelChannelResource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccMackerelChannel_Email(t *testing.T) {
	resourceName := "mackerel_channel.email"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-channel email %s", rand)
	nameUpdated := fmt.Sprintf("tf-channel email %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelChannelConfigEmail(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "email.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email.0.emails.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "email.0.user_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "email.0.events.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "slack.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelChannelConfigEmailUpdated(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "email.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email.0.emails.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "email.0.user_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "email.0.events.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "slack.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "0"),
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

func TestAccMackerelChannel_Slack(t *testing.T) {
	resourceName := "mackerel_channel.slack"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-channel slack %s", rand)
	nameUpdated := fmt.Sprintf("tf-channel slack %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelChannelConfigSlack(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "email.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "slack.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "slack.0.url"),
					resource.TestCheckResourceAttr(resourceName, "slack.0.mentions.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "slack.0.events.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelChannelConfigSlackUpdated(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "email.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "slack.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "slack.0.url"),
					resource.TestCheckResourceAttr(resourceName, "slack.0.mentions.ok", "OK!!!"),
					resource.TestCheckResourceAttr(resourceName, "slack.0.mentions.warning", "WARNING!!!"),
					resource.TestCheckResourceAttr(resourceName, "slack.0.mentions.critical", "CRITICAL!!!"),
					resource.TestCheckResourceAttr(resourceName, "slack.0.events.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "0"),
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

func TestAccMackerelChannel_Webhook(t *testing.T) {
	resourceName := "mackerel_channel.webhook"
	rand := acctest.RandString(5)
	name := fmt.Sprintf("tf-channel slack %s", rand)
	nameUpdated := fmt.Sprintf("tf-channel slack %s updated", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		CheckDestroy:             testAccCheckMackerelChannelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccMackerelChannelConfigWebhook(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "email.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "slack.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "webhook.0.url"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.events.#", "0"),
				),
			},
			// Test: Update
			{
				Config: testAccMackerelChannelConfigWebhookUpdated(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMackerelChannelExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceName, "email.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "slack.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "webhook.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "webhook.0.url"),
					resource.TestCheckResourceAttr(resourceName, "webhook.0.events.#", "6"),
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

func testAccCheckMackerelChannelDestroy(s *terraform.State) error {
	client := mackerelClient()
	for _, r := range s.RootModule().Resources {
		if r.Type != "mackerel_channel" {
			continue
		}
		channels, err := client.FindChannels()
		if err != nil {
			return err
		}
		for _, channel := range channels {
			if channel.ID == r.Primary.ID {
				return fmt.Errorf("channel still exists: %s", r.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckMackerelChannelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("channel not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no channel ID is set")
		}

		client := mackerelClient()
		channels, err := client.FindChannels()
		if err != nil {
			return err
		}

		for _, channel := range channels {
			if channel.ID == rs.Primary.ID {
				return nil
			}
		}

		return fmt.Errorf("channel not found from mackerel: %s", rs.Primary.ID)
	}
}

func testAccMackerelChannelConfigEmail(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "email" {
  name = "%s"
  email {}
}
`, name)
}

func testAccMackerelChannelConfigEmailUpdated(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "email" {
  name = "%s"
  email {
    emails = [
      "john.doe@example.test"]
    events = [
      "alert",
      "alertGroup"]
  }
}
`, name)
}

func testAccMackerelChannelConfigSlack(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "slack" {
  name = "%s"
  slack {
    url = "https://hooks.slack.com/services/xxx/yyy/zzz"
  }
}
`, name)
}

func testAccMackerelChannelConfigSlackUpdated(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "slack" {
  name = "%s"
  slack {
    url = "https://hooks.slack.com/services/xxx/yyy/zzz"
    mentions = {
      "ok" = "OK!!!"
      "warning" = "WARNING!!!"
      "critical" = "CRITICAL!!!"
    }
    enabled_graph_image = true
    events = [
      "alert",
      "alertGroup",
      "hostStatus",
      "hostRegister",
      "hostRetire",
      "monitor"]
  }
}
`, name)
}

func testAccMackerelChannelConfigWebhook(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "webhook" {
  name = "%s"
  webhook {
    url = "https://test.com/hook"
  }
}
`, name)
}

func testAccMackerelChannelConfigWebhookUpdated(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "webhook" {
  name = "%s"
  webhook {
    url = "https://test.com/hook"
    events = [
      "alert",
      "alertGroup",
      "hostStatus",
      "hostRegister",
      "hostRetire",
      "monitor"]
  }
}
`, name)
}
