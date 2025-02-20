package provider_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelDowntimeResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	provider.NewMackerelDowntimeResource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccCompat_MackerelDowntimeResource(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		config func(name string) string
	}{
		"undefined": {
			config: func(name string) string {
				return `
resource "mackerel_downtime" "foo" {
  name = "` + name + `"
  start = 2051254800
  duration = 3600
  recurrence {
    type = "weekly"
    interval = 2
  }
}`
			},
		},
		"null": {
			config: func(name string) string {
				return `
resource "mackerel_downtime" "foo" {
  name = "` + name + `"
  start = 2051254800
  duration = 3600
  service_scopes = null
  service_exclude_scopes = null
  role_scopes = null
  role_exclude_scopes = null
  monitor_scopes = null
  monitor_exclude_scopes = null
  recurrence {
    type = "weekly"
    interval = 2
    weekdays = null
  }
}`
			},
		},
		"empty": {
			config: func(name string) string {
				return `
resource "mackerel_downtime" "foo" {
  name = "` + name + `"
  start = 2051254800
  duration = 3600
  service_scopes = []
  service_exclude_scopes = []
  role_scopes = []
  role_exclude_scopes = []
  monitor_scopes = []
  monitor_exclude_scopes = []
  recurrence {
    type = "weekly"
    interval = 2
    weekdays = []
  }
}`
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			name := acctest.RandomWithPrefix("tf-compat-downtime")
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
