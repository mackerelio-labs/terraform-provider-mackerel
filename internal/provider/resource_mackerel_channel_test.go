package provider_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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

func TestAccCompat_MackerelChannelResource_Email(t *testing.T) {
	t.Parallel()

	resourceName := "mackerel_channel.email"
	name := acctest.RandomWithPrefix("tf-channel-email-compat")
	config := `
resource "mackerel_channel" "email" {
  name = "` + name + `"
  email {}
}`

	resource.Test(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("email").AtSliceIndex(0).AtMapKey("events"),
						knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("email").AtSliceIndex(0).AtMapKey("emails"),
						knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("email").AtSliceIndex(0).AtMapKey("user_ids"),
						knownvalue.Null()),
				},
			},
			stepNoPlanInFramework(config),
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("email").AtSliceIndex(0).AtMapKey("events"),
						knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("email").AtSliceIndex(0).AtMapKey("emails"),
						knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("email").AtSliceIndex(0).AtMapKey("user_ids"),
						knownvalue.ListSizeExact(0)),
				},
			},
			stepNoPlanInFramework(config),
		},
	})
}

func TestAccCompat_MackerelChannelResource_Slack(t *testing.T) {
	t.Parallel()

	resourceName := "mackerel_channel.slack"
	name := acctest.RandomWithPrefix("tf-channel-slack-compat")

	config := `
resource "mackerel_channel" "slack" {
  name = "` + name + `"
  slack {
    url = "https://hooks.slack.com/services/xxx/yyy/zzz"
  }
}`

	resource.Test(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		Steps: []resource.TestStep{
			// Test: SDK
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("slack").AtSliceIndex(0).AtMapKey("events"),
						knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("slack").AtSliceIndex(0).AtMapKey("mentions"),
						knownvalue.Null()),
				},
			},
			stepNoPlanInFramework(config),
			// Test: SDK
			// Apply config twice to normalize the state.
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				// { "slack": { "events": [], "menstions": {} }}
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("slack").AtSliceIndex(0).AtMapKey("events"),
						knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("slack").AtSliceIndex(0).AtMapKey("mentions"),
						knownvalue.MapSizeExact(0)),
				},
			},
			stepNoPlanInFramework(config),
		},
	})
}

func TestAccCompat_MackerelChannelResource_Webhook(t *testing.T) {
	t.Parallel()

	resourceName := "mackerel_channel.webhook"
	name := acctest.RandomWithPrefix("tf-channel-webhook-compat")

	config := `
resource "mackerel_channel" "webhook" {
  name = "` + name + `"
  webhook {
    url = "https://example.test"
  }
}`

	resource.Test(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("webhook").AtSliceIndex(0).AtMapKey("events"),
						knownvalue.Null()),
				},
			},
			stepNoPlanInFramework(config),
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("webhook").AtSliceIndex(0).AtMapKey("events"),
						knownvalue.ListSizeExact(0)),
				},
			},
			stepNoPlanInFramework(config),
		},
	})
}
