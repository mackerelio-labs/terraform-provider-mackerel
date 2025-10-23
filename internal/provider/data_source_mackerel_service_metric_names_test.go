package provider_test

import (
	"context"
	"fmt"
	"testing"

	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mackerelio-labs/terraform-provider-mackerel/internal/provider"
)

func Test_MackerelServiceMetricNamesDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := &fwdatasource.SchemaResponse{}
	provider.NewMackerelServiceMetricNamesDataSource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("schema method: %+v", resp.Diagnostics)
		return
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("schema validation: %+v", diags)
	}
}

func TestAccDataSourceMackerelServiceMetricNames(t *testing.T) {
	name := fmt.Sprintf("tf-service-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelServiceMetricNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mackerel_service_metric_names.foo", "id", name+":"),
					resource.TestCheckResourceAttr("data.mackerel_service_metric_names.foo", "name", name),
				),
			},
		},
	})
}

func testAccDataSourceMackerelServiceMetricNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
  memo = "This service is managed by Terraform."
}

data "mackerel_service_metric_names" "foo" {
  name = mackerel_service.foo.name
}
`, name)
}
