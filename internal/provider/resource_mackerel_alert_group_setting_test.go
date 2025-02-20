package provider_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
