package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelDashboardResource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwresource.SchemaRequest{}
	resp := fwresource.SchemaResponse{}
	if provider.NewMackerelDashboardResource().Schema(ctx, req, &resp); resp.Diagnostics.HasError() {
		t.Fatalf("schema method diagnostics: %+v", resp.Diagnostics)
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema validation diagnostics: %+v", diags)
	}
}

func TestAccMackerelDashboardGraphWithoutRange(t *testing.T) {
	resourceName := "mackerel_dashboard.graph"
	rand := acctest.RandString(5)
	title := fmt.Sprintf("tf-dashboard graph %s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: `
resource "mackerel_service" "include" {
  name = "tf-service-` + rand + `-include"
}

resource "mackerel_role" "include" {
  service = mackerel_service.include.name
  name    = "tf-role-` + rand + `-include"
}

resource "mackerel_dashboard" "graph" {
  title = "` + title + `"
  url_path = "` + rand + `"
  graph {
    title = "test graph role"
    role {
      role_fullname = "${mackerel_service.include.name}:${mackerel_role.include.name}"
      name = "loadavg5"
      is_stacked = true
    }
    layout {
      x = 2
      y = 12
      width = 10
      height = 8
    }
  }
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "url_path", rand),
					resource.TestCheckResourceAttr(resourceName, "graph.#", "1"),
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
