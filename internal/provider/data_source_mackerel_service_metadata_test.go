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

func Test_MackerelServiceMetadataDataSource_schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	req := fwdatasource.SchemaRequest{}
	resp := &fwdatasource.SchemaResponse{}
	provider.NewMackerelServiceMetadataDataSource().Schema(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("schema method: %+v", resp.Diagnostics)
		return
	}

	if diags := resp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Errorf("schema validation: %+v", diags)
	}
}

func TestAccDataSourceMackerelServiceMetadata(t *testing.T) {
	dsName := "data.mackerel_service_metadata.foo"
	rand := acctest.RandString(5)
	serviceName := fmt.Sprintf("tf-service-%s", rand)
	namespace := fmt.Sprintf("tf-namespace-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV5ProviderFactories: protoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMackerelServiceMetadataConfig(serviceName, namespace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "service", serviceName),
					resource.TestCheckResourceAttr(dsName, "namespace", namespace),
					resource.TestCheckResourceAttr(dsName, "metadata_json", `{"id":1}`),
				),
			},
		},
	})
}

func testAccDataSourceMackerelServiceMetadataConfig(serviceName, namespace string) string {
	return fmt.Sprintf(`
resource "mackerel_service" "foo" {
  name = "%s"
}

resource "mackerel_service_metadata" "foo" {
  service = mackerel_service.foo.name
  namespace = "%s"
  metadata_json = jsonencode({
    id = 1
  })
}

data "mackerel_service_metadata" "foo" {
  depends_on = [mackerel_service_metadata.foo]
  service = mackerel_service_metadata.foo.service
  namespace = mackerel_service_metadata.foo.namespace
}
`, serviceName, namespace)
}
