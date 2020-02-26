package mackerel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMackerelChannel_Email(t *testing.T) {
	channelName := fmt.Sprintf("tf-channel-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelChannelTypeEmailConfig(channelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_channel.foo", "name", channelName)),
			},
			{
				Config: testAccCheckMackerelChannelTypeEmailConfigUpdated(channelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_channel.foo", "name", channelName)),
			},
		},
	})
}

func TestAccMackerelChannel_Slack(t *testing.T) {
	channelName := fmt.Sprintf("tf-channel-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelChannelTypeSlackConfig(channelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_channel.foo", "name", channelName)),
			},
			{
				Config: testAccCheckMackerelChannelTypeSlackConfigUpdated(channelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_channel.foo", "name", channelName)),
			},
		},
	})
}

func TestAccMackerelChannel_Webhook(t *testing.T) {
	channelName := fmt.Sprintf("tf-channel-%s", acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMackerelChannelTypeWebhookConfig(channelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_channel.foo", "name", channelName)),
			},
			{
				Config: testAccCheckMackerelChannelTypeWebhookConfigUpdated(channelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mackerel_channel.foo", "name", channelName)),
			},
		},
	})
}

func testAccCheckMackerelChannelTypeEmailConfig(channelName string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "%s"
  email {
    events = ["alert"]
  }
}
`, channelName)
}

func testAccCheckMackerelChannelTypeEmailConfigUpdated(channelName string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "%s"
  email {
    emails = ["main.xcezx+mackerel@gmail.com"]
    user_ids = []
    events = ["alert", "alertGroup"]
  }
}
`, channelName)
}

func testAccCheckMackerelChannelTypeSlackConfig(channelName string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "%s"
  slack {
    url = "https://example.test/"
  }
}
`, channelName)
}

func testAccCheckMackerelChannelTypeSlackConfigUpdated(channelName string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "%s"
  slack {
    url = "https://example.test/"
    mentions = {
      ok = "ok message"
      warning = "warning message"
      critical = "critical message"
    }
    enabled_graph_image = true
    events = ["alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"]
  }
}`, channelName)
}

func testAccCheckMackerelChannelTypeWebhookConfig(channelName string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "%s"
  webhook {
    url = "https://example.test/"
  }
}
`, channelName)
}

func testAccCheckMackerelChannelTypeWebhookConfigUpdated(channelName string) string {
	// language=HCL
	return fmt.Sprintf(`
resource "mackerel_channel" "foo" {
  name = "%s"
  webhook {
    url = "https://example.test/"
    events = ["alert", "alertGroup", "hostStatus", "hostRegister", "hostRetire", "monitor"]
  }
}
`, channelName)
}
