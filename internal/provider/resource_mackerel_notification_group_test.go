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

func Test_MackerelNotificationGroupResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	provider.NewMackerelNotificationGroupResource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnotstics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnotstics: %+v", diags)
	}
}

func TestAccCompat_MackerelNotificationGroupResource(t *testing.T) {
	t.Parallel()

	resouceName := "mackerel_notification_group.default"
	name := acctest.RandomWithPrefix("tf-notification-group-compat")

	config := `
resource "mackerel_notification_group" "default" {
  name = "` + name + `"
  child_channel_ids = null
  child_notification_group_ids = []
}`

	resource.Test(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		Steps: []resource.TestStep{
			// Test: SDK
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resouceName,
						tfjsonpath.New("child_notification_group_ids"),
						knownvalue.Null()),
					statecheck.ExpectKnownValue(resouceName,
						tfjsonpath.New("child_channel_ids"),
						knownvalue.Null()),
				},
			},
			stepNoPlanInFramework(config),
			// Test: SDK
			{
				Config:                   config,
				ProtoV5ProviderFactories: protoV5SDKProviderFactories,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resouceName,
						tfjsonpath.New("child_notification_group_ids"),
						knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(resouceName,
						tfjsonpath.New("child_channel_ids"),
						knownvalue.ListSizeExact(0)),
				},
			},
			stepNoPlanInFramework(config),
		},
	})
}
