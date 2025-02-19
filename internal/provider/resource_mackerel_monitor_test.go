package provider_test

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelMonitorResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	provider.NewMackerelMonitorResource().Schema(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccCompat_MackerelMonitorResource(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		config func(name string) string
	}{
		"host_metric/undefined": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "host_metric" {
  name = "` + name + `"
  host_metric {
    metric = "cpu.sys"
    operator = ">"
    warning = "80"
    duration = 3
  }
}`
			},
		},
		"host_metric/null": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "host_metric" {
  name = "` + name + `"
  host_metric {
    metric = "cpu.sys"
    operator = ">"
    warning = "80"
    duration = 3
    scopes = null
    exclude_scopes = null
  }
}`
			},
		},
		"host_metric/empty": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "host_metric" {
  name = "` + name + `"
  host_metric {
    metric = "cpu.sys"
    operator = ">"
    warning = "80"
    duration = 3
    scopes = []
    exclude_scopes = []
  }
}`
			},
		},
		"connectivity/undefined": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "connectivity" {
  name = "` + name + `"
  connectivity {}
}`
			},
		},
		"connectivity/null": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "connectivity" {
  name = "` + name + `"
  connectivity {
    scopes = null
    exclude_scopes = null
  }
}`
			},
		},
		"connectivity/empty": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "connectivity" {
  name = "` + name + `"
  connectivity {
    scopes = []
	exclude_scopes = []
  }
}`
			},
		},
		"service_metric": {
			config: func(name string) string {
				return `
resource "mackerel_service" "svc" {
  name = "` + name + `-svc"
}
resource "mackerel_monitor" "service_metric" {
  name = "` + name + `"
  service_metric {
    service = mackerel_service.svc.name
    duration = 1	
    metric = "custom.access.2xx_ratio"
    operator = ">"
    warning = "99.9"
  }
}`
			},
		},
		"external": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "external" {
  name = "` + name + `"
  external {
    method = "GET"
    url = "https://terraform-provider-mackerel.test/"
  }
}`
			},
		},
		"external/null": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "external" {
  name = "` + name + `"
  external {
    method = "GET"
    url = "https://terraform-provider-mackerel.test/"
    headers = null
  }
}`
			},
		},
		"external/empty": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "external" {
  name = "` + name + `"
  external {
    method = "GET"
    url = "https://terraform-provider-mackerel.test/"
    headers = {}
  }
}`
			},
		},
		"expression": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "expression" {
  name = "` + name + `"
  expression {
    expression = "max(role(my-service:db, loadavg5))"
    operator = ">"
    warning = "0.7"
  }
}`
			},
		},
		"anomaly_detection": {
			config: func(name string) string {
				return `
resource "mackerel_service" "svc" {
  name = "` + name + `-svc"
}
resource "mackerel_role" "role" {
  service = mackerel_service.svc.name
  name = "` + name + `-role"
}
resource "mackerel_monitor" "anomaly_detection" {
  name = "` + name + `"
  anomaly_detection {
    warning_sensitivity = "insensitive"
    scopes = [mackerel_role.role.id]
  }
}`
			},
		},
		"query": {
			config: func(name string) string {
				return `
resource "mackerel_monitor" "foo" {
  name = "` + name + `"
  query {
    query = "container.cpu.utilization{k8s.deployment.name=\"httpbin\"}"
    legend = "cpu.utilization {{k8s.node.name}}"
    operator = ">"
    warning = "70"
  }
}`
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			name := acctest.RandomWithPrefix("tf-monitor-compat")
			config := tt.config(name)

			resource.Test(t, resource.TestCase{
				PreCheck: func() { preCheck(t) },
				Steps: []resource.TestStep{
					{
						Config:                   config,
						ProtoV5ProviderFactories: protoV5SDKProviderFactories,
					},
					stepNoPlanInFramework(config),
					{
						Config:                   config,
						ProtoV5ProviderFactories: protoV5SDKProviderFactories,
					},
					stepNoPlanInFramework(config),
				},
			})
		})
	}
}
