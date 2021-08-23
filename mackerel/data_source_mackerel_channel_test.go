package mackerel

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMackerelChannelEmail(t *testing.T) {
	dsName := "data.mackerel_channel.foo"
	name := fmt.Sprintf("tf-channel-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelChannelConfigEmail(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "email.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "email.0.emails.#", "1"),
						resource.TestCheckResourceAttr(dsName, "email.0.user_ids.#", "0"),
						resource.TestCheckResourceAttr(dsName, "email.0.events.#", "2"),
					),
					resource.TestCheckResourceAttr(dsName, "slack.#", "0"),
					resource.TestCheckResourceAttr(dsName, "webhook.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelChannelSlack(t *testing.T) {
	dsName := "data.mackerel_channel.foo"
	name := fmt.Sprintf("tf-channel-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelChannelConfigSlack(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "email.#", "0"),
					resource.TestCheckResourceAttr(dsName, "slack.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "slack.0.url", "https://hooks.slack.com/services/xxx/yyy/zzz"),
						resource.TestCheckResourceAttr(dsName, "slack.0.mentions.ok", "OK!!!"),
						resource.TestCheckResourceAttr(dsName, "slack.0.mentions.warning", "WARNING!!!"),
						resource.TestCheckResourceAttr(dsName, "slack.0.mentions.critical", "CRITICAL!!!"),
						resource.TestCheckResourceAttr(dsName, "slack.0.enabled_graph_image", "true"),
						resource.TestCheckResourceAttr(dsName, "slack.0.events.#", "6"),
					),
					resource.TestCheckResourceAttr(dsName, "webhook.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceMackerelChannelWebhook(t *testing.T) {
	dsName := "data.mackerel_channel.foo"
	name := fmt.Sprintf("tf-channel-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelChannelConfigWebhook(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "name", name),
					resource.TestCheckResourceAttr(dsName, "email.#", "0"),
					resource.TestCheckResourceAttr(dsName, "slack.#", "0"),
					resource.TestCheckResourceAttr(dsName, "webhook.#", "1"),
					resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dsName, "webhook.0.url", "https://test.com/hook"),
						resource.TestCheckResourceAttr(dsName, "webhook.0.events.#", "6"),
					),
				),
			},
		},
	})
}

func testAccDataSourceMackerelChannelConfigEmail(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "%s"
  email {
    emails = ["john.doe@example.test"]
    events = ["alert", "alertGroup"]
  }
}

data "mackerel_channel" "foo" {
  id = mackerel_channel.foo.id
}
`, name)
}

func testAccDataSourceMackerelChannelConfigSlack(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
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

data "mackerel_channel" "foo" {
  id = mackerel_channel.foo.id
}
`, name)
}

func testAccDataSourceMackerelChannelConfigWebhook(name string) string {
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
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

data "mackerel_channel" "foo" {
  id = mackerel_channel.foo.id
}
`, name)
}

func TestAccDataSourceMackerelChannelNotMatchAnyChannel(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `data "mackerel_channel" "foo" { id = "not-found" }`,
				ExpectError: regexp.MustCompile(`the ID 'not-found' does not match any channel in mackerel\.io`),
			},
		},
	})
}
